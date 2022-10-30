package main

import (
	"encoding/json"
	"log"
)

type TextMap map[int]string

//need avatarIds to simplify textmap
func generateAvatarNameMap(avatarIdToTextMapId map[int]int) map[int]string {
	enTextMapJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/TextMap/TextMapEN.json")
	if err != nil {
		log.Fatal(err)
	}
	var enTextMap TextMap
	if err := json.Unmarshal([]byte(enTextMapJson), &enTextMap); err != nil {
		log.Fatal(err)
	}

	avatarNameMap := make(map[int]string)
	for avatarId, textMapId := range avatarIdToTextMapId {
		avatarNameMap[avatarId] = enTextMap[textMapId]
	}
	return avatarNameMap
}
