package singleWordFinder

import (
	"bufio"
	"encoding/csv"
	"hash/fnv"
	"io"
	"os"
	"strconv"
)

const (
	MAXLINE = 64
	BASEDIR = "/home/lob/learning/singleWordFinder/tmp/"
)

// spiltFile read a file, logically spilt them into parts in chunkSize
// each part ends with delim, returns offsets of each part
// chunkSize < MAXLINE
func SpiltFile(file *os.File, chunkSize int64, delim byte) ([]int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	var offset int64
	segments := int((stat.Size() + chunkSize - 1) / chunkSize)
	starts := make([]int64, 0, segments)
	starts = append(starts, 0)

	buf := make([]byte, MAXLINE)

	var i, j int
	for i = 1; i < segments; i++ {
		_, err = file.ReadAt(buf, int64(i)*chunkSize)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, err
			}
		}
		for j = range buf {
			if buf[j] == delim {
				offset = int64(i)*chunkSize + int64(j+1)
				starts = append(starts, offset)
				break
			}
		}

		// if buf[j] != '\n' { // can't find delim within buf[MAXLINE]
		// 	// return nil, error
		// }
	}

	if starts[len(starts)-1] >= stat.Size() {
		starts = starts[:len(starts)-1]
	}
	return starts, nil
}

// reads file
func Mapper(inFilePath string, start int64, end int64, nReducer int) error {
	inFile, err := os.Open(inFilePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	stat, _ := inFile.Stat()
	if start >= stat.Size() {
		return nil
	}

	_, err = inFile.Seek(start, 0)
	if err != nil {
		return err
	}

	lreader := io.LimitReader(inFile, end-start)
	breader := bufio.NewReader(lreader)
	words := make(map[string]int64)

	// read words from original file segment
	var ok bool
	var pos int64 = start
	word, err := breader.ReadString('\n')
	for err == nil {
		if _, ok = words[word]; !ok {
			words[word] = pos
		} else {
			words[word] = -1
		}

		pos++
		word, err = breader.ReadString('\n')
	}
	if err != io.EOF {
		return err
	}

	// put unique words into files base on hashed value
	outFiles := make([]*bufio.Writer, nReducer)
	for i := range outFiles {
		file, err := os.OpenFile(BASEDIR+"Hashed_"+strconv.Itoa(i), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		defer file.Close()
		outFiles[i] = bufio.NewWriter(file)
	}

	var i int
	for word, pos := range words {
		if pos != -1 {
			i = int(hash(word)) % nReducer
			_, err := outFiles[i].WriteString(word[:len(word)-1] + "," + strconv.Itoa(int(pos)) + "\n")
			if err != nil {
				return err
			}
		}
	}

	for i := range outFiles {
		if err := outFiles[i].Flush(); err != nil {
			return err
		}
	}

	return nil
}

// Reducer return the first unique word in corresponding file
func Reducer(reduceNum int) (string, int, error) {
	filePath := BASEDIR + "Hashed_" + strconv.Itoa(reduceNum)
	file, err := os.Open(filePath)
	if err != nil {
		return "", -1, err
	}
	defer file.Close()

	breader := bufio.NewReader(file)
	words := make(map[string]int)
	reader := csv.NewReader(breader)

	var word string
	var pos int
	var ok bool
	record, err := reader.Read()
	for err == nil {
		word = record[0]
		if _, ok = words[word]; ok {
			words[word] = -1
		} else {
			pos, err = strconv.Atoi(record[1])
			if err != nil {
				return "", -1, err
			}
			words[word] = pos
		}
		record, err = reader.Read()
	}

	if err != io.EOF {
		return "", -1, err
	}

	// find first unique word
	first := int(^uint(0) >> 1)
	var firstWord string

	for word, pos = range words {
		if pos >= 0 && pos < first {
			first = pos
			firstWord = word
		}
	}

	return firstWord, first, nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
