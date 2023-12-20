// Package cooldown provides default implementation for SetCD, SetCDWithDelay, ResetActionCooldown, ReduceActionCooldown, ActionReady,
package character

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// SetCD takes two parameters:
//   - a action.Action: this is the action type we are triggering the cooldown for
//   - dur: duration in frames that the cooldown should last for
//
// It is assumed that AvailableCDCharges[a] > 0 (otherwise action should not have been allowed)
//
// SetCD works by adding the cooldown duration to a queue. This is because when there are
// multiple charges, the game will first finish recharging the first charge before starting
// the full cooldown for the second charge.
//
// When a cooldown is added to queue for the first time, a queue worker is started. This queue
// worker will check back at the cooldown specified for the first queued item, and if the queued
// cooldown did not change, it will increment the number of charges by 1, and reschedule itself
// to check back for the next item in queue
//
// Sometimes, the queued cooldown gets adjusted via ReduceActionCooldown or ResetActionCooldown.
// When this happens, the initial queued worker will check back at the wrong time. To prevent this,
// we use cdQueueWorkerStartedAt[a] which tracks the frame the worker started at. So when
// ReduceActionCooldown or ResetActionCooldown gets called, we start a new worker, updating
// cdQueueWorkerStartedAt[a] to represent the new worker start frame. This way the old worker can
// check this value first and then gracefully exit if it no longer matches its starting frame
func (c *Character) SetCD(a action.Action, dur int) {
	// setting cd is just adding a cd to the recovery queue
	// we need to check for cooldown reduction first to make sure the correct duration gets added
	modified := c.CDReduction(a, dur)
	// add current action and duration to the queue
	c.cdQueue[a] = append(c.cdQueue[a], modified)
	// if queue is zero before we added to it, then we'll start a cooldown queue worker
	if len(c.cdQueue[a]) == 1 {
		c.startCooldownQueueWorker(a)
	}
	// make sure to remove one from stack count
	c.AvailableCDCharge[a]--
	if c.AvailableCDCharge[a] < 0 {
		panic("unexpected charges less than 0")
	}
	c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index, a.String(), " cooldown triggered").
		Write("type", a.String()).
		Write("expiry", c.Cooldown(a)).
		Write("charges_remain", c.AvailableCDCharge).
		Write("original_cd", dur).
		Write("modified_cd_by_cdr", modified).
		Write("cooldown_queue", c.cdQueue)
}

func (c *Character) SetNumCharges(a action.Action, num int) {
	c.additionalCDCharge[a] = num - 1
	c.AvailableCDCharge[a] = num
}

func (c *Character) Charges(a action.Action) int {
	return c.AvailableCDCharge[a]
}

func (c *Character) SetCDWithDelay(a action.Action, dur, delay int) {
	if delay == 0 {
		c.SetCD(a, dur)
		return
	}
	c.Core.Tasks.Add(func() { c.SetCD(a, dur) }, delay)
}

func (c *Character) Cooldown(a action.Action) int {
	// remaining cooldown is src + first item in queue - current frame
	if c.AvailableCDCharge[a] > 0 {
		return 0
	}
	// otherwise check our queue; if zero then it's ready
	if len(c.cdQueue) == 0 {
		// panic("queue length is somehow 0??")
		return 0
	}
	return c.cdQueueWorkerStartedAt[a] + c.cdQueue[a][0] - c.Core.F
}

func (c *Character) ResetActionCooldown(a action.Action) {
	// if stacks already maxed then do nothing
	if c.AvailableCDCharge[a] == 1+c.additionalCDCharge[a] {
		return
	}
	// log.Printf("resetting; frame %v, queue %v\n", c.F, c.cdQueue[a])
	// otherwise add a stack && pop queue
	c.AvailableCDCharge[a]++
	c.Tags["skill_charge"]++
	c.cdQueue[a] = c.cdQueue[a][1:]
	// reset worker time
	c.cdQueueWorkerStartedAt[a] = c.Core.F
	c.cdCurrentQueueWorker[a] = nil
	c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index, a.String(), " cooldown forcefully reset").
		Write("type", a.String()).
		Write("charges_remain", c.AvailableCDCharge).
		Write("cooldown_queue", c.cdQueue)
	// check if anymore cd in queue
	if len(c.cdQueue) > 0 {
		c.startCooldownQueueWorker(a)
	}
}

