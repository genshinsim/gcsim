package stats

import "github.com/genshinsim/gcsim/pkg/core"

type CollectorCustomStats[T any] interface {
	Flush(core *core.Core) T
}

type NewStatsFuncCustomStats[T any] func(core *core.Core) (CollectorCustomStats[T], error)
