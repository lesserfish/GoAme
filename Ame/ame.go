package ame

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	modules    []module.Module
	ankiModule module.Module
}

func Initialize(config Configuration) (*AmeKanji, error) {
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
			return ameInstance, errors.New(errmsg)
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
					return ameInstance, errors.New(errmsg)
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
			return ameInstance, errors.New(errmsg)
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
					return ameInstance, errors.New(errmsg)
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
		examples_init.MaxExamples, _ = strconv.ParseUint(config["Examples"]["MaxExamples"], 10, 64)

		examples_mod, err := examples.Initialize(examples_init)

		if err != nil {
			errmsg := "Failed to initialize Examples module. Error: " + err.Error()
			return ameInstance, errors.New(errmsg)
		} else {
			ameInstance.modules = append(ameInstance.modules, examples_mod)
		}
	}

	anki_mod, err := anki.Initialize(anki.InitOptions{})
	if err != nil {
		errmsg := "Failed to initialize Anki module. Error: " + err.Error()
		return ameInstance, errors.New(errmsg)
	} else {
		ameInstance.ankiModule = anki_mod
	}

	return ameInstance, nil
}

type UpdateFunc func(float64)

func CleanInput(input map[string]string) string {
	copy := make(map[string]string)
	for key, value := range input {
		copy[key] = value
	}
	delete(copy, "savepath")
	return fmt.Sprint(copy)
}
func (ameKanji AmeKanji) URender(input Input, updatefunc UpdateFunc) (out string, errorlog string) {

	activeModules := []module.Module{}

	for _, module := range ameKanji.modules {
		if module.Active(input.Template.Fields) {
			activeModules = append(activeModules, module)
		}
	}

	for id := range input.Input {

		var progress float64 = 0.0
		progress = float64(id) / float64(len(input.Input))

		currentCard := input.Template.Copy()
		currentCSS := ""

		for _, mod := range activeModules {
			err := mod.Render(input.Input[id], &currentCard)
			currentCSS += mod.CSS()

			if err != nil {
				currentinput := CleanInput(input.Input[id])
				errmsg := fmt.Sprintf("Error rendering card %s.\nError: %s", currentinput, err.Error())
				errorlog += errmsg + "\n"
			}
		}
		// Render CSS
		CSSMap := make(map[string]string)
		CSSMap["CSS"] = currentCSS
		currentCard.Parse(CSSMap, false)

		// Anki Module

		err := ameKanji.ankiModule.Render(input.Input[id], &currentCard)
		if err != nil {
			currentinput := CleanInput(input.Input[id])
			errmsg := fmt.Sprintf("Error rendering card %s.\nError: %s", currentinput, err.Error())
			errorlog += errmsg + "\n"
		}

		out += currentCard.Render() + "\n"

		updatefunc(progress)
	}

	return out, errorlog
}

func (ameKanji AmeKanji) URenderAndSave(input Input, path string, errpath string, updatefunc UpdateFunc) (string, string) {
	content, err := ameKanji.URender(input, updatefunc)
	ioutil.WriteFile(path, []byte(content), 0666)
	ioutil.WriteFile(errpath, []byte(err), 0666)
	return content, err
}

func (ameKanji AmeKanji) Render(input Input) (string, string) {
	return ameKanji.URender(input, func(f float64) {})
}

func (AmeKanji AmeKanji) RenderAndSave(input Input, path string, errpath string) (string, string) {
	return AmeKanji.URenderAndSave(input, path, errpath, func(f float64) {})
}
