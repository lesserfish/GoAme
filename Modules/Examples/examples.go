package examples

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	module "github.com/lesserfish/GoAme/Modules"
	_ "github.com/mattn/go-sqlite3"
)

type InitOptions struct {
	DBPath  string
	Shuffle bool
	Seed    int64
}

type Example struct {
	JP  string
	ENG string
}
type ExampleModule struct {
	DB      *sql.DB
	Shuffle bool
	Seed    int64
}

func Initialize(options InitOptions) (out module.Module, err error) {
	newModule := new(ExampleModule)
	db, err := sql.Open("sqlite3", options.DBPath)
	if err != nil {
		return newModule, err
	}

	newModule.DB = db
	newModule.Shuffle = options.Shuffle
	newModule.Seed = options.Seed
	if options.Seed == 0 {
		newModule.Seed = time.Now().UnixMicro()
	}

	fmt.Println("Examples loaded!")
	return newModule, err
}

func (exampleModule ExampleModule) Close() {
	exampleModule.DB.Close()
}
func (exampleModule ExampleModule) Demo() {

}
func (exampleModule ExampleModule) Render(input module.Input, card *module.Card) (err error) {
	Kanji := input["kanjiword"]

	tx, err := exampleModule.DB.Begin()

	if err != nil {
		return err
	}

	idp, err := tx.Prepare("select id from wordmap where word == ?")

	if err != nil {
		return err
	}

	defer idp.Close()

	rowid, err := idp.Query(Kanji)

	if err != nil {
		return err
	}

	defer rowid.Close()

	ids := []int{}

	for rowid.Next() {
		var currentid int
		err = rowid.Scan(&currentid)

		if err != nil {
			return err
		}

		ids = append(ids, currentid)
	}

	exp, err := tx.Prepare("select sentence,translation from sentmap where ID == ?")

	if err != nil {
		return err
	}

	defer exp.Close()

	examples := []Example{}

	for _, id := range ids {
		rowex, err := exp.Query(id)

		if err != nil {
			return err
		}
		var JP string
		var ENG string
		for rowex.Next() {
			err = rowex.Scan(&JP, &ENG)
			if err != nil {
				return err
			}
			JP = strings.TrimSuffix(JP, "\n")
			ENG = strings.TrimSuffix(ENG, "\n")
			examples = append(examples, Example{JP, ENG})
		}
	}

	rand.Seed(exampleModule.Seed)
	if exampleModule.Shuffle && len(examples) > 1 {
		rand.Shuffle(len(examples), func(i, j int) {
			examples[i], examples[j] = examples[j], examples[i]
		})
	}
	keymap := KeymapFromEntry(examples)

	card.Render(keymap, false)

	return nil
}
func (ExampleModule ExampleModule) CSS(card *module.Card) {

}

func KeymapFromEntry(examples []Example) (out map[string]string) {
	fmt.Println(examples)
	out = make(map[string]string)
	if len(examples) == 0 {
		return out
	}

	canonicalvalue := "<div class = 'rexample' id = '0'><div class = 'JP'</div>" + examples[0].JP + "<div class = 'ENG'>" + examples[0].ENG + "</div></div>"
	out["Example"] = canonicalvalue
	for id, ex := range examples {
		key := "Example_" + strconv.Itoa(id)
		value := "<div class = 'rexample' id = '" + strconv.Itoa(id) + "'><div class = 'JP'>" + ex.JP + "</div><div class = 'ENG'>" + ex.ENG + "</div></div>"
		out[key] = value
		key = "Example_" + strconv.Itoa(id) + "_JP"
		value = "<div class = 'rexample' id = '" + strconv.Itoa(id) + "'><div class = 'JP'>" + ex.JP + "</div></div>"
		out[key] = value
		key = "Example_" + strconv.Itoa(id) + "_ENG"
		value = "<div class = 'rexample' id = '" + strconv.Itoa(id) + "'><div class = 'ENG'>" + ex.ENG + "</div></div>"
		out[key] = value

	}
	return out
}
