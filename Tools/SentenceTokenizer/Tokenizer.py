from janome.tokenizer import Tokenizer
t = Tokenizer()

sentence = "僕は子供からあなたが大好きだよ。"

for token in t.tokenize(sentence):
    print(token)

# Load Sentence Dictionary
# For each sentence, tokenize it, and for each word in the sentence, add the [word, id of sentence] to a DB