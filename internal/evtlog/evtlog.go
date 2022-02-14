package evtlog

import (
	"log"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core"
	easyjson "github.com/mailru/easyjson"
)

//Debugw
//Warnw

//easyjson:json
type keyVal struct {
	Key string      `json:"key,nocopy"`
	Val interface{} `json:"val"`
}

//easyjson:json
type Event struct {
	Typ     core.LogSource `json:"event"`
	F       int            `json:"frame"`
	Ended   int            `json:"ended"`
	SrcChar int            `json:"char_index"`
	Msg     string         `json:"msg,nocopy"`
	Logs    []keyVal       `json:"logs"`
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
		e.Logs = append(e.Logs, keyVal{
			Key: key,
			Val: keysAndValues[i],
		})
	}
}

func (e *Event) SetEnded(end int)          { e.Ended = end }
func (e *Event) LogSource() core.LogSource { return e.Typ }
func (e *Event) StartFrame() int           { return e.F }
func (e *Event) Src() int                  { return e.SrcChar }

type Ctrl struct {
	//keep it in an array so we can keep track order it occured
	// events []*Event
	events map[int]*Event
	count  int
	core   *core.Core
}

func NewCtrl(c *core.Core, size int) core.LogCtrl {
	ctrl := &Ctrl{
		events: make(map[int]*Event),
		core:   c,
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

func (c *Ctrl) NewEventBuildMsg(typ core.LogSource, srcChar int, msg ...string) core.LogEvent {
	if len(msg) == 0 {
		panic("no msg provided")
	}
	var sb strings.Builder
	for _, v := range msg {
		sb.WriteString(v)
	}
	return c.NewEvent(sb.String(), typ, srcChar)
}

func (c *Ctrl) NewEvent(msg string, typ core.LogSource, srcChar int, keysAndValues ...interface{}) core.LogEvent {
	e := &Event{
		Msg:     msg,
		F:       c.core.F,
		Ended:   c.core.F,
		Typ:     typ,
		SrcChar: srcChar,
		Logs:    make([]keyVal, 0, len(keysAndValues)+5), //+5 from default just in case we need to add in more keys
	}
	// c.events = append(c.events, e)
	c.events[c.count] = e
	c.count++
	if keysAndValues != nil {
		e.Write(keysAndValues...)
	}
	return e
}
