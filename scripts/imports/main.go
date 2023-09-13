package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type data struct {
	Art  []string
	Char []string
	Weap []string
}

func main() {
	var d data

	d.Art = walk("../../internal/artifacts")
	d.Char = walk("../../internal/characters")
	d.Weap = walk("../../internal/weapons")

	t, err := template.New("out").Parse(tmpl)
	if err != nil {
		log.Panic(err)
	}
	os.Remove("./out.txt")
	of, err := os.Create("./out.txt")
	if err != nil {
		log.Panic(err)
	}
	err = t.Execute(of, d)
	if err != nil {
		log.Panic(err)
	}
}

func walk(path string) []string {
	var r []string
	err := filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// fmt.Println(path, info.Size())
			// skip if any of the following is true
			switch {
			case !info.IsDir():
				return nil
			case strings.Contains(path, "/common"):
				return nil
			case strings.HasSuffix(path, "/artifacts"):
				return nil
			case strings.HasSuffix(path, "/characters"):
				return nil
			case strings.HasSuffix(path, "/weapons"):
				return nil
			case strings.HasSuffix(path, "/weapons/bow"):
				return nil
			case strings.HasSuffix(path, "/weapons/catalyst"):
				return nil
			case strings.HasSuffix(path, "/weapons/claymore"):
				return nil
			case strings.HasSuffix(path, "/weapons/spear"):
				return nil
			case strings.HasSuffix(path, "/weapons/sword"):
				return nil
			}

			r = append(r, strings.TrimPrefix(path, "../../"))

			return nil
		})
	if err != nil {
		panic(err)
	}
	return r
}

var tmpl = `package simulation

import (
	//artifacts
	{{- range $e := .Art}}
	_ "github.com/genshinsim/gcsim/{{$e}}"
	{{- end}}

	//char
	{{- range $e := .Char}}
	_ "github.com/genshinsim/gcsim/{{$e}}"
	{{- end}}

	//weapons
	{{- range $e := .Weap}}
	_ "github.com/genshinsim/gcsim/{{$e}}"
	{{- end}}
)

`
