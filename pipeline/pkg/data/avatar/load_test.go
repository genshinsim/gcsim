package avatar

import (
	"testing"
)

func TestLoadAvatarExcel(t *testing.T) {
	const src = "../../../data/ExcelBinOutput/" + AvatarExcelConfigData

	res, err := loadAvatarExcel(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}

	if _, ok := res[10000007]; !ok {
		t.Errorf("could not find id for lumine")
	}

	if _, ok := res[10000005]; !ok {
		t.Errorf("could not find id for aether")
	}
}

func TestLoadAvatarSkillDepot(t *testing.T) {
	const src = "../../../data/ExcelBinOutput/" + AvatarSkillDepotExcelConfigData

	res, err := loadAvatarSkillDepot(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}
}

func TestLoadAvatarSkillExcel(t *testing.T) {
	const src = "../../../data/ExcelBinOutput/" + AvatarSkillExcelConfigData

	res, err := loadAvatarSkillExcel(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}

}

func TestLoadFetterInfo(t *testing.T) {
	const src = "../../../data/ExcelBinOutput/" + FetterInfoExcelConfigData

	res, err := loadAvatarFetterInfo(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}

}

func TestLoadPromotData(t *testing.T) {
	const src = "../../../data/ExcelBinOutput/" + AvatarPromoteExcelConfigData

	res, err := loadAvatarPromoteData(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}

}
