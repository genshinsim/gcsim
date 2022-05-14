package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
)

//SetCD takes two parameters:
//	- a core.ActionType: this is the action type we are triggering the cooldown for
//  - dur: duration in frames that the cooldown should last for
//It is assumed that AvailableCDCharges[a] > 0 (otherwise action should not have been allowed)
//
//SetCD works by adding the cooldown duration to a queue. This is because when there are
//multiple charges, the game will first finish recharging the first charge before starting
//the full cooldown for the second charge.
//
//When a cooldown is added to queue for the first time, a queue worker is started. This queue
//worker will check back at the cooldown specified for the first queued item, and if the queued
//cooldown did not change, it will increment the number of charges by 1, and reschedule itself
//to check back for the next item in queue
//
//Sometimes, the queued cooldown gets adjusted via ReduceActionCooldown or ResetActionCooldown.
//When this happens, the initial queued worker will check back at the wrong time. To prevent this,
//we use cdQueueWorkerStartedAt[a] which tracks the frame the worker started at. So when
//ReduceActionCooldown or ResetActionCooldown gets called, we start a new worker, updating
//cdQueueWorkerStartedAt[a] to represent the new worker start frame. This way the old worker can
//check this value first and then gracefully exit if it no longer matches its starting frame
func (c *Tmpl) SetCD(a core.ActionType, dur int) {
	//setting cd is just adding a cd to the recovery queue
	//add current action and duration to the queue
	c.cdQueue[a] = append(c.cdQueue[a], dur)
	//if queue is zero before we added to it, then we'll start a cooldown queue worker
	if len(c.cdQueue[a]) == 1 {
		c.startCooldownQueueWorker(a, true)
	}
	//make sure to remove one from stack count
	c.AvailableCDCharge[a]--
	if c.AvailableCDCharge[a] < 0 {
		panic("unexpected charges less than 0")
	}
	if c.Tags["skill_charge"] > 0 {
		c.Tags["skill_charge"]--
	}
	c.Core.Log.NewEventBuildMsg(
		core.LogActionEvent,
		c.Index,
		a.String(), " cooldown triggered",
	).Write(
		"type", a.String(),
		"expiry", c.Cooldown(a),
		"charges_remain", c.AvailableCDCharge,
		"cooldown_queue", c.cdQueue,
	)
}

func (c *Tmpl) SetNumCharges(a core.ActionType, num int) {
	c.additionalCDCharge[a] = num - 1
	c.AvailableCDCharge[a] = num
}

func (c *Tmpl) Charges(a core.ActionType) int {
	return c.AvailableCDCharge[a]
}

//TODO: energy calc mode currently not working. we need to reset action or else we might get
//charges < 0 while trying to SetCD
func (c *Tmpl) ActionReady(a core.ActionType, p map[string]int) bool {
	//cannot execute actions if character is dead
	if c.HPCurrent <= 0 {
		return false
	}
	//up if energy is ready && stack > 0
	if a == core.ActionBurst && (c.Energy != c.EnergyMax) {
		return false
	}
	return c.AvailableCDCharge[a] > 0
}

func (c *Tmpl) SetCDWithDelay(a core.ActionType, dur int, delay int) {
	if delay == 0 {
		c.SetCD(a, dur)
		return
	}
	c.AddTask(func() { c.SetCD(a, dur) }, "set-cd", delay)
}

func (c *Tmpl) Cooldown(a core.ActionType) int {
	//remaining cooldown is src + first item in queue - current frame
	if c.AvailableCDCharge[a] > 0 {
		return 0
	}
	//otherwise check our queue; if zero then it's ready
	if len(c.cdQueue) == 0 {
		// panic("queue length is somehow 0??")
		return 0
	}
	return c.cdQueueWorkerStartedAt[a] + c.cdQueue[a][0] - c.Core.F
}

func (c *Tmpl) ResetActionCooldown(a core.ActionType) {
	//if stacks already maxed then do nothing
	if c.AvailableCDCharge[a] == 1+c.additionalCDCharge[a] {
		return
	}
	//log.Printf("resetting; frame %v, queue %v\n", c.Core.F, c.cdQueue[a])
	//otherwise add a stack && pop queue
	c.AvailableCDCharge[a]++
	c.Tags["skill_charge"]++
	c.cdQueue[a] = c.cdQueue[a][1:]
	//reset worker time
	c.cdQueueWorkerStartedAt[a] = c.Core.F
	c.cdCurrentQueueWorker[a] = nil
	c.Core.Log.NewEventBuildMsg(
		core.LogActionEvent,
		c.Index,
		a.String(), " cooldown forcefully reset",
	).Write(
		"type", a.String(),
		"charges_remain", c.AvailableCDCharge,
		"cooldown_queue", c.cdQueue,
	)
	//check if anymore cd in queue
	if len(c.cdQueue) > 0 {
		c.startCooldownQueueWorker(a, true)
	}
}

