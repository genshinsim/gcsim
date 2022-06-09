package glog

// Event describes one event to be logged
type Event interface {
	LogSource() Source              //returns the type of this log event i.e. character, sim, damage, etc...
	StartFrame() int                //returns the frame on which this event was started
	Src() int                       //returns the index of the character that triggered this event. -1 if it's not a character
	Write(keyAndVal ...interface{}) //write additional keyAndVal pairs to the event
	SetEnded(f int)
}

// Logger records LogEvents
type Logger interface {
	NewEvent(msg string, typ Source, srcChar int, keysAndValues ...interface{}) Event
	NewEventBuildMsg(typ Source, srcChar int, msg ...string) Event
	Dump() ([]byte, error) //print out all the logged events in array of JSON strings in the ordered they were added
}
