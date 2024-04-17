package kanjidic

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"os"
    "fmt"
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

	err = xml.Unmarshal(byteData, &parser.Dictionary)

	if err != nil {
		return err
	}

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
    kanji_count := 1

	for _, char := range *characters {

		instance_info := "<div class = 'kanji_instance'>"

        literal_info := "<div class = 'literal'>"
		literal_info += char.Literal
		literal_info += "</div>"

        instance_info += literal_info

        key := fmt.Sprintf("literal_%d", kanji_count)
        out[key] = literal_info

		instance_info += "<div class = meanings>"
		instance_info += "<ol>"
		for _, meaning := range char.ReadingMeaning.Rmgroup.Meaning {

			if meaning.MLang != "" { // Skip non-english meanings
				continue
			}
			text := "<li>" + meaning.Text + "</li>"
			instance_info += text
		}
		instance_info += "</ol>"
		instance_info += "</div>"

		instance_info += "<div class = readings>"
		instance_info += "<ul>"
		for _, reading := range char.ReadingMeaning.Rmgroup.Reading {
			if reading.RType != "ja_on" && reading.RType != "ja_kun" {
				continue
			}
			text := "<li>" + reading.Text + "</li>"
			instance_info += text
		}
		instance_info += "</ul>"
		instance_info += "</div>"

		kanji_info_basic += instance_info + "</div>"
        key = fmt.Sprintf("kanjiinfo_%d", kanji_count)
        out[key] = instance_info

		instance_info += "<div class='misc'>"
		instance_info += "<div class='grade'> Grade: " + char.Misc.Grade + "</div>"
		instance_info += "<div class='strokecount'> Stroke count: " + char.Misc.StrokeCount + "</div>"
		instance_info += "<div class='jlpt'> JLPT: " + char.Misc.Jlpt + "</div>"
		instance_info += "<div class='freq'> Frequency: " + char.Misc.Freq + "</div>"
		instance_info += "</div>"

		instance_info += "</div>"

		kanji_info_ext += instance_info
        key = fmt.Sprintf("kanjiinfoex_%d", kanji_count)
        out[key] = instance_info

        kanji_count += 1
	}
	kanji_info_basic += "</div>"
	kanji_info_ext += "</div>"

	out["kanjiinfo"] = kanji_info_basic
	out["kanjiinfoex"] = kanji_info_ext

	literal := ""
	for _, char := range *characters {
		literal += char.Literal
	}

	out["literal"] = "<div class = 'kliteral'>" + literal + "</div>"

	return out, err
}
