package mongo

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
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
			Database:    "testdb",
			Collection:  "testcol",
			BatchSize:   15 * 3,
			Iterations:  1000,
			CurrentHash: "ok",
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

	setValidView()
	setSubView()

	// insert some fake entries for testing
	err = insertFakeEntries()
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())

	// run tests
	code := m.Run()

	log.Println("test done - starting clean up")

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// disconnect mongodb client
	if err = s.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}

	log.Println("clean up done")

	os.Exit(code)
}

func Ptr[T any](v T) *T {
	return &v
}

func setValidView() {
	res := s.client.Database("testdb").RunCommand(
		context.Background(),
		bson.D{
			{Key: "create", Value: "testview"},
			{Key: "viewOn", Value: "testcol"},
			{
				Key: "pipeline",
				Value: bson.A{
					bson.D{
						{
							Key: "$match",
							Value: bson.D{
								{
									Key:   "is_db_valid",
									Value: true,
								},
							},
						},
					},
				},
			},
		},
	)
	if res.Err() != nil {
		log.Fatalf("creating view failed: %v", res)
	}
	s.cfg.ValidView = "testview"
}

func setSubView() {
	res := s.client.Database("testdb").RunCommand(
		context.Background(),
		bson.D{
			{Key: "create", Value: "testsubview"},
			{Key: "viewOn", Value: "testcol"},
			{
				Key: "pipeline",
				Value: bson.A{
					bson.D{
						{
							Key: "$match",
							Value: bson.D{
								{
									Key: "summary",
									Value: bson.D{
										{
											Key:   "$exists",
											Value: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	)
	if res.Err() != nil {
		log.Fatalf("creating view failed: %v", res)
	}
	s.cfg.SubView = "testsubview"
}

func makeEntry(id, src string, d, t bool) *db.Entry {
	e := &db.Entry{
		Id:        id,
		Config:    RandStringBytesMaskImpr(20),
		Submitter: src,
	}
	if d || t {
		e.Summary = &db.EntrySummary{
			TargetCount: 1,
		}
		if rand.Intn(3) == 1 {
			e.Hash = "ok"
		}
	}
	if t {
		e.AcceptedTags = append(e.AcceptedTags, model.DBTag_DB_TAG_GCSIM)
		e.IsDbValid = true
	}
	return e
}

var dbEntries map[string]*db.Entry = make(map[string]*db.Entry)
var subs map[string]*db.Entry = make(map[string]*db.Entry)
var dbNoTag map[string]*db.Entry = make(map[string]*db.Entry)

func insertFakeEntries() error {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	// entries with data no tag
	for i := 0; i < rand.Intn(10)+5; i++ {
		id := fmt.Sprintf("sample_db_no_tag_%v", i)
		e := makeEntry(id, "notag", true, false)
		_, err := col.InsertOne(context.TODO(), e)
		if err != nil {
			log.Fatal(err)
		}
		dbNoTag[id] = e
	}
	// entries with tag
	for i := 0; i < rand.Intn(10)+5; i++ {
		id := fmt.Sprintf("sample_db_approved_%v", i)
		e := makeEntry(id, "tag", true, true)
		_, err := col.InsertOne(context.TODO(), e)
		if err != nil {
			log.Fatal(err)
		}
		dbEntries[id] = e
	}
	// entries without
	for i := 0; i < rand.Intn(10)+5; i++ {
		id := fmt.Sprintf("sample_sub_only_%v", i)
		e := makeEntry(id, "sub", false, false)
		_, err := col.InsertOne(context.TODO(), e)
		if err != nil {
			log.Fatal(err)
		}
		subs[id] = e
	}

	return nil

}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandStringBytesMaskImpr(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func compareConfig(t *testing.T, expect, got *db.Entry) {
	if expect.Config != got.Config {
		t.Errorf("expecting config %v (id = %v), got %v (id = %v)", expect.Config, expect.Id, got.Config, got.Id)
	}
}

func compareHash(t *testing.T, expect string, got *db.Entry) {
	if expect != got.Hash {
		t.Errorf("expecting config %v, got %v (entry = %v)", expect, got.Hash, got.String())
	}
}
