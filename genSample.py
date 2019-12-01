#! /usr/bin/env python3

# generate sample file
# tested by: sort sample.csv| uniq -u | wc -l

# git clone https://github.com/dwyl/english-words.git
wordsFile = "english-words/words_alpha.txt"
sampleFile = "sample/sample.csv"

repeats = 5  # size(sampleFile) ≈ repeats * size(wordsFile)
uniqueRate = 10  # number of unique words in each repeat loop


def loadWords(path: str) -> set:
    with open(path) as file:
        words = set(file.read().split())
    return words


def writeSample(outFile: str, words: set, round: int, uniRate: int):
    unique = set()
    writtenUnique = set()
    for i in range(round*uniqueRate):
        unique.add(words.pop())

    with open(outFile, "w") as file:
        for i in range(round):
            for j in range(uniqueRate):
                u = unique.pop()
                words.add(u)
                writtenUnique.add(u)  # faster than union

            for word in words:
                file.write(word)
                file.write("\n")

            while len(writtenUnique) > 0:
                words.discard(writtenUnique.pop())


if __name__ == "__main__":
    words = loadWords(wordsFile)
    writeSample(sampleFile, words, repeats, uniqueRate)
