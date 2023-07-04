package attributes

import "encoding/json"

type ElementMap map[Element]float64

func (e ElementMap) MarshalJSON() ([]byte, error) {
	stringRep := make(map[string]float64)
	for key, value := range e {
		stringRep[ElementString[key]] = value
	}
	return json.Marshal(stringRep)
}
