from pickle import FALSE, TRUE
from pkgutil import iter_importers
from gtts import gTTS
from pydub import AudioSegment
import string
import xml.etree.ElementTree as ET
import requests
import hashlib
import argparse
import os
import shutil

#GLOBAL VARIABLES:

OUT_DELAY = 2000
IN_DELAY = 1000
REPEAT = 2
NUM_CHOICES = 3

parser = argparse.ArgumentParser()
parser.add_argument("--downloaddir", nargs='?', const="", default = "/tmp")
parser.add_argument("--outputdir", nargs='?', const="", default = "/tmp")
parser.add_argument("--dictpath", nargs='?', const="", default = "")
parser.add_argument("--audiouri", nargs='?', const="", default = "http://localhost:8000?kana={0}&kanji={1}")
parser.add_argument("--audioerr", nargs='?', const="", default = "./configuration/404.mp3")
args = parser.parse_args()

AUDIO_ERRCTNT = open(args.audioerr, "rb")
AUDIO_ERRMD5 = hashlib.md5(AUDIO_ERRCTNT.read()).digest()
AUDIO_ERRCTNT.close()

def GetMeaning(root , kanji : string, kana : string):

    IGNORE_KANJI = False
    IGNORE_KANA = False

    if len(kanji) == 0:
        IGNORE_KANJI = True
    if len(kana) == 0:
        IGNORE_KANA = True
    
    dentries = root.findall("entry")

    meanings = []
    for dentry in dentries:
        kele = dentry.find("k_ele")

        if kele is None:
            continue
        keb = kele.find("keb")

        if keb is None:
            continue
        
        rele = dentry.find("r_ele")

        if rele is None:
            continue
        
        reb = rele.find("reb")

        if reb is None:
            continue

        dkanji = keb.text.split()[0]
        dkana = reb.text.split()[0]

        if (dkanji == kanji or IGNORE_KANJI) and (dkana == kana or IGNORE_KANA):

            sense = dentry.find("sense")
            if(sense is None):
                continue

            gloss = sense.findall("gloss")
            for meaning in gloss:
                meanings.append(meaning.text)
            return meanings

    return []    

def DownloadPronunciation(kanji : string, kana : string, id : string):
    fullpath = args.audiouri.format(kana, kanji)
    audio = requests.get(fullpath)
    
    audiomd5 = hashlib.md5(audio.content).digest()
    
    path = args.downloaddir + "/" + id + "/tmp.mp3"

    if(AUDIO_ERRMD5 == audiomd5 or audio.status_code != 200):
        print("Audio is not available for the word: " + kanji + " : " + kana)
        tts = gTTS(kanji, lang='ja')
        tts.save(path)
    else:
        fp = open(path, "wb")
        fp.write(audio.content)
        fp.close()
    
    audio = AudioSegment.from_mp3(path)
    os.remove(path)
    return audio

def GenerateTTS(meanings, id):
    minimum = min(NUM_CHOICES, len(meanings))
    choices = range(0, minimum)
    
    sentence = ""

    faudio = AudioSegment.empty()
    silence = AudioSegment.silent(duration=IN_DELAY)
    i = 1
    for choice in choices:
        path = args.downloaddir + "/" + id + "/tmp.mp3"
        sentence = str(i) + ". " + meanings[choice]
        tts = gTTS(sentence, lang='en')
        tts.save(path)
        segment = AudioSegment.from_mp3(path)
        os.remove(path)
        faudio = faudio + segment + silence
        i = i + 1

    return faudio

def MergeAudio(Aaudio, Baudio):
    silence = AudioSegment.silent(duration=OUT_DELAY)
    faudio = Aaudio + silence + Baudio
    faudio = faudio * REPEAT
    return faudio



def HandleTask(kanji, kana, id):
    meanings = GetMeaning(root, kanji, kana)

    if len(meanings) == 0:
        return
    
    A = DownloadPronunciation(kanji, kana, id)
    B = GenerateTTS(meanings, id)
    out = MergeAudio(A, B)

    out.export(args.downloaddir + "/" + id + "/" + kanji + ".mp3", format="mp3")


tree = ET.parse(args.dictpath)
root = tree.getroot()

id = "1"
dpath = args.downloaddir + "/" + id
opath = args.outputdir + "/" + "out_" + id

os.mkdir(dpath)
HandleTask("応援", "おうえん", id)
shutil.make_archive(dpath, 'zip', opath)
shutil.rmtree(dpath)