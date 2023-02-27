package queue

import (
	context "context"
	"time"

	"github.com/genshinsim/gcsim/pkg/model"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotifyService interface {
	Notify(topic string, message string)
}

type Config struct {
	Timeout time.Duration
}

type Server struct {
	Config
	Log *zap.SugaredLogger
	UnimplementedWorkQueueServer
	addToQueue   chan addToQueueReq
	getWork      chan getWorkReq
	completeWork chan completeWorkReq
}

type getWorkReq struct {
	resp chan *model.ComputeWork
}

type completeWorkReq struct {
	id   string
	resp chan error
}

type addToQueueReq struct {
	w    []*model.ComputeWork
	resp chan []string
}

type Work struct {
	Key  string
	Work any
}

func NewQueue(cfg Config, cust ...func(*Server) error) (*Server, error) {
	s := &Server{
		Config:       cfg,
		addToQueue:   make(chan addToQueueReq),
		getWork:      make(chan getWorkReq),
		completeWork: make(chan completeWorkReq),
	}
	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}
	if s.Log == nil {
		logger, err := zap.NewProduction()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()

		s.Log = sugar
	}

	if s.Timeout <= 0 {
		s.Timeout = 5 * time.Minute //default 5 min
	}

	go s.queueCtrl()

	s.Log.Infow("queue service started")
	return s, nil
}

type queueBlock struct {
	work  []*model.ComputeWork //this is the priority queue
	index map[string]*model.ComputeWork
	wip   map[string]wipWork
}

type wipWork struct {
	w      *model.ComputeWork
	expiry time.Time
}

func (s *Server) Populate(ctx context.Context, req *PopulateReq) (*PopulateResp, error) {
	d := req.GetData()
	s.Log.Infow("populate queue req received", "length", len(d))
	if len(d) == 0 {
		return nil, status.Error(codes.InvalidArgument, "populate request data is blank")
	}
	resp := make(chan []string)
	s.addToQueue <- addToQueueReq{
		w:    d,
		resp: resp,
	}
	res := <-resp
	s.Log.Infow("populate done", "added", res)
	return &PopulateResp{
		Ids: res,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *GetReq) (*GetResp, error) {
	s.Log.Infow("get work request received")
	resp := make(chan *model.ComputeWork)
	s.getWork <- getWorkReq{
		resp: resp,
	}
	res := <-resp
	return &GetResp{
		Data: res,
	}, nil
}

func (s *Server) Complete(ctx context.Context, req *CompleteReq) (*CompleteResp, error) {
	id := req.GetId()
	s.Log.Infow("complete work request", "id", id)
	resp := make(chan error)
	s.completeWork <- completeWorkReq{
		resp: resp,
	}
	err := <-resp
	if err != nil {
		return nil, err
	}
	return &CompleteResp{}, nil
}

// TODO: we need to handle timeouts on these requests somehow???
func (s *Server) queueCtrl() {
	//single thread keep tracking of queue etc...
	q := queueBlock{
		index: make(map[string]*model.ComputeWork),
	}
	//when request for work comes in, pop off the work at the front and mark it as wip
	//as wip expires (no result came in), they should be appended again to the front of the queue
	//wip is scanned periodically for expiry
	wipTicker := time.NewTicker(10 * time.Second)

	for {
		select {
		case x := <-s.addToQueue:
			//this comes in periodically for the scheduled pulls
			//we need to resolve (and block) at this point which is alerady in queue (ignore)
			//and which are new (add to end)
			res := q.populate(x.w)
			x.resp <- res
		case req := <-s.getWork:
			w := q.pop()
			if w != nil {
				q.wip[w.GetId()] = wipWork{
					w:      w,
					expiry: time.Now().Add(s.Timeout),
				}
			}
			req.resp <- w
		case req := <-s.completeWork:
			//work should exist in wip; other wise NotFound
			if _, ok := q.wip[req.id]; !ok {
				req.resp <- status.Error(codes.NotFound, "work not found")
			}
			//ok to delete from wip
			delete(q.wip, req.id)
			req.resp <- nil
		case t := <-wipTicker.C:
			s.Log.Infow("purging expired wip...", "t", t.Format(time.RFC1123))
			expired := q.purge()
			s.Log.Infow("purge done", "expired", expired)
		}

	}
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
		if now.After(v.expiry) {
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
