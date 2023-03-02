package preview

import (
	context "context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"path"
	"regexp"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
	Files        embed.FS
	AssetsFolder string
}

type Store struct {
	Router *chi.Mux
	Log    *zap.SugaredLogger
	cfg    Config
	tmpl   *template.Template
	data   sync.Map
	UnimplementedEmbedServer
}

var re = regexp.MustCompile(`(?m)^\s*\<script\>[\s+\n\r]+var data = "\{[\S\s\n\r]*\}";[\s+\n\r]+\</script\>$`)

// serve 2 routes
// 1 is the index
// another is the data (JSON)
func New(cfg Config, cust ...func(*Store) error) (*Store, error) {
	s := &Store{
		cfg: cfg,
	}

	log.Println("setting up router....")

	s.Router = chi.NewRouter()
	for _, f := range cust {
		err := f(s)
		if err != nil {
			return nil, err
		}
	}

	if s.Log == nil {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, err := config.Build()
		if err != nil {
			return nil, err
		}
		sugar := logger.Sugar()
		sugar.Debugw("logger initiated")

		s.Log = sugar
	}

	s.routes()

	//build a template from index.html
	b, err := s.cfg.Files.ReadFile("dist/index.html")
	if err != nil {
		return s, fmt.Errorf("error reading index.html")
	}
	tmplStr := re.ReplaceAllString(string(b), "<script>var data = \"{{.Data}}\"</script>")
	tmpl, err := template.New("index").Parse(tmplStr)
	if err != nil {
		return s, fmt.Errorf("error compiling data template: %v", err)
	}

	s.tmpl = tmpl

	return s, nil

}

func (s *Store) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	data := req.GetData()
	key := req.GetId()
	s.Log.Infow("get request received on grpc", "key", key)
	if data == nil {
		return nil, status.Error(codes.InvalidArgument, "payload cannot be nil")
	}
	if key == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be blank")
	}
	s.data.Store(key, data)

	snap, err := s.generateSnapshot(key)
	if err != nil {
		return nil, status.Error(codes.Internal, "error getting snapshot: "+err.Error())
	}

	return &GetResponse{
		Data: snap,
	}, nil
}

func (s *Store) generateSnapshot(key string) ([]byte, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(540, 250),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(
		context.Background(),
		opts...,
	)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, s.fullScreenshot("http://localhost:3001/"+key, 100, &buf)); err != nil {
		return nil, err
	}

	return buf, nil
}

type myFS struct {
	content embed.FS
}

func (c myFS) Open(name string) (fs.File, error) {
	return c.content.Open(path.Join("dist", name))
}

func (s *Store) routes() {
	s.Log.Debugw("setting up server routes for preview generation server")
	s.Router.Use(middleware.Logger)

	fs := http.FileServer(http.FS(&myFS{content: s.cfg.Files}))

	//for requests to any assets
	s.Router.Handle("/assets/*", fs)

	imgFS := http.FileServer(http.Dir(s.cfg.AssetsFolder))

	//for images
	s.Router.Handle("/api/assets/*", http.StripPrefix("/api/assets/", imgFS))

	//root should serve index
	s.Router.Handle("/{key}", s.handleServeHTML())

}

func (s *Store) handleServeHTML() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.Log.Info("received request for embed html page")
		//pull data from result store, insert into template, and then server
		key := chi.URLParam(r, "key")
		var out struct {
			Data string
		}

		defer func() {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("content-type", "text/html")
			s.tmpl.Execute(w, out)
		}()

		d, ok := s.data.LoadAndDelete(key)
		if !ok {
			out.Data = `{"err":"unexpected error getting result; key not found"}`
			return
		}
		res, ok := d.(*model.SimulationResult)
		if !ok {
			out.Data = `{"err":"unexpected error getting result; bad data"}`
			return
		}
		data, err := res.MarshalJson()
		if err != nil {
			out.Data = `{"err":"unexpected error getting result: cannot convert to json"}`
			return
		}
		out.Data = string(data)
	}
}

func (s *Store) fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.ActionFunc(func(context.Context) error {
			s.Log.Info("chromedp task start")
			return nil
		}),
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(context.Context) error {
			s.Log.Info("waiting for card to be visible")
			return nil
		}),
		chromedp.WaitEnabled("#card"),
		chromedp.ActionFunc(func(context.Context) error {
			s.Log.Info("card ready")
			return nil
		}),
		chromedp.FullScreenshot(res, quality),
	}
}
