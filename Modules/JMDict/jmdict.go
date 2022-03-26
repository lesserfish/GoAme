package jmdict

import (
	"errors"

	module "github.com/lesserfish/GoAme/Modules"
)

type JMdictModule struct {
	DictionaryPath string
	FormatterPath  string
	Dictionary     JMdict
	Formatter      RegexFormatter
}
type InitOptions struct {
	DictionaryPath string
	FormatterPath  string
}

func Initialize(options InitOptions) (*JMdictModule, error) {
	newModule := new(JMdictModule)
	newModule.DictionaryPath = options.DictionaryPath
	newModule.FormatterPath = options.FormatterPath

	err := LoadDictionary(newModule)

	if err != nil {
		return newModule, err
	}

	if options.FormatterPath != "" {
		err = LoadFormatter(newModule)

		if err != nil {
			return newModule, err
		}

	}

	return newModule, nil
}
func (parser JMdictModule) Close() {

}
func (parser JMdictModule) Demo() {
}
func (parser JMdictModule) Render(input module.Input, card *module.Card) error {
	ignore_kanji := false
	ignore_kana := false

	if input["kanjiword"] == "" {
		ignore_kanji = true
	}
	if input["kanaword"] == "" {
		ignore_kana = true
	}

	if ignore_kana && ignore_kanji {
		return errors.New("No input given to JMdic module!")
	}

	kanji := input["kanjiword"]
	kana := input["kanaword"]

	entry, err := FindEntry(&parser.Dictionary, kanji, kana, ignore_kanji, ignore_kana)

	if err != nil {
		return err
	}

	err = CleanEntry(&entry, &parser.Formatter)

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
func (parser JMdictModule) CSS(card *module.Card) {

}
