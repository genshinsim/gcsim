package main

import "C"

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/genshinsim/gsim/pkg/combat"
	"github.com/genshinsim/gsim/pkg/parse"
)

func main() {}

//export Run
func Run(config string) *C.char {

	parser := parse.New("single", config)
	_, opts, err := parser.Parse()
	if err != nil {
		return C.CString(errToString("error parsing config"))
	}

	var data combat.AverageStats

	if opts.Debug {
		old := os.Stdout
		r, w, err := os.Pipe()
		if err != nil {
			return C.CString(errToString(err.Error()))
		}
		os.Stdout = w
		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			io.Copy(&buf, r)
			outC <- buf.String()
		}()
		defer func() {
			w.Close()
			os.Stdout = old
		}()
		opts.DebugPaths = []string{"stdout"}
		data, err = combat.Run(string(config), opts)
		if err != nil {
			return C.CString(errToString(err.Error()))
		}
		out := <-outC
		data.Debug = out
	} else {
		data, err = combat.Run(string(config), opts)
		if err != nil {
			return C.CString(errToString(err.Error()))
		}
	}

	result, _ := json.Marshal(data)
	return C.CString(string(result))

}

func errToString(s string) string {
	var r struct {
		Err string `json:"err"`
	}
	r.Err = s

	b, _ := json.Marshal(r)

	return string(b)
}
