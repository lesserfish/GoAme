package anki

import (
	"log"

	module "github.com/lesserfish/GoAme/Modules"
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
func (ankiModule AnkiModule) Demo() {
}
func (ankiModule AnkiModule) Render(input module.Input, card *module.Card) (err error) {

	//for id, field := range card.Fields {
	//updated_field := strings.ReplaceAll(field, "\"", "\"\"")
	//card.Fields[id] = updated_field
	//}

	//	updated_tag := strings.ReplaceAll(card.Tag, "\"", "\"\"")
	//	card.Tag = updated_tag

	out := make(map[string]string)
	card.Parse(out, true)
	return nil
}
func (ankiModule AnkiModule) CSS() string {
	return ""
}
