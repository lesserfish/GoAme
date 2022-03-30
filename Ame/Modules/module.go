package module

type Input map[string]string

type Card struct {
	Fields []string
	Tag    string
}
type Module interface {
	Demo()
	Render(Input, *Card) error
	CSS() string
}

func (card Card) Parse(keymap map[string]string, clear_unused bool) {
	for id, field := range card.Fields {

		rendered := ParseString(field, keymap, clear_unused)
		card.Fields[id] = rendered
	}

	tag := ParseString(card.Tag, keymap, clear_unused)

	card.Tag = tag
}
func (card Card) AddToFields(content string) {
	for id, field := range card.Fields {
		newfield := field + content
		card.Fields[id] = newfield
	}
}

func (card Card) Render() (out string) {
	for _, field := range card.Fields {
		out += "\"" + field + "\"" + ";"
	}
	out += card.Tag

	return out
}

func (card Card) Copy() Card {
	newcard := Card{}
	for _, field := range card.Fields {
		newcard.Fields = append(newcard.Fields, field)
	}
	newcard.Tag = card.Tag
	return newcard
}
func ParseString(input string, keymap map[string]string, clear_unused bool) string {
	for i, c := range input {
		if c == '}' {
			end := i

			for start := end; start >= 0; start-- {
				if start+1 >= len(input) {
					continue
				}
				if input[start+1] == '{' && input[start] == '@' {
					lowsegment := input[0:start]
					highsegment := ""

					if end+1 < len(input) {
						highsegment = input[end+1:]
					}

					key := input[start+2 : end]

					value, ok := keymap[key]

					if !ok {
						if !clear_unused {
							break
						}
						value = ""
					}

					translation := lowsegment + value + highsegment
					return ParseString(translation, keymap, clear_unused)
				}
			}
		}
	}

	return input
}
