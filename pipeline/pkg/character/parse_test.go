package character

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestParseCharConfigs(t *testing.T) {
	err := os.RemoveAll("./test")
	if err != nil {
		t.Fatal(err)
	}
	//write 2 config yaml to file, read it back
	err = os.Mkdir("./test", 0755)
	if err != nil {
		t.Fatal(err)
	}
	r := []Config{
		{
			PackageName: "a",
		},
		{
			PackageName: "b",
		},
	}
	for _, v := range r {
		err := os.Mkdir("./test/"+v.PackageName, 0755)
		if err != nil {
			t.Fatal(err)
		}
		data, err := yaml.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(fmt.Sprintf("./test/%v/config.yml", v.PackageName), data, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	cfgs, err := ParseCharConfig("./test")
	if err != nil {
		t.Errorf("error encountered parsing config: %v", err)
		t.FailNow()
	}
	if len(cfgs) == 0 {
		t.Error("configs read should not be 0")
		t.FailNow()
	}
	for i, v := range cfgs {
		if v.PackageName != r[i].PackageName {
			t.Errorf("data not matching, expecting pkg name %v, got pkg name %v", r[i].PackageName, v.PackageName)
		}
	}

}
