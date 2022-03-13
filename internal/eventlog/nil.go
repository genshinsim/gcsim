package eventlog

import "github.com/genshinsim/gcsim/pkg/coretype"

type NilLogEvent struct{}

func (n *NilLogEvent) LogSource() coretype.LogSource  { return coretype.LogSimEvent }
func (n *NilLogEvent) StartFrame() int                { return -1 }
func (n *NilLogEvent) Src() int                       { return 0 }
func (n *NilLogEvent) Write(keyAndVal ...interface{}) {}
func (n *NilLogEvent) SetEnded(f int)                 {}

type NilLogger struct{}

func (n *NilLogger) Dump() ([]byte, error) { return []byte{}, nil }
func (n *NilLogger) NewEventBuildMsg(typ coretype.LogSource, srcChar int, msg ...string) coretype.LogEvent {
	return &NilLogEvent{}
}
func (n *NilLogger) NewEvent(msg string, typ coretype.LogSource, srcChar int, keysAndValues ...interface{}) coretype.LogEvent {
	return &NilLogEvent{}
}
