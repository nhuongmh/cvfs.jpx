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