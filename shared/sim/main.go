package main

import "C"

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/genshinsim/gcsim/pkg/parse"
)

func main() {}

//export Run
func Run(config string) *C.char {

	parser := parse.New("single", config)
	_, opts, err := parser.Parse()
	if err != nil {
		return C.CString(errToString("error parsing config"))
	}

	opts.LogDetails = true

	var data gsim.Result

	if opts.Debug {
		//make a log file, write to it, read it back and store as string
		now := time.Now().Unix()

		// r, w, err := os.Pipe()
		// if err != nil {
		// 	return C.CString(errToString(err.Error()))
		// }
		// defer w.Close()
		// defer r.Close()

		// outC := make(chan string)
		// // copy the output in a separate goroutine so printing can't block indefinitely
		// go func() {
		// 	var buf bytes.Buffer
		// 	io.Copy(&buf, r)
		// 	outC <- buf.String()
		// }()
		// zap.RegisterSink("gsim", func(url *url.URL) (zap.Sink, error) {
		// 	return w, nil
		// })
		file := "./" + strconv.FormatInt(now, 10) + ".log"
		opts.DebugPaths = []string{file}

		data, err = gsim.Run(string(config), opts)
		if err != nil {
			return C.CString(errToString(err.Error()))
		}

		l, err := ioutil.ReadFile(file)
		if err != nil {
			return C.CString(errToString(err.Error()))
		}

		data.Debug = string(l)

		//remove the file
		os.Remove(file)
	} else {
		data, err = gsim.Run(string(config), opts)
		if err != nil {
			return C.CString(errToString(err.Error()))
		}
	}

	data.Text = data.PrettyPrint()

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
