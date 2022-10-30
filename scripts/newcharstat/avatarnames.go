package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type TextMap map[int]string

//need avatarIds to simplify textmap
func generateAvatarNameMap(textMapIds []int) map[int]string {
	enTextMapJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/TextMap/TextMapEN.json")
	if err != nil {
		log.Fatal(err)
	}
	var enTextMap TextMap
	if err := json.Unmarshal([]byte(enTextMapJson), &enTextMap); err != nil {
		log.Fatal(err)
	}

	avatarNameMap := make(map[int]string)
	for _, textMapId := range textMapIds {
		avatarNameMap[textMapId] = enTextMap[textMapId]
	}
	fmt.Printf("avatarNameMap: %+v", avatarNameMap)
	return avatarNameMap
}
