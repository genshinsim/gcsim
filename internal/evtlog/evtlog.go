package evtlog

import (
	"log"

	"github.com/genshinsim/gcsim/pkg/core"
)

//Debugw
//Warnw

//easyjson:json
type keyVal struct {
	Key string      `json:"key"`
	Val interface{} `json:"val"`
}

//easyjson:json
type Event struct {
	Typ     core.LogSource `json:"type"`
	F       int            `json:"frame"`
	Ended   int            `json:"ended"`
	SrcChar int            `json:"char_index"`
	Msg     string         `json:"msg"`
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

func (e *Event) SetEnded(end int) {
	e.Ended = end
}

type Ctrl struct {
	//keep it in an array so we can keep track order it occured
	events []*Event
	core   *core.Core
}

func (c *Ctrl) NewEvent(msg string, typ core.LogSource, srcChar int) *Event {
	e := &Event{
		Msg:     msg,
		F:       c.core.F,
		Ended:   c.core.F,
		Typ:     typ,
		SrcChar: srcChar,
		Logs:    make([]keyVal, 0, 10),
	}
	c.events = append(c.events, e)
	return e
}

func (c *Ctrl) NewEventWithFields(msg string, typ core.LogSource, srcChar int, keysAndValues ...interface{}) *Event {
	e := c.NewEvent(msg, typ, srcChar)
	e.Write(keysAndValues...)
	return e
}

//temporary work around
func (c *Ctrl) Debugw(msg string, keysAndValues ...interface{}) {
	c.NewEventWithFields(msg, core.LogSimEvent, 0, keysAndValues...)
}
