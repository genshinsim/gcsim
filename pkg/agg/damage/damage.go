package damage

import (
	calc "github.com/aclements/go-moremath/stats"
	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/stats"
)

// 30 = .5s
const bucketSize uint32 = 30

func init() {
	agg.Register(NewAgg)
}

// TODO: We need to populate targetDPS with 0s if damage wasn't done that iteration
// for an accurate measure. The problem is that we need target keys to be decided at the cfg level
// not the core level.
// We also have no guarantee that targets will have the same key across iterations. This will solve
// the problem.
type buffer struct {
	elementDPS            map[string]*calc.StreamStats
	targetDPS             map[int]*calc.StreamStats
	sourceDPS             []map[string]*calc.StreamStats
	sourceDamageInstances []map[string]*calc.StreamStats
	characterDPS          []*calc.StreamStats // i = char
	dpsByElement          []map[string]*calc.StreamStats
	dpsByTarget           []map[int]*calc.StreamStats

	damageBuckets     []*calc.StreamStats
	cumulativeContrib [][]*calc.StreamStats
}

func NewAgg(cfg *info.ActionList) (agg.Aggregator, error) {
	out := buffer{
		elementDPS:            make(map[string]*calc.StreamStats),
		targetDPS:             make(map[int]*calc.StreamStats),
		sourceDPS:             make([]map[string]*calc.StreamStats, len(cfg.Characters)),
		sourceDamageInstances: make([]map[string]*calc.StreamStats, len(cfg.Characters)),
		characterDPS:          make([]*calc.StreamStats, len(cfg.Characters)),
		dpsByElement:          make([]map[string]*calc.StreamStats, len(cfg.Characters)),
		dpsByTarget:           make([]map[int]*calc.StreamStats, len(cfg.Characters)),
		cumulativeContrib:     make([][]*calc.StreamStats, len(cfg.Characters)),
		damageBuckets:         make([]*calc.StreamStats, 0),
	}

	// start with single entry
	out.damageBuckets = append(out.damageBuckets, &calc.StreamStats{})

	for i := 0; i < len(cfg.Characters); i++ {
		out.sourceDPS[i] = make(map[string]*calc.StreamStats)
		out.sourceDamageInstances[i] = make(map[string]*calc.StreamStats)
		out.characterDPS[i] = &calc.StreamStats{}
		out.dpsByElement[i] = make(map[string]*calc.StreamStats)
		out.dpsByTarget[i] = make(map[int]*calc.StreamStats)
		out.cumulativeContrib[i] = make([]*calc.StreamStats, 0)
		out.cumulativeContrib[i] = append(out.cumulativeContrib[i], &calc.StreamStats{})
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result) {
	time := 60 / float64(result.Duration)
	targetDPS := make(map[int]float64)
	elementDPS := makeElementMap()

	b.damageBuckets = expandBuckets(
		b.damageBuckets, max(len(b.damageBuckets), len(result.DamageBuckets)))
	for i, stat := range b.damageBuckets {
		var val float64
		if i < len(result.DamageBuckets) {
			val = result.DamageBuckets[i]
		} else {
			val = 0
		}
		stat.Add(val)
	}

	for i := range result.Characters {
		var charDPS float64
		charElementDPS := makeElementMap()
		charTargetDPS := make(map[int]float64)
		sourceDPS := make(map[string]float64)
		sourceDamageInstances := make(map[string]float64)

		b.cumulativeContrib[i] = expandCumu(
			b.cumulativeContrib[i],
			max(len(b.cumulativeContrib[i]), len(result.Characters[i].DamageCumulativeContrib)))
		var prev float64
		for j, stat := range b.cumulativeContrib[i] {
			var val float64
			if j < len(result.Characters[i].DamageCumulativeContrib) {
				val = result.Characters[i].DamageCumulativeContrib[j]
			} else {
				val = prev
			}
			prev = val
			stat.Add(val)
		}

		for _, ev := range result.Characters[i].DamageEvents {
			if _, ok := charTargetDPS[ev.Target]; !ok {
				charTargetDPS[ev.Target] = 0
			}
			charTargetDPS[ev.Target] += ev.Damage
			charElementDPS[ev.Element] += ev.Damage
			charDPS += ev.Damage
			sourceDPSKey := ev.Source
			reactModifier := string(ev.ReactionModifier)
			if reactModifier != "" {
				sourceDPSKey += " (" + reactModifier + ")"
			}
			sourceDPS[sourceDPSKey] += ev.Damage
			sourceDamageInstances[sourceDPSKey] += 1
		}

		b.characterDPS[i].Add(charDPS * time)
		for k, v := range charElementDPS {
			if _, ok := elementDPS[k]; !ok {
				elementDPS[k] = 0
			}
			elementDPS[k] += v

			if _, ok := b.dpsByElement[i][k]; !ok {
				b.dpsByElement[i][k] = &calc.StreamStats{}
			}
			b.dpsByElement[i][k].Add(v * time)
		}

		for k, v := range charTargetDPS {
			if _, ok := targetDPS[k]; !ok {
				targetDPS[k] = 0
			}
			targetDPS[k] += v

			if _, ok := b.dpsByTarget[i][k]; !ok {
				b.dpsByTarget[i][k] = &calc.StreamStats{}
			}
			b.dpsByTarget[i][k].Add(v * time)
		}

		for k, v := range sourceDPS {
			if _, ok := b.sourceDPS[i][k]; !ok {
				b.sourceDPS[i][k] = &calc.StreamStats{}
			}
			b.sourceDPS[i][k].Add(v * time)
		}

		for k, v := range sourceDamageInstances {
			if _, ok := b.sourceDamageInstances[i][k]; !ok {
				b.sourceDamageInstances[i][k] = &calc.StreamStats{}
			}
			b.sourceDamageInstances[i][k].Add(v)
		}
	}

	for k, v := range targetDPS {
		if _, ok := b.targetDPS[k]; !ok {
			b.targetDPS[k] = &calc.StreamStats{}
		}
		b.targetDPS[k].Add(v * time)
	}

	for k, v := range elementDPS {
		if _, ok := b.elementDPS[k]; !ok {
			b.elementDPS[k] = &calc.StreamStats{}
		}
		b.elementDPS[k].Add(v * time)
	}
}

func (b *buffer) Flush(result *model.SimulationStatistics) {
	result.ElementDps = make(map[string]*model.DescriptiveStats)
	for k, v := range b.elementDPS {
		if v.Mean() > 0 {
			result.ElementDps[k] = agg.ToDescriptiveStats(v)
		}
	}

	result.TargetDps = make(map[int32]*model.DescriptiveStats)
	for k, v := range b.targetDPS {
		result.TargetDps[int32(k)] = agg.ToDescriptiveStats(v)
	}

	result.CharacterDps = make([]*model.DescriptiveStats, len(b.characterDPS))
	for i, v := range b.characterDPS {
		result.CharacterDps[i] = agg.ToDescriptiveStats(v)
	}

	result.BreakdownByElementDps = make([]*model.ElementStats, len(b.dpsByElement))
	for i, em := range b.dpsByElement {
		elements := make(map[string]*model.DescriptiveStats)
		for k, v := range em {
			if v.Mean() > 0 {
				elements[k] = agg.ToDescriptiveStats(v)
			}
		}

		result.BreakdownByElementDps[i] = &model.ElementStats{
			Elements: elements,
		}
	}

	result.BreakdownByTargetDps = make([]*model.TargetStats, len(b.dpsByTarget))
	for i, t := range b.dpsByTarget {
		targets := make(map[int32]*model.DescriptiveStats)
		for k, v := range t {
			targets[int32(k)] = agg.ToDescriptiveStats(v)
		}

		result.BreakdownByTargetDps[i] = &model.TargetStats{
			Targets: targets,
		}
	}

	damageBuckets := make([]*model.DescriptiveStats, len(b.damageBuckets))
	for i, v := range b.damageBuckets {
		damageBuckets[i] = agg.ToDescriptiveStats(v)
	}
	result.DamageBuckets = &model.BucketStats{
		BucketSize: bucketSize,
		Buckets:    damageBuckets,
	}

	characterBuckets := make([]*model.CharacterBuckets, len(b.cumulativeContrib))
	for i, c := range b.cumulativeContrib {
		buckets := make([]*model.DescriptiveStats, len(c))
		for j, v := range c {
			buckets[j] = agg.ToDescriptiveStats(v)
		}
		characterBuckets[i] = &model.CharacterBuckets{
			Buckets: buckets,
		}
	}

	result.SourceDps = make([]*model.SourceStats, len(b.sourceDPS))
	for i, c := range b.sourceDPS {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.SourceDps[i] = &model.SourceStats{
			Sources: source,
		}
	}

	result.SourceDamageInstances = make([]*model.SourceStats, len(b.sourceDamageInstances))
	for i, c := range b.sourceDamageInstances {
		source := make(map[string]*model.DescriptiveStats)
		for k, s := range c {
			source[k] = agg.ToDescriptiveStats(s)
		}

		result.SourceDamageInstances[i] = &model.SourceStats{
			Sources: source,
		}
	}

	result.CumulativeDamageContribution = &model.CharacterBucketStats{
		BucketSize: bucketSize,
		Characters: characterBuckets,
	}
}

func makeElementMap() map[string]float64 {
	out := make(map[string]float64)
	for _, ele := range attributes.ElementString {
		out[ele] = 0
	}
	return out
}

func expandCumu(arr []*calc.StreamStats, size int) []*calc.StreamStats {
	last := arr[len(arr)-1]
	for size > len(arr) {
		cpy := *last
		arr = append(arr, &cpy)
	}
	return arr
}

func expandBuckets(arr []*calc.StreamStats, size int) []*calc.StreamStats {
	last := arr[len(arr)-1]
	for size > len(arr) {
		newStat := &calc.StreamStats{}
		for i := 0; i < int(last.Count); i++ {
			newStat.Add(0)
		}
		arr = append(arr, newStat)
	}
	return arr
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
