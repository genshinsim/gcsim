package core

type CommandHandler interface {
	Exec(n Command) (frames int, done bool, err error) //return frames, if executed, any errors
}

type CommandType int

const (
	CommandTypeAction CommandType = iota
	CommandTypeWait
	CommandTypeNoSwap
	CommandTypeResetLimit
)

//Command is what gets executed by the sim.
type Command interface {
	Type() CommandType
}

type CmdResetLimit struct {
}

func (c CmdResetLimit) Type() CommandType { return CommandTypeResetLimit }

type CmdWaitType int

const (
	CmdWaitTypeInvalid CmdWaitType = iota
	CmdWaitTypeTimed
	CmdWaitTypeParticle
	CmdWaitTypeMods
)

type CmdWait struct {
	For        CmdWaitType
	Max        int //cannot be 0 if type is timed
	Source     string
	Conditions Condition
	FillAction ActionItem
}

func (c *CmdWait) Clone() CmdWait {
	next := *c
	next.Conditions = c.Conditions.Clone()
	next.FillAction = c.FillAction.Clone()
	return next
}

type CmdCalcWait struct {
	Frames bool
	Val    int
}

func (c *CmdWait) Type() CommandType { return CommandTypeWait }

type CmdNoSwap struct {
	Val int
}

func (c *CmdNoSwap) Type() CommandType { return CommandTypeNoSwap }
