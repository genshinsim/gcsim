package main

import (
	"encoding/json"
	"errors"
	"strconv"
	"syscall/js"
	"time"

	"github.com/genshinsim/gcsim/pkg/agg"
	"github.com/genshinsim/gcsim/pkg/gcs/ast"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/simulator"
	"github.com/genshinsim/gcsim/pkg/stats"
)

const DefaultBufferLength = 1024 * 10

// assigned by compiler
var shareKey string

// shared variables
var cfg string
var simcfg *ast.ActionList
var buffer []byte

// Aggregator variables
var aggregators []agg.Aggregator
var cachedResult *model.SimulationResult

func main() {
	//GOOS=js GOARCH=wasm go build -o main.wasm
	ch := make(chan struct{}, 0)

	// Helper Functions (stateless, no init call needed)
	js.Global().Set("sample", js.FuncOf(doSample))
	js.Global().Set("validateConfig", js.FuncOf(validateConfig))

	// Worker Functions
	js.Global().Set("initializeWorker", js.FuncOf(initializeWorker))
	js.Global().Set("simulate", js.FuncOf(simulate))

	// Aggregator Functions
	js.Global().Set("initializeAggregator", js.FuncOf(initializeAggregator))
	js.Global().Set("aggregate", js.FuncOf(aggregate))
	js.Global().Set("flush", js.FuncOf(flush))

	<-ch
}

// static helper functions (stateless)

// sample(cfg: string, seed: string) -> string
func doSample(this js.Value, args []js.Value) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = errorRecover(r)
		}
	}()

	opts := simulator.Options{
		GZIPResult:       false,
		ResultSaveToPath: "",
		ConfigPath:       "",
	}

	cfg := args[0].String()
	seed, _ := strconv.ParseUint(args[1].String(), 10, 64)

	data, err := simulator.GenerateSampleWithSeed(cfg, seed, opts)
	if err != nil {
		return marshal(err)
	}

	marshalled, err := data.MarshalJson()
	if err != nil {
		return marshal(err)
	}

	return string(marshalled)
}

// validateConfig(cfg: string) -> string
func validateConfig(this js.Value, args []js.Value) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = errorRecover(r)
		}
	}()

	in := args[0].String()

	cfg, err := simulator.Parse(in)
	if err != nil {
		return marshal(err)
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return marshal(err)
	}
	return string(data)
}

// worker functions

// initializeWorker(cfg: string)
func initializeWorker(this js.Value, args []js.Value) interface{} {
	in := args[0].String()
	if err := initialize(in); err != nil {
		return marshal(err)
	}
	return nil
}

// simulate() -> js Uint8Array
func simulate(this js.Value, args []js.Value) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = errorRecover(r)
		}
	}()

	cpycfg := simcfg.Copy()
	core, err := simulation.NewCore(simulator.CryptoRandSeed(), false, cpycfg)
	if err != nil {
		return marshal(err)
	}

	sim, err := simulation.New(cpycfg, core)
	if err != nil {
		return marshal(err)
	}

	result, err := sim.Run()
	if err != nil {
		return marshal(err)
	}

	buffer, err = result.MarshalMsg(buffer[:0])
	if err != nil {
		return marshal(err)
	}

	dst := js.Global().Get("Uint8Array").New(len(buffer))
	copyLen := js.CopyBytesToJS(dst, buffer)
	if copyLen != len(buffer) {
		marshal(errors.New("BytesToJS: copied array was the incorrect size!"))
	}
	return dst
}

// aggregator functions

// initializeAggregator(cfg: string) -> string
func initializeAggregator(this js.Value, args []js.Value) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = errorRecover(r)
		}
	}()

	in := args[0].String()
	if err := initialize(in); err != nil {
		return marshal(err)
	}

	aggregators = aggregators[:0]
	for _, aggregator := range agg.Aggregators() {
		a, err := aggregator(simcfg)
		if err != nil {
			return marshal(err)
		}
		aggregators = append(aggregators, a)
	}

	opts := simulator.Options{
		GZIPResult:       false,
		ResultSaveToPath: "",
		ConfigPath:       "",
	}
	result, err := simulator.GenerateResult(cfg, simcfg, opts)
	if err != nil {
		return marshal(err)
	}

	// // store the result for reuse
	cachedResult = result

	marshalled, err := result.MarshalJson()
	if err != nil {
		return marshal(err)
	}
	return string(marshalled)
}

// aggregate(src: Uint8Array)
func aggregate(this js.Value, args []js.Value) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = errorRecover(r)
		}
	}()

	src := args[0]
	var err error

	// golang wasm copy requires src and destination length to have enough capacity to copy
	// should be enforced by capacity and not length but whatev....
	// have to waste cycles here to make sure buffer has the right length
	length := src.Get("length").Int()
	if length > cap(buffer) {
		buffer = make([]byte, length)
	} else {
		buffer = buffer[:length]
	}

	copyLen := js.CopyBytesToGo(buffer, src)
	if copyLen != len(buffer) {
		marshal(errors.New("BytesToGo: copied array was the incorrect size!"))
	}

	result := stats.Result{}
	buffer, err = result.UnmarshalMsg(buffer)
	if err != nil {
		return marshal(err)
	}

	for _, a := range aggregators {
		a.Add(result)
	}
	return nil
}

// flush(startTime: int) -> string
func flush(this js.Value, args []js.Value) (out interface{}) {
	defer func() {
		if r := recover(); r != nil {
			out = errorRecover(r)
		}
	}()

	startTime := args[0].Int()

	stats := &model.SimulationStatistics{}
	for _, a := range aggregators {
		a.Flush(stats)
	}
	stats.Runtime = float64(time.Now().Nanosecond() - startTime)

	// build full result from cache and sign
	cachedResult.Statistics = stats
	hash, _ := cachedResult.Sign(shareKey)

	signedResults := &model.SignedSimulationStatistics{
		Stats: stats,
		Hash:  hash,
	}

	marshalled, err := signedResults.MarshalJson()
	if err != nil {
		return marshal(err)
	}
	return string(marshalled)
}

// internal helper functions

func initialize(raw string) error {
	parser := ast.New(raw)
	out, err := parser.Parse()
	if err != nil {
		return err
	}

	if cap(buffer) < DefaultBufferLength {
		buffer = make([]byte, 0, DefaultBufferLength)
	}

	cfg = raw
	simcfg = out
	return nil
}

func marshal(err error) string {
	d := struct {
		Err string `json:"error"`
	}{
		Err: err.Error(),
	}
	b, _ := json.Marshal(d)
	return string(b)
}

func errorRecover(r interface{}) string {
	var err error
	switch x := r.(type) {
	case string:
		err = errors.New(x)
	case error:
		err = x
	default:
		err = errors.New("unknown error")
	}
	return marshal(err)
}
