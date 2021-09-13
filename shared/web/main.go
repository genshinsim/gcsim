package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"syscall/js"

	"github.com/genshinsim/gsim"
	"github.com/genshinsim/gsim/pkg/core"
)

func main() {
	//GOOS=js GOARCH=wasm go build -o ../../app/public/sim.wasm
	done := make(chan struct{}, 0)

	global := js.Global()

	runSimFunc := js.FuncOf(runSim)
	defer runSimFunc.Release()
	global.Set("sim", runSimFunc)

	<-done
}

type runConfig struct {
	Options core.RunOpt `json:"options"`
	Config  string      `json:"config"`
}

func runSim(this js.Value, args []js.Value) interface{} {
	in := args[0].String()
	callback := args[1]

	var cfg runConfig
	err := json.Unmarshal([]byte(in), &cfg)
	if err != nil {
		callback.Invoke(err.Error(), nil)
		return js.Undefined()
	}

	cfg.Options.Iteration = 100

	cfg.Options.Debug = false

	if cfg.Options.Debug {

		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			callback.Invoke(err.Error(), nil)
			return js.Undefined()
		}
		os.Stdout = w

		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()

		cfg.Options.DebugPaths = []string{"stdout"}

		result, err := gsim.Run(cfg.Config, cfg.Options)
		if err != nil {
			callback.Invoke(err.Error(), nil)
			return js.Undefined()
		}

		w.Close()
		os.Stdout = old
		out := <-outC

		result.Text = result.PrettyPrint()
		result.Debug = out

		data, _ := json.Marshal(result)

		callback.Invoke(nil, string(data))

	} else {
		result, err := gsim.Run(cfg.Config, cfg.Options)
		if err != nil {
			callback.Invoke(err.Error(), nil)
			return js.Undefined()
		}

		result.Text = result.PrettyPrint()

		data, _ := json.Marshal(result)
		callback.Invoke(nil, string(data))
	}

	return js.Undefined()
}
