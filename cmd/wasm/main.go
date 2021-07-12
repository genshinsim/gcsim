package main

import (
	"encoding/json"
	"syscall/js"
)

func main() {
	c := make(chan struct{}, 0)

	println("Go WebAssembly Initialized")

	js.Global().Set("runSim", js.FuncOf(runSim))

	<-c
}

type runConfig struct {
	LogLvl    string   `json:"log"`
	Seconds   int      `json:"seconds"`
	Config    string   `json:"config"`
	HP        float64  `json:"hp"`
	AvgMode   bool     `json:"avg_mode"`
	Iter      int      `json:"iter"`
	NoSeed    bool     `json:"noseed"`
	LogEvents []string `json:"logs"`
}

func runSim(this js.Value, args []js.Value) interface{} {

	var cfg runConfig
	err := json.Unmarshal([]byte(r.Payload), &cfg)

	return js.ValueOf(colors)
}
