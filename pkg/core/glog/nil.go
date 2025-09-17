package glog

// NilLogEvent implements LogEvent and is used when no logger is needed
type NilLogEvent struct{}

func (n *NilLogEvent) LogSource() Source                    { return LogSimEvent }
func (n *NilLogEvent) StartFrame() int                      { return -1 }
func (n *NilLogEvent) Src() int                             { return 0 }
func (n *NilLogEvent) WriteBuildMsg(keyAndVal ...any) Event { return n }
func (n *NilLogEvent) Write(key string, val any) Event      { return n }
func (n *NilLogEvent) SetEnded(f int) Event                 { return n }

// NilLogger implements Logger and is used when no logger is needed
type NilLogger struct{}

func (n *NilLogger) Dump() ([]byte, error) { return []byte{}, nil }
func (n *NilLogger) NewEventBuildMsg(typ Source, srcChar int, msg ...string) Event {
	return &NilLogEvent{}
}

func (n *NilLogger) NewEvent(msg string, typ Source, srcChar int) Event {
	return &NilLogEvent{}
}
