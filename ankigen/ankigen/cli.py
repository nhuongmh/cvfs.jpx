from ankigen.ankigen_main import AnkiGen

def main():  # pragma: no cover
    print("Generating Anki deck...")
    ankigen = AnkiGen()
    ankigen.execute()