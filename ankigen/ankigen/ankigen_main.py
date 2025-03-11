import dotenv
import os
import base64
import json
import genanki
import csv
from ankigen.google_sheet import GoogleSheet
import time
SCOPES = ["https://www.googleapis.com/auth/spreadsheets.readonly"]

class CardData:
    def __init__(self, lesson, sentence, hv, meaning, structure):
        self.lesson = lesson
        self.structure = structure
        self.sentence = sentence
        self.meaning = meaning
        self.hv = hv

class AnkiGen:
    def __init__(self):  # pragma: no cover
        pass
    
    def execute(self):
        card_datas = self.load_from_google()
        if card_datas:
            # group cards by lesson
            lesson_cards = {}
            for card_data in card_datas:
                if card_data.lesson not in lesson_cards:
                    lesson_cards[card_data.lesson] = []
                lesson_cards[card_data.lesson].append(card_data)
            # generate a deck for each lesson
            # self.generate_anki_deck(f"Minna no Nihongo Grammar", [], "output/testm.apkg")
            for lesson, cards in lesson_cards.items():
                outfile = f'output/{lesson.replace(" ", "_").replace("/", "_")}.apkg'
                self.generate_anki_deck(f"Minna no Nihongo Grammar::{lesson}", cards, outfile)
                print(f"Deck {outfile} generated.")

    def load_from_google(self) -> list[CardData]:
        config = {
        **dotenv.dotenv_values(".env"),
        **dotenv.dotenv_values(".private.env"),
        **os.environ
        }
        encoded_creds = config.get("GOOGLE_API_KEY_BASE64", None)
        creds_data = None
        if encoded_creds:
            decoded_creds = base64.b64decode(encoded_creds)
            creds_data = json.loads(decoded_creds)
        else:
            print("No credentials found.")
            return
        spreadsheet_id = config.get("GOOGLE_SPREADSHEET_ID", None)
        sheet_name = config.get("GOOGLE_GRAMMAR_SHEET_NAME", None)
        if not spreadsheet_id or not sheet_name:
            print("No spreadsheet ID or sheet name found.")
            return
        google_sheet = GoogleSheet(creds_data, spreadsheet_id, sheet_name)
        rows = google_sheet.load()
        print(f"Number of rows: {len(rows)}")
        # convert rows to CardData objects
        card_datas = []
        if rows:
            for row in rows:
                if len(row) >=4:
                    card_datas.append(CardData(lesson=row[0], structure=row[1], sentence=row[2], meaning=row[3], hv='' if len(row) < 5 else row[4]))
        return card_datas    

    def generate_anki_deck(self, desk_name, card_datas: list[CardData], output_file):
        fields=[
            {'name': 'Sentence'},
            {'name': 'HV'},
            {'name': 'Meaning'},
            {'name': 'GrammarStructure'},
            {'name': 'PitchAccent'},
            {'name': 'Audio'}
        ]
        templates=[
            {
            'name': 'Card 1',
            'qfmt': '''<div class=jp; style='font-family: STKaiti; font-size: 100px;'> 
<input type="checkbox" id="check"/>
<label for="check">{{furigana:Sentence}}</label> 
</div><br>''',
            'afmt': '''<div id="back">{{FrontSide}}</div>
<hr id=answer>
<div style='font-family: Times New Roman; font-size: 50px;'>{{Meaning}}</div>
<div style='font-family: Times New Roman; font-size: 30px;'>{{HV}}</div>
<div style='font-family: Times New Roman; font-size: 30px;'>{{GrammarStructure}}</div>
<div style='font-family: STKaiti; font-size: 20px;'>{{Audio}}</div>
<div style='font-family: Times New Roman; font-size: 30px;'>{{PitchAccent}}</div>''',
            },
            {
            'name': 'Card 2',
            'qfmt': '''<div class=jp; style='font-family: STKaiti; font-size: 100px;'> 
<input type="checkbox" id="check"/>
<label for="check">{{Meaning}}</label> 
</div><br>''',
            'afmt': '''<div id="back"> {{FrontSide}} </div>
<hr id=answer>
<div style='font-family: Times New Roman; font-size: 50px;'>{{furigana:Sentence}}</div>
<div style='font-family: Times New Roman; font-size: 30px;'>{{HV}}</div>
<div style='font-family: STKaiti; font-size: 20px;'>{{Audio}}</div>
<div style='font-family: Times New Roman; font-size: 30px;'>{{GrammarStructure}}</div>
<div style='font-family: Times New Roman; font-size: 30px;'>{{PitchAccent}}</div>''',
            },
        ]
        css = '''
@import url("_ajt_japanese_24.10.8.1.css");
.card {
 font-family: arial;
 font-size: 20px;
 text-align: center;
 color: black;
 background-color: white;
}

.jp { font-size: 100px }
.win .jp { font-family: "MS Mincho", "ＭＳ 明朝"; }
.mac .jp { font-family: "Hiragino Mincho Pro", "ヒラギノ明朝 Pro"; }
.linux .jp { font-family: "Kochi Mincho", "東風明朝"; }
.mobile .jp { font-family: "Hiragino Mincho ProN"; }

/* hide the checkbox */
#check{
	display:none;
}

/* when checkbox is active, show for furigana */
#check:checked + label ruby rt {
    visibility: visible; 
}

/* furigana hidden by default*/
ruby rt { 
	visibility: hidden; 
}

#back #check + label ruby rt {
	visibility: visible !important; 
}'''
        deck = genanki.Deck(1417337802, desk_name)
        model = genanki.Model(
            1963116713,
            'MinnaNoGrammar Model',
            fields=fields,
            templates=templates,
            css=css
        )
        for card_data in card_datas:
            note = genanki.Note(
                model=model,
                fields=[card_data.sentence, card_data.hv, card_data.meaning, card_data.structure, '', ''])
            deck.add_note(note)
        genanki.Package(deck).write_to_file(output_file)