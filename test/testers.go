package main

import (
	"bufio"
	"io"

	"../../singleWordFinder"

	"os"
)

const (
	Sample = "/home/lob/learning/singleWordFinder/sample/sample.csv"
	Chunk  = 1024 * 1024
	Delim  = '\n'
)

func testSpiltFile() {
	file, _ := os.Open("/home/lob/learning/singleWordFinder/sample/sample.csv")
	defer file.Close()

	starts, _ := singleWordFinder.SpiltFile(file, 1024*1024, '\n')
	buf := make([]byte, 1)
	for _, offset := range starts[1:] {
		file.ReadAt(buf, offset-1)
		if buf[0] != '\n' {
			panic("wrong!!!!!!")
		}
	}
}

// another approach, must be correct
func testFirstUnique() {
	dictFile, err := os.Open("/home/lob/learning/singleWordFinder/sample/sample.csv")
	must(err)
	defer dictFile.Close()
	r := bufio.NewReader(dictFile)

	var word string
	pos := 0
	dict := make(map[string]int)
	for {
		word, err = r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}
		word = word[:len(word)-1]
		if _, ok := dict[word]; !ok {
			dict[word] = pos
		} else {
			dict[word] = -1
		}
		pos++
	}

	first := 0xffffffff
	var firstWord string
	for w, p := range dict {
		if p >= 0 && p < first {
			first = p
			firstWord = w
		}
	}
	for w, p := range dict {
		if p >= 0 {
			println(w, p)
		}
	}

	println("pos:", first)
	println("word:", firstWord)
	// file, _ := os.Open("../sample/sample.csv")
	// defer file.Close()

	// reader := bufio.NewReader(file)

}

func testMapper() {
	file, err := os.Open(Sample)
	must(err)

	starts, err := singleWordFinder.SpiltFile(file, Chunk, Delim)
	must(err)
	err = singleWordFinder.Mapper(Sample, starts[0], starts[1], len(starts))
	must(err)
}

func testALL() {
	file, err := os.Open(Sample)
	must(err)

	starts, err := singleWordFinder.SpiltFile(file, Chunk, Delim)
	must(err)

	for i := 0; i < len(starts)-1; i++ {
		err = singleWordFinder.Mapper(Sample, starts[i], starts[i+1], len(starts))
		must(err)
	}

	firsts := make(map[string]int)
	var word string
	var pos int
	for i := 0; i < len(starts); i++ {
		word, pos, err = singleWordFinder.Reducer(i)
		must(err)
		if _, ok := firsts[word]; ok {
			panic("NOTTTTTTTTTT UNIQUE")
		}
		firsts[word] = pos
	}

	var firstWord string
	firstPos := 0xfffffffff
	for word, pos := range firsts {
		if pos >= 0 && pos < firstPos {
			firstPos = pos
			firstWord = word
		}
	}

	println("first:", firstWord, firstPos)
}

func main() {
	testALL()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
