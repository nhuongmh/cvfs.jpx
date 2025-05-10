from flask import Flask, request, jsonify
from flask_cors import CORS
import traceback
from cvfspy.usecases.cambrige_dict_fetcher import get_dictionary
from cvfspy.usecases.article_fetcher import fetch_article
from cvfspy.usecases.vocab_analyzer import article_vocab_extractor
app = Flask(__name__)
CORS(app, resources={r"/*": {"origins": "*"}})


# @app.route('/')
# def index():
#     return send_from_directory('.', 'index.html')

@app.route('/api/dictionary/<language>/<entry>')
def get_english_dictionary(language, entry):
    try:
        rep =  get_dictionary(language, entry)
        return jsonify(rep), 200
    except Exception as e:
        traceback.print_exc()
        return jsonify({"error": str(e)}), 500
    
@app.route('/api/dictionary/<language>/', methods=['POST'])
def fetch_english_multiple_vocab(language):
    data = request.json
    if not data:
        return jsonify({"error": "Invalid request"}), 400
    def empty_word(word):
        return {
            "word": word,
            "position": [],
            "verbs": [],
            "pronunciation": [],
            "definition": [],
        }
    try:
        processed_vocab = {}
        for entry in data:
            if not entry:
                print("Empty entry found, skipping...")
                continue
            try:
                vocab =  get_dictionary(language, entry)
            except Exception as e:
                print(f"Error fetching dictionary for {entry}: {e}")
                vocab = empty_word(entry)
            if not vocab:
                print(f"No vocabulary found for entry: {entry}")
                vocab = empty_word(entry)
            print(f"Fetched dictionary for {entry}: {vocab}")
            processed_vocab[entry] = vocab
        return jsonify(processed_vocab), 200
    except Exception as e:
        traceback.print_exc()
        return jsonify({"error": str(e)}), 500

@app.route('/api/article/')
def get_article():
    try:
        url = request.args.get('url')
        if not url:
            return jsonify({"error": "URL is required"}), 400
        article = fetch_article(url)
        return jsonify(article.model_dump()), 200
    except Exception as e:
        traceback.print_exc()
        return jsonify({"error": str(e)}), 500
    
@app.route('/api/vocab_extractor/', methods=['POST'])
def extract_vocab():
    data = request.json
    if not data:
        return jsonify({"error": "Invalid request"}), 400
    text = data.get('content')
    if not text:
        text = data.get('text')
    if not text:
        return jsonify({"error": "Content is required"}), 400
    threshold = data.get('threshold', 4.5)
    try:
        words = article_vocab_extractor(text, threshold)
        return jsonify([word.model_dump() for word in words]), 200
    except Exception as e:
        traceback.print_exc()
        return jsonify({"error": str(e)}), 500

def serve():
    app.run(debug=True, port=5000)