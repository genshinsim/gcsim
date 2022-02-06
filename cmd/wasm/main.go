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
	global.Set("sim", runSimFunc)
	global.Set("setcfg", setConfigFunc)

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
	c := simulation.NewDefaultCoreWithDefaultLogger(int64(seed))
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

func marshalErr(err error) string {
	d := struct {
		Err string `json:"err"`
	}{
		Err: err.Error(),
	}
	b, _ := json.Marshal(d)
	return string(b)
}
