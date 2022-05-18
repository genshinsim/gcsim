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
