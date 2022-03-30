package examples

import (
	"bytes"
	"database/sql"
	"io/ioutil"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	module "github.com/lesserfish/GoAme/Ame/Modules"
	_ "github.com/mattn/go-sqlite3"
)

type InitOptions struct {
	DBPath  string
	Shuffle bool
	Seed    int64
	CSSPath string
}

type Example struct {
	JP  string
	ENG string
}
type ExampleModule struct {
	DB         *sql.DB
	Shuffle    bool
	Seed       int64
	CSSContent string
}

func Initialize(options InitOptions) (*ExampleModule, error) {
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

	CSSdata, err := ioutil.ReadFile(options.CSSPath)

	if err != nil {
		return newModule, err
	}

	newModule.CSSContent = strings.ReplaceAll(bytes.NewBuffer(CSSdata).String(), "\n", "")
	log.Println("Example Module initialized!")
	return newModule, nil
}

func (exampleModule ExampleModule) Close() {
	exampleModule.DB.Close()
}
func (exampleModule ExampleModule) Demo() {

}
func (exampleModule ExampleModule) Render(input module.Input, card *module.Card) (err error) {
	Word := input["kanjiword"]

	if Word == "" {
		Word = input["kanaword"]
	}

	tx, err := exampleModule.DB.Begin()

	if err != nil {
		return err
	}

	idp, err := tx.Prepare("select id from wordmap where word == ?")

	if err != nil {
		return err
	}

	defer idp.Close()

	rowid, err := idp.Query(Word)

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

	card.Parse(keymap, false)

	return nil
}
func (exampleModule ExampleModule) CSS() string {
	return exampleModule.CSSContent
}

func KeymapFromEntry(examples []Example) (out map[string]string) {
	out = make(map[string]string)

	out["example"] = ""
	for i := 0; i < 100; i++ {
		out["example_"+strconv.Itoa(i)] = ""
		out["example_"+strconv.Itoa(i)+"_eng"] = ""
		out["example_"+strconv.Itoa(i)+"_jp"] = ""
	}

	if len(examples) == 0 {
		return out
	}

	canonicalvalue := "<div class = 'rexample' id = '0'><div class = 'JP'>" + examples[0].JP + "</div><div class = 'ENG'>" + examples[0].ENG + "</div></div>"
	out["example"] = canonicalvalue
	for id, ex := range examples {
		key := "example_" + strconv.Itoa(id)
		value := "<div class = 'rexample' id = '" + strconv.Itoa(id) + "'><div class = 'JP'>" + ex.JP + "</div><div class = 'ENG'>" + ex.ENG + "</div></div>"
		out[key] = value
		key = "example_" + strconv.Itoa(id) + "_jp"
		value = "<div class = 'rexample' id = '" + strconv.Itoa(id) + "'><div class = 'JP'>" + ex.JP + "</div></div>"
		out[key] = value
		key = "example_" + strconv.Itoa(id) + "_eng"
		value = "<div class = 'rexample' id = '" + strconv.Itoa(id) + "'><div class = 'ENG'>" + ex.ENG + "</div></div>"
		out[key] = value

	}
	return out
}
