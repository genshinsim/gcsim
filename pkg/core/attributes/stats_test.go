package attributes

import (
	"slices"
	"testing"
)

func TestPrettyPrintStatsPrecision(t *testing.T) {
	cases := []struct {
		v    float64
		want string
	}{
		{0.225, "dmg%: 0.225"},
		{0.075, "dmg%: 0.075"},
		{0.12, "dmg%: 0.12"},
		{1.1234, "dmg%: 1.1234"},
		{0.1 + 0.2, "dmg%: 0.3"},
		{0.1234567, "dmg%: 0.123457"},
		{311, "dmg%: 311"},
	}
	for _, c := range cases {
		stats := make([]float64, EndStatType)
		stats[DmgP] = c.v

		if got := PrettyPrintStatsSlice(stats); !slices.Equal(got, []string{c.want}) {
			t.Errorf("PrettyPrintStatsSlice(%v) = %q, want %q", c.v, got, c.want)
		}
		if got := PrettyPrintStats(stats); got != c.want {
			t.Errorf("PrettyPrintStats(%v) = %q, want %q", c.v, got, c.want)
		}
	}
}

func TestPrettyPrintStatsSkipsZero(t *testing.T) {
	stats := make([]float64, EndStatType)
	stats[ATK] = 311
	stats[CR] = 0.05

	want := []string{"atk: 311", "cr: 0.05"}
	if got := PrettyPrintStatsSlice(stats); !slices.Equal(got, want) {
		t.Errorf("PrettyPrintStatsSlice = %q, want %q", got, want)
	}
	if got, want := PrettyPrintStats(stats), "atk: 311 cr: 0.05"; got != want {
		t.Errorf("PrettyPrintStats = %q, want %q", got, want)
	}
}
