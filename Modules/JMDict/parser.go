package jmdict

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type JMdict struct {
	XMLName xml.Name `xml:"JMdict"`
	Text    string   `xml:",chardata"`
	Entry   []struct {
		Text   string `xml:",chardata"`
		EntSeq string `xml:"ent_seq"`
		REle   []struct {
			Text      string   `xml:",chardata"`
			Reb       string   `xml:"reb"`
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
	} `xml:"entry"`
}

func LoadDictionary(parser *Parser) error {
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
