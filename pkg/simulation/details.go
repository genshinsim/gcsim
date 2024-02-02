package simulation

import (
	"fmt"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/model"
)

func (s *Simulation) CharacterDetails() []*model.Character {
	out := make([]*model.Character, len(s.C.Player.Chars()))

	for i := range s.cfg.Characters {
		m := make(map[string]int32)
		for k, v := range s.cfg.Characters[i].Sets {
			m[k.String()] = int32(v)
		}

		char := &model.Character{
			Name:     s.cfg.Characters[i].Base.Key.String(),
			Element:  s.cfg.Characters[i].Base.Element.String(),
			Level:    int32(s.cfg.Characters[i].Base.Level),
			MaxLevel: int32(s.cfg.Characters[i].Base.MaxLevel),
			Cons:     int32(s.cfg.Characters[i].Base.Cons),
			Weapon: &model.Weapon{
				Name:     s.cfg.Characters[i].Weapon.Key.String(),
				Refine:   int32(s.cfg.Characters[i].Weapon.Refine),
				Level:    int32(s.cfg.Characters[i].Weapon.Level),
				MaxLevel: int32(s.cfg.Characters[i].Weapon.MaxLevel),
			},
			Talents: &model.CharacterTalents{
				Attack: int32(s.cfg.Characters[i].Talents.Attack),
				Skill:  int32(s.cfg.Characters[i].Talents.Skill),
				Burst:  int32(s.cfg.Characters[i].Talents.Burst),
			},
			Sets:  m,
			Stats: s.cfg.Characters[i].Stats,
		}
		out[i] = char
	}

	// grab a snapshot for each char
	for i, c := range s.C.Player.Chars() {
		snap := c.Snapshot(&combat.AttackInfo{
			Abil:      "stats-check",
			AttackTag: attacks.AttackTagNone,
		})
		// convert all atk%, def% and hp% into flat amounts by tacking on base
		snap.Stats[attributes.HP] += c.Base.HP * (1 + snap.Stats[attributes.HPP])
		snap.Stats[attributes.DEF] += c.Base.Def * (1 + snap.Stats[attributes.DEFP])
		snap.Stats[attributes.ATK] += (c.Base.Atk + c.Weapon.BaseAtk) * (1 + snap.Stats[attributes.ATKP])
		snap.Stats[attributes.HPP] = 0
		snap.Stats[attributes.DEFP] = 0
		snap.Stats[attributes.ATKP] = 0
		if s.C.Combat.Debug {
			evt := s.C.Log.NewEvent(
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
