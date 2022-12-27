module github.com/genshinsim/gcsim/backend

go 1.19

replace github.com/genshinsim/gcsim/pkg/model => ../
replace github.com/genshinsim/gcsim/pkg/simulator => ../

require (
	github.com/aidarkhanov/nanoid/v2 v2.0.5
	github.com/chromedp/chromedp v0.8.6
	github.com/dgraph-io/badger/v3 v3.2103.4
	github.com/dgraph-io/ristretto v0.1.1
	github.com/genshinsim/gcsim v1.3.6
	github.com/genshinsim/gcsim/pkg/model v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi v1.5.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/jaevor/go-nanoid v1.3.0
	github.com/mattn/go-sqlite3 v1.14.16
	go.mongodb.org/mongo-driver v1.11.0
	go.uber.org/zap v1.23.0
	golang.org/x/oauth2 v0.2.0
	google.golang.org/grpc v1.50.1
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chromedp/cdproto v0.0.0-20220924210414-0e3390be1777 // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.1.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/flatbuffers v22.10.26+incompatible // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.15.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/montanaflynn/stats v0.0.0-20171201202039-1bf9dbcd8cbe // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tinylib/msgp v1.1.6 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/text v0.4.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221114212237-e4508ebdbee1 // indirect
)
