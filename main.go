package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	ame "github.com/lesserfish/GoAme/Ame"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Usage:\n\tame configuration.json input.json")
		return
	}

	config_file := args[0]
	config_content, err := ioutil.ReadFile(config_file)

	if err != nil {
		log.Println(err)
		return
	}

	var config ame.Configuration
	json.Unmarshal(config_content, &config)

	ame.Initialize(config)
}
