import json
import logging

logger = logging.getLogger("idf_logger")
logger.setLevel(logging.DEBUG)

word_dict = {}
article_count = 0

def process_data(data_dict):
    doc_dict = []
    global article_count, word_dict
    for doc in data_dict["documents"]:
        for para in doc["segmented_paragraphs"]:
            for word in para:
                if word not in doc_dict: doc_dict.append(word)
    for word in doc_dict:
        if word in word_dict:
            word_dict[word] += 1
        else:
            word_dict[word] = 1
    article_count += 1
    if not article_count % 200:
        logging.warn("end article" + str(article_count))
        logging.warn("word_dict size is " + str(len(word_dict)))

def reader(filepath, start_line=0, end_line=None):
    with open(filepath) as f:
        for _ in range(start_line): next(f)
        for i, l in enumerate(f):
            if end_line is not None and i+start_line == end_line: return
            process_data(json.loads(l))

def saver(filepath="idf.dict"):
    with open(filepath, "w") as f:
        f.write(str(article_count) + "\n")
        for k, v in word_dict.items():
            f.write(k + " " + str(v) + "\n")
def main():
    reader("search.train.json")
    saver()

if __name__ == "__main__":
    main()
