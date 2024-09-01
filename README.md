# AmeKanji

AmeKanji is a web service and Firefox extension for tracking Japanese vocabulary from Jisho and automatically creating Anki cards.

You can access it by visiting [AmeKanji.com](https://amekanji.com/)

## Usage

1. Go to the AmeKanji website.
2. Enter your Japanese vocabulary. AmeKanji will auto-select the kana reading, which you can adjust if needed.
3. Click on the arrow to generate your cards.

![image](https://github.com/user-attachments/assets/cd9b4d9c-5671-42e9-860b-df1ef37a9ae7)


![image](https://github.com/user-attachments/assets/f3836f6b-42ef-4e36-9744-39f32adbc399)


AmeKanji will provide:

  - An .apkg file with all your Anki cards (Japanese → English, English → Japanese, Kanji → Meaning).
  - An anki_deck.txt file for manual Anki import.
  - A Media/ directory with audio and kanji stroke GIFs.

Simply import the .apkg file to Anki and begin your studies:

![image](https://github.com/user-attachments/assets/41601485-93b4-4bef-b2e0-3cd96cd68d6c)



##  Extension

You can also download an [extension](https://addons.mozilla.org/en-US/firefox/addon/amekanji/). 

This extension will add an "Add to memory" button to all Jisho entries:

![image](https://github.com/user-attachments/assets/837917c3-2200-4988-9d3d-e2fcc15fd907)

Access your stored vocabulary by clicking the extension icon:

![image](https://github.com/user-attachments/assets/3ba49cb3-4cfe-45a9-bf47-5d5edc6d89e7)


Download your data directly from the extension!

## GoAme

GoAme includes a standalone software. The web-app had its backend written in GoLang, powered by RabbitMQ, with a frontend built in Svelte.
