package yaemiko

import "github.com/genshinsim/gcsim/pkg/core"

func (c *char) ActionReady(a core.ActionType, p map[string]int) bool {
	//up if energy is ready && stack > 0
	if a == core.ActionBurst && (c.Energy != c.EnergyMax) && !c.Core.Flags.EnergyCalcMode {
		return false
	}
	return c.availableCDCharge[a] > 0
}

func (c *char) SetCD(a core.ActionType, dur int) {
	//setting cd is just adding a cd to the recovery queue
	//add current action and duration to the queue
	c.cdQueue[a] = append(c.cdQueue[a], dur)
	//if queue is zero before we added to it, then we'll start a cooldown queue worker
	if len(c.cdQueue[a]) == 1 {
		c.startCooldownQueueWorker(a, true)
	}
	//make sure to remove one from stack count
	c.availableCDCharge[a]--
	if c.availableCDCharge[a] < 0 {
		panic("unexpected charges less than 0")
	}
	c.Core.Log.Debugw(
		a.String()+" cooldown triggered",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"char", c.Index,
		"type", a.String(),
		"charges_remain", c.availableCDCharge,
		"cooldown_queue", c.cdQueue,
	)
}

func (c *char) Cooldown(a core.ActionType) int {
	//remaining cooldown is src + first item in queue - current frame
	if c.availableCDCharge[a] > 0 {
		return 0
	}
	//otherwise check our queue
	if len(c.cdQueue) == 0 {
		panic("queue length is somehow 0??")
	}
	return c.cdQueueWorkerStartedAt[a] + c.cdQueue[a][0] - c.Core.F
}

func (c *char) ResetActionCooldown(a core.ActionType) {
	//if stacks already maxed then do nothing
	if c.availableCDCharge[a] == 1+c.additionalCDCharge[a] {
		return
	}
	//otherwise add a stack && pop queue
	c.availableCDCharge[a]++
	c.cdQueue[a] = c.cdQueue[a][1:]
	//reset worker time
	c.cdQueueWorkerStartedAt[a] = c.Core.F
	c.Core.Log.Debugw(
		a.String()+" cooldown forcefully reset",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"char", c.Index,
		"type", a.String(),
		"charges_remain", c.availableCDCharge,
		"cooldown_queue", c.cdQueue,
	)
	//check if anymore cd in queue
	if len(c.cdQueue) > 0 {
		c.startCooldownQueueWorker(a, true)
	}
}

func (c *char) ReduceActionCooldown(a core.ActionType, v int) {
	//do nothing if stacks already maxed
	if c.availableCDCharge[a] == 1+c.additionalCDCharge[a] {
		return
	}
	//check if reduction > time remaing? if so then call reset cd
	remain := c.cdQueueWorkerStartedAt[a] + c.cdQueue[a][0] - c.Core.F
	if v > remain {
		c.ResetActionCooldown(a)
		return
	}
	//otherwise reduce remain and restart queue
	c.cdQueue[a][0] = remain - v
	c.Core.Log.Debugw(
		a.String()+" cooldown forcefully reduced",
		"frame", c.Core.F,
		"event", core.LogCharacterEvent,
		"char", c.Index,
		"type", a.String(),
		"charges_remain", c.availableCDCharge,
		"cooldown_queue", c.cdQueue,
	)
	c.startCooldownQueueWorker(a, false)
}

func (c *char) startCooldownQueueWorker(a core.ActionType, cdReduct bool) {
	//check the length of the queue for action a, if there's nothing then there's
	//nothing to start
	if len(c.cdQueue[a]) == 0 {
		return
	}
	//set the time we starter this worker at
	c.cdQueueWorkerStartedAt[a] = c.Core.F
	src := c.Core.F

	//reduce the first item by the current cooldown reduction
	if cdReduct {
		c.cdQueue[a][0] = c.cdReduction(a, c.cdQueue[a][0])
	}

	//wait for c.cooldownQueue[a][0], then add a stack
	c.AddTask(func() {
		//check if src changed; if so do nothing
		if src != c.cdQueueWorkerStartedAt[a] {
			// c.Core.Log.Debugw("src changed", "frame", c.Core.F, "src", src, "new", c.cdQueueWorkerStartedAt[a])
			return
		}
		//check to make sure queue is not 0
		if len(c.cdQueue[a]) == 0 {
			//this should never happen
			panic("charges > max??")
			// return
		}
		//otherwise add a stack and pop first item in queue
		c.availableCDCharge[a]++
		c.cdQueue[a] = c.cdQueue[a][1:]

		// c.Core.Log.Debugw("stack restored", "frame", c.Core.F, "avail", c.availableCDCharge[a], "queue", c.cdQueue)

		if c.availableCDCharge[a] > 1+c.additionalCDCharge[a] {
			//sanity check, this should never happen
			panic("charges > max??")
			// c.availableCDCharge[a] = 1 + c.additionalCDCharge[a]
			// return
		}

		c.Core.Log.Debugw(
			a.String()+" cooldown ready",
			"frame", c.Core.F,
			"event", core.LogCharacterEvent,
			"char", c.Index,
			"type", a.String(),
			"charges_remain", c.availableCDCharge,
			"cooldown_queue", c.cdQueue,
		)

		//if queue still has len > 0 then call start queue again
		if len(c.cdQueue) > 0 {
			c.startCooldownQueueWorker(a, true)
		}

	}, "cooldown-worker-"+a.String(), c.cdQueue[a][0])
}

func (c *char) cdReduction(a core.ActionType, dur int) int {
	var cd float64 = 1
	n := 0
	for _, v := range c.CDReductionFuncs {
		//if not expired
		if v.Expiry == -1 || v.Expiry > c.Core.F {
			amt := v.Amount(a)
			c.Core.Log.Debugw("applying cooldown modifier", "frame", c.Core.F, "event", core.LogCharacterEvent, "char", c.Index, "key", v.Key, "modifier", amt, "expiry", v.Expiry)
			cd += amt
			c.CDReductionFuncs[n] = v
			n++
		}
	}
	c.CDReductionFuncs = c.CDReductionFuncs[:n]

	return int(float64(dur) * cd)
}
