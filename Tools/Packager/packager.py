import genanki
import csv
import os
import sys
import argparse


MODEL_ID = 1443365057
DECK_ID = 2093107000
FIELDS = [
        {'name': 'Kana'},
        {'name': 'Kanji'},
        {'name': 'Audio'},
        {'name': 'Sense'},
        {'name': 'Kanjiinfo'},
        {'name': 'Kanjiinfo_1'},
        {'name': 'Kanjiinfo_2'},
        {'name': 'Kanjiinfo_3'},
        {'name': 'Kanjiinfo_4'},
        {'name': 'Kanjiinfo_5'},
        {'name': 'Kanjiinfo_6'},
        {'name': 'Kanjiinfo_7'},
        {'name': 'Kanjiinfo_8'},
        {'name': 'Kanjiinfo_9'},
        {'name': 'Kanjiinfo_10'},
        {'name': 'Kanjiinfoex'},
        {'name': 'Kanjiinfoex_1'},
        {'name': 'Kanjiinfoex_2'},
        {'name': 'Kanjiinfoex_3'},
        {'name': 'Kanjiinfoex_4'},
        {'name': 'Kanjiinfoex_5'},
        {'name': 'Kanjiinfoex_6'},
        {'name': 'Kanjiinfoex_7'},
        {'name': 'Kanjiinfoex_8'},
        {'name': 'Kanjiinfoex_9'},
        {'name': 'Kanjiinfoex_10'},
        {'name': 'Stroke'},
        {'name': 'Stroke_1'},
        {'name': 'Stroke_2'},
        {'name': 'Stroke_3'},
        {'name': 'Stroke_4'},
        {'name': 'Stroke_5'},
        {'name': 'Stroke_6'},
        {'name': 'Stroke_7'},
        {'name': 'Stroke_8'},
        {'name': 'Stroke_9'},
        {'name': 'Stroke_10'},
        {'name': 'Literal'},
        {'name': 'Literal_1'},
        {'name': 'Literal_2'},
        {'name': 'Literal_3'},
        {'name': 'Literal_4'},
        {'name': 'Literal_5'},
        {'name': 'Literal_6'},
        {'name': 'Literal_7'},
        {'name': 'Literal_8'},
        {'name': 'Literal_9'},
        {'name': 'Literal_10'},
        {'name': 'Example_1'},
        {'name': 'Example_2'},
        {'name': 'Example_3'},
        {'name': 'Example_4'},
        {'name': 'Example_5'},
        {'name': 'Example_6'},
        {'name': 'Example_7'},
        {'name': 'Example_8'},
        {'name': 'Example_9'},
        {'name': 'Example_10'},
]


