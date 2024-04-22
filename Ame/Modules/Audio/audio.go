package audio

import (
	"errors"
    "fmt"
	"io"
	"log"
	"os"

	module "github.com/lesserfish/GoAme/Ame/Modules"
	jmdict "github.com/lesserfish/GoAme/Ame/Modules/JMDict"
)

type InitOptions struct {
	AudioPath   string
	JMdictMod *jmdict.JMdictModule
}

type AudioModule struct {
	AudioPath        string
	JMdictMod  *jmdict.JMdictModule
}

func Initialize(options InitOptions) (*AudioModule, error) {
	newModule := new(AudioModule)
	newModule.AudioPath = options.AudioPath
	newModule.JMdictMod = options.JMdictMod

	log.Println("Audio Module initialized!")
	return newModule, nil
}

func (audioModule AudioModule) Render(input module.Input, card *module.Card) (err error) {

	kanji := input["kanjiword"]
	kana := input["kanaword"]
	path := input["savepath"]

	if kana == "" && kanji == "" {
		return nil
	}
	if kana == "" {
		kana, err = GetKana(kanji, &audioModule.JMdictMod.Dictionary)
		if err != nil {
			return err
		}
	}

	filename := GetFilename(kana, kanji, path)
    filepath := fmt.Sprintf("%s/%s", audioModule.AudioPath, filename)

    _, err = os.Stat(filepath);

    if err != nil {
        return nil
    }

    output_path := fmt.Sprintf("%s/%s", path, filename)

    err = CopyFile(filepath, output_path)

    if err != nil {
        return errors.New("Internal audio error")
    }

    keymap := KeymapFromEntry(kanji, kana)
    card.AddToFields("Audio", keymap["audio"])
     
	return nil
}
func GetFilename(kana string, kanji string, path string) (out string) {
    if kanji == "" {
        return fmt.Sprintf("audio_%s.mp3", kana)
    }
	return fmt.Sprintf("audio_%s_%s.mp3", kana, kanji)
}
func GetKana(kanji string, dict *jmdict.JMdict) (string, error) {

	entry, err := jmdict.FindEntry(dict, kanji, "", false, true)

	if err != nil {
		return "", err
	}

	if len(entry.REle) == 0 {
		return "", errors.New("no kana expression")
	}
	rele := entry.REle[0]

	if len(rele.Reb) == 0 {
		return "", errors.New("no kana expression")
	}
	kana := rele.Reb[0]
	return kana, nil
}

func CopyFile(source string, target string) error {
	// Open the source file for reading
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents of the source file to the destination file
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}
	return nil
}
func KeymapFromEntry(kana string, kanji string) (out map[string]string) {
	out = make(map[string]string)

	filename := GetFilename(kanji, kana, "")

	value := "<div class = 'audio'>"
	value += "[sound:" + filename + "]"
	value += "</div>"
	out["audio"] = value
	return out
}
