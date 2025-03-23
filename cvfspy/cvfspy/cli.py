import argparse
from cvfspy.ankigen_main import AnkiGen
from cvfspy.wserver import serve

argparser = argparse.ArgumentParser(description="CVFSPY")
argparser.add_argument(
    "-m",
    "--mode",
    dest="mode",
    default="server",
    help="Mode: anki/server",
)
args = argparser.parse_args()

def main():  # pragma: no cover
    if args.mode == "anki":
        print("Generating Anki deck...")
        ankigen = AnkiGen()
        ankigen.execute()
    else:
        print("Starting server...")
        serve()