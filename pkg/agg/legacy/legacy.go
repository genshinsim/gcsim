package legacy

import (
	"strconv"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/agg/util"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/stats"
)

func init() {
	agg.Register(NewAgg)
}

type buffer struct {
	iterations int

	// damage (w/ SD)
	damageOverTime map[string]*util.FloatBuffer
	dpsByTarget    map[int]*util.FloatBuffer

	// damage (no SD)
	damageByChar          []map[string]*util.FloatBuffer
	damageInstancesByChar []map[string]*util.IntBuffer
	damageByCharByTargets []map[int]*util.FloatBuffer

	// metadata per char
	charActiveTime       []*util.IntBuffer
	abilUsageCountByChar []map[string]*util.IntBuffer

	// metadata per enemy
	elementUptime []map[string]*util.IntBuffer

	// overall metadata
	particleCount      map[string]*util.FloatBuffer
	reactionsTriggered map[string]*util.IntBuffer
}

func NewAgg(cfg *ast.ActionList) (agg.Aggregator, error) {
	out := buffer{
		iterations:            cfg.Settings.Iterations,
		damageOverTime:        make(map[string]*util.FloatBuffer),
		dpsByTarget:           make(map[int]*util.FloatBuffer),
		damageByChar:          make([]map[string]*util.FloatBuffer, len(cfg.Characters)),
		damageInstancesByChar: make([]map[string]*util.IntBuffer, len(cfg.Characters)),
		damageByCharByTargets: make([]map[int]*util.FloatBuffer, len(cfg.Characters)),
		abilUsageCountByChar:  make([]map[string]*util.IntBuffer, len(cfg.Characters)),
		charActiveTime:        make([]*util.IntBuffer, len(cfg.Characters)),
		elementUptime:         make([]map[string]*util.IntBuffer, len(cfg.Targets)),
		particleCount:         make(map[string]*util.FloatBuffer),
		reactionsTriggered:    make(map[string]*util.IntBuffer),
	}

	for i := 0; i < len(cfg.Characters); i++ {
		out.damageByChar[i] = make(map[string]*util.FloatBuffer)
		out.damageInstancesByChar[i] = make(map[string]*util.IntBuffer)
		out.damageByCharByTargets[i] = make(map[int]*util.FloatBuffer)
		out.abilUsageCountByChar[i] = make(map[string]*util.IntBuffer)
		out.charActiveTime[i] = util.NewIntBuffer(out.iterations)
	}

	for i := 0; i < len(cfg.Targets); i++ {
		out.elementUptime[i] = make(map[string]*util.IntBuffer)
	}

	return &out, nil
}

func (b *buffer) Add(result stats.Result, itr int) {
	dd := float64(result.Duration) / 60

	for k, v := range result.Legacy.DamageOverTime {
		if _, ok := b.damageOverTime[k]; !ok {
			b.damageOverTime[k] = util.NewFloatBuffer(b.iterations)
		}
		b.damageOverTime[k].Add(v, itr)
	}

	// loop over all chars
	damagePerTarget := make(map[int]float64)
	for i := 0; i < len(b.charActiveTime); i++ {
		for k, v := range result.Legacy.DamageByChar[i] {
			if _, ok := b.damageByChar[i][k]; !ok {
				b.damageByChar[i][k] = util.NewFloatBufferNoSD(b.iterations)
			}
			b.damageByChar[i][k].Add(v/dd, itr)
		}

		for k, v := range result.Legacy.DamageByCharByTargets[i] {
			key, _ := strconv.Atoi(k)
			damagePerTarget[key] += v

			if _, ok := b.damageByCharByTargets[i][key]; !ok {
				b.damageByCharByTargets[i][key] = util.NewFloatBufferNoSD(b.iterations)
			}
			b.damageByCharByTargets[i][key].Add(v/dd, itr)
		}

		for k, v := range result.Legacy.DamageInstancesByChar[i] {
			if _, ok := b.damageInstancesByChar[i][k]; !ok {
				b.damageInstancesByChar[i][k] = util.NewIntBuffer(b.iterations)
			}
			b.damageInstancesByChar[i][k].Add(v)
		}

		for k, v := range result.Legacy.AbilUsageCountByChar[i] {
			if _, ok := b.abilUsageCountByChar[i][k]; !ok {
				b.abilUsageCountByChar[i][k] = util.NewIntBuffer(b.iterations)
			}
			b.abilUsageCountByChar[i][k].Add(v)
		}

		b.charActiveTime[i].Add(result.Legacy.CharActiveTime[i])
	}

	for k, v := range damagePerTarget {
		if _, ok := b.dpsByTarget[k]; !ok {
			b.dpsByTarget[k] = util.NewFloatBuffer(b.iterations)
		}
		b.dpsByTarget[k].Add(v/dd, itr)
	}

	for i := 0; i < len(result.Legacy.ElementUptime); i++ {
		for k, v := range result.Legacy.ElementUptime[i] {
			if _, ok := b.elementUptime[i][k]; !ok {
				b.elementUptime[i][k] = util.NewIntBuffer(b.iterations)
			}
			b.elementUptime[i][k].Add(v)
		}
	}

	for k, v := range result.Legacy.ParticleCount {
		if _, ok := b.particleCount[k]; !ok {
			b.particleCount[k] = util.NewFloatBufferNoSD(b.iterations)
		}
		b.particleCount[k].Add(v, itr)
	}

	for k, v := range result.Legacy.ReactionsTriggered {
		if _, ok := b.reactionsTriggered[k]; !ok {
			b.reactionsTriggered[k] = util.NewIntBuffer(b.iterations)
		}
		b.reactionsTriggered[k].Add(v)
	}
}

