package queue

import (
	"context"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type queueBlock struct {
	work  []*model.ComputeWork //this is the priority queue
	index map[string]*model.ComputeWork
	wip   map[string]wipWork
	log   *zap.SugaredLogger
}

// TODO: we need to handle timeouts on these requests somehow???
func (s *Server) queueCtrl() {
	//single thread keep tracking of queue etc...
	q := queueBlock{
		index: make(map[string]*model.ComputeWork),
		wip:   make(map[string]wipWork),
		log:   s.Log,
	}
	//when request for work comes in, pop off the work at the front and mark it as wip
	//as wip expires (no result came in), they should be appended again to the front of the queue
	//wip is scanned periodically for expiry
	wipTicker := time.NewTicker(10 * time.Second)

	//regular ticker to poll for work;
	pollTicker := time.NewTicker(30 * time.Second)

	//populate at start
	s.populate(&q)

	for {
		select {
		case req := <-s.getWork:
			w := q.pop()
			if w != nil {
				q.wip[w.GetId()] = wipWork{
					w:      w,
					Expiry: time.Now().Add(s.Timeout),
				}
			}
			s.Log.Infow("found next work", "id", w.GetId(), "wip", q.wip)
			req.resp <- w
		case req := <-s.completeWork:
			s.Log.Infow("complete work request", "id", req.id, "wip", q.wip)
			//work should exist in wip; other wise NotFound
			if _, ok := q.wip[req.id]; !ok {
				req.resp <- status.Error(codes.NotFound, "work not found")
			} else {
				//ok to delete from wip
				delete(q.wip, req.id)
				req.resp <- nil
			}
		case <-wipTicker.C:
			s.Log.Info("purging expired wip...")
			expired := q.purge()
			s.Log.Infow("purge done", "expired", expired, "wip", q.wip)
		case <-pollTicker.C:
			s.populate(&q)
		}

	}
}

func (s *Server) populate(q *queueBlock) {
	s.Log.Info("polling for work")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

	w, err := s.DBWork.GetWork(ctx)
	if err != nil {
		s.Log.Infow("error getting work from db", "err", err)
	} else {
		res := q.populate(w)
		s.Log.Infow("db work populated", "res", res)
	}

	w, err = s.SubWork.GetWork(ctx)
	if err != nil {
		s.Log.Infow("error getting work from sub", "err", err)
	} else {
		res := q.populate(w)
		s.Log.Infow("sub work populated", "res", res)
	}
	cancel()
}

func (q *queueBlock) purge() []string {
	var expired []string
	now := time.Now()
	//may not need to do it this way? just to be safe i guess
	var keys []string
	//this is a slow operation...
	for k := range q.wip {
		keys = append(keys, k)

	}
	for _, k := range keys {
		v := q.wip[k]
		if now.After(v.Expiry) {
			//add it back
			q.insert(v.w)
			delete(q.wip, k)
			expired = append(expired, k)
		}
	}
	return expired
}

func (q *queueBlock) populate(next []*model.ComputeWork) []string {
	var added []string
	//check both wip and index; only add to end of work if not in either
	for _, w := range next {
		id := w.GetId()
		if id == "" {
			//this shouldnt be happening?
			continue
		}
		if _, ok := q.index[id]; ok {
			continue
		}
		if _, ok := q.wip[id]; ok {
			continue
		}
		q.append(w)
		added = append(added, id)

	}
	return added
}

func (q *queueBlock) insert(w *model.ComputeWork) {
	q.work = append([]*model.ComputeWork{w}, q.work...)
	q.index[w.GetId()] = w
}

func (q *queueBlock) append(w *model.ComputeWork) {
	q.work = append(q.work, w)
	q.index[w.GetId()] = w
}

func (q *queueBlock) pop() *model.ComputeWork {
	if len(q.work) == 0 {
		return nil
	}
	x := q.work[0]
	q.work = q.work[1:]
	//purge it from index; if popping for getting job make sure to add to wip
	delete(q.index, x.GetId())
	return x
}
