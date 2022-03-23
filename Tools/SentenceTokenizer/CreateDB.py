from janome.tokenizer import Tokenizer
import sqlite3

DB = "../../Database/Sentences.db"
DICT = "../../Repository/Sentences/jpn_sentences.tsv"

tokenizer = Tokenizer()
conn = sqlite3.connect(DB)
print("Creating table...")
try:
    conn.execute("CREATE TABLE wordmap (word TEXT, id INTEGER);")
    conn.execute("CREATE TABLE sentmap (id INTEGER, sentence TEXT);")
    print("Success!")
except sqlite3.OperationalError:
    conn.execute("DROP TABLE wordmap")
    conn.execute("DROP TABLE sentmap")
    conn.execute("CREATE TABLE wordmap (word TEXT, id INTEGER);")
    conn.execute("CREATE TABLE sentmap (id INTEGER, sentence TEXT);")
    print("Table already exists. Recreated!")

sfile = open(DICT, "r")
sentences = sfile.readlines()

id = 0
for sentence in sentences:
    sep = sentence.split("\t")
    id = id + 1
    result = str(sep[2])
    
    conn.execute("INSERT INTO sentmap VALUES (?, ?)", (id, result))

    for token in tokenizer.tokenize(result, wakati=True):
        conn.execute("INSERT INTO wordmap VALUES (?, ?)", (str(token), id))
    
conn.commit()
conn.close()
