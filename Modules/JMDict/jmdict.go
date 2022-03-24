package jmdict

import (
	"errors"
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
)

type Parser struct {
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
	newParser := new(Parser)
	newParser.DictionaryPath = options.DictionaryPath
	newParser.FormatterPath = options.FormatterPath

	err = LoadDictionary(newParser)

	if err != nil {
		return out, err
	}

	if options.FormatterPath != "" {
		err = LoadFormatter(newParser)

		if err != nil {
			return out, err
		}

	}
	out = *newParser

	return out, nil
}
func (parser Parser) Demo() {
	entry, _ := FindEntry(&parser.dictionary, "警察", "")
	CleanEntry(&entry, &parser.formatter)
	km, _ := KeymapFromEntry(&entry)

	fmt.Println(km)

}
func (parser Parser) Render(input module.Input, card *module.Card) error {
	if len(input) < 1 {
		return errors.New("No input given to JMdict module!")
	}

	kanji := input[0]
	kana := ""

	if len(input) > 1 {
		kana = input[1]
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

	err = card.Render(keymap)

	if err != nil {
		return err
	}

	return nil
}
