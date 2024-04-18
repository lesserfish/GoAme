package anki

import (
	"log"
	"strings"

	module "github.com/lesserfish/GoAme/Ame/Modules"
)

type InitOptions struct {
}

type AnkiModule struct {
}

func Initialize(options InitOptions) (*AnkiModule, error) {
	newModule := new(AnkiModule)
	log.Println("Anki Module initialized!")
	return newModule, nil
}

func (ankiModule AnkiModule) Close() {
}
func (ankiModule AnkiModule) Render(input module.Input, card *module.Card) (err error) {

    // TODO: Do the rest
	card.Tag = strings.ReplaceAll(card.Tag, "\"", "\"\"")
	return nil
}
func (ankiModule AnkiModule) CSS() string {
	return ""
}
func (ankiModule AnkiModule) Active(Fields []string) (out bool) {
	return true
}
