package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

//Provides mapping in JSON between gcsim keys and character data such as element, weapon type, etc...
func main() {
	//we pull the information from the curves pkg
	//but in reality this should all just be piped from dimbreath's repo instead
	res := make(map[string]curves.CharBase)

	//we start at TravelerDelim because traveler has to be handled as special case
	for i := keys.TravelerDelim + 1; i < keys.EndCharKeys; i++ {
		res[i.String()] = curves.CharBaseMap[i]
	}

	//handle traveler
	for i := keys.NoChar + 1; i < keys.TravelerDelim-2; i++ {
		key := i.String()
		//odd is aether, even is lumine
		name := "aether"
		if !strings.HasPrefix(key, name) {
			name = "lumine"
		}
		element := attributes.StringToEle(strings.TrimPrefix(key, name))

		res[key] = curves.CharBase{
			Rarity:     5,
			Body:       profile.BodyBoy,
			Element:    element,
			WeaponType: weapon.WeaponClassSword,
		}

		fmt.Println(res[key])
	}

	//write to file
	out, _ := json.MarshalIndent(res, "", " ")
	ioutil.WriteFile("./character_data.json", out, 0644)

}
