package audio

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"

	module "github.com/lesserfish/GoAme/Modules"
	jmdict "github.com/lesserfish/GoAme/Modules/JMDict"
)

type InitOptions struct {
	URI       string
	JMdictMod *jmdict.JMdictModule
}

type AudioModule struct {
	URI       string
	JMdictMod *jmdict.JMdictModule
}

func Initialize(options InitOptions) (*AudioModule, error) {
	newModule := new(AudioModule)
	newModule.URI = options.URI
	newModule.JMdictMod = options.JMdictMod

	return newModule, nil
}

func (audioModule AudioModule) Close() {
}
func (audioModule AudioModule) Demo() {

}
func (audioModule AudioModule) Render(input module.Input, card *module.Card) (err error) {

	kanji := input["kanjiword"]
	kana := input["kanaword"]
	path := input["savepath"]

	if kana == "" && kanji == "" {
		return errors.New("No input given to JMdic module!")
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

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
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
		return errors.New("Error. Status code: " + strconv.Itoa(resp.StatusCode))
	}

	defer resp.Body.Close()
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		return err
	}

	card.Parse(KeymapFromEntry(kanji, kana), false)

	return nil

}
func (audioModule AudioModule) CSS(card *module.Card) {
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
		return "", errors.New("No kana expression.")
	}
	rele := entry.REle[0]

	if len(rele.Reb) == 0 {
		return "", errors.New("No kana expression.")
	}
	kana := rele.Reb[0]
	return kana, nil
}
func KeymapFromEntry(kana string, kanji string) (out map[string]string) {
	out = make(map[string]string)

	filename := GetFilename(kana, kanji, "")

	value := "<div class = 'audio'>"
	value += "[sound:{{" + filename + "}}]"
	value += "</div>"
	out["Audio"] = value
	return out
}
