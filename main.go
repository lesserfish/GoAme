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
	if len(args) < 2 {
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

	ameinstance := ame.Initialize(config)

	input_file := args[1]

	input_content, err := ioutil.ReadFile(input_file)

	if err != nil {
		log.Print(err)
		return
	}
	var input ame.Input

	err = json.Unmarshal(input_content, &input)

	if err != nil {
		log.Print(err)
		return
	}

	out := ameinstance.Render(input)

	ioutil.WriteFile("/home/lesserfish/test.txt", []byte(out), 0666)
	fmt.Println(out)
}
