from janome.tokenizer import Tokenizer
import sqlite3

DB = "../../Database/Sentences_2.db"

search = "脇"

conn = sqlite3.connect(DB)
cur = conn.cursor()
cur.execute("SELECT * FROM dict where word == ?", (search,))
while x := cur.fetchone():
    print(x)
conn.close()
