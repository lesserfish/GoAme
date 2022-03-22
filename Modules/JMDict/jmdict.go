package jmdict

import (
	"errors"
	"fmt"

	module "github.com/lesserfish/GoAme/Modules"
)

type Parser struct {
	DictionaryPath string
	dictionary     JMdict
}
type InitOptions struct {
	DictionaryPath string
}

func Initialize(options InitOptions) (out module.Module, err error) {
	newParser := new(Parser)
	newParser.DictionaryPath = options.DictionaryPath

	err = LoadDictionary(newParser)

	out = *newParser

	return out, err
}
func (parser Parser) Demo() {
	entry, err := FindEntry(&parser.dictionary, "食べる", "")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print((entry))
	}

	CleanEntry(&entry)
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

	err = CleanEntry(&entry)

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
