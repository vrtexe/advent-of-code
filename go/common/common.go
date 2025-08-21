package common

import (
	"bufio"
	"errors"
	"fmt"
	"math"
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

func IntPow(value int, exponent int) int {
	return int(math.Pow(float64(value), float64(exponent)))
}

func ParseInt64(number string) int64 {
	if i, err := strconv.ParseInt(number, 10, 64); err == nil {
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

func Prepend(x []int, y int) []int {
	x = append(x, 0)
	copy(x[1:], x)
	x[0] = y
	return x
}

func Pop[S ~[]E, E any](s S) (E, S) {
	lastIndex := len(s) - 1
	last := s[lastIndex]
	return last, s[:lastIndex]
}

func Shift[S ~[]E, E any](s S) (E, S) {
	last := s[0]
	return last, s[1:]
}

func Det(mat [][]float64) float64 {
	result, err := det(mat)
	if err != nil {
		fmt.Println(err)
	}
	return result
}

func det(mat [][]float64) (float64, error) {
	// Base cases and rules
	if len(mat) != len(mat[0]) {
		return 0.0, errors.New("determinant can only be performed on square matrices")
	}

	if len(mat) == 1 {
		return (mat[0][0]), nil
	}

	if len(mat) == 2 {
		return (mat[0][0] * mat[1][1]) - (mat[0][1] * mat[1][0]), nil
	}

	s := 0.0 // accumulator
	for i := 0; i < len(mat[0]); i++ {

		sm := subMat(mat[1:][:], i) // peel off top row before passing
		z, err := det(sm)           // get determinant of sub-matrix

		if err == nil {
			if i%2 != 0 {
				s -= mat[0][i] * z
			} else {
				s += mat[0][i] * z
			}
		}
	}
	return s, nil
}

func subMat(mat [][]float64, p int) [][]float64 {
	stacks := make([]stack, len(mat))
	for n := range mat {
		stacks[n] = stack{}
		for j := range mat[n] {
			if j != p {
				stacks[n].push(mat[n][j])
			}
		}
	}
	out := make([][]float64, len(mat))
	for k := range stacks {
		out[k] = stacks[k].ToSlice()
	}
	return out
}

type stack []float64

func (s *stack) isEmpty() bool {
	return len(*s) == 0
}
func (s *stack) push(n float64) {
	*s = append(*s, n)
}
func (s *stack) pop() (float64, bool) {
	if s.isEmpty() {
		return 0, false
	}
	i := len(*s) - 1
	n := (*s)[i]
	*s = (*s)[:i]
	return n, true
}
func (s *stack) ToSlice() []float64 {
	return *s
}