TEMPLATES=[
    {
      'name': 'Japanese -> English',
      'qfmt': '{{Kanji}}',
      'afmt': '{{FrontSide}}\n\n<hr id=answer>\n\n{{Kana}}\n{{Sense}}\n{{Audio}}\n\n{{Example_1}}\n{{Example_2}}\n{{Example_3}}'
    },
    {
      'name': 'English -> Japanese',
      'qfmt': '{{Sense}}',
      'afmt': '{{FrontSide}}\n\n<hr id=answer>\n\n{{Kanji}}\n{{Kana}}\n{{Audio}}\n\n{{Example_1}}\n{{Example_2}}\n{{Example_3}}'
    },
    {
      'name': 'Kanji -> Meaning',
      'qfmt': '{{Literal}}',
      'afmt': '{{FrontSide}}\n\n<hr id=answer>\n\n{{Kanjiinfo}}\n{{Stroke}}'
    },
]
CSS = ".card {\nfont-family: arial;\nfont-size: 20px;\ntext-align: center;\ncolor: black;\nbackground-color: white;\n}\n\n.rexample {\ntext-align: left;\nmargin-left: 10%;\n}\n\n.rexample > .JP {\nfont-size: 100%;\ncolor: #1d1d29;\nmargin-top: 1em;\n}\n\n.rexample > .ENG {\nfont-size: 75%;\ncolor: #303030;\nmargin-bottom: 1em;\n}\n\n.Kele {\ntext-align: center;\nlist-style-position: inside;\n}\n\n.Kele > ol > li {\ncolor: #333;\nfont-size: 100%;\n}\n\n.Kele > ol > li:nth-child(1) {\ncolor: #000;\nfont-size: 150%;\n}\n\n.Kele > ol > li:nth-child(2) {\ncolor: #131313;\nfont-size: 125%;\n}\n\n.Keb {\ndisplay: inline-block;\ntext-align: left;\n}\n\n.Rele {\ntext-align: center;\nlist-style-position: inside;\n}\n\n.Rele > ol > li {\ncolor: #333;\nfont-size: 100%;\n}\n\n.Rele > ol > li:nth-child(1) {\ncolor: #000;\nfont-size: 150%;\n}\n\n.Rele > ol > li:nth-child(2) {\ncolor: #131313;\nfont-size: 125%;\n}\n\n.Reb {\ndisplay: inline-block;\ntext-align: left;\n}\n\n.Sense > ol > li {\ntext-align: left;\nmargin-left: 10%;\n}\n\n.pos > ul {\nlist-style-type: disc;\n}\n\n.pos > ul > li {\ncolor: gray;\nfont-size: 75%;\ndisplay: inline-flex;\nmargin-left: 1em;\nmargin-right: 1em;\n}\n\n.gloss > ol > li {\nmargin-top: 1em;\nmargin-bottom: 1em;\n}\n\n.gloss > ol > li:nth-child(1) {\nfont-size: 115%;\n}\n\n.example > .lang_jpn {\nfont-size: 100%;\ncolor: #1d1d29;\nmargin-top: 1em;\n}\n\n.example > .lang_eng {\nfont-size: 75%;\ncolor: #303030;\nmargin-bottom: 1em;\n}\n\n.example {\nmargin-top: 1em;\nmargin-bottom: 1em;\n}\n\n.kanji_info {\ntext-align: center;\nlist-style-position: inside;\n}\n\n.kanji_info > ol > li {\ncolor: #333;\nfont-size: 100%;\n}\n\n.kanji_info > .kanji_instance > .literal {\ncolor: #000;\nfont-size: 150%;\ntext-align: center;\n}\n\n.kanji_info > .kanji_instance > .meanings {\ntext-align: center;\nmargin-right: 3%;\n}\n\n.kanji_info > .kanji_instance > .meanings > ol > li {\nfont-size: 100%;\n}\n\n.kanji_info > .kanji_instance > .meanings > ol > li:nth-child(1) {\nfont-size: 125%;\n}\n\n.kanji_info > .kanji_instance > .readings > ul > li {\ndisplay: inline;\nmargin-left: 4%;\nmargin-right: 4%;\n}\n\n.kanji_info > .kanji_instance > .misc > * {\nmargin-top: 3%;\nmargin-bottom: 3%;\n}\n\n.kliteral {\ntext-align: center;\ncolor: #000;\nfont-size: 150%;\n}\n"

model = genanki.Model(
        MODEL_ID,
        'Ame Model',
        fields=FIELDS,
        templates=TEMPLATES,
        css=CSS)


def SafeCheck(PATH):
    if not os.path.isdir(PATH): sys.exit(2)
    if not os.path.isdir(PATH + "/Media"): sys.exit(3)
    if not os.path.isfile(PATH + "/anki_deck.txt"): sys.exit(4)

def ListDirectory(dir):
    return [file for file in os.listdir(dir) if os.path.isfile(os.path.join(dir, file))]

def Package(PATH):
    SafeCheck(PATH)
    deck = genanki.Deck(
        DECK_ID,
        "AmeKanji")

    # Add cards
    with open(PATH + "/anki_deck.txt", "r") as file:
        csvfile = csv.reader(file, delimiter=';')
        for row in csvfile:
            cfields = row[:-1]
            ctags = row[-1:]
            note = genanki.Note(
                    model=model,
                    fields=cfields,
                    tags=ctags)
            deck.add_note(note)
    # Add media
    package = genanki.Package(deck)
    files = [PATH + "/Media/" + filename for filename in ListDirectory(PATH + "/Media")]
    package.media_files = files
    package.write_to_file(PATH + '/anki_deck.apkg')

def main():
    # Create ArgumentParser object
    parser = argparse.ArgumentParser(description="Combines the output of GoAme to a .apkg file")

    # Add arguments
    parser.add_argument("-p", type=str, help="Path to the directory where the files are contained")
    # Add more arguments as needed

    # Parse the command-line arguments
    args = parser.parse_args()
    PATH = args.p

    if PATH is None:
        sys.exit(1)

    Package(PATH)
    sys.exit(0)

if __name__ == "__main__":
    main()
