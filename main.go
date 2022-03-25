package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	kanjidic "github.com/lesserfish/GoAme/Modules/Kanjidic"
	strokes "github.com/lesserfish/GoAme/Modules/Strokes"
)

var modules []module.Module

func main() {

	_, _, kmod := kanjidic.Initialize(kanjidic.InitOptions{"Repository/Kanji/kanjidic2.xml"})
	smod, _ := strokes.Initialize(strokes.InitOptions{"Repository/Strokes/sodzip", kmod, false})

	card := module.Card{[]string{"@{Stroke}"}, ""}
	input := module.Input{"literal": "単語", "savepath": "/home/lesserfish/Documents/tmp"}

	err := smod.Render(input, &card)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(card)
}
