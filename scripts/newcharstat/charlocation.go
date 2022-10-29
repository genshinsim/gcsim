package main

import (
	"encoding/json"
	"log"
)

type FetterInfo []struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int    `json:"avatarId"`
}

func getCharLocationMap() map[int]string {
	fetterInfoJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/FetterInfoExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}

	var fetterInfo FetterInfo
	if err := json.Unmarshal([]byte(fetterInfoJson), &fetterInfo); err != nil {
		log.Fatal(err)
	}

	// reshape fetterInfo to map of avatarId to AvatarAssocType
	locationMap := make(map[int]string)
	for _, v := range fetterInfo {
		locationMap[v.AvatarId] = v.AvatarAssocType
	}

	return locationMap
}
