import requests
import os
import sys
import json
import hashlib

APIHelper = "https://www.amekanji.com/api/help"
Rank = 0
FORCE_DOWNLOAD = False

def calculate_hash(file_contents):
    sha256_hash = hashlib.sha256()
    sha256_hash.update(file_contents)
    return sha256_hash.hexdigest()

ERRFile = "./err.mp3"
with open(ERRFile, "rb") as file:
    ERRHash = calculate_hash(file.read())

def char_is_hiragana(c) -> bool:
    return u'\u3040' <= c <= u'\u309F'
def char_is_katakana(c) -> bool:
    return u'\u30A0' <= c <= u'\u30FF'

def IsKanaChar(c):
    return (char_is_katakana(c) or char_is_hiragana(c))

def IsKana(string):
    for char in string:
        if not IsKanaChar(char):
            return False
    return True

def SplitMixed(words):
    KanaWords = []
    MixedWords = []

    for word in words:
        if IsKanaChar(word):
            KanaWords.append(word)
        else:
            MixedWords.append(word)
    return (KanaWords, MixedWords)

def ParseDownload(response):
    content = response["Response"]
    content_group = [(key, content[key]) for key in list(content.keys())]
    output = [(key, reading[Rank]) for (key, reading) in content_group if len(reading) > Rank]
    return output

def DownloadReadings(mixed):
    response = requests.post(APIHelper, data=json.dumps(mixed))
    return (response.status_code, json.loads(response.content.decode('utf-8')))

def LoadWordlist(id):
    wordlist = "wordlist_{:02d}".format(id)
    print(wordlist)
    with open(wordlist, "r") as file:
        content = [line.strip() for line in file.readlines()]
    return content

def ConstructURI(Kana, Kanji):
    if Kanji is None:
        URI = "http://assets.languagepod101.com/dictionary/japanese/audiomp3.php?kana={}".format(Kana)
    else:
        URI = "http://assets.languagepod101.com/dictionary/japanese/audiomp3.php?kana={}&kanji={}".format(Kana, Kanji)

    return URI

def ValidateAudio(content):
    audio_hash = calculate_hash(content)
    return audio_hash != ERRHash

def GetName(Kana, Kanji):
    if Kanji is None:
        Name = "audio_{}.mp3".format(Kana)
    else:
        Name = "audio_{}_{}.mp3".format(Kana, Kanji)
    return Name

def DownloadWord(Kana, Kanji, ID):
    filepath = "Audio/{:02d}/".format(ID) + GetName(Kana, Kanji)

    if os.path.exists(filepath) and not FORCE_DOWNLOAD:
        print(",", end="")
        sys.stdout.flush()
        return 0

    missing_filepath = "Audio/Missing/" + GetName(Kana, Kanji)
    if os.path.exists(missing_filepath) and not FORCE_DOWNLOAD:
        print(",", end="")
        sys.stdout.flush()
        return 0

    URI = ConstructURI(Kana, Kanji)
    headers = {
        'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3'
    }

    response = requests.get(URI, headers=headers, stream=True)

    if response.status_code != 200:
        with open(missing_filepath, "wb") as file:
            file.write(b"")
        return 0

    if not ValidateAudio(response.content):
        with open(missing_filepath, "wb") as file:
            file.write(b"")
        return 0

    with open(filepath, "wb") as file:
        file.write(response.content)
    print(".", end="") 
    sys.stdout.flush()
    
    return 1


def DownloadWordlist(kana_only, mixed_only, IDD):
    downloads = 0
    for kana in kana_only:
        downloads = downloads + DownloadWord(kana, None, IDD)
    for (kanji, kana) in mixed_only:
        downloads = downloads + DownloadWord(kana, kanji, IDD)
    print("")
    return downloads

def Main(IDD):
    wl = LoadWordlist(IDD)
    (k, m) = SplitMixed(wl)
    (status, response) = DownloadReadings(m)
    mixed_readings = ParseDownload(response)
    total_downloads = DownloadWordlist(k, mixed_readings, IDD)
    return total_downloads


for i in range(0, 86):
    print("Downloading file {}".format(i))
    downloads = Main(i)
    print("Downloaded {} files.".format(downloads))
