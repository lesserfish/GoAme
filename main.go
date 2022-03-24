package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
)

var modules []module.Module

func main() {
	init := jmdict.InitOptions{}
	init.DictionaryPath = "./Repository/Vocabulary/JMdict_e_examp.xml"
	init.FormatterPath = "./Tools/POLXML/out.xml"

	mod, err := jmdict.Initialize(init)

	if err != nil {
		fmt.Println(err)
		return
	}
	card := module.Card{[]string{"@{Kanji}", "@{Kanji}<br>@{Kana}<br><br>@{Sense}"}, "test"}
	input := module.Input{"警察"}

	err = mod.Render(input, &card)

	if err != nil {
		fmt.Println((err))
	}

	fmt.Println(card)
}
