package main

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/keys"
)

func main() {

	for _, set := range keys.SetNames {
		fmt.Println(set)
	}
}

func writeSetNameJSON(outputPath string) error {
	//delete existing

}
