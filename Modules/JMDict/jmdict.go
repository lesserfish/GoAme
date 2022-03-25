package jmdict

import (
	"errors"
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
)

type JMdictModule struct {
	DictionaryPath string
	FormatterPath  string
	dictionary     JMdict
	formatter      RegexFormatter
}
type InitOptions struct {
	DictionaryPath string
	FormatterPath  string
}

func Initialize(options InitOptions) (out module.Module, err error) {
	newModule := new(JMdictModule)
	newModule.DictionaryPath = options.DictionaryPath
	newModule.FormatterPath = options.FormatterPath

	err = LoadDictionary(newModule)

	if err != nil {
		return out, err
	}

	if options.FormatterPath != "" {
		err = LoadFormatter(newModule)

		if err != nil {
			return out, err
		}

	}
	out = *newModule

	return out, nil
}
func (parser JMdictModule) Close() {

}
func (parser JMdictModule) Demo() {
	entry, _ := FindEntry(&parser.dictionary, "警察", "")
	CleanEntry(&entry, &parser.formatter)
	km, _ := KeymapFromEntry(&entry)

	fmt.Println(km)

}
func (parser JMdictModule) Render(input module.Input, card *module.Card) error {
	if len(input) < 1 {
		return errors.New("No input given to JMdict module!")
	}

	kanji := input["kanjiword"]
	kana := ""

	if len(input) > 1 {
		kana = input["kanaword"]
	}

	entry, err := FindEntry(&parser.dictionary, kanji, kana)

	if err != nil {
		return err
	}

	err = CleanEntry(&entry, &parser.formatter)

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
func (parser JMdictModule) CSS(card *module.Card) {

}
