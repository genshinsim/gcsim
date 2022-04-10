package calcqueue

import "github.com/genshinsim/gcsim/pkg/core"

type Queuer struct {
	core *core.Core
	pq   []core.ActionBlock
	ind  int
	wait int
}

func New(c *core.Core) *Queuer {
	return &Queuer{
		core: c,
		wait: -1,
	}
}

func (q *Queuer) SetActionList(a []core.ActionBlock) error {
	q.pq = a
	q.core.Log.NewEvent(
		"priority queued received",
		core.LogQueueEvent,
		-1,
		"pq", a,
	)
	return nil
}

func (q *Queuer) Next() ([]core.Command, bool, error) {

	var r []core.Command
	active := q.core.Chars[q.core.ActiveChar].Key()
	for {
		if q.ind >= len(q.pq) {
			q.core.Log.NewEvent(
				"no more actions",
				core.LogQueueEvent,
				-1,
			)
			return nil, false, nil
		}
		//we only care for either wait or sequence; macro and anything else can be ignored
		//also, go down the list 1 at a time
		v := q.pq[q.ind]

		//check if we're on wait
		if q.wait > q.core.F {
			q.core.Log.NewEvent(
				"queuer on wait",
				core.LogQueueEvent,
				-1,
				"wait", q.wait,
				"index", q.ind,
			)
			return nil, false, nil
		}

		switch v.Type {
		case core.ActionBlockTypeCalcRestart:
			q.ind = 0
			return nil, false, nil
		case core.ActionBlockTypeCalcWait:
			//depending on the type of wait here
			if v.CalcWait.Frames {
				q.wait = v.CalcWait.Val
			} else {
				q.wait = q.core.F + v.CalcWait.Val
			}
		case core.ActionBlockTypeSequence:
			//check if we need to swap
			if active != v.SequenceChar {
				r = append(r, &core.ActionItem{
					Typ:    core.ActionSwap,
					Target: v.SequenceChar,
				})
			}
			//add the rest
			for i := 0; i < len(v.Sequence); i++ {
				r = append(r, &v.Sequence[i])
			}
			// for _, v := range v.Sequence {
			// 	r = append(r, &v)
			// }
		default:
			//ignore and move on
			q.ind++
			q.core.Log.NewEvent(
				"queuer skipping non sequence options",
				core.LogQueueEvent,
				-1,
				"index", q.ind,
				"type", v.Type,
			)
			continue
		}

		// q.core.Log.Debugw(
		// 	"item queued",
		// 	"frame", q.core.F,
		// 	core.LogQueueEvent,
		// 	"index", q.ind,
		// 	"name", v.Name,
		// 	"target", v.Target,
		// 	"is seq", v.IsSeq,
		// 	"strict", v.IsStrict,
		// 	"exec", v.Exec,
		// 	"once", v.Once,
		// 	"post", v.PostAction.String(),
		// 	"swap_to", v.SwapTo,
		// 	"raw", v.Raw,
		// )

		q.ind++
		return r, false, nil
	}

}
