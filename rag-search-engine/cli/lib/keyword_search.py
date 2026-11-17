import os
import pickle
import string
from collections import Counter, defaultdict
import math

from nltk.stem import PorterStemmer

from .search_utils import *

class InvertedIndex:
    def __init__(self) -> None:
        self.index = defaultdict(set)
        self.docmap: dict[int, Movie] = {}
        self.index_path = os.path.join(CACHE_DIR, "index.pkl")
        self.docmap_path = os.path.join(CACHE_DIR, "docmap.pkl")
        self.tf_path = os.path.join(CACHE_DIR, "term_frequencies.pkl")
        self.term_frequencies = defaultdict(Counter)
        self.doc_lengths = {}
        self.doc_lengths_path = os.path.join(CACHE_DIR, "doc_lengths.pkl")

    def build(self) -> None:
        movies = load_movies()
        for m in movies:
            doc_id = m["id"]
            doc_description = f"{m['title']} {m['description']}"
            self.docmap[doc_id] = m
            self.__add_document(doc_id, doc_description)

    def save(self) -> None:
        os.makedirs(CACHE_DIR, exist_ok=True)
        with open(self.index_path, "wb") as f:
            pickle.dump(self.index, f)
        with open(self.docmap_path, "wb") as f:
            pickle.dump(self.docmap, f)
        with open(self.tf_path, "wb") as f:
            pickle.dump(self.term_frequencies, f)
        with open(self.doc_lengths_path, "wb") as f: # "a"
            pickle.dump(self.doc_lengths, f)

    def load(self) -> None:
        with open(self.index_path, "rb") as f:
            self.index = pickle.load(f)
        with open(self.docmap_path, "rb") as f:
            self.docmap = pickle.load(f)
        with open(self.tf_path, "rb") as f:
            self.term_frequencies = pickle.load(f)
        with open(self.doc_lengths_path, "rb") as f:
            self.doc_lengths = pickle.load(f)

    def get_documents(self, term: str) -> list[int]:
        doc_ids = self.index.get(term, set())
        return sorted(doc_ids)

    def __add_document(self, doc_id: int, text: str) -> None:
        tokens = tokenize_text(text)
        for token in set(tokens):
            self.index[token].add(doc_id)
        self.term_frequencies[doc_id].update(tokens)
        self.doc_lengths[doc_id] = len(tokens)

    def get_tf(self, doc_id: int, term: str) -> int:
        return self.term_frequencies[doc_id][term]

    def get_bm25_idf(self, term: str) -> float:
        idx = self
        n = len(idx.docmap)
        df = len(idx.get_documents(term))
        bm25_idf_score = math.log((n - df + 0.5) / (df + 0.5) + 1)
        return bm25_idf_score

    def get_bm25_tf(self, doc_id: int, term: str, k1=BM25_K1, b=BM25_B) -> int:
        tf = self.get_tf(doc_id, term)
        doc_length = self.doc_lengths[doc_id]
        avg_doc_length = self.__get_avg_doc_length()
        length_norm = 1 - b + b * (doc_length / avg_doc_length)
        return (tf * (k1 + 1)) / (tf + k1 * length_norm)

    def __get_avg_doc_length(self) -> float:
        lengths = []
        if len(self.docmap) == 0:
            return 0.0
        else:
            for length in self.doc_lengths.values():
                lengths.append(length)
                
        average_doc_length = sum(lengths)/len(lengths)
        return average_doc_length

    def bm25(self, doc_id: int, term: str):
        bm25_tf = self.get_bm25_tf(doc_id, term)
        bm25_idf = self.get_bm25_idf(term)
        return bm25_tf * bm25_idf

    def bm25_search(self, query, limit=5) -> dict[int, int]:
        tokens = tokenize_text(query)
        scores = {}
        for doc_id in self.docmap:
            total_bm25 = 0
            for token in tokens:
                total_bm25 += self.bm25(doc_id, token)
            scores[doc_id] = total_bm25
        return dict(sorted(scores.items(), key=lambda item: item[1], reverse=True)[:limit])

def build_command() -> None:
    idx = InvertedIndex()
    idx.build()
    idx.save()

def search_command(query: str, limit: int = DEFAULT_SEARCH_LIMIT) -> list[Movie]:
    idx = InvertedIndex()
    idx.load()
    query_tokens = tokenize_text(query)
    seen, results = set(), []
    for query_token in query_tokens:
        matching_doc_ids = idx.get_documents(query_token)
        for doc_id in matching_doc_ids:
            if doc_id in seen:
                continue
            seen.add(doc_id)
            doc = idx.docmap[doc_id]
            results.append(doc)
            if len(results) >= limit:
                return results

    return results


def preprocess_text(text: str) -> str:
    text = text.lower()
    text = text.translate(str.maketrans("", "", string.punctuation))
    return text


def load_stopwords() -> list[str]:
    with open(STOPWORDS_PATH, "r") as f:
        return [preprocess_text(word) for word in f.read().splitlines()]


STOPWORDS = load_stopwords()


def tokenize_text(text: str) -> list[str]:
    text = preprocess_text(text)
    tokens = text.split()
    valid_tokens = []
    for token in tokens:
        if token:
            valid_tokens.append(token)
    filtered_words = []
    for word in valid_tokens:
        if word not in STOPWORDS:
            filtered_words.append(word)
    stemmer = PorterStemmer()
    stemmed_words = []
    for word in filtered_words:
        stemmed_words.append(stemmer.stem(word))
    return stemmed_words


def tokenize_single_term(term: str) -> str:
    tokens = tokenize_text(term)
    if len(tokens) != 1:
        raise ValueError("term must be a single token")
    return tokens[0]


def tf_command(doc_id: int, term: str) -> int:
    idx = InvertedIndex()
    idx.load()
    return idx.get_tf(doc_id, tokenize_single_term(term))


def idf_command(term: str) -> int:
    idx = InvertedIndex()
    idx.load()
    tokenized_term = tokenize_single_term(term)
    total_doc_count = len(idx.docmap)
    term_match_doc_count = len(idx.get_documents(tokenized_term))
    idf = math.log((total_doc_count + 1) / (term_match_doc_count + 1))
    return idf

def tf_idf_command(doc_id: int, term: str) -> int:
    idx = InvertedIndex()
    idx.load()
    tf = tf_command(doc_id, term)
    idf = idf_command(term)
    return tf * idf

def bm25_idf_command(term: str) -> int:
    idx = InvertedIndex()
    idx.load()
    tokenized_term = tokenize_single_term(term)
    bm25_idf = idx.get_bm25_idf(tokenized_term)
    return float(bm25_idf)


def bm25_tf_command(doc_id: int, term: str, k1=BM25_K1, b=BM25_B):
    idx = InvertedIndex()
    idx.load()
    tokenized_term = tokenize_single_term(term)
    bm25_tf = idx.get_bm25_tf(doc_id, tokenized_term, k1, b)
    return bm25_tf