package module

import (
    "strconv"
    "fmt"
)

type Input map[string]string

type Card struct {
    Kanaword string
    Kanjiword string
    Audio string
    Sense string
    Kanjiinfo string
    Kanjisinfo []string
    Kanjiinfoex string
    Kanjisinfoex []string
    Stroke string
    Strokes []string
    Literal string
    Literals []string
    Examples []string
	Tag    string
}

type Module interface {
	Render(Input, *Card) error
}

func PrepareArrays(keyword string, max_selection int) []string {
    var output []string
    for i := 1; i <= max_selection; i++ {
        output = append(output, "@{" + keyword + "_" + strconv.Itoa(i) + "}")
    }
    return output
}
func NewCard(max_selection int) Card {
    card := Card{
        Kanaword:     "",
        Kanjiword:    "",
        Audio:        "",
        Sense:        "",
        Kanjiinfo:    "",
        Kanjisinfo:   make([]string, 0),
        Kanjiinfoex:  "",
        Kanjisinfoex: make([]string, 0),
        Stroke:       "",
        Strokes:      make([]string, 0),
        Literal:      "",
        Literals:     make([]string, 0),
        Examples:     make([]string, 0),
        Tag:          "",
    }
    return card
}

func (card *Card) AddToFields(field string, content string) {
    switch field {
    case "Kanaword":
        card.Kanaword = content
    case "Kanjiword":
        card.Kanjiword = content
    case "Audio":
        card.Audio = content
    case "Sense":
        card.Sense = content
    case "Kanjiinfo":
        card.Kanjiinfo = content
    case "Kanjisinfo":
        card.Kanjisinfo = append(card.Kanjisinfo, content)
    case "Kanjiinfoex":
        card.Kanjiinfoex = content
    case "Kanjisinfoex":
        card.Kanjisinfoex = append(card.Kanjisinfoex, content)
    case "Stroke":
        card.Stroke = content
    case "Strokes":
        card.Strokes = append(card.Strokes, content)
    case "Literal":
        card.Literal = content
    case "Literals":
        card.Literals = append(card.Literals, content)
    case "Examples":
        card.Examples = append(card.Examples, content)
    case "Tag":
        card.Tag = content
    }
}

func StandardizeStrings (input []string, max_selection int) string {
    out := ""
    for i := 0; i < max_selection; i++ {
        if len(input) > i {
            fmt.Printf(">>> %s", input[i])
            out += input[i] + ";"
        } else {
            out += ";"
        }
    }
    return out
}
func (card Card) Render(max_selection int) (out string) {
    out += card.Kanaword + ";"
    out += card.Kanjiword + ";"
    out += card.Audio + ";"
    out += card.Sense + ";"
    out += card.Kanjiinfo + ";"
    out += StandardizeStrings(card.Kanjisinfo, max_selection)
    out += card.Kanjiinfoex + ";"
    out += StandardizeStrings(card.Kanjisinfoex, max_selection)
    out += card.Stroke + ";"
    out += StandardizeStrings(card.Strokes, max_selection)
    out += card.Literal + ";"
    out += StandardizeStrings(card.Literals, max_selection)
    out += StandardizeStrings(card.Examples, max_selection)
    out += card.Tag
    return out
}

func (original Card) Copy() Card {
    // Create a new Card instance
    copiedCard := Card{
        Kanaword:     original.Kanaword,
        Kanjiword:    original.Kanjiword,
        Audio:        original.Audio,
        Sense:        original.Sense,
        Kanjiinfo:    original.Kanjiinfo,
        Kanjisinfo:   make([]string, len(original.Kanjisinfo)),
        Kanjiinfoex:  original.Kanjiinfoex,
        Kanjisinfoex: make([]string, len(original.Kanjisinfoex)),
        Stroke:       original.Stroke,
        Strokes:      make([]string, len(original.Strokes)),
        Literal:      original.Literal,
        Literals:     make([]string, len(original.Literals)),
        Examples:     make([]string, len(original.Literals)),
        Tag:          original.Tag,
    }

    // Copy slices to avoid modifying the original card's slices
    copy(copiedCard.Kanjisinfo, original.Kanjisinfo)
    copy(copiedCard.Kanjisinfoex, original.Kanjisinfoex)
    copy(copiedCard.Strokes, original.Strokes)
    copy(copiedCard.Literals, original.Literals)
    copy(copiedCard.Examples, original.Examples)

    return copiedCard
}

