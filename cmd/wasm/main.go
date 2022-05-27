package main

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"syscall/js"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/pkg/simulation"
)

var (
	sha1ver   string // sha1 revision used to build the program
	buildTime string // when the executable was built
)

func main() {
	//GOOS=js GOARCH=wasm go build -o main.wasm
	done := make(chan struct{}, 0)

	global := js.Global()

	setConfigFunc := js.FuncOf(setConfig)
	defer setConfigFunc.Release()

	checkConfigFunc := js.FuncOf(checkConfig)
	defer checkConfigFunc.Release()

	runSimFunc := js.FuncOf(run)
	defer runSimFunc.Release()

	debugFunc := js.FuncOf(debug)
	defer debugFunc.Release()

	collectFunc := js.FuncOf(collect)
	defer collectFunc.Release()

	versionFunc := js.FuncOf(version)
	defer versionFunc.Release()

	global.Set("sim", runSimFunc)
	global.Set("setcfg", setConfigFunc)
	global.Set("checkcfg", checkConfigFunc)
	global.Set("debug", debugFunc)
	global.Set("collect", collectFunc)
	global.Set("version", versionFunc)

	<-done
}

var cfg core.SimulationConfig
var cfgStr string

func setConfig(this js.Value, args []js.Value) interface{} {
	in := args[0].String()
	//parse this
	parser := parse.New("single", in)
	var err error
	cfg, err = parser.Parse()
	if err != nil {
		return marshalErr(err)
	}
	cfgStr = in

	data, err := json.Marshal(cfg)
	if err != nil {
		return marshalErr(err)
	}
	return string(data)
}

func checkConfig(this js.Value, args []js.Value) interface{} {
	in := args[0].String()
	//parse this
	parser := parse.New("single", in)
	var err error
	cfg, err = parser.Parse()
	if err != nil {
		return marshalErr(err)
	}

	for i, v := range cfg.Characters.Profile {
		log.Printf("%v: %v\n", i, v.Base.Key.String())
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return marshalErr(err)
	}
	return string(data)
}

func version(this js.Value, args []js.Value) interface{} {
	return fmt.Sprintf(`{"hash":"%v","date":"%v"}`, sha1ver, buildTime)
}

//run simulation once
func run(this js.Value, args []js.Value) interface{} {
	//seed this with now
	t := cfg.Clone()
	c, err := simulation.NewCore(cryptoRandSeed(), false, t.Settings)
	if err != nil {
		return marshalErr(err)
	}
	s, err := simulation.New(t, c)
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

//debug generates the debug log (does not track dps value)
func debug(this js.Value, args []js.Value) interface{} {
	c, err := simulation.NewCore(cryptoRandSeed(), true, cfg.Settings)
	if err != nil {
		return marshalErr(err)
	}
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

func collect(this js.Value, args []js.Value) interface{} {
	var in []simulation.Result
	s := args[0].String()
	err := json.Unmarshal([]byte(s), &in)
	if err != nil {
		log.Println(err)
		return marshalErr(err)
	}

	chars := make([]string, len(cfg.Characters.Profile))
	for i, v := range cfg.Characters.Profile {
		chars[i] = v.Base.Key.String()
	}

	r := result.CollectResult(
		in,
		cfg.DamageMode,
		chars,
		true,
		false,
	)

	r.Iterations = cfg.Settings.Iterations
	r.ActiveChar = cfg.Characters.Initial.String()
	if cfg.DamageMode {
		r.Duration.Mean = float64(cfg.Settings.Duration)
		r.Duration.Min = float64(cfg.Settings.Duration)
		r.Duration.Max = float64(cfg.Settings.Duration)
	}

	r.NumTargets = len(cfg.Targets)
	r.CharDetails = in[0].CharDetails
	for i := range r.CharDetails {
		r.CharDetails[i].Stats = cfg.Characters.Profile[i].Stats
	}
	r.TargetDetails = cfg.Targets
	r.Text = r.PrettyPrint()
	r.Config = cfgStr

	out, err := json.Marshal(r)
	if err != nil {
		return marshalErr(err)
	}

	return string(out)
}

//results aggregates all the results

func cryptoRandSeed() int64 {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
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
