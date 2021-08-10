package combat

import (
	"log"

	"github.com/genshinsim/gsim/pkg/core"
)

func (s *Sim) ApplyDamage(ds *core.Snapshot) {
	died := false
	for i, v := range s.targets {
		d := ds.Clone()
		dmg, crit := v.Attack(&d)
		s.stats.Damage += dmg
		s.stats.DamageByChar[d.ActorIndex][d.Abil] += dmg
		//check if target is dead
		if s.cfg.Mode.HPMode && v.HP() <= 0 {
			died = true
			s.OnTargetDefeated(s.targets[i])
			s.targets[i] = nil
			log.Println("target died", i, dmg)
		}

		amp := ""
		if d.IsMeltVape {
			amp = string(d.ReactionType)
		}

		s.log.Debugw(
			d.Abil,
			"frame", s.f,
			"event", core.LogDamageEvent,
			"char", d.ActorIndex,
			"target", i,
			"attack_tag", d.AttackTag,
			"damage", dmg,
			"crit", crit,
			"amp", amp,
			"abil", d.Abil,
		)

	}
	if died {
		//wipe out nil entries
		n := 0
		for _, v := range s.targets {
			if v != nil {
				s.targets[n] = v
				s.targets[n].SetIndex(n)
				n++
			}
		}
		s.targets = s.targets[:n]
	}
}

func (s *Sim) execQueue() (int, error) {
	//if length of q is 0, search for next
	if len(s.queue) == 0 {
		next, err := s.querer.Next(s.chars[s.active].Name())
		if err != nil {
			return 0, err
		}
		if len(next) == 0 {
			return 0, nil
		}
		s.queue = append(s.queue, next...)
	}

	willWait := false
	n := s.queue[0]
	//otherwise pop first item on queue and execute

	c := s.chars[s.active]
	f := 0

	s.log.Debugw(
		"attempting to execute "+n.Typ.String(),
		"frame", s.f,
		"event", core.LogActionEvent,
		"char", s.active,
		"action", n.Typ.String(),
		"target", n.Target,
		"swap_cd_pre", s.swapCD,
		"stam_pre", s.stam,
		"animation", f,
	)

	//do one last ready check
	if !c.ActionReady(n.Typ, n.Param) {
		s.log.Warnw("queued action is not ready, should not happen; skipping frame")
		return 0, nil
	}
	switch n.Typ {
	case core.ActionSwapLock:
		s.swapCD += n.SwapLock
		// return 0
	case core.ActionSkill:
		s.executeEventHook(core.PreSkillHook)
		f = c.Skill(n.Param)
		s.executeEventHook(core.PostSkillHook)
		s.ResetAllNormalCounter()
	case core.ActionBurst:
		s.executeEventHook(core.PreBurstHook)
		f = c.Burst(n.Param)
		s.executeEventHook(core.PostBurstHook)
		s.ResetAllNormalCounter()
	case core.ActionAttack:
		f = c.Attack(n.Param)
		s.executeEventHook(core.PostAttackHook)
	case core.ActionCharge:
		req := s.StamPercentMod(core.ActionCharge) * c.ActionStam(core.ActionCharge, n.Param)
		if s.stam <= req {
			f = 90 - (s.f - s.lastStamUse)
			s.log.Warnw("insufficient stam: charge attack", "have", s.stam, "last", s.lastStamUse, "recharge", f)
			willWait = true
		} else {
			s.stam -= req
			f += c.ChargeAttack(n.Param)
			s.ResetAllNormalCounter()
			s.lastStamUse = s.f
		}
	case core.ActionHighPlunge:
		f = c.HighPlungeAttack(n.Param)
	case core.ActionLowPlunge:
		f = c.LowPlungeAttack(n.Param)
	case core.ActionAim:
		f = c.Aimed(n.Param)
		s.ResetAllNormalCounter()
	case core.ActionSwap:
		s.executeEventHook(core.PreSwapHook)
		f = swapFrames
		//if we're still in cd then forcefully wait up the cd
		if s.swapCD > 0 {
			f += s.swapCD
		}
		s.swapCD = swapCDFrames
		// s.Log.Debugw("swap", "frame", s.F, "event", LogActionEvent, "from", s.ActiveChar, "to", n.Target)
		ind := s.charPos[n.Target]
		s.active = ind
		s.ResetAllNormalCounter()
		s.executeEventHook(core.PostSwapHook)
		//this duration reset needs to be after the hook for spine to behave properly
		s.charActiveDuration = 0
	case core.ActionCancellable:
	case core.ActionDash:
		//check if enough stam
		stam := s.StamPercentMod(core.ActionDash) * c.ActionStam(core.ActionDash, n.Param)
		if s.stam <= stam {
			f = 90 - (s.f - s.lastStamUse)
			// s.Log.Warnw("insufficient stam: dash", "have", s.Stam, "last", s.LastStamUse, "recharge", f)
			willWait = true
		} else {
			s.stam -= stam
			f = c.Dash(n.Param)
			s.ResetAllNormalCounter()
			s.lastStamUse = s.f
			s.executeEventHook(core.PostDashHook)
		}
	case core.ActionJump:
		f = jumpFrames
		s.ResetAllNormalCounter()
	}

	if willWait {
		s.log.Debugw(
			"execution will wait "+n.Typ.String(),
			"frame", s.f,
			"event", core.LogActionEvent,
			"char", s.active,
			"action", n.Typ.String(),
			"target", n.Target,
			"will_wait", willWait,
			"animation", f,
		)
		return f, nil
	}

	s.queue = s.queue[1:]

	s.stats.AbilUsageCountByChar[s.active][n.Typ.String()]++

	// s.Log.Infof("[%v] %v executing %v", s.Frame(), s.ActiveChar, a.Action)
	s.log.Debugw(
		"executed "+n.Typ.String(),
		"frame", s.f,
		"event", core.LogActionEvent,
		"char", s.active,
		"action", n.Typ.String(),
		"target", n.Target,
		"swap_cd_post", s.swapCD,
		"stam_post", s.stam,
		"animation", f,
	)

	return f, nil
}

func (s *Sim) StamPercentMod(a core.ActionType) float64 {
	var m float64 = 1
	for _, f := range s.stamModifier {
		m += f(a)
	}
	return m
}

func (s *Sim) AddStamMod(f func(a core.ActionType) float64) {
	s.stamModifier = append(s.stamModifier, f)
}

func (s *Sim) ResetAllNormalCounter() {
	for _, c := range s.chars {
		c.ResetNormalCounter()
	}
}
