package simulator

import (
	"bytes"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/simulation"
	"go.uber.org/zap"
)

//GenerateDebugLogWithSeed will run one simulation with debug enabled using the given seed and output
//the debug log. Used for generating debug for min/max runs
func GenerateDebugLogWithSeed(cfg core.SimulationConfig, seed int64) (string, error) {
	//parse the config
	r, w, err := os.Pipe()
	if err != nil {
		log.Println(err)
		return "", err
	}
	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	zap.RegisterSink("gsim", func(url *url.URL) (zap.Sink, error) {
		return w, nil
	})
	//set up core to use logger with custom path in order to capture debug
	logger, err := core.NewDefaultLogger(true, true, []string{"gsim://"})
	if err != nil {
		return "", err
	}
	c := simulation.NewDefaultCoreWithCustomLogger(seed, logger)
	c.Flags.LogDebug = true
	//create a new simulation and run
	s, err := simulation.New(cfg, c)
	if err != nil {
		return "", err
	}
	_, err = s.Run()
	if err != nil {
		return "", err
	}
	//capture the log
	out := <-outC
	return out, nil
}

//GenerateDebugLog will run one simulation with debug enabled using a random seed
func GenerateDebugLog(cfg core.SimulationConfig) (string, error) {
	return GenerateDebugLogWithSeed(cfg, cryptoRandSeed())
}
