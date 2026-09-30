package shield

import (
	"reflect"
	"testing"

	"github.com/genshinsim/gcsim/pkg/stats"
)

func hp(v float64) map[string]float64 {
	out := map[string]float64{normalized: v}
	for _, e := range elements {
		out[e.String()] = v
	}
	return out
}

func effective(intervals ...stats.ShieldSingleInterval) map[string][]stats.ShieldSingleInterval {
	out := map[string][]stats.ShieldSingleInterval{normalized: intervals}
	for _, e := range elements {
		out[e.String()] = intervals
	}
	return out
}

func TestComputeEffective(t *testing.T) {
	tests := []struct {
		name    string
		shields map[string][]stats.ShieldInterval
		want    map[string][]stats.ShieldSingleInterval
	}{
		{
			// zeta and alpha tie once big ends, as crystallize shields of different elements do
			// on every other element
			name: "tie goes to the shield that started first",
			shields: map[string][]stats.ShieldInterval{
				"big":   {{Start: 0, End: 30, HP: hp(1000)}},
				"zeta":  {{Start: 10, End: 50, HP: hp(100)}},
				"alpha": {{Start: 12, End: 80, HP: hp(100)}},
			},
			want: effective(
				stats.ShieldSingleInterval{Start: 0, End: 30, HP: 1000},
				stats.ShieldSingleInterval{Start: 30, End: 50, HP: 100},
				stats.ShieldSingleInterval{Start: 50, End: 80, HP: 100},
			),
		},
		{
			name: "a shield ending on the same frame doesn't take over",
			shields: map[string][]stats.ShieldInterval{
				"big": {{Start: 0, End: 30, HP: hp(1000)}},
				"z":   {{Start: 10, End: 30, HP: hp(500)}},
				"c":   {{Start: 10, End: 80, HP: hp(100)}},
			},
			want: effective(
				stats.ShieldSingleInterval{Start: 0, End: 30, HP: 1000},
				stats.ShieldSingleInterval{Start: 30, End: 80, HP: 100},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// the result used to follow map iteration order, so one run could pass by chance
			for range 100 {
				if got := computeEffective(tt.shields); !reflect.DeepEqual(got, tt.want) {
					t.Fatalf("computeEffective() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
