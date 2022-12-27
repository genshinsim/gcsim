package mongo

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestCreate(t *testing.T) {
	const colstr = "test_collection"
	const dbstr = "test"
	srv, err := NewServer(Config{
		URL:        "mongodb://localhost:27017",
		Database:   dbstr,
		Collection: colstr,
		Username:   "root",
		Password:   "example",
	})

	if err != nil {
		t.Error(err)
	}

	//clean up everything first
	col := srv.client.Database(dbstr).Collection(colstr)
	res, err := col.DeleteMany(context.TODO(), bson.D{})

	if err != nil {
		t.Error(err)
	}

	srv.Log.Infow("delete done", "deleted_count", res.DeletedCount)

	e := &model.DBEntry{
		Key:        "blah",
		CreateDate: uint64(time.Now().Unix()),
		RunDate:    uint64(time.Now().Unix()),
		SimDuration: &model.DescriptiveStats{
			Min:  0,
			Max:  100,
			Mean: 50,
		},
		Config: "blah",
		Hash:   "blah",
	}

	_, err = srv.Create(context.TODO(), e)
	if err != nil {
		t.Error(err)
	}

	var qs = &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"key": {
				Kind: &structpb.Value_StringValue{
					StringValue: e.Key,
				},
			},
		},
	}

	entries, err := srv.Get(context.TODO(), qs, 30, 1)
	if err != nil {
		t.Error(err)
	}

	log.Println(entries)

}
