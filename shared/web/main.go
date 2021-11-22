package main

import (
	"encoding/json"
	"net/url"
	"syscall/js"

	"github.com/genshinsim/gcsim"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
	"go.uber.org/zap"
)

func main() {
	//GOOS=js GOARCH=wasm go build -o ../../../gsimweb/public/sim.wasm
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

type CustomSink struct {
	WriteTo func(msg string)
}

func (*CustomSink) Sync() error  { return nil }
func (*CustomSink) Close() error { return nil }

func (c *CustomSink) Write(p []byte) (int, error) {
	c.WriteTo(string(p))
	return len(p), nil
}

var sink CustomSink

func init() {
	zap.RegisterSink("gsim", func(url *url.URL) (zap.Sink, error) {
		return &sink, nil
	})
}

func runSim(this js.Value, args []js.Value) interface{} {
	in := args[0].String()
	callback := args[1]
	update := args[2]
	debug := args[3]

	var runCfg runConfig
	err := json.Unmarshal([]byte(in), &runCfg)
	if err != nil {
		callback.Invoke(err.Error(), nil)
		return js.Undefined()
	}

	// runCfg.Options.Iteration = 100
	// runCfg.Options.Debug = false
	opt := runCfg.Options

	var data []gcsim.Stats

	parser := parse.New("single", runCfg.Config)
	cfg, _, err := parser.Parse()
	if err != nil {
		callback.Invoke(err.Error(), nil)
		return js.Undefined()
	}

	charCount := len(cfg.Characters.Profile)

	if charCount > 4 {
		callback.Invoke("cannot have more than 4 characters in a team", nil)
		return js.Undefined()
	}

	chars := make([]string, len(cfg.Characters.Profile))
	for i, v := range cfg.Characters.Profile {
		chars[i] = v.Base.Name
	}

	count := opt.Iteration
	if count == 0 {
		count = 1000
	}

	if opt.Debug {
		count--
	}

	t := opt.Debug
	opt.Debug = false

	for i := 0; i < count; i++ {
		s, err := gcsim.NewSim(cfg, opt)
		if err != nil {
			callback.Invoke(err.Error(), nil)
			return js.Undefined()
		}
		v, err := s.Run()
		if err != nil {
			callback.Invoke(err.Error(), nil)
			return js.Undefined()
		}
		data = append(data, v)

		update.Invoke(i + 1)
	}

	var out string
	if t {

		opt.Debug = true

		// sink := CustomSink{}
		sink.WriteTo = func(msg string) {
			debug.Invoke(msg)
		}

		zap.RegisterSink("gsim", func(url *url.URL) (zap.Sink, error) {
			return &sink, nil
		})

		opt.DebugPaths = []string{"gsim://"}

		s, err := gcsim.NewSim(cfg, opt)
		if err != nil {
			callback.Invoke(err.Error(), nil)
			return js.Undefined()
		}
		v, err := s.Run()
		if err != nil {
			callback.Invoke(err.Error(), nil)
			return js.Undefined()
		}

		// log.Println(v)
		data = append(data, v)
	}

	result := gcsim.CollectResult(data, cfg.DamageMode, chars, opt.LogDetails)
	result.Iterations = opt.Iteration
	if !cfg.DamageMode {
		result.Duration.Mean = float64(opt.Duration)
		result.Duration.Min = float64(opt.Duration)
		result.Duration.Max = float64(opt.Duration)
	}
	result.Text = result.PrettyPrint()
	result.Debug = out

	resultj, _ := json.Marshal(result)
	callback.Invoke(nil, string(resultj))

	return js.Undefined()
}
