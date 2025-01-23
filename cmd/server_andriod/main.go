// build for andriod with
// fyne package -os android/arm64 -appID com.gcsim.server -icon ../../ui/packages/ui/src/Images/logo.png --release

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sort"
	"time"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/genshinsim/gcsim/pkg/servermode"
)

var (
	shareKey     string
	version      string
	simulatorURL = must(url.Parse("https://gcsim.app/simulator"))
	viewerURL    = must(url.Parse("https://gcsim.app/web"))
)

type opts struct {
	host        string
	port        string
	shareKey    string
	timeout     int
	update      bool
	showVersion bool
}

func must[T any](x T, _ error) T {
	return x
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
	flag.BoolVar(&opt.update, "update", false, "run autoupdater (default: false)")
	flag.BoolVar(&opt.showVersion, "version", false, "show currrent version")
	flag.Parse()

	if opt.showVersion {
		fmt.Println("Running version: ", version)
		return
	}

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

	a := app.New()
	w := a.NewWindow("gcsim server mode")

	w.SetContent(widget.NewButton("Head to https://gcsim.app/simulator to use", func() { a.OpenURL(simulatorURL) }))
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", opt.host, opt.port), server.Router))
	}()
	go func() {
		bars := make(map[string]*widget.ProgressBar)
		for {
			progress := server.Progress()
			cont := container.NewVBox()
			for id := range bars {
				if _, ok := progress[id]; !ok {
					cont.Add(widget.NewButton("Head to https://gcsim.app/web to view results", func() { a.OpenURL(viewerURL) }))
					delete(bars, id)
				}
			}
			ids := make([]string, 0, len(progress))
			for k := range progress {
				ids = append(ids, k)
			}
			sort.Strings(ids)
			for _, id := range ids {
				prog := progress[id]
				progBar, ok := bars[id]
				if !ok {
					progBar = widget.NewProgressBar()
					progBar.Max = float64(prog.Max)
					bars[id] = progBar
				}
				progBar.SetValue(float64(prog.Curr))
				cont.Add(progBar)
			}
			if len(cont.Objects) > 0 {
				w.SetContent(cont)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	w.ShowAndRun()
}
