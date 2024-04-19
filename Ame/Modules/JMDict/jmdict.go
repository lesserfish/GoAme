package jmdict

import (
	"log"
	module "github.com/lesserfish/GoAme/Ame/Modules"
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

	log.Println("JMdict Module initialized!")
	return newModule, nil
}
func (parser JMdictModule) Close() {

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
		return nil
	}

	kanji := input["kanjiword"]
	kana := input["kanaword"]

	entry, err := FindEntry(&parser.Dictionary, kanji, kana, ignore_kanji, ignore_kana)

	if err != nil {
		return err
	}

	CleanEntry(&entry, &parser.Formatter)

	if err != nil {
		return err
	}

	keymap, err := KeymapFromEntry(&entry)

	if err != nil {
		return err
	}

    card.AddToFields("Kanjiword", keymap["kanjiword"])
    card.AddToFields("Kanaword", keymap["kanaword"])
    card.AddToFields("Sense", keymap["sense"])

	if err != nil {
		return err
	}

	return nil
}

