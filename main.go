package main

import (
	"fmt"
	"os"

	"github.com/armosec/armotoy/controller"
)

func main() {

	argsWithoutProg := os.Args[1:]

	srcType := "JSON"
	jsonPath := ""
	version := "v2"
	if len(argsWithoutProg) == 1 {
		jsonPath = argsWithoutProg[0]
	}

	ctrler, err := controller.InitController(jsonPath, srcType, version, "./config.json")
	if err != nil {
		fmt.Printf("unable to initialize controller due to: %s", err.Error())
		os.Exit(1)
	}
	ctrler.Start()
}
