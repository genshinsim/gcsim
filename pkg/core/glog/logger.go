package glog

import (
	"log"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

// Debugw
// Warnw

// https://github.com/dominikh/go-tools/issues/836

//nolint:staticcheck
type keyVal struct {
	Key string      `json:"key,nocopy"`
	Val interface{} `json:"val"`
}

// https://github.com/dominikh/go-tools/issues/836

//nolint:staticcheck
//easyjson:json
type LogEvent struct {
	Typ      Source                 `json:"event"`
	F        int                    `json:"frame"`
	Ended    int                    `json:"ended"`
	SrcChar  int                    `json:"char_index"`
	Msg      string                 `json:"msg,nocopy"`
	Logs     map[string]interface{} `json:"logs"`
	Ordering map[string]int         `json:"ordering"`
	counter  int
}

//easyjson:json
type EventArr []*LogEvent

func (e *LogEvent) Write(key string, value interface{}) Event {
	e.Logs[key] = value
	e.Ordering[key] = e.counter
	e.counter++

	return e
}

func (e *LogEvent) WriteBuildMsg(keysAndValues ...interface{}) Event {
	// should be even number
	var key string
	var ok bool
	for i := 0; i < len(keysAndValues); i++ {
		key, ok = keysAndValues[i].(string)
		if !ok {
			log.Panicf("invalid key %v, expected type to be string", keysAndValues[i].(string))
		}
		// make sure there's a corresponding val
		i++
		if i == len(keysAndValues) {
			log.Panicf("expected an associated value after key %v, got nothing", key)
		}
		// e.Logs = append(e.Logs, keyVal{
		// 	Key: key,
		// 	Val: keysAndValues[i],
		// })
		e.Logs[key] = keysAndValues[i]
		e.Ordering[key] = e.counter
		e.counter++
	}
	return e
}

func (e *LogEvent) SetEnded(end int) Event {
	e.Ended = end
	return e
}

func (e *LogEvent) LogSource() Source { return e.Typ }
func (e *LogEvent) StartFrame() int   { return e.F }
func (e *LogEvent) Src() int          { return e.SrcChar }

type Ctrl struct {
	// keep it in an array so we can keep track order it occured
	// events []*Event
	events map[int]*LogEvent
	count  int
	f      *int
}

func New(f *int, size int) Logger {
	ctrl := &Ctrl{
		events: make(map[int]*LogEvent),
		f:      f,
	}
	return ctrl
}

func (c *Ctrl) Dump() ([]byte, error) {
	r := make(EventArr, 0, c.count)
	for i := 0; i < c.count; i++ {
		v, ok := c.events[i]
		if ok {
			r = append(r, v)
		}
	}
	return easyjson.Marshal(r)
}

func (c *Ctrl) NewEventBuildMsg(typ Source, srcChar int, msg ...string) Event {
	if len(msg) == 0 {
		panic("no msg provided")
	}
	var sb strings.Builder
	for _, v := range msg {
		sb.WriteString(v)
	}
	return c.NewEvent(sb.String(), typ, srcChar)
}

func (c *Ctrl) NewEvent(msg string, typ Source, srcChar int) Event {
	e := &LogEvent{
		Msg:      msg,
		F:        *c.f,
		Ended:    *c.f,
		Typ:      typ,
		SrcChar:  srcChar,
		Logs:     make(map[string]interface{}), //+5 from default just in case we need to add in more keys
		Ordering: make(map[string]int),
	}
	// c.events = append(c.events, e)
	c.events[c.count] = e
	c.count++
	return e
}
