package result

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	os.RemoveAll("./testdb")

	store, err := New(Config{
		DBPath: "./testdb",
	})
	if err != nil {
		panic(err)
	}

	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()

	RegisterResultStoreServer(s, store)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestResultStore(t *testing.T) {

	client := &Client{}

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client.srvClient = NewResultStoreClient(conn)

	uuid, err := client.Create([]byte("test"), context.TODO())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	data, _, err := client.Read(uuid, context.TODO())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(data) != "test" {
		t.Errorf("expecting result to be test, got  %v", data)
	}

	err = client.Update(uuid, []byte("next"), context.TODO())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	data, _, err = client.Read(uuid, context.TODO())

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(data) != "next" {
		t.Errorf("expecting result to be next, got  %v", data)
	}

	err = client.Delete(uuid, context.TODO())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	_, _, err = client.Read(uuid, context.TODO())

	if err != api.ErrKeyNotFound {
		t.Errorf("expecting key to be gone, but got err %v", err)
		t.FailNow()
	}

}
