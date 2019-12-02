#! /usr/bin/env python3

# generate sample file
# tested by: sort sample.csv| uniq -u | wc -l

# git clone https://github.com/dwyl/english-words.git
wordsFile = "english-words/words_alpha.txt"
sampleFile = "sample/sample.csv"

repeats = 256  # size(sampleFile) ≈ repeats * size(wordsFile)
uniqueRate = 50  # number of unique words in each repeat loop

#
# parameter for genBig
bigFile = "sample/big.csv"
maxNum = 1024 * 1024 * 4
rounds = 16
uniqueInRound = 10


def loadWords(path: str) -> set:
    with open(path) as file:
        words = set(file.read().split())
    return words


def writeSample(outFile: str, words: set, nround: int, uniRate: int):
    r"""generate sample file with normal words"""
    unique = set()
    writtenUnique = set()
    for i in range(nround*uniqueRate):
        unique.add(words.pop())

    with open(outFile, "w") as file:
        for i in range(nround):
            for j in range(uniqueRate):
                u = unique.pop()
                words.add(u)
                writtenUnique.add(u)  # faster than union

            for word in words:
                file.write(word)
                file.write("\n")

            while len(writtenUnique) > 0:
                words.discard(writtenUnique.pop())


def genBig(outFile: str):
    r"""can generate file that almost all words are unique"""
    unique = set()
    for u in range(rounds*uniqueInRound*2):  # overfill the set
        unique.add(str(maxNum+u))

    interval = round(maxNum / uniqueInRound)

    with open(outFile, "w") as file:
        for r in range(rounds):
            for i in range(maxNum):
                file.write(str(i))
                file.write("\n")
                if i % interval == 0:  # the second word in file must be unique
                    file.write(unique.pop())
                    file.write("\n")


if __name__ == "__main__":
    # words = loadWords(wordsFile)
    # writeSample(sampleFile, words, repeats, uniqueRate)

    genBig(bigFile)
