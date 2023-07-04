package mongo

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
)

func TestApproveTag(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	//insert new entry so not to pollute other tests
	e := makeEntry("approve_test", "poop", true, false)
	_, err := col.InsertOne(context.Background(), e)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = s.ApproveTag(context.Background(), e.Id, model.DBTag_DB_TAG_TESTING)
	if err != nil {
		t.Error(err)
	}

	res, err := s.getOne(context.Background(), col, bson.M{"_id": e.Id})
	if err != nil {
		t.Error(err)
	}

	count := 0
	for _, v := range res.AcceptedTags {
		if v == model.DBTag_DB_TAG_TESTING {
			count++
		}
	}

	if count != 1 {
		t.Errorf("could not find 1 count of testing tag in result: %v", res.String())
	}

	if !res.IsDbValid {
		t.Error("result should be db valid")
	}

	//try adding same tag again
	err = s.ApproveTag(context.Background(), e.Id, model.DBTag_DB_TAG_TESTING)
	if err != nil {
		t.Error(err)
	}

	res, err = s.getOne(context.Background(), col, bson.M{"_id": e.Id})
	if err != nil {
		t.Error(err)
	}

	count = 0
	for _, v := range res.AcceptedTags {
		if v == model.DBTag_DB_TAG_TESTING {
			count++
		}
	}

	if count != 1 {
		t.Errorf("could not find 1 count of testing tag in result: %v", res.String())
	}

	if !res.IsDbValid {
		t.Error("result should be db valid")
	}
}

func TestRejectTag(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	//insert new entry so not to pollute other tests
	e := makeEntry("deny_test", "poop", true, true)
	e.AcceptedTags = append(e.AcceptedTags, model.DBTag_DB_TAG_TESTING)
	_, err := col.InsertOne(context.Background(), e)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = s.RejectTag(context.Background(), e.Id, model.DBTag_DB_TAG_GCSIM)
	if err != nil {
		t.Error(err)
	}

	res, err := s.getOne(context.Background(), col, bson.M{"_id": e.Id})
	if err != nil {
		t.Error(err)
	}

	count := 0
	for _, v := range res.AcceptedTags {
		if v == model.DBTag_DB_TAG_GCSIM {
			count++
		}
	}

	if count != 0 {
		t.Errorf("tag count not 0 in result: %v", res.String())
	}

	//there's still one more tag to remove
	if !res.IsDbValid {
		t.Error("result should db valid still")
	}

	err = s.RejectTag(context.Background(), e.Id, model.DBTag_DB_TAG_TESTING)
	if err != nil {
		t.Error(err)
	}

	res, err = s.getOne(context.Background(), col, bson.M{"_id": e.Id})
	if err != nil {
		t.Error(err)
	}

	count = 0
	for _, v := range res.AcceptedTags {
		if v == model.DBTag_DB_TAG_TESTING {
			count++
		}
	}

	if count != 0 {
		t.Errorf("tag count not 0 in result: %v", res.String())
	}

	//there's still one more tag to remove
	if res.IsDbValid {
		t.Error("result should not be db valid")
	}
}

func TestRejectAllUnapprovedTag(t *testing.T) {
	col := s.client.Database(s.cfg.Database).Collection(s.cfg.Collection)

	//insert new entries so not to pollute other tests
	e := makeEntry("reject_all_unapproved", "poop", true, true)
	e.AcceptedTags = append(e.AcceptedTags, model.DBTag_DB_TAG_TESTING)
	_, err := col.InsertOne(context.Background(), e)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	total := rand.Intn(25)
	//add a bunch of random entries
	for i := 0; i < total; i++ {
		e := makeEntry(fmt.Sprintf("reject_all_unapproved_%v", i), "poop", true, false)
		_, err := col.InsertOne(context.Background(), e)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}

	_, err = s.RejectTagAllUnapproved(context.Background(), model.DBTag_DB_TAG_TESTING)
	if err != nil {
		t.Error(err)
	}

	//we should have total of reject_all_unapproved with testing as rejected tag
	results, err := s.get(context.Background(), col, bson.D{})
	if err != nil {
		t.Error(err)
	}
	rejectCount := 0
	acceptCount := 0
	for _, v := range results {
		//we only care about reject_all_unapproved
		if strings.HasPrefix(v.Id, "reject_all_unapproved") {
			for _, x := range v.RejectedTags {
				if x == model.DBTag_DB_TAG_TESTING {
					rejectCount++
				}
			}
			for _, x := range v.AcceptedTags {
				if x == model.DBTag_DB_TAG_TESTING {
					acceptCount++
				}
			}
		}
	}

	if acceptCount != 1 {
		t.Errorf("expecting %v accepted, got %v", 1, acceptCount)
	}

	if rejectCount != total {
		t.Errorf("expecting %v rejected, got %v", total, rejectCount)
	}
}
