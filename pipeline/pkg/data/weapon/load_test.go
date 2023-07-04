package weapon

import (
	"testing"
)

func TestLoadWeaponExcel(t *testing.T) {
	const src = "../../../data/ExcelBinOutput/" + WeaponExcelConfigData

	res, err := loadWeaponExcel(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}

	if _, ok := res[14405]; !ok {
		t.Errorf("could not find id for Solar Pearl")
	}

	if _, ok := res[12502]; !ok {
		t.Errorf("could not find id for Wolf's Gravestone")
	}
}

func TestLoadPromotData(t *testing.T) {
	const src = "../../../data/ExcelBinOutput/" + WeaponPromoteExcelConfigData

	res, err := loadWeaponPromoteData(src)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Error("res length cannot be 0")
	}

}
