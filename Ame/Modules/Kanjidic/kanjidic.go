package kanjidic

import (
	"bytes"
	"errors"
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
func (parser Kanjidic_Module) Active(Fields []string) (out bool) {
	keywords := []string{"kanjiinfo", "kaniinfoex"}

	out = false
keyword_search:
	for _, keyword := range keywords {
		key := fmt.Sprintf("@{%s}", keyword)

		for _, field := range Fields {
			if strings.Contains(field, key) {
				out = true
				break keyword_search
			}
		}
	}

	return out

}
