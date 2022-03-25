package strokes

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	module "github.com/lesserfish/GoAme/Modules"
	kanjidic "github.com/lesserfish/GoAme/Modules/Kanjidic"
)

type InitOptions struct {
	StrokePath string
	Kanjimod   *kanjidic.Kanjidic_Module
	PreferJIS  bool
}

type StrokeModule struct {
	Path      string
	kanjimod  *kanjidic.Kanjidic_Module
	PreferJIS bool
}

type StrokeOutput struct {
	IdealPath string
	JISPath   string
	ALSASPath string
}

func Initialize(options InitOptions) (out module.Module, err error) {
	newModule := new(StrokeModule)
	newModule.Path = options.StrokePath
	newModule.kanjimod = options.Kanjimod
	newModule.PreferJIS = options.PreferJIS

	_, err = os.Stat(newModule.Path)

	if os.IsNotExist(err) {
		return newModule, err
	}

	return newModule, nil
}

func (strokeModule StrokeModule) Close() {
}
func (strokeModule StrokeModule) Demo() {

}
func (strokeModule StrokeModule) Render(input module.Input, card *module.Card) (err error) {
	literals := input["literal"]
	savepath := input["savepath"]

	output := []StrokeOutput{}

	if savepath == "" {
		return errors.New("Unspecified output file path!")
	}

	characters, err := kanjidic.FindEntry(&strokeModule.kanjimod.Dictionary, literals)

	if err != nil {
		return err
	}

	for _, character := range characters {
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

		if strokeModule.PreferJIS && JIS != "" {
			currentstroke.IdealPath = JIS
		} else if JIS != "" && ANDAS == "" {
			currentstroke.IdealPath = JIS
		} else {
			currentstroke.IdealPath = ANDAS
		}

		output = append(output, currentstroke)
	}

	err = CopyOutput(output, strokeModule.Path, savepath)

	if err != nil {
		return err
	}
	card.Parse(KeymapFromEntry(output), false)
	return nil
}
func (strokeModule StrokeModule) CSS(card *module.Card) {
}
func CopyOutput(output []StrokeOutput, inpath string, outpath string) (out error) {
	for _, file := range output {
		fullinpath := inpath + "/" + file.IdealPath
		fulloutpath := outpath + "/" + file.IdealPath
		inpathstat, err := os.Stat(fullinpath)

		if err != nil {
			return err
		}
		if !inpathstat.Mode().IsRegular() {
			return errors.New("Fullpath is not a file")
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
func KeymapFromEntry(output []StrokeOutput) (out map[string]string) {
	out = make(map[string]string)

	value := "<div class = 'stroke_set'>"
	for _, out := range output {
		value += "<div class = 'stroke'>" + "<img src='" + out.IdealPath + "'>" + "</div>"
	}
	value += "</div>"
	out["Stroke"] = value
	return out
}
