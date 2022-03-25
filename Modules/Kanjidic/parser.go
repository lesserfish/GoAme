package kanjidic

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Kanjidic struct {
	XMLName xml.Name `xml:"kanjidic2"`
	Text    string   `xml:",chardata"`
	Header  struct {
		Text            string `xml:",chardata"`
		FileVersion     string `xml:"file_version"`
		DatabaseVersion string `xml:"database_version"`
		DateOfCreation  string `xml:"date_of_creation"`
	} `xml:"header"`
	Characters []Character `xml:"character"`
}

type Character struct {
	XMLName   xml.Name `xml:"character"`
	Text      string   `xml:",chardata"`
	Literal   string   `xml:"literal"`
	Codepoint struct {
		Text    string `xml:",chardata"`
		CpValue []struct {
			Text   string `xml:",chardata"`
			CpType string `xml:"cp_type,attr"`
		} `xml:"cp_value"`
	} `xml:"codepoint"`
	Radical struct {
		Text     string `xml:",chardata"`
		RadValue []struct {
			Text    string `xml:",chardata"`
			RadType string `xml:"rad_type,attr"`
		} `xml:"rad_value"`
	} `xml:"radical"`
	Misc struct {
		Text        string `xml:",chardata"`
		Grade       string `xml:"grade"`
		StrokeCount string `xml:"stroke_count"`
		Variant     struct {
			Text    string `xml:",chardata"`
			VarType string `xml:"var_type,attr"`
		} `xml:"variant"`
		Freq string `xml:"freq"`
		Jlpt string `xml:"jlpt"`
	} `xml:"misc"`
	DicNumber struct {
		Text   string `xml:",chardata"`
		DicRef []struct {
			Text   string `xml:",chardata"`
			DrType string `xml:"dr_type,attr"`
			MVol   string `xml:"m_vol,attr"`
			MPage  string `xml:"m_page,attr"`
		} `xml:"dic_ref"`
	} `xml:"dic_number"`
	QueryCode struct {
		Text  string `xml:",chardata"`
		QCode []struct {
			Text   string `xml:",chardata"`
			QcType string `xml:"qc_type,attr"`
		} `xml:"q_code"`
	} `xml:"query_code"`
	ReadingMeaning struct {
		Text    string `xml:",chardata"`
		Rmgroup struct {
			Text    string `xml:",chardata"`
			Reading []struct {
				Text  string `xml:",chardata"`
				RType string `xml:"r_type,attr"`
			} `xml:"reading"`
			Meaning []struct {
				Text  string `xml:",chardata"`
				MLang string `xml:"m_lang,attr"`
			} `xml:"meaning"`
		} `xml:"rmgroup"`
		Nanori string `xml:"nanori"`
	} `xml:"reading_meaning"`
}

func LoadDictionary(parser *Kanjidic_Module) (err error) {
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
func FindEntry(dict *Kanjidic, literals string) (out []Character, err error) {
	for _, literal := range literals {
	char_search:
		for _, char := range dict.Characters {
			if string(literal) == char.Literal {
				out = append(out, char)
				break char_search
			}
		}
	}
	if len(out) == 0 {
		errorstr := "Failed to find any literal in " + literals + "."
		return out, errors.New(errorstr)
	}

	return out, nil
}
func KeymapFromEntry(characters *[]Character) (out map[string]string, err error) {
	out = make(map[string]string)

	kanji_info_ext := "<div class = 'kanji_info'>"
	kanji_info_basic := "<div class = 'kanji_info'>"
	for _, char := range *characters {
		kanji_info := "<div class = 'kanji_instance'>"

		kanji_info += "<div class = 'literal'>"
		kanji_info += char.Literal
		kanji_info += "</div>"

		kanji_info += "<div class = meanings>"
		kanji_info += "<ol>"
		for _, meaning := range char.ReadingMeaning.Rmgroup.Meaning {

			if meaning.MLang != "" { // Skip non-english meanings
				continue
			}
			text := "<li>" + meaning.Text + "</li>"
			kanji_info += text
		}
		kanji_info += "</ol>"
		kanji_info += "</div>"

		kanji_info += "<div class = readings>"
		kanji_info += "<ul>"
		for _, reading := range char.ReadingMeaning.Rmgroup.Reading {
			if reading.RType != "ja_on" && reading.RType != "ja_kun" {
				continue
			}
			text := "<li>" + reading.Text + "</li>"
			kanji_info += text
		}
		kanji_info += "</ul>"
		kanji_info += "</div>"

		kanji_info_basic += kanji_info + "</div>"

		kanji_info += "<div class='misc'>"
		kanji_info += "<div class='grade'> Grade: " + char.Misc.Grade + "</div>"
		kanji_info += "<div class='strokecount'> Stroke count: " + char.Misc.StrokeCount + "</div>"
		kanji_info += "<div class='jlpt'> JLPT: " + char.Misc.Jlpt + "</div>"
		kanji_info += "<div class='freq'> Frequency: " + char.Misc.Freq + "</div>"
		kanji_info += "</div>"

		kanji_info += "</div>"

		kanji_info_ext += kanji_info
	}
	kanji_info_basic += "</div>"
	kanji_info_ext += "</div>"

	out["kanjiinfo"] = kanji_info_basic
	out["kanjiinfoex"] = kanji_info_ext
	return out, err
}
