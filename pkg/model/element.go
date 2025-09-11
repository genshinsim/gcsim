package model

var elementoToString = map[Element]string{
	Element_Electric:    "electro",
	Element_Fire:        "pyro",
	Element_Water:       "hydro",
	Element_Grass:       "dendro",
	Element_Wind:        "anemo",
	Element_Ice:         "cryo",
	Element_Rock:        "geo",
	Element_Frozen:      "frozen",
	Element_Overdose:    "quicken",
	Element_Burning:     "burning",
	Element_BurningFuel: "dendro-fuel",
}

func ElementToString(e Element) string {
	return elementoToString[e]
}
