package mongo

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInsert(t *testing.T) {
	e := &db.Entry{
		Config:      "blah",
		Description: "test",
		Submitter:   "poop",
		Hash:        "ok",
		Summary: &db.EntrySummary{
			TargetCount: 2,
		},
	}

	id, err := s.Create(context.Background(), e)
	if err != nil {
		t.Error(err)
	}

	if id == "" {
		t.Error("id shouldn't be blank")
	}
	log.Println(id)
}

func TestGet(t *testing.T) {
	var e *db.Entry
	for _, v := range dbEntries {
		e = v
		break
	}

	if e == nil {
		t.Fatalf("entries shouldnt be empty")
	}

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	results, err := s.get(context.Background(), col, bson.M{"_id": e.Id})
	if err != nil {
		t.Error(err)
	}

	if len(results) != 1 {
		t.Errorf("expected 1 result, got %v", len(results))
	}

	res := results[0]

	if res.Config != e.Config {
		t.Errorf("expecting config to be %v, got %v", e.Config, res.Config)
	}
}

func TestGetOne(t *testing.T) {
	var e *db.Entry
	for _, v := range dbEntries {
		e = v
		break
	}

	if e == nil {
		t.Fatalf("entries shouldnt be empty")
	}

	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	res, err := s.getOne(context.Background(), col, bson.M{"_id": e.Id})
	if err != nil {
		t.Error(err)
	}

	if res.Config != e.Config {
		t.Errorf("expecting config to be %v, got %v", e.Config, res.Config)
	}
}

func TestGetMany(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	results, err := s.get(context.Background(), col, bson.D{})
	if err != nil {
		t.Error(err)
	}
	for _, v := range results {
		// we only care about results with sample prefix
		switch {
		case strings.HasPrefix(v.Id, "sample_db_no_tag"):
			compareConfig(t, dbNoTag[v.Id], v)
		case strings.HasPrefix(v.Id, "sample_db_approved"):
			compareConfig(t, dbEntries[v.Id], v)
		case strings.HasPrefix(v.Id, "sample_sub_only"):
			compareConfig(t, subs[v.Id], v)
		}
	}
}

func TestGetValid(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.ValidView)
	results, err := s.get(context.Background(), col, bson.M{})
	if err != nil {
		t.Error(err)
	}
	if len(results) != len(dbEntries) {
		t.Errorf("expecting %v entries, got %v", len(dbEntries), len(results))
		t.FailNow()
	}
	for _, v := range results {
		compareConfig(t, dbEntries[v.Id], v)
	}
}

func TestGetSubs(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.SubView)
	results, err := s.get(context.Background(), col, bson.M{})
	if err != nil {
		t.Error(err)
	}
	if len(results) != len(subs) {
		t.Errorf("expecting %v entries, got %v", len(subs), len(results))
		t.FailNow()
	}
	for _, v := range results {
		compareConfig(t, subs[v.Id], v)
	}
}

func TestDelete(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	// insert new entry so not to pollute other tests
	e := makeEntry("delete_test", "poop", true, false)
	_, err := col.InsertOne(context.Background(), e)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = s.Delete(context.Background(), e.Id)
	if err != nil {
		t.Error(err)
	}
}

func TestReplace(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)
	// insert new entry so not to pollute other tests
	e := makeEntry("update_test", "poop", true, false)
	_, err := col.InsertOne(context.Background(), e)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	e.Config = "cool stuff"

	err = s.Replace(context.Background(), e)
	if err != nil {
		t.Error(err)
	}

	res, err := s.getOne(context.Background(), col, bson.M{"_id": e.Id})
	if err != nil {
		t.Error(err)
	}

	if res.Config != "cool stuff" {
		t.Error("update failed; config incorrect")
	}

}
