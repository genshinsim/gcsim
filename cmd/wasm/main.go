package main

import (
	"encoding/json"
	"log"
	"syscall/js"
	"time"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

func main() {
	//GOOS=js GOARCH=wasm go build -o main.wasm
	done := make(chan struct{}, 0)

	global := js.Global()

	setConfigFunc := js.FuncOf(setConfig)
	defer setConfigFunc.Release()

	runSimFunc := js.FuncOf(run)
	defer runSimFunc.Release()

	runSimCalcFunc := js.FuncOf(runCalcMode)
	defer runSimCalcFunc.Release()

	debugFunc := js.FuncOf(debug)
	defer debugFunc.Release()

	debugCalcFunc := js.FuncOf(debugCalcMode)
	defer debugFunc.Release()

	global.Set("sim", runSimFunc)
	global.Set("simcalc", runSimCalcFunc)
	global.Set("setcfg", setConfigFunc)
	global.Set("debug", debugFunc)
	global.Set("debugcalc", debugCalcFunc)

	<-done
}

var cfg core.SimulationConfig

func setConfig(this js.Value, args []js.Value) interface{} {
	in := args[0].String()
	//parse this
	parser := parse.New("single", in)
	var err error
	cfg, err = parser.Parse()
	if err != nil {
		return err.Error()
	}
	return "ok"
}

//run runs simulation once
func run(this js.Value, args []js.Value) interface{} {
	//seed this with now
	seed := time.Now().Nanosecond()
	c := simulation.NewDefaultCore(int64(seed))
	s, err := simulation.New(cfg, c)
	if err != nil {
		return marshalErr(err)
	}
	res, err := s.Run()
	if err != nil {
		return marshalErr(err)
	}
	// log.Println(res.DPS)
	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func runCalcMode(this js.Value, args []js.Value) interface{} {
	//seed this with now
	seed := time.Now().Nanosecond()
	c := simulation.NewDefaultCoreWithCalcQueue(int64(seed))
	s, err := simulation.New(cfg, c)
	if err != nil {
		return marshalErr(err)
	}
	res, err := s.Run()
	if err != nil {
		return marshalErr(err)
	}
	// log.Println(res.DPS)
	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
	}
	return string(b)
}

func debug(this js.Value, args []js.Value) interface{} {
	seed := int64(time.Now().Nanosecond())
	c := simulation.NewDefaultCoreWithDebug(seed)
	c.Flags.LogDebug = true
	//create a new simulation and run
	s, err := simulation.New(cfg, c)
	if err != nil {
		return marshalErr(err)
	}
	_, err = s.Run()
	if err != nil {
		return marshalErr(err)
	}
	//capture the log
	out, err := c.Log.Dump()
	if err != nil {
		return marshalErr(err)
	}
	return string(out)
}

func debugCalcMode(this js.Value, args []js.Value) interface{} {
	seed := int64(time.Now().Nanosecond())
	c := simulation.NewDefaultCoreWithDebugCalcQueue(seed)
	c.Flags.LogDebug = true
	//create a new simulation and run
	s, err := simulation.New(cfg, c)
	if err != nil {
		return marshalErr(err)
	}
	_, err = s.Run()
	if err != nil {
		return marshalErr(err)
	}
	//capture the log
	out, err := c.Log.Dump()
	if err != nil {
		return marshalErr(err)
	}
	return string(out)
}

func marshalErr(err error) string {
	d := struct {
		Err string `json:"err"`
	}{
		Err: err.Error(),
	}
	b, _ := json.Marshal(d)
	return string(b)
}
