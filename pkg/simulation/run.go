package simulation

import "github.com/genshinsim/gcsim/pkg/simulation/queue"

func (s *Simulation) Run() (Result, error) {
	return s.stats, nil
}

func (s *Simulation) AdvanceFrame() error {
	var err error
	var done bool
	for !done {
		done, err = s.queueAndExec()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Simulation) queueAndExec() (bool, error) {
	//TODO: we should do some optimization here to at least skip some frames
	//if we know for sure the action won't be ready

	//if queue empty, grab item from queue
	if len(s.queue) == 0 {
		next, dropIfFailed, err := s.queuer.Next()
		if err != nil {
			return false, err
		}
		//if queue is empty, then there's nothing to do
		//so we skip this frame
		if len(next) == 0 {
			return true, nil
		}
		s.dropQueueIfFailed = dropIfFailed
		s.queue = append(s.queue, next...)
	}

	//try executing first item in queue, if failed b/c not ready, skip frame
	a, isAction := s.queue[0].(*queue.ActionItem)

	//check if we need to purge queue

	return true, nil
}

//executes the provided command, returns 2 booleans:
//	- executed if this action was successfully executed and should be purged from queue
//	- frameDone if this action consumes a frame
func (s *Simulation) execCommand(c queue.Command) (executed, frameDone bool) {

	return
}

func (s *Simulation) QueueNext() {

}
