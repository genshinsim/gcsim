//bin provides the methods required to run simulations; the cmd line tools should be a wrapper
//around this
package simulator

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"log"
	"net/url"
	"os"

	"go.uber.org/zap"
)

//GenerateDebugLogWithSeed will run one simulation with debug enabled using the given seed and output
//the debug log. Used for generating debug for min/max runs
func GenerateDebugLogWithSeed(cfg string, seed int64) (string, error) {
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

}

//GenerateDebugLog will run one simulation with debug enabled using a random seed
func GenerateDebugLog(cfg string) (string, error) {
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		log.Panic("cannot seed math/rand package with cryptographically secure random number generator")
	}
	seed := int64(binary.LittleEndian.Uint64(b[:]))
	return GenerateDebugLogWithSeed(cfg, seed)
}

//RunOnce provide convenience wrapper around Run
func RunOnce() {

}

//Run will run the simulation given number of times
func Run(opts SimulatorOptions) {

}
