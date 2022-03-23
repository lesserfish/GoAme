from janome.tokenizer import Tokenizer
import sqlite3

DB = "../../Database/Sentences_2.db"
DICT = "../../Repository/Sentences/jpn_sentences.tsv"

tokenizer = Tokenizer()
conn = sqlite3.connect(DB)
print("Creating table...")
try:
    conn.execute("CREATE TABLE dict (word TEXT, id INTEGER, sentence TEXT);")
    print("Success!")
except sqlite3.OperationalError:
    conn.execute("DROP TABLE dict")
    conn.execute("CREATE TABLE dict (word TEXT, id INTEGER, sentence TEXT);")
    print("Table already exists. Recreated!")

sfile = open(DICT, "r")
sentences = sfile.readlines()

i = 0
for sentence in sentences:
    sep = sentence.split("\t")
    id = int(sep[0])
    result = str(sep[2])

    for token in tokenizer.tokenize(result, wakati=True):
        conn.execute("INSERT INTO dict VALUES (?, ?, ?)", (str(token), id, str(result)))
    
conn.commit()
conn.close()
