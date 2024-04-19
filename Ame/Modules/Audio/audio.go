package audio

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	module "github.com/lesserfish/GoAme/Ame/Modules"
	jmdict "github.com/lesserfish/GoAme/Ame/Modules/JMDict"
)

type InitOptions struct {
	URI       string
	JMdictMod *jmdict.JMdictModule
}

type AudioModule struct {
	URI        string
	JMdictMod  *jmdict.JMdictModule
}

func Initialize(options InitOptions) (*AudioModule, error) {
	newModule := new(AudioModule)
	newModule.URI = options.URI
	newModule.JMdictMod = options.JMdictMod

	log.Println("Audio Module initialized!")
	return newModule, nil
}

func (audioModule AudioModule) Close() {
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

	params := url.Values{}
	params.Add("kana", kana)
	params.Add("kanji", kanji)
	URI := audioModule.URI + params.Encode()

	filepath := GetFilename(kana, kanji, path)

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	req, err := http.NewRequest("GET", URI, nil)

	if err != nil {
		return err
	}

	SetHeaders(&req.Header)

	resp, err := client.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("could not find audio file")
	}

	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

    keymap := KeymapFromEntry(kanji, kana)
    card.AddToFields("Audio", keymap["audio"])

	return nil

}
func SetHeaders(header *http.Header) {
	header.Set("charset", "utf-8")
}
func GetFilename(kana string, kanji string, path string) (out string) {
	out = path + "/" + kana + ".mp3"
	return out
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
func KeymapFromEntry(kana string, kanji string) (out map[string]string) {
	out = make(map[string]string)

	filename := GetFilename(kanji, kana, "")
	split := strings.Split(filename, "/")
	filename = split[1]

	value := "<div class = 'audio'>"
	value += "[sound:" + filename + "]"
	value += "</div>"
	out["audio"] = value
	return out
}
