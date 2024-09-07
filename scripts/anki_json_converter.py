import os
import json
import argparse


parser = argparse.ArgumentParser(description='Anki JSON to Tsunami converter')
parser.add_argument('-f', dest='jsonfile')           # positional argument

args = parser.parse_args()

def process(input_file):
    output_file = os.path.splitext(input_file)[0] + '_processed.json'
    with open(input_file, 'r') as infile, open(output_file, 'w', newline='') as outfile:
        data = json.load(f)
        