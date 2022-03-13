package eventlog

import (
	"log"
	"strings"

	"github.com/genshinsim/gcsim/pkg/coretype"
	easyjson "github.com/mailru/easyjson"
)

//Debugw
//Warnw

type keyVal struct {
	Key string      `json:"key,nocopy"`
	Val interface{} `json:"val"`
}

//easyjson:json
type Event struct {
	Typ      coretype.LogSource     `json:"event"`
	F        int                    `json:"frame"`
	Ended    int                    `json:"ended"`
	SrcChar  int                    `json:"char_index"`
	Msg      string                 `json:"msg,nocopy"`
	Logs     map[string]interface{} `json:"logs"`
	Ordering map[string]int         `json:"ordering"`
	counter  int
}

//easyjson:json
type EventArr []*Event

func (e *Event) Write(keysAndValues ...interface{}) {
	//should be even number
	var key string
	var ok bool
	for i := 0; i < len(keysAndValues); i++ {
		key, ok = keysAndValues[i].(string)
		if !ok {
			log.Panicf("invalid key %v, expected type to be string", keysAndValues[i].(string))
		}
		//make sure there's a corresponding val
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
}

func (e *Event) SetEnded(end int)              { e.Ended = end }
func (e *Event) LogSource() coretype.LogSource { return e.Typ }
func (e *Event) StartFrame() int               { return e.F }
func (e *Event) Src() int                      { return e.SrcChar }

type Ctrl struct {
	//keep it in an array so we can keep track order it occured
	// events []*Event
	events map[int]*Event
	count  int
	f      *int
}

func NewCtrl(f *int, size int) coretype.Logger {
	ctrl := &Ctrl{
		events: make(map[int]*Event),
		f:      f,
	}
	return ctrl
}

func (c *Ctrl) Dump() ([]byte, error) {
	var r EventArr = make(EventArr, 0, c.count)
	for i := 0; i < c.count; i++ {
		v, ok := c.events[i]
		if ok {
			r = append(r, v)
		}
	}
	return easyjson.Marshal(r)
}

func (c *Ctrl) NewEventBuildMsg(typ coretype.LogSource, srcChar int, msg ...string) coretype.LogEvent {
	if len(msg) == 0 {
		panic("no msg provided")
	}
	var sb strings.Builder
	for _, v := range msg {
		sb.WriteString(v)
	}
	return c.NewEvent(sb.String(), typ, srcChar)
}

func (c *Ctrl) NewEvent(msg string, typ coretype.LogSource, srcChar int, keysAndValues ...interface{}) coretype.LogEvent {
	e := &Event{
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
	if keysAndValues != nil {
		e.Write(keysAndValues...)
	}
	return e
}
