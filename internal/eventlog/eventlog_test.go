package eventlog

import (
	"testing"

	"github.com/genshinsim/gcsim/pkg/core"
	easyjson "github.com/mailru/easyjson"
)

func TestEventWriteKeyOnlyPanic(t *testing.T) {
	e := &Event{
		Msg:     "test",
		F:       1,
		Typ:     core.LogCharacterEvent,
		SrcChar: 0,
		Logs:    map[string]interface{}{},
	}
	//test writing
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	//this should panic
	e.Write("keyonly")

}

func TestEventWriteNonStringKeyPanic(t *testing.T) {
	e := &Event{
		Msg:     "test",
		F:       1,
		Typ:     core.LogCharacterEvent,
		SrcChar: 0,
		Logs:    map[string]interface{}{},
	}
	//test writing
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	//this should panic
	e.Write(1)

}

func TestEventWriteKeyVal(t *testing.T) {
	e := &Event{
		Msg:     "test",
		F:       1,
		Typ:     core.LogCharacterEvent,
		SrcChar: 0,
		Logs:    map[string]interface{}{},
	}

	//this should be ok no panic
	e.Write("stuff", 1, "goes", true, "here", "two")

}

func BenchmarkEasyJSONSerialization(b *testing.B) {
	//generate roughly 2 lines of debug per frame over 90s
	//each line should be roughly 10 fields
	//so that's 10800 events
	count := 10800
	var testdata EventArr
	testdata = make([]*Event, 0, count)
	for i := 0; i < count; i++ {
		e := &Event{
			Msg:     "test",
			F:       1,
			Typ:     core.LogCharacterEvent,
			SrcChar: 0,
			Logs:    map[string]interface{}{},
		}
		e.Write("a", 1, "b", true, "c", "stuff", "e", 123, "f", "boo", "g", 111)
		testdata = append(testdata, e)
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		easyjson.Marshal(testdata)
	}

}
