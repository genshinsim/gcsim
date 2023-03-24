package character

import "testing"

func TestGenerateUIJSON(t *testing.T) {
	g, err := NewGenerator(GeneratorConfig{
		Root:   "../../../internal/characters",
		Excels: "../../data/ExcelBinOutput",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = g.DumpUIJSON(".")
	if err != nil {
		t.Fatal(err)
	}

}
