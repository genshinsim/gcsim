package mongo

import (
	"context"
	"strings"
	"testing"
)

func TestGetWork(t *testing.T) {
	work, err := s.GetWork(context.Background())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	for _, v := range work {
		//we only care about results with sample prefix
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
