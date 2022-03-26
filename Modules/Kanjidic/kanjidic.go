package kanjidic

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"strings"

	module "github.com/lesserfish/GoAme/Modules"
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

	newModule.CSSContent = strings.TrimSpace(bytes.NewBuffer(CSSdata).String())

	log.Println("Kanjidic Module initialized!")
	return newModule, nil
}
func (parser Kanjidic_Module) Close() {

}
func (parser Kanjidic_Module) Demo() {

}
func (parser Kanjidic_Module) Render(input module.Input, card *module.Card) error {
	if input["literal"] == "" {
		return errors.New("No input given to Kanjidic module!")
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

	card.Parse(keymap, false)

	if err != nil {
		return err
	}

	return nil
}
func (parser Kanjidic_Module) CSS() string {
	return parser.CSSContent
}
