import os
import tempfile
import genanki
import requests
import shutil
from .utils import format_filename

'''
sample card data input:
{
            "word": "pitfall",
            "pronunciation": [
                {
                    "lang": "us",
                    "pron": "/ˈpɪt.fɑːl/",
                    "url": "https://dictionary.cambridge.org/us/media/english/us_pron/p/pit/pitfa/pitfall.mp3"
                }
            ],
            "definition": [
                {
                    "text": "a likely mistake or problem in a situation: ",
                    "pos": "",
                    "example": [
                        {
                            "text": "The store fell into one of the major pitfalls of small business, borrowing from suppliers by paying bills late."
                        },
                        {
                            "text": "There's a video that tells new students about pitfalls to avoid."
                        }
                    ]
                },
                {
                    "text": "an unexpected danger or difficulty: ",
                    "pos": "",
                    "example": [
                        {
                            "text": "Who knows what kind of pitfalls they’re going to run into."
                        }
                    ]
                },
            ],
            "properties": null,
            "context": "The store fell into one of the major pitfalls of small business, borrowing from suppliers by paying bills late.",
            "wordfreq": 1.234
        }
'''

fields=[
            {'name': 'Word'},
            {'name': 'Definition'},
            {'name': 'Pronunciation'},
            {'name': 'Ipa'},
            {'name': 'Context'},
            {'name': 'WordFreq'},
        ]
templates=[
        {
        'name': 'Card 1',
        'qfmt': '''<div class=jp; style='font-size: 60px;'> 
{{Word}}
</div><br>''',
        'afmt': '''<div id="back">{{FrontSide}}</div>
<hr id=answer>
<div class="definition">{{Definition}}</div>
<div class="ipa">{{Ipa}}</div>
<div class="pron">{{Pronunciation}}</div>
<div class="context">{{Context}}</div>''',
        },
        {
        'name': 'Card 2',
        'qfmt': '''<div class="context"> 
{{Context}}
</div><br>''',
        'afmt': '''<div id="back"> {{FrontSide}} </div>
<hr id=answer>
<div class="word">{{Word}}</div>
<div class="definition">{{Definition}}</div>
<div class="ipa">{{Ipa}}</div>
<div class="pron">{{Pronunciation}}</div>''',
        },
    ]

css = '''
.card {
font-family: arial;
font-size: 20px;
text-align: center;
color: black;
background-color: white;
}
.definition {
font-size: 15px;
text-align: left;
}
.pron {
font-size: 20px;
}
.context {
font-size: 20px;
text-align: left;
color: #dcdcdc;;
}
.word {
font-size: 60px;
}
.ipa {
font-size: 20px;
}
'''

http_headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
        "Connection": "keep-alive",
        "Upgrade-Insecure-Requests": "1",
    }
class IeCardData:
    def __init__(self, word: str, definition: str, pron_file_path: str, ipa: str, context: str, wordfreq: str):
        self.word = word
        self.definition = definition
        self.pron_file_path = pron_file_path
        self.ipa = ipa
        self.context = context
        self.wordfreq = wordfreq

def _parse_card_data(card_datas: list, media_dir) -> list[IeCardData]:
    parsed_data = []
    for card_data in card_datas:
        word = card_data.get('word', '')
        if not word:
            continue
        definition = card_data.get('definition', [])
        pronunciation = card_data.get('pronunciation', [])
        ipa = ''
        context = card_data.get('context_sentence', '')
        wordfreq = card_data.get('wordfreq', '')

        # Convert definition and pronunciation to string
        definition_combined = ''
        for idx, deff in enumerate(definition):
            text = deff.get('text', '')
            if not text:
                continue
            pos = deff.get('pos', '')
            examples = deff.get('example', [])
            example_str = ''
            if examples:
                example_str = examples[0].get('text', '')
            
            if pos:
                text = f"{text} <{pos}>"
            definition_combined += f"""
            <ul>
                <li>
                    <b>{idx+1}.</b> {text}
                    <br>
                    <span style="font-size: smaller; font-style: italic;">e.g: {example_str}</span>
                </li>
            </ul>
            """

        # only get the first pronunciation, download audio file to local temporary folder
        if pronunciation:
            audio_url = pronunciation[0].get('url', '')
            ipa = pronunciation[0].get('pron', '')
            audio_name = os.path.basename(audio_url)
            audio_path = os.path.join(media_dir, audio_name)
            if not os.path.exists(audio_path):
                # Download the audio file
                with open(audio_path, 'wb') as f:
                    with requests.get(audio_url, headers=http_headers, stream=True) as r:
                        r.raise_for_status()
                        for chunk in r.iter_content(chunk_size=8192):
                            f.write(chunk)
        else:
            audio_path = ''

        # highlight the word in context
        if context:
            context = context.replace(word, f"<span style='font-size: larger; font-weight: bold; color: red;'>{word}</span>")

        parsed_data.append(IeCardData(word, definition_combined, audio_path, ipa, context, wordfreq))

    return parsed_data

def gen_ie_anki_cards(card_datas: list, desk_name: str) -> str:

    media_dir = tmpdir = tempfile.mkdtemp()
    processed_card_datas = _parse_card_data(card_datas, media_dir)
    output_dir = "data/anki_ie_export/"
    if not os.path.exists(output_dir):
        os.makedirs(output_dir)
    output_file = os.path.join(output_dir, f"{format_filename(desk_name)}.apkg")

    deck = genanki.Deck(1672199389, desk_name)
    model = genanki.Model(
        1573999279,
        'English Vocab Model',
        fields=fields,
        templates=templates,
        css=css
    )
    
    media_files = []
    for card_data in processed_card_datas:
        if card_data.pron_file_path:
            media_files.append(card_data.pron_file_path)
        note = genanki.Note(
            model=model,
            fields=[card_data.word, 
                    card_data.definition, 
                    f"[sound:{os.path.basename(card_data.pron_file_path)}]",
                    card_data.ipa, 
                    card_data.context,
                    card_data.wordfreq],)
        deck.add_note(note)
    try:
        package = genanki.Package(deck)
        package.media_files = media_files
        package.write_to_file(output_file)
        print(f"Anki deck generated: {output_file}")
        return output_file
    finally:
        shutil.rmtree(tmpdir)