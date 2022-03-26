package kanjidic

import (
	"errors"

	module "github.com/lesserfish/GoAme/Modules"
)

type Kanjidic_Module struct {
	DictionaryPath string
	Dictionary     Kanjidic
}
type InitOptions struct {
	DictionaryPath string
}

func Initialize(options InitOptions) (*Kanjidic_Module, error) {
	newModule := new(Kanjidic_Module)
	newModule.DictionaryPath = options.DictionaryPath

	err := LoadDictionary(newModule)

	if err != nil {
		return newModule, err
	}

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
func (parser Kanjidic_Module) CSS(card *module.Card) {

}
