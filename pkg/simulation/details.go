package simulation

import (
	"fmt"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type CharacterDetail struct {
	Name          string         `json:"name"`
	Element       string         `json:"element"`
	Level         int            `json:"level"`
	MaxLevel      int            `json:"max_level"`
	Cons          int            `json:"cons"`
	Weapon        WeaponDetail   `json:"weapon"`
	Talents       TalentDetail   `json:"talents"`
	Sets          map[string]int `json:"sets"`
	Stats         []float64      `json:"stats"`
	SnapshotStats []float64      `json:"snapshot"`
}

type WeaponDetail struct {
	Name     string `json:"name"`
	Refine   int    `json:"refine"`
	Level    int    `json:"level"`
	MaxLevel int    `json:"max_level"`
}

type TalentDetail struct {
	Attack int `json:"attack"`
	Skill  int `json:"skill"`
	Burst  int `json:"burst"`
}

func (sim *Simulation) CharacterDetails() []CharacterDetail {
	out := make([]CharacterDetail, len(sim.C.Player.Chars()))

	for i, v := range sim.cfg.Characters {
		m := make(map[string]int)
		for k, v := range v.Sets {
			m[k.String()] = v
		}

		char := CharacterDetail{
			Name:     v.Base.Key.String(),
			Element:  v.Base.Element.String(),
			Level:    v.Base.Level,
			MaxLevel: v.Base.MaxLevel,
			Cons:     v.Base.Cons,
			Weapon: WeaponDetail{
				Name:     v.Weapon.Key.String(),
				Refine:   v.Weapon.Refine,
				Level:    v.Weapon.Level,
				MaxLevel: v.Weapon.MaxLevel,
			},
			Talents: TalentDetail{
				Attack: v.Talents.Attack,
				Skill:  v.Talents.Skill,
				Burst:  v.Talents.Burst,
			},
			Sets:  m,
			Stats: v.Stats,
		}
		out[i] = char
	}

	//grab a snapshot for each char
	for i, c := range sim.C.Player.Chars() {
		snap := c.Snapshot(&combat.AttackInfo{
			Abil:      "stats-check",
			AttackTag: attacks.AttackTagNone,
		})
		//convert all atk%, def% and hp% into flat amounts by tacking on base
		snap.Stats[attributes.HP] += c.Base.HP * (1 + snap.Stats[attributes.HPP])
		snap.Stats[attributes.DEF] += c.Base.Def * (1 + snap.Stats[attributes.DEFP])
		snap.Stats[attributes.ATK] += (c.Base.Atk + c.Weapon.Atk) * (1 + snap.Stats[attributes.ATKP])
		snap.Stats[attributes.HPP] = 0
		snap.Stats[attributes.DEFP] = 0
		snap.Stats[attributes.ATKP] = 0
		if sim.C.Combat.Debug {
			evt := sim.C.Log.NewEvent(
				fmt.Sprintf("%v final stats", c.Base.Key.Pretty()),
				glog.LogCharacterEvent,
				i,
			)
			for i, v := range snap.Stats {
				if v != 0 {
					evt.Write(attributes.StatTypeString[i], strconv.FormatFloat(v, 'f', 3, 32))
				}
			}
		}
		out[i].SnapshotStats = snap.Stats[:]
	}

	return out
}
