package jmdict

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
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

type RegexInstance struct {
	Operations []struct {
		Find    string `xml:"find"`
		Replace string `xml:"replace"`
	} `xml:"operation"`
}
type RegexOrder struct {
	Text   RegexInstance `xml:",chardata"`
	EntSeq RegexInstance `xml:"ent_seq"`
	REle   struct {
		Text      RegexInstance `xml:",chardata"`
		Reb       RegexInstance `xml:"reb"`
		ReRestr   RegexInstance `xml:"re_restr"`
		RePri     RegexInstance `xml:"re_pri"`
		ReNokanji RegexInstance `xml:"re_nokanji"`
		ReInf     RegexInstance `xml:"re_inf"`
	} `xml:"r_ele"`
	Sense struct {
		Text  RegexInstance `xml:",chardata"`
		Pos   RegexInstance `xml:"pos"`
		Gloss struct {
			Text RegexInstance `xml:",chardata"`
			Lang RegexInstance `xml:"lang"`
		} `xml:"gloss"`
		Example struct {
			Text   RegexInstance `xml:",chardata"`
			ExSrce struct {
				Text      RegexInstance `xml:",chardata"`
				ExsrcType RegexInstance `xml:"exsrc_type"`
			} `xml:"ex_srce"`
			ExText RegexInstance `xml:"ex_text"`
			ExSent struct {
				Text RegexInstance `xml:",chardata"`
				Lang RegexInstance `xml:"lang"`
			} `xml:"ex_sent"`
		} `xml:"example"`
		Xref  RegexInstance `xml:"xref"`
		Ant   RegexInstance `xml:"ant"`
		Misc  RegexInstance `xml:"misc"`
		Dial  RegexInstance `xml:"dial"`
		Stagr RegexInstance `xml:"stagr"`
	} `xml:"sense"`
	KEle struct {
		Text  RegexInstance `xml:",chardata"`
		Keb   RegexInstance `xml:"keb"`
		KePri RegexInstance `xml:"ke_pri"`
		KeInf RegexInstance `xml:"ke_inf"`
	} `xml:"k_ele"`
	Info struct {
		Text  RegexInstance `xml:",chardata"`
		Audit struct {
			Text    RegexInstance `xml:",chardata"`
			UpdDate RegexInstance `xml:"upd_date"`
			UpdDetl RegexInstance `xml:"upd_detl"`
		} `xml:"audit"`
	} `xml:"info"`
}

func LoadDictionary(parser *Parser) (err error) {
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

	return nil
}
func LoadFormatter(parser *Parser) (err error) {
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
func CleanEntry(entry *Entry, order *RegexOrder) (out error) {
	return out
}

func KeymapFromEntry(entry *Entry) (out map[string]string, err error) {
	return out, err
}
