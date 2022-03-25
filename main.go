package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	examples "github.com/lesserfish/GoAme/Modules/Examples"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
	kanjidic "github.com/lesserfish/GoAme/Modules/Kanjidic"
)

var modules []module.Module

func main() {

	exmod, err := examples.Initialize(examples.InitOptions{"Database/Sentences.db", true, 0})
	mod, err := kanjidic.Initialize(kanjidic.InitOptions{"Repository/Kanji/kanjidic2.xml"})
	jmmod, err := jmdict.Initialize(jmdict.InitOptions{"Repository/Vocabulary/JMdict_e_examp.xml", "Tools/POLXML/out.xml"})

	if err != nil {
		fmt.Println(err)
	}

	card := module.Card{[]string{"@{KanjiWord} @{KanaWord} @{Sense} @{kanjiinfoex} @{Example}"}, ""}
	input := module.Input{"literal": "単語", "kanaword": "キリン"}

	err = mod.Render(input, &card)
	exmod.Render(input, &card)
	jmmod.Render(input, &card)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(card)
}
