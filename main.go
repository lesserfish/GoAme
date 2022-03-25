package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	kanjidic "github.com/lesserfish/GoAme/Modules/Kanjidic"
)

var modules []module.Module

func main() {

	mod, err := kanjidic.Initialize(kanjidic.InitOptions{"Repository/Kanji/kanjidic2.xml"})

	if err != nil {
		fmt.Println(err)
	}

	card := module.Card{[]string{"@{kanjiinfoex}"}, ""}
	input := module.Input{"literal": "警察"}

	err = mod.Render(input, &card)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(card)
}
