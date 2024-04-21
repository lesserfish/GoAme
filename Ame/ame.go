package ame

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
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
type Input []module.Input

type AmeKanji struct {
	modules    []module.Module
	ankiModule module.Module
}

func ParseEnv(input string) string {
	for id, char := range input {
		if char == ')' {
			end := id
			if end > 0 {
				for start := end - 1; start >= 0; start-- {
					if input[start] == ' ' {
						break
					}
					if input[start] == '$' && input[start+1] == '(' {
						prelude := input[0:start]
						prologue := input[end+1:]
						env := input[start+2 : end]
						envval := os.Getenv(env)
						newstring := prelude + envval + prologue
						return ParseEnv(newstring)
					}
				}
			}
		}
	}
	return input
}

func ParseConfiguration(Input Configuration) Configuration {
	output := make(Configuration)

	for modname, modconf := range Input {
		outconf := make(map[string]string)
		for key, val := range modconf {
			updated_val := ParseEnv(val)
			outconf[key] = updated_val
		}
		output[modname] = outconf
	}

	return output
}
func Initialize(config Configuration) (*AmeKanji, error) {
	ameInstance := new(AmeKanji)

	config = ParseConfiguration(config)

	_, jmdict_ok := config["JMdict"]
	if jmdict_ok {
		jmdict_init := jmdict.InitOptions{}
		jmdict_init.DictionaryPath = config["JMdict"]["DictionaryPath"]
		jmdict_init.FormatterPath = config["JMdict"]["FormatterPath"]

		jmdict_mod, err := jmdict.Initialize(jmdict_init)

		if err != nil {
			errmsg := "Failed to initialize JMDict module. Error: " + err.Error()
			return ameInstance, errors.New(errmsg)
		} else {

			ameInstance.modules = append(ameInstance.modules, jmdict_mod)

			_, audio_ok := config["Audio"]
			if audio_ok {
				audio_init := audio.InitOptions{}
				audio_init.AudioPath = config["Audio"]["Path"]
				audio_init.JMdictMod = jmdict_mod

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

func CharIsHiragana(c rune) bool {
	return '\u3040' <= c && c <= '\u309F'
}

func CharIsKatakana(c rune) bool {
	return '\u30A0' <= c && c <= '\u30FF'
}

func CharIsKanji (c rune) bool {
	return ('\u4E00' <= c && c <= '\u9FaF') || ('\u3400' <= c && c <= '\u4DBf')
}

func CleanKanji(input string) string {
    for _, c:= range input {
        if CharIsKanji(c) {
            return input
        }
    }
    return ""
}

func CleanInput (input map[string]string) {
    kanjiword, exists := input["kanjiword"]
    if exists {
        input["kanjiword"] = CleanKanji(kanjiword)
    }
}
func CleanRender(input map[string]string) string {
	copy := make(map[string]string)
	for key, value := range input {
		copy[key] = value
	}

    kanjiword, exists := input["kanjiword"]
    if exists {
        newkw := CleanKanji(kanjiword)
        fmt.Printf("Changed %s to %s\n", kanjiword, newkw);
        copy["kanjiword"] = CleanKanji(kanjiword)
        fmt.Printf("Entry after update: %s\n", copy["kanjiword"])
    }
	delete(copy, "savepath")
	return fmt.Sprint(copy)
}

func ValidateString(input string) bool {
    for _, r := range input {
        if !CharIsHiragana(r) && !CharIsKatakana(r) && !CharIsKanji(r) {
            return false
        }
    }
    return true
}
func ValidateInput(input map[string] string) bool {
    
    literal, exists := input["literal"]
    if exists {
        if !ValidateString(literal) {
            return false
        }
    }
    kanaword, exists := input["kanaword"]
    if exists {
        if !ValidateString(kanaword) {
            return false
        }
    }
    kanjiword, exists := input["kanjiword"]
    if exists {
        if !ValidateString(kanjiword) {
            return false
        }
    }
    return true
}
func (ameKanji AmeKanji) URender(input Input, updatefunc UpdateFunc) (out string, errorlog string) {

	activeModules := []module.Module{}

	for _, module := range ameKanji.modules {
		activeModules = append(activeModules, module)
	}

	for id := range input {

        entry := input[id]
        CleanInput(entry);

        if !ValidateInput(entry){
            continue
        }


		var progress float64 = 0.0
		progress = float64(id) / float64(len(input))

		currentCard := module.NewCard()

		for _, mod := range activeModules {

			err := mod.Render(entry, &currentCard)

			if err != nil {
				currentinput := CleanRender(entry)
				errmsg := fmt.Sprintf("Error rendering card %s.\nError: %s", currentinput, err.Error())
				errorlog += errmsg + "\n"
			}
		}
		// Anki Module

		err := ameKanji.ankiModule.Render(entry, &currentCard)
		if err != nil {
			currentinput := CleanRender(entry)
			errmsg := fmt.Sprintf("Error rendering card %s.\nError: %s", currentinput, err.Error())
			errorlog += errmsg + "\n"
		}

        // TODO: CHANGE 10 TO A VARIABLE
		out += currentCard.Render(10) + "\n"

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
