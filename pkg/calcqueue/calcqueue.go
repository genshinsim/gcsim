package calcqueue

import "github.com/genshinsim/gcsim/pkg/core"

type Queuer struct {
	core *core.Core
	prio []core.Action
	ind  int
}

func New(c *core.Core) *Queuer {
	return &Queuer{
		core: c,
	}
}

func (q *Queuer) SetActionList(a []core.Action) {
	q.prio = a
}

func (q *Queuer) Next() ([]core.ActionItem, error) {

	var r []core.ActionItem
	active := q.core.Chars[q.core.ActiveChar].Key()
	for {
		if q.ind >= len(q.prio) {
			q.core.Log.Debugw(
				"no more actions",
				"frame", q.core.F,
				"event", core.LogQueueEvent,
			)
			return nil, nil
		}
		//we only accept action+=?? target=character wait=150
		//also, go down the list 1 at a time
		v := q.prio[q.ind]

		if v.IsSeq {
			//ignore and move on
			q.ind++
			continue
		}

		//check wait
		if v.Wait > q.core.F {
			q.core.Log.Debugw(
				"item on wait",
				"frame", q.core.F,
				"event", core.LogQueueEvent,
				"wait", v.Wait,
				"index", q.ind,
				"name", v.Name,
				"target", v.Target,
				"is seq", v.IsSeq,
				"strict", v.IsStrict,
				"exec", v.Exec,
				"once", v.Once,
				"post", v.PostAction.String(),
				"swap_to", v.SwapTo,
				"raw", v.Raw,
			)
			return nil, nil
		}

		//check if we need to swap
		if active != v.Target {
			r = append(r, core.ActionItem{
				Target: v.Target,
				Typ:    core.ActionSwap,
			})
		}

		r = append(r, v.Exec[0])

		q.core.Log.Debugw(
			"item queued",
			"frame", q.core.F,
			"event", core.LogQueueEvent,
			"index", q.ind,
			"name", v.Name,
			"target", v.Target,
			"is seq", v.IsSeq,
			"strict", v.IsStrict,
			"exec", v.Exec,
			"once", v.Once,
			"post", v.PostAction.String(),
			"swap_to", v.SwapTo,
			"raw", v.Raw,
		)

		q.ind++
		return r, nil
	}

}
