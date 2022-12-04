package damage

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	elementDPS   map[string]*calc.StreamStats
	characterDPS []*calc.StreamStats // i = char
	targetDPS    []*calc.StreamStats // i = target
	dpsByElement []map[string]*calc.StreamStats
	dpsByTarget  [][]*calc.StreamStats // i = char, j = target
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		elementDPS:   make(map[string]*calc.StreamStats),
		characterDPS: make([]*calc.StreamStats, len(cfg.Characters)),
		targetDPS:    make([]*calc.StreamStats, len(cfg.Targets)),
		dpsByElement: make([]map[string]*calc.StreamStats, len(cfg.Characters)),
		dpsByTarget:  make([][]*calc.StreamStats, len(cfg.Characters)),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.characterDPS[i] = &calc.StreamStats{}
		out.dpsByElement[i] = make(map[string]*calc.StreamStats)

		out.dpsByTarget[i] = make([]*calc.StreamStats, len(cfg.Targets))
		for j := 0; j < len(cfg.Targets); j++ {
			out.dpsByTarget[i][j] = &calc.StreamStats{}
		}
	}

	for i := 0; i < len(cfg.Targets); i++ {
		out.targetDPS[i] = &calc.StreamStats{}
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	time := 60 / float64(result.Duration)
	targetDPS := make([]float64, len(b.targetDPS))
	elementDPS := makeElementMap()

	for i, char := range result.Characters {
		var charDPS float64
		charElementDPS := makeElementMap()
		charTargetDPS := make([]float64, len(b.targetDPS))

		for _, ev := range char.DamageEvents {
			charElementDPS[ev.Element] += ev.Damage
			charTargetDPS[ev.Target-1] += ev.Damage
			charDPS += ev.Damage
		}

		b.characterDPS[i].Add(charDPS * time)
		for k, v := range charElementDPS {
			if _, ok := b.dpsByElement[i][k]; !ok {
				b.dpsByElement[i][k] = &calc.StreamStats{}
			}
			b.dpsByElement[i][k].Add(v * time)
			elementDPS[k] += v
		}

		for j, v := range charTargetDPS {
			b.dpsByTarget[i][j].Add(v * time)
			targetDPS[j] += v
		}
	}

	for i, v := range targetDPS {
		b.targetDPS[i].Add(v * time)
	}

	for k, v := range elementDPS {
		if _, ok := b.elementDPS[k]; !ok {
			b.elementDPS[k] = &calc.StreamStats{}
		}
		b.elementDPS[k].Add(v * time)
	}
}

func (b *buffer) Flush(result *agg.Result) {
	result.ElementDPS = make(map[string]agg.FloatStat)
	for k, v := range b.elementDPS {
		if v.Min > 0 {
			result.ElementDPS[k] = agg.ConvertToFloatStat(v)
		}
	}

	result.TargetDPS = make([]agg.FloatStat, len(b.targetDPS))
	for i, v := range b.targetDPS {
		result.TargetDPS[i] = agg.ConvertToFloatStat(v)
	}

	result.CharacterDPS = make([]agg.FloatStat, len(b.characterDPS))
	for i, v := range b.characterDPS {
		result.CharacterDPS[i] = agg.ConvertToFloatStat(v)
	}

	result.BreakdownByElementDPS = make([]map[string]agg.FloatStat, len(b.dpsByElement))
	for i, em := range b.dpsByElement {
		result.BreakdownByElementDPS[i] = make(map[string]agg.FloatStat)
		for k, v := range em {
			if v.Min > 0 {
				result.BreakdownByElementDPS[i][k] = agg.ConvertToFloatStat(v)
			}
		}
	}

	result.BreakdownByTargetDPS = make([][]agg.FloatStat, len(b.dpsByTarget))
	for i, t := range b.dpsByTarget {
		result.BreakdownByTargetDPS[i] = make([]agg.FloatStat, len(t))
		for j, v := range t {
			result.BreakdownByTargetDPS[i][j] = agg.ConvertToFloatStat(v)
		}
	}
}

func makeElementMap() map[string]float64 {
	out := make(map[string]float64)
	for _, ele := range attributes.ElementString {
		out[ele] = 0
	}
	return out
}
