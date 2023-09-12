package mongo

import (
	"context"
	"log"
	"strings"
	"testing"

	"github.com/genshinsim/gcsim/backend/pkg/services/db"
	"github.com/genshinsim/gcsim/pkg/model"
)

func TestGetWork(t *testing.T) {
	work, err := s.GetWork(context.Background())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, v := range work {
		// we only care about results with sample prefix
		switch {
		case strings.HasPrefix(v.Id, "sample_db_no_tag"):
			compareHash(t, "", dbNoTag[v.Id])
		case strings.HasPrefix(v.Id, "sample_db_approved"):
			compareHash(t, "", dbEntries[v.Id])
		case strings.HasPrefix(v.Id, "sample_sub_only"):
			compareHash(t, "", subs[v.Id])
		}
	}
}

func TestGetAllEntriesWithoutTag(t *testing.T) {
	res, err := s.GetAllEntriesWithoutTag(context.Background(), model.DBTag_DB_TAG_GCSIM, &db.QueryOpt{})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	for _, v := range res {
		log.Println(v.String())
	}
}