func (c *Character) ReduceActionCooldown(a action.Action, v int) {
	// do nothing if stacks already maxed
	if c.AvailableCDCharge[a] == 1+c.additionalCDCharge[a] {
		return
	}
	// check if reduction > time remaing? if so then call reset cd
	remain := c.cdQueueWorkerStartedAt[a] + c.cdQueue[a][0] - c.Core.F
	// log.Printf("hello reducing; reduction %v, remaining %v, frame %v, old queue %v\n", v, remain, c.F, c.cdQueue[a])
	if v >= remain {
		c.ResetActionCooldown(a)
		return
	}
	// otherwise reduce remain and restart queue
	c.cdQueue[a][0] = remain - v
	c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index, a.String(), " cooldown forcefully reduced").
		Write("type", a.String()).
		Write("expiry", c.Cooldown(a)).
		Write("charges_remain", c.AvailableCDCharge).
		Write("cooldown_queue", c.cdQueue)
	c.startCooldownQueueWorker(a)
	// log.Printf("started: %v, new queue: %v, worker frame: %v\n", c.cdQueueWorkerStartedAt[a], c.cdQueue[a], c.cdQueueWorkerStartedAt[a])
}

func (c *Character) startCooldownQueueWorker(a action.Action) {
	// check the length of the queue for action a, if there's nothing then there's
	// nothing to start
	if len(c.cdQueue[a]) == 0 {
		return
	}

	// set the time we starter this worker at
	c.cdQueueWorkerStartedAt[a] = c.Core.F
	var src *func()

	worker := func() {
		// check if src changed; if so do nothing
		if src != c.cdCurrentQueueWorker[a] {
			// c.Log.Debugw("src changed",  "src", src, "new", c.cdQueueWorkerStartedAt[a])
			return
		}
		// log.Printf("cd worker triggered, started; %v, queue: %v\n", c.cdQueueWorkerStartedAt[a], c.cdQueue[a])
		// check to make sure queue is not 0
		if len(c.cdQueue[a]) == 0 {
			// this should never happen
			panic(fmt.Sprintf(
				"queue is empty? index :%v, frame : %v, worker src: %v, started: %v",
				c.Index,
				c.Core.F,
				src,
				c.cdQueueWorkerStartedAt[a],
			))
			// return
		}
		// otherwise add a stack and pop first item in queue
		c.AvailableCDCharge[a]++
		c.Tags["skill_charge"]++
		c.cdQueue[a] = c.cdQueue[a][1:]

		// c.Log.Debugw("stack restored",  "avail", c.availableCDCharge[a], "queue", c.cdQueue)

		if c.AvailableCDCharge[a] > 1+c.additionalCDCharge[a] {
			// sanity check, this should never happen
			panic(fmt.Sprintf("charges > max? index :%v, frame : %v", c.Index, c.Core.F))
		}

		c.Core.Log.NewEventBuildMsg(glog.LogCooldownEvent, c.Index, a.String(), " cooldown ready").
			Write("type", a.String()).
			Write("charges_remain", c.AvailableCDCharge).
			Write("cooldown_queue", c.cdQueue)

		// if queue still has len > 0 then call start queue again
		if len(c.cdQueue) > 0 {
			c.startCooldownQueueWorker(a)
		}
	}

	c.cdCurrentQueueWorker[a] = &worker
	src = &worker

	// wait for c.cooldownQueue[a][0], then add a stack
	c.Core.Tasks.Add(worker, c.cdQueue[a][0])
}
