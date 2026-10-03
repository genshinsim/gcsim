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
	github.com/creativeprojects/go-selfupdate v1.6.0
	github.com/davecgh/go-spew v1.1.1
	github.com/fatih/color v1.16.0
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/cors v1.2.2
	github.com/mailru/easyjson v0.9.2
	github.com/sanity-io/litter v1.5.9-0.20260504104730-2ddefc21bc33
	github.com/schollz/progressbar/v3 v3.19.1
	github.com/shizukayuki/excel-hk4e v0.0.0-20260717230206-c93b17a7e33b
	github.com/tinylib/msgp v1.6.4
	github.com/urfave/cli/v3 v3.13.0
	go.mongodb.org/mongo-driver v1.17.10
	go.uber.org/zap v1.28.0
	go.yaml.in/yaml/v4 v4.0.0-rc.6
	golang.org/x/tools v0.49.0
	google.golang.org/protobuf v1.36.12
	mvdan.cc/gofumpt v0.12.0
)

require (
	code.gitea.io/sdk/gitea v0.23.2 // indirect
	github.com/42wim/httpsig v1.2.4 // indirect
	github.com/Masterminds/semver/v3 v3.5.0 // indirect
	github.com/davidmz/go-pageant v1.0.2 // indirect
	github.com/dmarkham/enumer v1.6.3 // indirect
	github.com/go-fed/httpsig v1.1.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-github/v86 v86.0.0 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.8 // indirect
	github.com/hashicorp/go-version v1.9.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pascaldekloe/name v1.0.0 // indirect
	github.com/philhofer/fwd v1.2.0 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/ulikunitz/xz v0.5.15 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	gitlab.com/gitlab-org/api/client-go v1.46.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.53.0 // indirect
	golang.org/x/mod v0.40.0 // indirect
	golang.org/x/oauth2 v0.36.0 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.44.0 // indirect
	golang.org/x/text v0.38.0 // indirect
	golang.org/x/time v0.15.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/imdario/mergo => github.com/imdario/mergo v0.3.16
