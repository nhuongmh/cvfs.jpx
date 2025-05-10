from bs4 import BeautifulSoup
import requests
from flask import jsonify
import json
import traceback


def fetch_verbs(wiki):
    try:
        response = requests.get(wiki)
        soup = BeautifulSoup(response.text, "html.parser")
        verb_elements = soup.select("tr > td > p")
        verb_text = ""
        for element in verb_elements:
            verb_text += element.text

        lines = [line.strip() for line in verb_text.split("\n") if line.strip()]

        verbs = []
        for i in range(0, len(lines), 2):
            # Check if we have both type and text
            if i + 1 < len(lines):
                verb_type = lines[i]
                verb_text = lines[i + 1]

                # Check if this entry already exists
                exists = False
                for verb in verbs:
                    if verb.get("type") == verb_type and verb.get("text") == verb_text:
                        exists = True
                        break

                if exists:
                    break

                if verb_type and verb_text:
                    verbs.append(
                        {"id": len(verbs), "type": verb_type, "text": verb_text}
                    )

        return verbs
    except Exception as e:
        print(f"Error fetching verbs: {e}")
        return []


def get_dictionary(language, entry) -> dict:
    slug_language = language
    nation = "us"

    if slug_language == "en":
        language = "english"
    elif slug_language == "uk":
        language = "english"
        nation = "uk"

    url = f"https://dictionary.cambridge.org/{nation}/dictionary/{language}/{entry}"
    wiki = f"https://simple.wiktionary.org/wiki/{entry}"

    headers = {
        "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Accept-Language": "en-US,en;q=0.5",
        "Connection": "keep-alive",
        "Upgrade-Insecure-Requests": "1",
    }
    response = requests.get(url, headers=headers, timeout=30)

    if response.status_code != 200:
        raise Exception(f"Failed to fetch data from {url}")

    soup = BeautifulSoup(response.text, "html.parser")
    site_url = "https://dictionary.cambridge.org"

    # Get verbs
    verbs = fetch_verbs(wiki)

    # Basic info
    word = soup.select_one(".hw.dhw").text if soup.select_one(".hw.dhw") else ""

    # Part of speech
    pos_elements = soup.select(".pos.dpos")
    all_pos = [pos.text for pos in pos_elements]
    pos = list(dict.fromkeys(all_pos))  # Remove duplicates while preserving order

    # Phonetics audios
    audio = []
    pos_headers = soup.select(".pos-header.dpos-h")

    for header in pos_headers:
        pos_node = header.select_one(".dpos-g")
        if not pos_node or not pos_node.contents:
            continue

        p = pos_node.contents[0].text if pos_node.contents else ""
        pron_nodes = header.select("span.dpron-i")

        if not pron_nodes:
            continue

        for node in pron_nodes:
            if len(node.contents) < 3:
                continue

            lang = node.contents[0].text
            audio_element = node.select_one("audio")

            if not audio_element:
                continue

            source = audio_element.select_one("source")

            if not source:
                continue

            url = site_url + source.get("src")
            pron = node.contents[2].text if len(node.contents) > 2 else ""

            audio.append({"pos": p, "lang": lang, "url": url, "pron": pron})

    # Definition & example
    def_bodies = soup.select(".def-body.ddef_b")
    example_count = [len(body.select(".examp.dexamp")) for body in def_bodies]

    # Calculate cumulative counts
    for i in range(1, len(example_count)):
        example_count[i] = example_count[i] + example_count[i - 1]

    # Get examples
    example_elements = soup.select(".examp.dexamp > .eg.deg")
    example_trans = soup.select(".examp.dexamp > .trans.dtrans.dtrans-se.hdb.break-cj")

    examples = []
    for idx, element in enumerate(example_elements):
        translation = ""
        if idx < len(example_trans):
            translation = example_trans[idx].text

        examples.append({"id": idx, "text": element.text, "translation": translation})

    # Helper functions for definitions
    def get_source(element):
        parent = element.find_parent(".pr.dictionary")
        return parent.get("data-id") if parent else ""

    def get_pos(element):
        parent = element.find_parent(".pr.entry-body__el")
        if parent:
            pos_element = parent.select_one(".pos.dpos")
            return pos_element.text if pos_element else ""
        return ""

    def get_examples(element):
        examples = []
        example_elements = element.select(".def-body.ddef_b > .examp.dexamp")

        for idx, ex in enumerate(example_elements):
            text_element = ex.select_one(".eg.deg")
            trans_element = ex.select_one(".trans.dtrans")

            examples.append(
                {
                    "id": idx,
                    "text": text_element.text if text_element else "",
                    "translation": trans_element.text if trans_element else "",
                }
            )

        return examples

    # Get definitions
    definition_elements = soup.select(".def-block.ddef_block")
    definitions = []

    for idx, element in enumerate(definition_elements):
        def_text_element = element.select_one(".def.ddef_d.db")
        def_trans_element = element.select_one(".def-body.ddef_b > span.trans.dtrans")

        definitions.append(
            {
                "id": idx,
                "position": get_pos(element),
                "source": get_source(element),
                "text": def_text_element.text if def_text_element else "",
                "translation": def_trans_element.text if def_trans_element else "",
                "example": get_examples(element),
            }
        )

    # API response
    if not word:
        raise Exception("Word not found")
    else:
        return {
            "word": word,
            "position": pos,
            "verbs": verbs,
            "pronunciation": audio,
            "definition": definitions,
        }
