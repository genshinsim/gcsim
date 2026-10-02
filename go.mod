module github.com/genshinsim/gcsim

go 1.27.0

ignore (
	./ui
	node_modules
)

tool (
	github.com/dmarkham/enumer
	github.com/tinylib/msgp
)

require (
	github.com/aclements/go-moremath v0.0.0-20210112150236-f10218a38794
	github.com/adrg/xdg v0.5.3
	github.com/caarlos0/env/v10 v10.0.0
	github.com/creativeprojects/go-selfupdate v1.1.3
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.15.0
	github.com/go-chi/chi v1.5.5
	github.com/go-chi/cors v1.2.1
	github.com/mailru/easyjson v0.7.7
	github.com/sanity-io/litter v1.5.9-0.20260504104730-2ddefc21bc33
	github.com/schollz/progressbar/v3 v3.18.0
	github.com/shizukayuki/excel-hk4e v0.0.0-20260717230206-c93b17a7e33b
	github.com/tinylib/msgp v1.1.9
	github.com/urfave/cli/v3 v3.10.1
	go.mongodb.org/mongo-driver v1.12.1
	go.uber.org/zap v1.26.0
	go.yaml.in/yaml/v4 v4.0.0-rc.6
	golang.org/x/tools v0.44.0
	google.golang.org/protobuf v1.36.12
	mvdan.cc/gofumpt v0.10.0
)

require (
	code.gitea.io/sdk/gitea v0.17.1 // indirect
	github.com/Masterminds/semver/v3 v3.2.1 // indirect
	github.com/davidmz/go-pageant v1.0.2 // indirect
	github.com/dmarkham/enumer v1.6.3 // indirect
	github.com/go-fed/httpsig v1.1.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-github/v30 v30.1.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.4 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pascaldekloe/name v1.0.0 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/xanzy/go-gitlab v0.95.2 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/mod v0.35.0 // indirect
	golang.org/x/oauth2 v0.22.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/term v0.42.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	golang.org/x/time v0.3.0 // indirect
)

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.16
