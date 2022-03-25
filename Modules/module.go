package module

type Input map[string]string

type Card struct {
	Fields []string
	Tag    string
}
type Module interface {
	Demo()
	Render(Input, *Card) error
	CSS(*Card)
	Close()
}

func (card Card) Parse(keymap map[string]string, clear_unused bool) {
	for id, field := range card.Fields {

		rendered := ParseString(field, keymap, clear_unused)
		card.Fields[id] = rendered
	}

	tag := ParseString(card.Tag, keymap, clear_unused)

	card.Tag = tag
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
