package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// var weapons []curves.WeaponBase
	weaponConfigArray := getWeapon()

	for _, weaponConfig := range weaponConfigArray {
		// var weapon curves.WeaponBase
		// weapon. = weaponConfig.Name
		fmt.Println(weaponConfig.ID)
		fmt.Printf("%+v\n", weaponConfig.WeaponProp)
		// weapons = append(weapons, weapon)
	}

}

func getJsonFromFile[V WeaponConfigs](path string) V {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var configData V
	json.Unmarshal(byteValue, &configData)

	return configData
}
