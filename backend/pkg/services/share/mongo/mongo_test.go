package mongo

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/services/share"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var s *Server

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "5.0",
		Env: []string{
			// username and password for mongodb superuser
			"MONGO_INITDB_ROOT_USERNAME=root",
			"MONGO_INITDB_ROOT_PASSWORD=password",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	s = &Server{
		cfg: Config{
			Database:   "testdb",
			Collection: "testcol",
		},
	}
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("creating logger failed: %s", err)
	}
	sugar := logger.Sugar()
	sugar.Debugw("logger initiated")

	s.Log = sugar

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		s.client, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://root:password@localhost:%s", resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		return s.client.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// run tests
	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err = s.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func TestReadWrite(t *testing.T) {
	share := &share.ShareEntry{
		Result: &model.SimulationResult{
			Config: "blah",
		},
		Submitter: "tester",
	}

	id, err := s.Create(context.TODO(), share)
	if err != nil {
		t.Fatal(err)
	}

	res, err := s.Read(context.TODO(), id)
	if err != nil {
		t.Fatal(err)
	}

	if res.Result.Config != "blah" {
		t.Errorf("expecting config to equal blah, got %v", res.Result.Config)
	}

	if res.Submitter != "tester" {
		t.Errorf("expecting tester for submitter, got %v", res.Submitter)
	}
}

func TestUpdate(t *testing.T) {
	share := &share.ShareEntry{
		Result: &model.SimulationResult{
			Config: "blah",
		},
		Submitter: "tester",
	}

	id, err := s.Create(context.TODO(), share)
	if err != nil {
		t.Fatal(err)
	}

	share.Id = id
	share.Result.Config = "poop"

	_, err = s.Update(context.TODO(), share)
	if err != nil {
		t.Fatal(err)
	}

	res, err := s.Read(context.TODO(), id)
	if err != nil {
		t.Fatal(err)
	}

	if res.Result.Config != "poop" {
		t.Errorf("expecting config to equal poop, got %v", res.Result.Config)
	}

	if res.Submitter != "tester" {
		t.Errorf("expecting tester for submitter, got %v", res.Submitter)
	}
}

func TestDelete(t *testing.T) {
	share := &share.ShareEntry{
		Result: &model.SimulationResult{
			Config: "blah",
		},
		Submitter: "tester",
	}

	id, err := s.Create(context.TODO(), share)
	if err != nil {
		t.Fatal(err)
	}

	err = s.Delete(context.TODO(), id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Read(context.TODO(), id)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expecting status error, got %v", err)
	}
	if st.Code() != codes.NotFound {
		t.Fatalf("expecting not found error, got %v", st.Code())
	}
}

func TestRandom(t *testing.T) {
	share := &share.ShareEntry{
		Result: &model.SimulationResult{
			Config: "blah",
		},
		Submitter: "tester",
	}

	_, err := s.Create(context.TODO(), share)
	if err != nil {
		t.Fatal(err)
	}

	id, err := s.Random(context.TODO())
	if err != nil {
		t.Error(err)
	}

	log.Println(id)
}
