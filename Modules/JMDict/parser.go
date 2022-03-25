package jmdict

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

type JMdict struct {
	XMLName xml.Name `xml:"JMdict"`
	Text    string   `xml:",chardata"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Text   string `xml:",chardata"`
	EntSeq string `xml:"ent_seq"`
	REle   []struct {
		Text      string   `xml:",chardata"`
		Reb       []string `xml:"reb"`
		ReRestr   string   `xml:"re_restr"`
		RePri     []string `xml:"re_pri"`
		ReNokanji string   `xml:"re_nokanji"`
		ReInf     string   `xml:"re_inf"`
	} `xml:"r_ele"`
	Sense []struct {
		Text  string   `xml:",chardata"`
		Pos   []string `xml:"pos"`
		Gloss []struct {
			Text string `xml:",chardata"`
			Lang string `xml:"lang,attr"`
		} `xml:"gloss"`
		Example []struct {
			Text   string `xml:",chardata"`
			ExSrce struct {
				Text      string `xml:",chardata"`
				ExsrcType string `xml:"exsrc_type,attr"`
			} `xml:"ex_srce"`
			ExText string `xml:"ex_text"`
			ExSent []struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"ex_sent"`
		} `xml:"example"`
		Xref  []string `xml:"xref"`
		Ant   string   `xml:"ant"`
		Misc  string   `xml:"misc"`
		Dial  string   `xml:"dial"`
		Stagr []string `xml:"stagr"`
	} `xml:"sense"`
	KEle []struct {
		Text  string   `xml:",chardata"`
		Keb   []string `xml:"keb"`
		KePri []string `xml:"ke_pri"`
		KeInf string   `xml:"ke_inf"`
	} `xml:"k_ele"`
	Info struct {
		Text  string `xml:",chardata"`
		Audit []struct {
			Text    string `xml:",chardata"`
			UpdDate string `xml:"upd_date"`
			UpdDetl string `xml:"upd_detl"`
		} `xml:"audit"`
	} `xml:"info"`
}
type Operation struct {
	Find    string `xml:"find"`
	Replace string `xml:"replace"`
}
type RegexFormatter struct {
	XMLName xml.Name    `xml:"RegexFormatter"`
	Pos     []Operation `xml:"pos"`
}

func LoadDictionary(parser *JMdictModule) (err error) {
	DictionaryPath := parser.DictionaryPath

	xmlFile, err := os.Open(DictionaryPath)

	if err != nil {
		return err
	}

	byteData, _ := ioutil.ReadAll(xmlFile)

	err = xml.Unmarshal(byteData, &parser.dictionary)

	if err != nil {
		return err
	}

	fmt.Println("Dictionary loaded!")
	return nil
}
func LoadFormatter(parser *JMdictModule) (err error) {
	FormatterPath := parser.FormatterPath
	xmlFile, err := os.Open(FormatterPath)

	if err != nil {
		return err
	}

	byteData, _ := ioutil.ReadAll(xmlFile)

	err = xml.Unmarshal(byteData, &parser.formatter)

	if err != nil {
		return err
	}
	fmt.Println("Formatter loaded!")
	return nil
}
func FindEntry(dict *JMdict, kanji string, kana string) (out Entry, err error) {

	err = errors.New("Failed to find entry in dictionary.")

	search_kana := false
	if kana != "" {
		search_kana = true
	}

entry_search:
	for _, entry := range dict.Entries {
		match_kanji := false
		match_kana := false

	kanji_search:
		for _, kele := range entry.KEle {
			for _, keb := range kele.Keb {
				if keb == kanji {
					match_kanji = true
					break kanji_search
				}
			}
		}
		if search_kana {
		kana_search:
			for _, rele := range entry.REle {
				for _, reb := range rele.Reb {
					if reb == kana {
						match_kana = true
						break kana_search
					}
				}
			}
		} else {
			match_kana = true
		}

		if match_kanji && match_kana {
			out = entry
			err = nil
			break entry_search
		}
	}

	return out, err
}
func CleanEntry(entry *Entry, order *RegexFormatter) (out error) {
	for senseid, sense := range entry.Sense {
		for posid, pos := range sense.Pos {
			out_string := pos
			for _, instruction := range order.Pos {
				regex := regexp.MustCompile(instruction.Find)
				newstring := regex.ReplaceAllString(out_string, instruction.Replace)
				out_string = newstring
			}
			entry.Sense[senseid].Pos[posid] = out_string
		}
	}

	return out
}

func KeymapFromEntry(entry *Entry) (out map[string]string, err error) {
	out = make(map[string]string)

	Kanji := ""
	Kanji += "<div class='Kele'><ol>"
	for _, kele := range entry.KEle {
		Kanji += "<li>"
		for _, keb := range kele.Keb {
			Kanji += "<div class='Keb'>"
			Kanji += keb
			Kanji += "</div>"
		}
		Kanji += "</li>"
	}
	Kanji += "</ol></div>"
	Kana := ""
	Kana += "<div class='Rele'><ol>"
	for _, rele := range entry.REle {
		Kana += "<li>"
		for _, reb := range rele.Reb {
			Kana += "<div class='Reb'>"
			Kana += reb
			Kana += "</div>"
		}
		Kana += "</li>"
	}
	Kana += "</div>"
	Sense := ""
	Sense += "<div class='Sense'><ol>"
	for _, sense := range entry.Sense {
		Sense += "<li>"
		Sense += "<div class='pos'><ul>"
		for _, pos := range sense.Pos {
			Sense += "<li>"
			Sense += pos
			Sense += "</li>"
		}
		Sense += "</ul></div>"
		Sense += "<div class = 'gloss'><ol>"
		for _, gloss := range sense.Gloss {
			Sense += "<li>"
			Sense += gloss.Text
			Sense += "</li>"
		}
		Sense += "</ol></div>"
		Sense += "<div class='example'>"
		for _, example := range sense.Example {
			for _, lang := range example.ExSent {
				Sense += "<div class='lang " + lang.Lang + "'>"
				Sense += lang.Text
				Sense += "</div>"
			}
		}
		Sense += "</div>"
		Sense += "</li>"
	}
	Sense += "</div>"

	out["KanaWord"] = Kana
	out["KanjiWord"] = Kanji
	out["Sense"] = Sense

	return out, err
}
