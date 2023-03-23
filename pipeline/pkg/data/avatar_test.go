package data

import "testing"

func TestLoadAvatarExcel(t *testing.T) {
	const src = "../../data/ExcelBinOutput/AvatarExcelConfigData.json"

	res, err := loadAvatarExcel(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}

	//map the ids, there are some ids we know for sure should exist
	ids := make(map[int]AvatarExcel)
	for _, v := range res {
		//mhy can break this...
		if _, ok := ids[v.ID]; ok {
			t.Errorf("duplicate found for id %v", v.ID)
		}
		ids[v.ID] = v
	}

	if _, ok := ids[10000007]; !ok {
		t.Errorf("could not find id for lumine")
	}

	if _, ok := ids[10000005]; !ok {
		t.Errorf("could not find id for aether")
	}
}
