package kanjidic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	module "github.com/lesserfish/GoAme/Ame/Modules"
)

type Kanjidic_Module struct {
	DictionaryPath string
	Dictionary     Kanjidic
	CSSContent     string
}
type InitOptions struct {
	DictionaryPath string
	CSSPath        string
}

func Initialize(options InitOptions) (*Kanjidic_Module, error) {
	newModule := new(Kanjidic_Module)
	newModule.DictionaryPath = options.DictionaryPath

	err := LoadDictionary(newModule)

	if err != nil {
		return newModule, err
	}

	CSSdata, err := ioutil.ReadFile(options.CSSPath)

	if err != nil {
		return newModule, err
	}

	newModule.CSSContent = strings.ReplaceAll(bytes.NewBuffer(CSSdata).String(), "\n", "")

	log.Println("Kanjidic Module initialized!")
	return newModule, nil
}
func (parser Kanjidic_Module) Close() {

}
func (parser Kanjidic_Module) Render(input module.Input, card *module.Card) error {
	if input["literal"] == "" {
		return nil
	}

	literal := input["literal"]

	entry, err := FindEntry(&parser.Dictionary, literal)

	if err != nil {
		return err
	}

	keymap, err := KeymapFromEntry(&entry)

	if err != nil {
		return err
	}

    card.AddToFields("Kanjiinfo", keymap["kanjiinfo"])
    card.AddToFields("Kanjiinfoex", keymap["kanjiinfoex"])
    card.AddToFields("Literal", keymap["literal"])

    // Not optimal, but cleaner, I think...
    for i := 1; i < 100; i++ {

        key := fmt.Sprintf("literal_%d", i)
        value, exists := keymap[key]

        if !exists {
            break
        }

        card.AddToFields("Literals", value)
    }
    for i := 1; i < 100; i++ {

        key := fmt.Sprintf("kanjiinfo_%d", i)
        value, exists := keymap[key]

        if !exists {
            break
        }

        card.AddToFields("Kanjisinfo", value)
    }

    for i := 1; i < 100; i++ {

        key := fmt.Sprintf("kanjiinfoex_%d", i)
        value, exists := keymap[key]

        if !exists {
            break
        }

        card.AddToFields("Kanjisinfoex", value)
    }

	if err != nil {
		return err
	}

	return nil
}
