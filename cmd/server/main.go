package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/genshinsim/gcsim/pkg/servermode"
)

var (
	shareKey string
)

type opts struct {
	host     string
	port     string
	shareKey string
	timeout  int
}

func main() {
	if shareKey == "" {
		shareKey = os.Getenv("GCSIM_SHARE_KEY")
	}

	var opt opts
	flag.StringVar(&opt.host, "host", "localhost", "host to listen to (default: localhost)")
	flag.StringVar(&opt.port, "port", "54321", "port to listen on (default: 54321)")
	flag.StringVar(&opt.shareKey, "sharekey", "", "share key to use (default: build flag OR GCSIM_SHARE_KEY env variable if not available)")
	flag.IntVar(&opt.timeout, "timeout", 5*60, "how long to run each sim for in seconds before timing out (default: 300s)")
	flag.Parse()

	if opt.shareKey != "" {
		shareKey = opt.shareKey
	}

	server, err := servermode.New(
		servermode.WithDefaults(),
		servermode.WithShareKey(shareKey),
		servermode.WithTimeout(time.Duration(opt.timeout)*time.Second),
	)

	if err != nil {
		panic(err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", opt.host, opt.port), server.Router))
}
