package character

import (
	"encoding/json"
	"os"
)

func (g *Generator) WriteFieldDocs(path string) error {
	data := make(map[string][]FieldDocData)

	for i := range g.chars {
		if len(g.chars[i].Documentation.FieldsData) == 0 {
			continue
		}
		data[g.chars[i].Key] = g.chars[i].Documentation.FieldsData
	}

	f, err := os.OpenFile(path, os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	e := json.NewEncoder(f)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	return e.Encode(data)
}
