package share

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestMain(m *testing.M) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	x, err := newMock()
	if err != nil {
		log.Fatal(err)
	}

	srv, err := New(Config{
		Store: x,
	})
	if err != nil {
		log.Fatal(err)
	}

	RegisterShareStoreServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	code := m.Run()
	os.Exit(code)
}

func TestCreate(t *testing.T) {
	c := &Client{}
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	c.srvClient = NewShareStoreClient(conn)

	n := time.Now().Unix() + 100

	id, err := c.Create(context.TODO(), &model.SimulationResult{
		Config: "blah",
	}, uint64(n), "tester")
	if err != nil {
		t.Error(err)
	}

	// expecting expiry to be <= n + 2 (assuming test takes 1s to run; shouldnt)

	res, expiry, err := c.Read(context.TODO(), id)
	if err != nil {
		t.Error(err)
	}

	if expiry != uint64(n) {
		t.Errorf("expiry expected to be %v, got %v", n+2, expiry)

	}

	if res.GetConfig() != "blah" {
		t.Errorf("expecting config to be blah, got %v", res.GetConfig())
	}

	log.Println(res, expiry)

}

func TestReplace(t *testing.T) {
	c := &Client{}
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	c.srvClient = NewShareStoreClient(conn)

	res := &model.SimulationResult{
		Config: "blah",
	}

	id, err := c.Create(context.TODO(), res, 0, "tester")
	if err != nil {
		t.Error(err)
	}

	res.Config = "boo"
	err = c.Replace(context.TODO(), id, res)
	if err != nil {
		t.Error(err)
	}

	res, expiry, err := c.Read(context.TODO(), id)
	if err != nil {
		t.Error(err)
	}

	if expiry != 0 {
		t.Errorf("expiry expected to be %v, got %v", 0, expiry)
	}

	if res.GetConfig() != "boo" {
		t.Errorf("expecting config to be boo, got %v", res.GetConfig())
	}

	log.Println(res, expiry)

}

func TestDelete(t *testing.T) {
	c := &Client{}
	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	c.srvClient = NewShareStoreClient(conn)

	res := &model.SimulationResult{
		Config: "blah",
	}

	id, err := c.Create(context.TODO(), res, 0, "tester")
	if err != nil {
		t.Error(err)
	}

	err = c.Delete(context.TODO(), id)
	if err != nil {
		t.Error(err)
	}

	_, _, err = c.Read(context.TODO(), id)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expecting status error, got %v", err)
	}
	if st.Code() != codes.NotFound {
		t.Fatalf("expecting not found error, got %v", st.Code())
	}
}
