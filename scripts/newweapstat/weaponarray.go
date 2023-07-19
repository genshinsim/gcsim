package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type WeaponConfig []struct {
	WeaponType string `json:"weaponType"`
	RankLevel  int    `json:"rankLevel"`
	// WeaponBaseExp int    `json:"weaponBaseExp"`
	SkillAffix []int `json:"skillAffix"`
	WeaponProp []struct {
		PropType  string  `json:"propType,omitempty"`
		InitValue float64 `json:"initValue,omitempty"`
		Type      string  `json:"type"`
	} `json:"weaponProp"`
	// AwakenTexture              string        `json:"awakenTexture"`
	// AwakenLightMapTexture      string        `json:"awakenLightMapTexture"`
	// AwakenIcon                 string        `json:"awakenIcon"`
	WeaponPromoteID int `json:"weaponPromoteId"`
	// StoryID                    int           `json:"storyId"`
	// AwakenCosts                []interface{} `json:"awakenCosts"`
	// GachaCardNameHashSuffix    int64         `json:"gachaCardNameHashSuffix"`
	// DestroyRule                string        `json:"destroyRule"`
	// DestroyReturnMaterial      []int         `json:"destroyReturnMaterial"`
	// DestroyReturnMaterialCount []int         `json:"destroyReturnMaterialCount"`
	ID              int   `json:"id"`
	NameTextMapHash int64 `json:"nameTextMapHash"`
	// DescTextMapHash            int64         `json:"descTextMapHash"`
	Icon     string `json:"icon"`
	ItemType string `json:"itemType"`
	// Weight                     int           `json:"weight"`
	Rank int `json:"rank"`
	// GadgetID                   int           `json:"gadgetId"`
}

func getWeapon() WeaponConfig {
	weaponDataJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/WeaponExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}

	var weapons WeaponConfig
	if err := json.Unmarshal([]byte(weaponDataJson), &weapons); err != nil {
		log.Fatal(err)
	}

	fmt.Println(weapons)

	return weapons
}

func fetchJsonFromUrl(path string) (string, error) {

	resp, err := http.Get(path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%v: %v", resp.Status, path)
	}

	out, err := io.ReadAll(resp.Body)
	return string(out), err
}
