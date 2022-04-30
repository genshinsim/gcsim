package character

import (
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

// AttackModFunc returns an array containing the stats boost and whether mod applies
type AttackModFunc func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool)

type attackMod struct {
	Amount AttackModFunc
	modTmpl
}

func (c *CharWrapper) AddAttackMod(key string, dur int, f AttackModFunc) {
	expiry := *c.f + dur
	if dur < 0 {
		expiry = -1
	}
	mod := attackMod{
		Amount: f,
		modTmpl: modTmpl{
			key:    key,
			expiry: expiry,
		},
	}
	addMod(c, c.attackMods, &mod)
}

func (c *CharWrapper) DeleteAttackMod(key string) {
	deleteMod(c, c.attackMods, key)
}

func (c *CharWrapper) ApplyAttackMods(a *combat.AttackEvent, t combat.Target) []interface{} {
	//skip if this is reaction damage
	if a.Info.AttackTag >= combat.AttackTagNoneStat {
		return nil
	}

	var sb strings.Builder
	var logDetails []interface{}

	if c.debug {
		logDetails = make([]interface{}, 0, len(c.attackMods))
	}

	n := 0
	for _, m := range c.attackMods {

		if m.expiry > *c.f || m.expiry == -1 {

			amt, ok := m.Amount(a, t)
			if ok {
				for k, v := range amt {
					a.Snapshot.Stats[k] += v
				}
			}
			c.attackMods[n] = m
			n++

			if c.debug {
				modStatus := make([]string, 0)
				if ok {
					sb.WriteString(m.key)
					modStatus = append(
						modStatus,
						"status: added",
						"expiry_frame: "+strconv.Itoa(m.expiry),
					)
					modStatus = append(
						modStatus,
						attributes.PrettyPrintStatsSlice(amt)...,
					)
					logDetails = append(logDetails, sb.String(), modStatus)
					sb.Reset()
				} else {
					sb.WriteString(m.key)
					modStatus = append(
						modStatus,
						"status: rejected",
						"reason: conditions not met",
					)
					logDetails = append(logDetails, sb.String(), modStatus)
					sb.Reset()
				}
			}
		}
	}
	c.attackMods = c.attackMods[:n]
	return logDetails
}
