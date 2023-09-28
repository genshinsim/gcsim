package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"
	"text/template"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	names := readNameMap()
	// fmt.Println(d)
	writeTmpl(tmplKeys, "./keys.txt", names)
	writeTmpl(tmplShortcuts, "./shortcuts.txt", names)
}

func writeTmpl(tmplStr, outFile string, d map[string]string) {
	t, err := template.New("out").Parse(tmplStr)
	if err != nil {
		log.Panic(err)
	}
	os.Remove(outFile)
	of, err := os.Create(outFile)
	if err != nil {
		log.Panic(err)
	}
	err = t.Execute(of, d)
	if err != nil {
		log.Panic(err)
	}
}

type namemap struct {
	Names map[string]string `json:"names"`
}

var re = regexp.MustCompile(`(?i)[^0-9a-z]`)

func readNameMap() map[string]string {
	f, err := os.ReadFile("./artifacts.json")
	if err != nil {
		log.Panic(err)
	}
	var m namemap
	err = json.Unmarshal(f, &m)
	if err != nil {
		log.Panic(err)
	}

	for k, v := range m.Names {
		// strip out any none word characters
		v = strings.ReplaceAll(v, "'", "")
		m.Names[k] = re.ReplaceAllString(cases.Title(language.Und, cases.NoLower).String(v), "")
	}
	return m.Names
}

var tmplKeys = `package keys

import (
	"encoding/json"
	"errors"
	"strings"
)

type Set int

func (s *Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(weaponNames[*s])
}

func (c *Set) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	s = strings.ToLower(s)
	for i, v := range setNames {
		if v == s {
			*c = Set(i)
			return nil
		}
	}
	return errors.New("unrecognized set key")
}

func (c Set) String() string {
	return setNames[c]
}


var setNames = []string{
	"",
	{{- range $key, $value := . }}
	"{{$key}}",
	{{- end }}
}

const (
	NoSet Set = iota
	{{- range $key, $value := . }}
	{{$value}}
	{{- end }}
)
`

var tmplShortcuts = `package shortcut

import "github.com/genshinsim/gcsim/pkg/core/keys"

var SetNameToKey = map[string]keys.Set{
	{{- range $key, $value := . }}
	"{{$key}}": keys.{{$value}},
	{{- end }}
}
`
