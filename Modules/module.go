package module

import "errors"

type Input []string

type Card struct {
	Fields []string
	Tag    []string
}
type Module interface {
	Demo()
	Render(Input, *Card) error
}

func (card Card) Render(vars map[string]string) (out error) {
	return out
}

func RenderString(input string, keymap map[string]string) (string, error) {
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
						err_string := "Failed to find variable " + key + " in map"
						return input, errors.New(err_string)
					}

					translation := lowsegment + value + highsegment
					return RenderString(translation, keymap)
				}
			}
		}
	}

	return input, nil
}
