package singleWordFinder

import (
	"os"
	"testing"
)

func TestSpiltFile(T *testing.T) {
	// 1. normal file  # done
	// 2. can't find delim in buf[maxline]
	// 3. starts[-1] == size(file)

	file, _ := os.Open("./sample/sample.csv")
	defer file.Close()
	starts, _ := SpiltFile(file, 1024*1024, '\n')

	buf := make([]byte, 1)
	for _, offset := range starts[1:] {
		file.ReadAt(buf, offset-1)
		if buf[0] != '\n' {
			T.Fatal("wrong start")
		}
	}
}
