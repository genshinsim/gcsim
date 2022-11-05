package main

type WeaponConfigs []struct {
	// WeaponType string `json:"weaponType"`
	// RankLevel  int    `json:"rankLevel"`
	// WeaponBaseExp int    `json:"weaponBaseExp"`
	// SkillAffix []int `json:"skillAffix"`
	// AwakenTexture              string        `json:"awakenTexture"`
	// AwakenLightMapTexture      string        `json:"awakenLightMapTexture"`
	// AwakenIcon                 string        `json:"awakenIcon"`
	// StoryID                    int           `json:"storyId"`
	// AwakenCosts                []interface{} `json:"awakenCosts"`
	// GachaCardNameHashSuffix    int64         `json:"gachaCardNameHashSuffix"`
	// DestroyRule                string        `json:"destroyRule"`
	// DestroyReturnMaterial      []int         `json:"destroyReturnMaterial"`
	// DestroyReturnMaterialCount []int         `json:"destroyReturnMaterialCount"`
	// DescTextMapHash            int64         `json:"descTextMapHash"`
	// Icon     string `json:"icon"`
	// ItemType string `json:"itemType"`
	// Weight                     int           `json:"weight"`
	// Rank int `json:"rank"`
	// GadgetID                   int           `json:"gadgetId"`
	ID         int `json:"id"`
	WeaponProp []struct {
		PropType  string  `json:"propType,omitempty"`
		InitValue float64 `json:"initValue,omitempty"`
		Type      string  `json:"type"`
	} `json:"weaponProp"`
	WeaponPromoteID int   `json:"weaponPromoteId"`
	NameTextMapHash int64 `json:"nameTextMapHash"`
}

func getWeapon() WeaponConfigs {
	// weaponDataJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/WeaponExcelConfigData.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	weaponDataJson := getJsonFromFile("../ExcelBinOutput/WeaponExcelConfigData.json")
	// fmt.Println(weaponDataJson)
	return weaponDataJson
}

// func fetchJsonFromUrl(path string) (string, error) {

// 	resp, err := http.Get(path)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != 200 {
// 		return "", fmt.Errorf("%v: %v", resp.Status, path)
// 	}

// 	out, err := io.ReadAll(resp.Body)
// 	return string(out), err
// }
