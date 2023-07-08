package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	files, err := ioutil.ReadDir("../../../internal/artifacts")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
