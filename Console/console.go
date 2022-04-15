package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	ame "github.com/lesserfish/GoAme/Ame"
)

func main() {

	var config_file string
	var input_file string
	var output_file string
	var log_file string

	flag.StringVar(&config_file, "c", "", "path of configuration json file")
	flag.StringVar(&input_file, "i", "", "path of input json file")
	flag.StringVar(&output_file, "o", "", "path of the output file, containing the anki deck in .txt format")
	flag.StringVar(&log_file, "log", "./log.txt", "path of the log file.")

	flag.Parse()

	if len(config_file) == 0 || len(input_file) == 0 || len(output_file) == 0 {
		flag.PrintDefaults()
		return
	}

	config_content, err := ioutil.ReadFile(config_file)

	if err != nil {
		log.Println(err)
		return
	}

	var config ame.Configuration
	json.Unmarshal(config_content, &config)

	ameinstance, err := ame.Initialize(config)

	if err != nil {
		log.Fatalln(err)
	}

	input_content, err := ioutil.ReadFile(input_file)

	if err != nil {
		log.Fatalln(err)
	}
	var input ame.Input

	err = json.Unmarshal(input_content, &input)

	if err != nil {
		log.Fatalln(err)
	}

	out, log := ameinstance.URenderAndSave(input, output_file, log_file, func(progress float64) { fmt.Println("Progress: " + fmt.Sprint(progress*100) + "%") })
	fmt.Println(out)
	fmt.Println(log)
}
