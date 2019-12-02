package main

import (
	"os"
	"strconv"

	"../../singleWordFinder"
)

func main() {
	if len(os.Args) != 3 {
		println("usage: unique [file path] [chunk size(MB)]")
		return
	}
	path := os.Args[1]
	chunk, err := strconv.Atoi(os.Args[2])
	if err != nil {
		println("invalid chunk size")
		return
	}

	word, err := singleWordFinder.FindUnique(path, int64(chunk*1024*1024))
	if err != nil {
		panic(err)
	}
	println("first unique:", word)
}
