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

func FormatField(input string) string {
    return strings.ReplaceAll(input, "\"", "\"\"")
}
func FormatFields(input []string) []string {
    output := make([]string, 0)
    for _, value := range input {
        new_value := strings.ReplaceAll(value, "\"", "\"\"")
        output = append(output, new_value)
    }
    return output
}
func (ankiModule AnkiModule) Render(input module.Input, card *module.Card) (err error) {

    // TODO: Do the rest
	card.Kanaword = FormatField(card.Kanaword)
	card.Kanjiword = FormatField(card.Kanjiword)
	card.Audio = FormatField(card.Audio)
	card.Sense = FormatField(card.Sense)
	card.Kanjiinfo = FormatField(card.Kanjiinfo)
	card.Kanjisinfo = FormatFields(card.Kanjisinfo)
	card.Kanjiinfoex = FormatField(card.Kanjiinfoex)
	card.Kanjisinfoex = FormatFields(card.Kanjisinfoex)
	card.Stroke = FormatField(card.Stroke)
	card.Strokes = FormatFields(card.Strokes)
	card.Literal = FormatField(card.Literal)
	card.Literals = FormatFields(card.Literals)
    card.Examples = FormatFields(card.Examples)
	card.Tag = FormatField(card.Tag)
	card.Tag = FormatField(card.Tag)
	return nil
}
func (ankiModule AnkiModule) CSS() string {
	return ""
}
func (ankiModule AnkiModule) Active(Fields []string) (out bool) {
	return true
}
