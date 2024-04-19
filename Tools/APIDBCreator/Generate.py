import xml.etree.cElementTree as ET
import sqlite3

dict_file = "../../Resources/Repository/Vocabulary/JMdict_e_examp.xml"
out_file = "../../Resources/Database/API.sqlite"

conn = sqlite3.connect(out_file)
print("Creating tables...")
try:
    conn.execute("CREATE TABLE kanjikana (kanji TEXT NOT NULL, kana TEXT NOT NULL);")
except sqlite3.OperationalError:
    conn.execute("DROP TABLE kanjikana")
    conn.execute("CREATE TABLE kanjikana (kanji TEXT NOT NULL, kana TEXT NOT NULL);")

try:
	conn.execute("CREATE TABLE clients (ip TEXT NOT NULL, date TEXT, reqsize INTEGER);")
except sqlite3.OperationalError:
    conn.execute("DROP TABLE clients")
    conn.execute("CREATE TABLE clients (ip TEXT NOT NULL, date TEXT, reqsize INTEGER);")

print("OK!")

tree = ET.parse(dict_file)
root = tree.getroot()

KKMAP = dict()

for entry in root:
    Kanji = []
    Kana = []
    keles = entry.findall("k_ele")
    for kele in keles:
        kebs = kele.findall("keb")
        for keb in kebs:
            Kanji.append(keb.text)
    reles = entry.findall("r_ele")
    for rele in reles:
            rebs = rele.findall("reb")
            for reb in rebs:
                Kana.append(reb.text)
    
    for kanji in Kanji:
        if not kanji in KKMAP:
            KKMAP[kanji] = []
        for kana in Kana:
            if not kana in KKMAP[kanji]:
                KKMAP[kanji].append(kana)
    

for kanji in KKMAP.keys():
    for kana in KKMAP[kanji]:
        conn.execute("INSERT INTO kanjikana VALUES (?, ?);", (kanji, kana))

conn.commit()  
