package kanjidic

import (
	"errors"

	module "github.com/lesserfish/GoAme/Modules"
)

type Kanjidic_Module struct {
	DictionaryPath string
	dictionary     Kanjidic
}
type InitOptions struct {
	DictionaryPath string
}

func Initialize(options InitOptions) (out module.Module, err error) {
	newModule := new(Kanjidic_Module)
	newModule.DictionaryPath = options.DictionaryPath

	err = LoadDictionary(newModule)

	if err != nil {
		return out, err
	}

	out = *newModule

	return out, nil
}
func (parser Kanjidic_Module) Close() {

}
func (parser Kanjidic_Module) Demo() {

}
func (parser Kanjidic_Module) Render(input module.Input, card *module.Card) error {
	if len(input) < 1 {
		return errors.New("No input given to Kanjidic module!")
	}

	literal := input["literal"]

	entry, err := FindEntry(&parser.dictionary, literal)

	if err != nil {
		return err
	}

	keymap, err := KeymapFromEntry(&entry)

	if err != nil {
		return err
	}

	card.Render(keymap, false)

	if err != nil {
		return err
	}

	return nil
}
func (parser Kanjidic_Module) CSS(card *module.Card) {

}
