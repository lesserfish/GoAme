package strokes

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	module "github.com/lesserfish/GoAme/Ame/Modules"
	kanjidic "github.com/lesserfish/GoAme/Ame/Modules/Kanjidic"
)

type InitOptions struct {
	StrokePath string
	Kanjimod   *kanjidic.Kanjidic_Module
	PreferJIS  bool
}

type StrokeModule struct {
	Path       string
	kanjimod   *kanjidic.Kanjidic_Module
	PreferJIS  bool
}

type StrokeOutput struct {
	Path string
	Type string
    ID int
}

func Initialize(options InitOptions) (*StrokeModule, error) {
	newModule := new(StrokeModule)
	newModule.Path = options.StrokePath
	newModule.kanjimod = options.Kanjimod
	newModule.PreferJIS = options.PreferJIS

	_, err := os.Stat(newModule.Path)

	if os.IsNotExist(err) {
		return newModule, err
	}

	log.Println("Stroke Module initialized!")

	return newModule, nil
}

func (strokeModule StrokeModule) Close() {
}
func (strokeModule StrokeModule) Render(input module.Input, card *module.Card) (err error) {
    
    if input["literal"] == "" {
		return nil
	}
    if input["savepath"] == "" {
		return errors.New("Unspecified output file path!")
	}


	literals := input["literal"]
	savepath := input["savepath"]

	output := []StrokeOutput{}

	characters, err := kanjidic.FindEntry(&strokeModule.kanjimod.Dictionary, literals)

	if err != nil {
		return err
	}

	for id, character := range characters {
		currentstroke := StrokeOutput{}
		ANDAS := ""
		JIS := ""
		for _, code := range character.DicNumber.DicRef {
			if code.DrType == "halpern_njecd" {
				ANDAS = code.Text
				break
			}
		}

		RAWJIS := ""
		for _, code := range character.Codepoint.CpValue {
			if code.CpType == "jis208" {
				RAWJIS = code.Text
				break
			}
		}

		JIScomponents := strings.Split(RAWJIS, "-")
		if len(JIScomponents) < 3 {
			break
		} else {
			JISrow, errr := strconv.Atoi(JIScomponents[1])
			JIScolumn, errc := strconv.Atoi(JIScomponents[2])

			if errr != nil || errc != nil {
				break
			}
			JISrow += 0x20
			JIScolumn += 0x20

			JIS = fmt.Sprintf("%x", JISrow) + fmt.Sprintf("%x", JIScolumn)
		}

		ANDAS = "ANDAS" + ANDAS + ".gif"
		JIS = JIS + ".gif"

		ANDASFP := strokeModule.Path + "/" + ANDAS
		JISFP := strokeModule.Path + "/" + JIS

		ANDASstat, err := os.Stat(ANDASFP)
		if err != nil || !ANDASstat.Mode().IsRegular() {
			ANDAS = ""
		}
		JISstat, err := os.Stat(JISFP)
		if err != nil || !JISstat.Mode().IsRegular() {
			JIS = ""
		}

		if strokeModule.PreferJIS && JIS != "" {
			currentstroke.Path = JIS
			currentstroke.Type = "JIS"
		} else if JIS != "" && ANDAS == "" {
			currentstroke.Path = JIS
			currentstroke.Type = "JIS"
		} else if ANDAS != "" {
			currentstroke.Path = ANDAS
			currentstroke.Type = "ANDAS"
		} else {
			continue
		}

        currentstroke.ID = id

		output = append(output, currentstroke)

	}

	err = CopyOutput(output, strokeModule.Path, savepath)

	if err != nil {
		return err
	}

    keymap := KeymapFromEntry(output)
    card.AddToFields("Stroke", keymap["stroke"])
    
    for i := 1; i < 10; i++ {

        key := fmt.Sprintf("stroke_%d", i)
        value, exists := keymap[key]
        if exists {
            card.AddToFields("Strokes", value)
        }
    }

	return nil
}
func CopyOutput(output []StrokeOutput, inpath string, outpath string) (out error) {
	for _, file := range output {
		fullinpath := inpath + "/" + file.Path
		fulloutpath := outpath + "/" + file.Path
		inpathstat, err := os.Stat(fullinpath)

		if err != nil {
			return err
		}
		if !inpathstat.Mode().IsRegular() {
			return errors.New("fullpath is not a file")
		}
		source, err := os.Open(fullinpath)
		if err != nil {
			return err
		}
		defer source.Close()

		destination, err := os.Create(fulloutpath)
		if err != nil {
			return err
		}
		defer destination.Close()

		_, err = io.Copy(destination, source)

		if err != nil {
			return err
		}
	}
	return out
}
func KeymapFromEntry(stroke_output []StrokeOutput) (output map[string]string) {
	output = make(map[string]string)

	value := "<div class = 'stroke_set'>"
	for _, stroke := range stroke_output {
        instance := "<div class = 'stroke " + stroke.Type + "'>" + "<img src='" + stroke.Path + "'>" + "</div>"
		value += instance

        key := fmt.Sprintf("stroke_%d", stroke.ID)
        output[key] = instance
    }
	value += "</div>"
	output["stroke"] = value
	return output
}
