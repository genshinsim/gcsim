package weapon

import (
	"fmt"
	"os"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestParseCharConfigs(t *testing.T) {
	// write 2 config yaml to file, read it back
	dir := t.TempDir()
	r := []Config{
		{
			PackageName: "a",
		},
		{
			PackageName: "b",
		},
	}
	for _, v := range r {
		err := os.Mkdir(dir+"/"+v.PackageName, 0o755)
		if err != nil {
			t.Fatal(err)
		}
		data, err := yaml.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(fmt.Sprintf("%v/%v/config.yml", dir, v.PackageName), data, 0o644)
		if err != nil {
			t.Fatal(err)
		}
	}

	cfgs, err := ParseWeaponConfig(dir)
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
