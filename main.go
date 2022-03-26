package main

import (
	"fmt"
	"io/ioutil"

	module "github.com/lesserfish/GoAme/Modules"
	anki "github.com/lesserfish/GoAme/Modules/Anki"
	audio "github.com/lesserfish/GoAme/Modules/Audio"
	examples "github.com/lesserfish/GoAme/Modules/Examples"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
	kanjidic "github.com/lesserfish/GoAme/Modules/Kanjidic"
	strokes "github.com/lesserfish/GoAme/Modules/Strokes"
)

var modules []module.Module

func main() {

	modj, e1 := jmdict.Initialize(jmdict.InitOptions{"Repository/Vocabulary/JMdict_e_examp.xml", "Tools/POLXML/out.xml", "Modules/JMDict/Static/default.css"})
	modk, e2 := kanjidic.Initialize(kanjidic.InitOptions{"Repository/Kanji/kanjidic2.xml", "Modules/Kanjidic/Static/default.css"})
	mods, e3 := strokes.Initialize(strokes.InitOptions{"Repository/Strokes/sodzip", modk, false, "Modules/Strokes/Static/default.css"})
	mode, e4 := examples.Initialize(examples.InitOptions{"Database/Sentences.db", true, 0, "Modules/Examples/Static/default.css"})
	moda, e5 := audio.Initialize(audio.InitOptions{"http://localhost:8000/?", modj, "Modules/Audio/Static/default.css"})
	end, _ := anki.Initialize(anki.InitOptions{})

	for _, e := range []error{e1, e2, e3, e4, e5} {
		if e != nil {
			fmt.Println(e)
		}
	}
	modules = []module.Module{modj, modk, mods, mode, moda, end}

	card := module.Card{[]string{"@{asdasd}<b>@{KanjiWord}</b> @{KanaWord} @{Sense}", "@{KanjiInfoEx}", "@{Example} @{Example_1}", "@{Stroke}", " @{Audio}"}, "tag"}
	input := module.Input{
		"kanjiword": "警察",
		"literal":   "警察",
		"savepath":  "/home/lesserfish/.local/share/Anki2/Dev/collection.media/"}

	for _, mod := range modules {
		e := mod.Render(input, &card)
		if e != nil {
			fmt.Println(e)
		}
	}

	CSS := "<style>"
	for _, mod := range modules {
		CSS += mod.CSS()
	}
	CSS += "</style>"
	card.AddToFields(CSS)

	out := card.Render()

	fmt.Println(out)

	ioutil.WriteFile("/home/lesserfish/Documents/tmp/test.txt", []byte(out), 0644)
}
