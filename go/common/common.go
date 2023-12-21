package common

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func Split2(s string, sep string) (s1, s2 string) {
	slices := strings.SplitN(s, sep, 2)
	return strings.Trim(slices[0], " "), strings.Trim(slices[1], " ")
}

func ParseInt(number string) int {
	if i, err := strconv.Atoi(number); err == nil {
		return i
	} else {
		panic(err)
	}
}

func ParseUint64(number string) uint64 {
	if i, err := strconv.ParseUint(number, 10, 64); err == nil {
		return i
	} else {
		panic(err)
	}
}

func ReadFile(path string, handle func(line string)) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewScanner(file)

	for reader.Scan() {
		line := reader.Text()
		handle(line)
	}
}

func ReadFileLines(path string, handle func(line string, index int)) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewScanner(file)

	index := 0
	for reader.Scan() {
		line := reader.Text()
		handle(line, index)
		index++
	}
}

func IntCompare(a, b int) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func Uint64Compare(a, b uint64) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func CreateFileWriter(path string) FileWriter {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	writer := bufio.NewWriter(file)

	return FileWriter{
		Write: func(line string) {
			writer.WriteString(line)
			if err := writer.Flush(); err != nil {
				panic(err)
			}
		},
		Close: func() {
			file.Close()
		},
	}
}

type FileWriter struct {
	Write func(line string)
	Close func()
}
