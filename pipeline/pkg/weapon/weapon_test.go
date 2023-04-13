package weapon

import "testing"

func TestGenerateUIJSON(t *testing.T) {
	g, err := NewGenerator(GeneratorConfig{
		Root:   "../../../internal/weapons",
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