func (c *Tmpl) ReduceActionCooldown(a core.ActionType, v int) {
	//do nothing if stacks already maxed
	if c.AvailableCDCharge[a] == 1+c.additionalCDCharge[a] {
		return
	}
	//check if reduction > time remaing? if so then call reset cd
	remain := c.cdQueueWorkerStartedAt[a] + c.cdQueue[a][0] - c.Core.F
	//log.Printf("hello reducing; reduction %v, remaining %v, frame %v, old queue %v\n", v, remain, c.Core.F, c.cdQueue[a])
	if v >= remain {
		c.ResetActionCooldown(a)
		return
	}
	//otherwise reduce remain and restart queue
	c.cdQueue[a][0] = remain - v
	c.Core.Log.NewEventBuildMsg(
		core.LogActionEvent,
		c.Index,
		a.String(), " cooldown forcefully reduced",
	).Write(
		"type", a.String(),
		"expiry", c.Cooldown(a),
		"charges_remain", c.AvailableCDCharge,
		"cooldown_queue", c.cdQueue,
	)
	c.startCooldownQueueWorker(a, false)
	//log.Printf("started: %v, new queue: %v, worker frame: %v\n", c.cdQueueWorkerStartedAt[a], c.cdQueue[a], c.cdQueueWorkerStartedAt[a])
}

func (c *Tmpl) startCooldownQueueWorker(a core.ActionType, cdReduct bool) {
	//check the length of the queue for action a, if there's nothing then there's
	//nothing to start
	if len(c.cdQueue[a]) == 0 {
		return
	}

	//set the time we starter this worker at
	c.cdQueueWorkerStartedAt[a] = c.Core.F
	var src *func()

	//reduce the first item by the current cooldown reduction
	if cdReduct {
		c.cdQueue[a][0] = c.cdReduction(a, c.cdQueue[a][0])
	}

	worker := func() {
		//check if src changed; if so do nothing
		if src != c.cdCurrentQueueWorker[a] {
			// c.Core.Log.Debugw("src changed",  "src", src, "new", c.cdQueueWorkerStartedAt[a])
			return
		}
		//log.Printf("cd worker triggered, started; %v, queue: %v\n", c.cdQueueWorkerStartedAt[a], c.cdQueue[a])
		//check to make sure queue is not 0
		if len(c.cdQueue[a]) == 0 {
			//this should never happen
			panic(fmt.Sprintf("queue is empty? character :%v, frame : %v, worker src: %v, started: %v", c.Name(), c.Core.F, src, c.cdQueueWorkerStartedAt[a]))
			// return
		}
		//otherwise add a stack and pop first item in queue
		c.AvailableCDCharge[a]++
		c.Tags["skill_charge"]++
		c.cdQueue[a] = c.cdQueue[a][1:]

		// c.Core.Log.Debugw("stack restored",  "avail", c.availableCDCharge[a], "queue", c.cdQueue)

		if c.AvailableCDCharge[a] > 1+c.additionalCDCharge[a] {
			//sanity check, this should never happen
			panic(fmt.Sprintf("charges > max? character :%v, frame : %v", c.Name(), c.Core.F))
			// c.availableCDCharge[a] = 1 + c.additionalCDCharge[a]
			// return
		}

		c.Core.Log.NewEventBuildMsg(
			core.LogActionEvent,
			c.Index,
			a.String(), " cooldown ready",
		).Write(
			"type", a.String(),
			"charges_remain", c.AvailableCDCharge,
			"cooldown_queue", c.cdQueue,
		)

		//if queue still has len > 0 then call start queue again
		if len(c.cdQueue) > 0 {
			c.startCooldownQueueWorker(a, true)
		}

	}

	c.cdCurrentQueueWorker[a] = &worker
	src = &worker

	//wait for c.cooldownQueue[a][0], then add a stack
	c.AddTask(worker, "cooldown-worker-"+a.String(), c.cdQueue[a][0])

}

func (c *Tmpl) cdReduction(a core.ActionType, dur int) int {
	var cd float64 = 1
	n := 0
	for _, v := range c.CDReductionFuncs {
		//if not expired
		if v.Expiry == -1 || v.Expiry > c.Core.F {
			amt := v.Amount(a)
			c.Core.Log.NewEvent("applying cooldown modifier", core.LogActionEvent, c.Index, "key", v.Key, "modifier", amt, "expiry", v.Expiry)
			cd += amt
			c.CDReductionFuncs[n] = v
			n++
		}
	}
	c.CDReductionFuncs = c.CDReductionFuncs[:n]

	return int(float64(dur) * cd)
}

func (c *Tmpl) AddCDAdjustFunc(rd core.CDAdjust) {
	ind := -1
	for i, v := range c.CDReductionFuncs {
		//if expired already, set to nil and ignore
		if v.Key == rd.Key {
			ind = i
		}
	}
	if ind > -1 {
		c.CDReductionFuncs[ind] = rd
	} else {
		c.CDReductionFuncs = append(c.CDReductionFuncs, rd)
	}
}