func (b *buffer) Flush(result *agg.Result) {
	numChars := len(b.damageByChar)
	numEnemies := len(b.elementUptime)

	result.Legacy.DamageByChar = make([]map[string]agg.FloatStat, numChars)
	result.Legacy.DamageInstancesByChar = make([]map[string]agg.IntStat, numChars)
	result.Legacy.DamageByCharByTargets = make([]map[int]agg.FloatStat, numChars)
	result.Legacy.CharActiveTime = make([]agg.IntStat, numChars)
	result.Legacy.AbilUsageCountByChar = make([]map[string]agg.IntStat, numChars)
	for i := 0; i < numChars; i++ {
		result.Legacy.DamageByChar[i] = make(map[string]agg.FloatStat)
		for k, v := range b.damageByChar[i] {
			result.Legacy.DamageByChar[i][k] = v.Flush()
		}

		result.Legacy.DamageInstancesByChar[i] = make(map[string]agg.IntStat)
		for k, v := range b.damageInstancesByChar[i] {
			result.Legacy.DamageInstancesByChar[i][k] = v.Flush()
		}

		result.Legacy.DamageByCharByTargets[i] = make(map[int]agg.FloatStat)
		for k, v := range b.damageByCharByTargets[i] {
			result.Legacy.DamageByCharByTargets[i][k] = v.Flush()
		}

		result.Legacy.CharActiveTime[i] = b.charActiveTime[i].Flush()

		result.Legacy.AbilUsageCountByChar[i] = make(map[string]agg.IntStat)
		for k, v := range b.abilUsageCountByChar[i] {
			result.Legacy.AbilUsageCountByChar[i][k] = v.Flush()
		}
	}

	result.Legacy.ElementUptime = make([]map[string]agg.IntStat, numEnemies)
	for i := 0; i < numEnemies; i++ {
		result.Legacy.ElementUptime[i] = make(map[string]agg.IntStat)
		for k, v := range b.elementUptime[i] {
			result.Legacy.ElementUptime[i][k] = v.Flush()
		}
	}

	result.Legacy.ParticleCount = make(map[string]agg.FloatStat)
	for k, v := range b.particleCount {
		result.Legacy.ParticleCount[k] = v.Flush()
	}

	result.Legacy.ReactionsTriggered = make(map[string]agg.IntStat)
	for k, v := range b.reactionsTriggered {
		result.Legacy.ReactionsTriggered[k] = v.Flush()
	}

	result.Legacy.DPSByTarget = make(map[int]agg.FloatStat)
	for k, v := range b.dpsByTarget {
		result.Legacy.DPSByTarget[k] = v.Flush()
	}

	result.Legacy.DamageOverTime = make(map[string]agg.FloatStat)
	for k, v := range b.damageOverTime {
		result.Legacy.DamageOverTime[k] = v.Flush()
	}
}
