package agg

import (
	"math"

	"github.com/aclements/go-moremath/stats"
)

func ConvertToFloatStat(ss *stats.StreamStats) FloatStat {
	sd := ss.StdDev()
	if math.IsNaN(sd) {
		sd = 0
	}

	return FloatStat{
		Min:  ss.Min,
		Max:  ss.Max,
		Mean: ss.Mean(),
		SD:   sd,
	}
}
