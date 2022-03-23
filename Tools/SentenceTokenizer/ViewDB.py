from janome.tokenizer import Tokenizer
import sqlite3

DB = "../../Database/Sentences.db"

search = "脇"

conn = sqlite3.connect(DB)
cur = conn.cursor()
cut = conn.cursor()
cur.execute("SELECT id FROM wordmap where word == ?", (search,))
while x := cur.fetchone():
    id = int(x[0])
    cut.execute("SELECT sentence from sentmap where ID == ?", (id, ))
    while y:= cut.fetchone():
        print(y[0])
    
conn.close()
