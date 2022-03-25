package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	examples "github.com/lesserfish/GoAme/Modules/Examples"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
)

var modules []module.Module

func main() {

	jsmod, err := jmdict.Initialize(jmdict.InitOptions{"Repository/Vocabulary/JMdict_e_examp.xml", "Tools/POLXML/out.xml"})

	if err != nil {
		fmt.Println(err)
	}

	exmod, err := examples.Initialize(examples.InitOptions{"Database/Sentences.db", true})

	if err != nil {
		fmt.Println(err)
	}

	card := module.Card{[]string{"@{Kanji}", "@{Kana} @{Sense} <h2>Examples:</h2> @{Example_1} @{Example_1_ENG} @{Example_1_JP}"}, ""}
	input := module.Input{"kanji": "警察"}

	err = jsmod.Render(input, &card)

	if err != nil {
		fmt.Println(err)
	}

	err = exmod.Render(input, &card)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(card)
}
