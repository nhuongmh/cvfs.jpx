import nltk
from wordfreq import zipf_frequency
from nltk.corpus import words, stopwords
from nltk.tokenize import word_tokenize, sent_tokenize, SpaceTokenizer
from nltk.stem import WordNetLemmatizer
from cvfspy.usecases.common import BaseSchema

nltk.download('punkt')
nltk.download('punkt_tab')
nltk.download('words')
nltk.download('stopwords')
nltk.download('wordnet')

wnl = WordNetLemmatizer()

class Vocab(BaseSchema):
    word: str 
    context_sentence: str = ""
    freq: float = 0.0


def article_vocab_extractor(text: str, freq_threshold: float = 4.5) -> list[Vocab]:
    sentences = sent_tokenize(text)
    # Get English dictionary words and common words to exclude
    english_vocab = set(words.words())
    common_words = set(stopwords.words('english'))

    learning_words: dict[str, Vocab] = {}
    for sentence in sentences:
        tk = SpaceTokenizer() 
        tokens = tk.tokenize(sentence)
        # Clean tokens (remove punctuation, numbers, etc.)
        tokens = [wnl.lemmatize(word) for word in tokens if word.isalpha()]
        for word in tokens:
            if len(word) < 3:
                continue
            word_freq = zipf_frequency(word, 'en')
            if word_freq == 0.0 or word_freq > freq_threshold:
                continue

            # if word start with capital letter and not in english_vocab, skip
            if word[0].isupper() and word.lower() not in english_vocab:
                continue
            # if word not in english_vocab:
            #     continue

            if word in common_words:
                continue

            if word in learning_words:
                continue

            learning_words[word] = Vocab(word=word, context_sentence=sentence, freq=word_freq)
    words_list = learning_words.values()
    #sort learning_words by freq in descending order
    words_list = sorted(words_list, key=lambda x: x.freq)

    return words_list
    