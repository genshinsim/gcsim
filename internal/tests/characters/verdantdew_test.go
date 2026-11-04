package characters

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func TestSetVerdantDewLog(t *testing.T) {
	for _, tc := range []struct {
		amount int
		want   int
	}{
		{amount: 2, want: 2},
		{amount: -1, want: 0},
		{amount: 10, want: 3},
	} {
		t.Run(fmt.Sprintf("amount=%d", tc.amount), func(t *testing.T) {
			c, _ := makeCore(1)
			c.Player.SetMoonridgeDew(1)
			c.Player.SetVerdantDew(tc.amount)
			if got := c.Player.VerdantDew(); got != tc.want {
				t.Fatalf("VerdantDew() = %d, want %d", got, tc.want)
			}
			if got := c.Player.MoonridgeDew(); got != 1 {
				t.Fatalf("MoonridgeDew() = %d, want 1", got)
			}
			data, err := c.Log.Dump()
			if err != nil {
				t.Fatal(err)
			}
			var events []glog.LogEvent
			if err := json.Unmarshal(data, &events); err != nil {
				t.Fatal(err)
			}
			if len(events) == 0 {
				t.Fatal("missing Verdant Dew log event")
			}
			if got, want := events[len(events)-1].Msg, fmt.Sprintf("verdant dew set to %d", tc.want); got != want {
				t.Errorf("log message = %q, want %q", got, want)
			}
		})
	}
}
