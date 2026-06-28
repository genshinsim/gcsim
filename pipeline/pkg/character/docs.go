package character

import (
	"bytes"
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

	b := bytes.NewBuffer(nil)
	e := json.NewEncoder(b)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	if err := e.Encode(data); err != nil {
		return err
	}
	return os.WriteFile(path, b.Bytes(), 0o644)
}
