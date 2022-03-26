package main

import (
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
	examples "github.com/lesserfish/GoAme/Modules/Examples"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
	kanjidic "github.com/lesserfish/GoAme/Modules/Kanjidic"
	strokes "github.com/lesserfish/GoAme/Modules/Strokes"
)

var modules []module.Module

func main() {

	modj, e1 := jmdict.Initialize(jmdict.InitOptions{"Repository/Vocabulary/JMdict_e_examp.xml", "Tools/POLXML/out.xml"})
	modk, e2 := kanjidic.Initialize(kanjidic.InitOptions{"Repository/Kanji/kanjidic2.xml"})
	mods, e3 := strokes.Initialize(strokes.InitOptions{"Repository/Strokes/sodzip", modk, false})
	mode, e4 := examples.Initialize(examples.InitOptions{"Database/Sentences.db", true, 0})

	for _, e := range []error{e1, e2, e3, e4} {
		if e != nil {
			fmt.Println(e)
		}
	}

	card := module.Card{[]string{"<b>@{KanjiWord}</b> @{KanaWord} @{Sense} @{KanjiInfoEx} @{Example} @{Example_1} @{Stroke}"}, ""}
	input := module.Input{
		"kanjiword": "警察",
		"literal":   "警察",
		"savepath":  "/home/lesserfish/Documents/tmp"}

	e1 = modj.Render(input, &card)
	e2 = modk.Render(input, &card)
	e3 = mods.Render(input, &card)
	e4 = mode.Render(input, &card)

	for _, e := range []error{e1, e2, e3, e4} {
		if e != nil {
			fmt.Println(e)
		}
	}
	fmt.Println(card)
}
