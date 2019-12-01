#! /usr/bin/env python3

PATH = "./sample/sample.csv"

words = dict()
count = 0
with open(PATH) as file:
    for line in file:
        if line in words:
            words[line] = -1
        else:
            words[line] = count
        count += 1


first = 0xfffffffff
for k, v in words.items():
    if v > 0 and v < first:
        first = v
        word = k

unique = 0
for k, v in words.items():
    if v != -1:
        unique += 1
        print(k[:-1], v)

print("unique:", unique)
print("word:", word, "pos:", first)
