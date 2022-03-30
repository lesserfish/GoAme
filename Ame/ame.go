package ame

import (
	"fmt"
	"log"
	"strconv"

	module "github.com/lesserfish/GoAme/Ame/Modules"
	anki "github.com/lesserfish/GoAme/Ame/Modules/Anki"
	audio "github.com/lesserfish/GoAme/Ame/Modules/Audio"
	examples "github.com/lesserfish/GoAme/Ame/Modules/Examples"
	jmdict "github.com/lesserfish/GoAme/Ame/Modules/JMDict"
	kanjidic "github.com/lesserfish/GoAme/Ame/Modules/Kanjidic"
	strokes "github.com/lesserfish/GoAme/Ame/Modules/Strokes"
)

type Configuration map[string]map[string]string
type Input struct {
	Template module.Card
	Input    []module.Input
}

type AmeKanji struct {
	modules []module.Module
}

func Initialize(config Configuration) *AmeKanji {
	ameInstance := new(AmeKanji)

	_, jmdict_ok := config["JMdict"]
	if jmdict_ok {
		jmdict_init := jmdict.InitOptions{}
		jmdict_init.DictionaryPath = config["JMdict"]["DictionaryPath"]
		jmdict_init.FormatterPath = config["JMdict"]["FormatterPath"]
		jmdict_init.CSSPath = config["JMdict"]["CSSPath"]

		jmdict_mod, err := jmdict.Initialize(jmdict_init)

		if err != nil {
			errmsg := "Failed to initialize JMDict module. Error: " + err.Error()
			log.Println(errmsg)
		} else {

			ameInstance.modules = append(ameInstance.modules, jmdict_mod)

			_, audio_ok := config["Audio"]
			if audio_ok {
				audio_init := audio.InitOptions{}
				audio_init.URI = config["Audio"]["URI"]
				audio_init.JMdictMod = jmdict_mod
				audio_init.CSSPath = config["Audio"]["CSSPath"]

				audio_mod, err := audio.Initialize(audio_init)

				if err != nil {
					errmsg := "Failed to initialize Audio module. Error: " + err.Error()
					log.Println(errmsg)
				} else {
					ameInstance.modules = append(ameInstance.modules, audio_mod)
				}
			}
		}

	}

	_, kanjidic_ok := config["Kanjidic"]
	if kanjidic_ok {
		kanjidic_init := kanjidic.InitOptions{}
		kanjidic_init.DictionaryPath = config["Kanjidic"]["DictionaryPath"]
		kanjidic_init.CSSPath = config["Kanjidic"]["CSSPath"]

		kanjidic_mod, err := kanjidic.Initialize(kanjidic_init)

		if err != nil {
			errmsg := "Failed to initialize Kanjidic module. Error: " + err.Error()
			log.Println(errmsg)
		} else {

			ameInstance.modules = append(ameInstance.modules, kanjidic_mod)

			_, strokes_ok := config["Strokes"]
			if strokes_ok {
				strokes_init := strokes.InitOptions{}
				strokes_init.StrokePath = config["Strokes"]["StrokePath"]
				strokes_init.Kanjimod = kanjidic_mod
				strokes_init.CSSPath = config["Strokes"]["CSSPath"]
				strokes_init.PreferJIS, _ = strconv.ParseBool(config["Strokes"]["PreferJIS"])

				strokes_mod, err := strokes.Initialize(strokes_init)

				if err != nil {
					errmsg := "Failed to initialize Strokes module. Error: " + err.Error()
					log.Println(errmsg)
				} else {
					ameInstance.modules = append(ameInstance.modules, strokes_mod)
				}
			}
		}

	}

	_, examples_ok := config["Examples"]
	if examples_ok {
		examples_init := examples.InitOptions{}
		examples_init.DBPath = config["Examples"]["DBPath"]
		examples_init.CSSPath = config["Examples"]["CSSPath"]
		examples_init.Seed, _ = strconv.ParseInt(config["Examples"]["Seed"], 10, 64)
		examples_init.Shuffle, _ = strconv.ParseBool(config["Examples"]["Shuffle"])

		examples_mod, err := examples.Initialize(examples_init)

		if err != nil {
			errmsg := "Failed to initialize Examples module. Error: " + err.Error()
			log.Println(errmsg)
		} else {
			ameInstance.modules = append(ameInstance.modules, examples_mod)
		}
	}

	_, anki_ok := config["Anki"]
	if anki_ok {
		anki_init := anki.InitOptions{}

		anki_mod, err := anki.Initialize(anki_init)

		if err != nil {
			errmsg := "Failed to initialize Anki module. Error: " + err.Error()
			log.Println(errmsg)
		} else {
			ameInstance.modules = append(ameInstance.modules, anki_mod)
		}
	}

	return ameInstance
}

func (ameKanji AmeKanji) Render(input Input) (out string) {

	for id, _ := range input.Input {

		fmt.Println(strconv.Itoa(id+1) + " / " + strconv.Itoa(len(input.Input)))

		currentCard := input.Template.Copy()
		currentCSS := ""

		for _, mod := range ameKanji.modules {

			err := mod.Render(input.Input[id], &currentCard)
			currentCSS += mod.CSS()

			if err != nil {
				errmsg := "Error rendering card. Error: " + err.Error()
				log.Println(errmsg)
			}

		}
		currentCard.AddToFields(currentCSS)
		out += currentCard.Render() + "\n"
	}

	return out
}
