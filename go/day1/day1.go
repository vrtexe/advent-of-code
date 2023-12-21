package day1

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

func Solution() {
	file, err := os.Open("day1/day1.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewScanner(file)

	result := 0
	for reader.Scan() {
		line := reader.Text()
		result += findNumber(line)
	}
}

func findNumber(line string) int {
	firstNumber, lastNumber := findNumberOccurrenceIndex(line)
	number := fmt.Sprint(firstNumber.number) + fmt.Sprint(lastNumber.number)
	if value, err := strconv.Atoi(number); err != nil {
		panic(err)
	} else {
		return value
	}
}

func findNumberOccurrenceIndex(line string) (firstValue, lastIndex StringNumber) {

	minNumber := findStringDigit(line)
	maxNumber := findLastStringDigit(line)
	actualMinDigitIndex := strings.IndexFunc(line, isNumber)
	actualMaxDigitIndex := strings.LastIndexFunc(line, isNumber)

	if actualMinDigitIndex < 0 || actualMaxDigitIndex < 0 {
		return *minNumber, *maxNumber
	}

	val, _ := strconv.Atoi(string(line[actualMinDigitIndex]))
	actualDigit := StringNumber{actualMinDigitIndex, val}

	maxVal, _ := strconv.Atoi(string(line[actualMaxDigitIndex]))
	actualMaxDigit := StringNumber{actualMaxDigitIndex, maxVal}

	if minNumber == nil || maxNumber == nil {
		if actualDigit.number == actualMaxDigit.number {
			fmt.Println(actualDigit.number, actualMaxDigit.number, line)
		}
		return actualDigit, actualMaxDigit
	}

	return getMin(actualDigit, *minNumber), getMax(actualMaxDigit, *maxNumber)
}

func getMin(a, b StringNumber) StringNumber {
	if minIndex := min(a.index, b.index); minIndex == a.index {
		return a
	} else {
		return b
	}
}

func getMax(a, b StringNumber) StringNumber {
	if maxIndex := max(a.index, b.index); maxIndex == a.index {
		return a
	} else {
		return b
	}
}

type StringNumber struct {
	index  int
	number int
}

func findStringDigit(str string) *StringNumber {
	numberStrings := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	numbers := []StringNumber{}
	for index, value := range numberStrings {
		if foundIndex := strings.Index(str, value); foundIndex < 3 && foundIndex >= 0 {
			return &StringNumber{foundIndex, index + 1}
		} else if foundIndex >= 0 {
			numbers = append(numbers, StringNumber{foundIndex, index + 1})
		}
	}

	if len(numbers) <= 0 {
		return nil
	}

	result := slices.MinFunc(numbers, func(a, b StringNumber) int {
		return intCompare(a.index, b.index)
	})

	return &result
}

func findLastStringDigit(str string) *StringNumber {
	numberStrings := []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	numbers := []StringNumber{}
	for index, value := range numberStrings {
		if foundIndex := strings.LastIndex(str, value); foundIndex > len(str)-3 {
			return &StringNumber{foundIndex, index + 1}
		} else if foundIndex >= 0 {
			numbers = append(numbers, StringNumber{foundIndex, index + 1})
		}
	}

	if len(numbers) <= 0 {
		return nil
	}

	result := slices.MaxFunc(numbers, func(a, b StringNumber) int {
		return intCompare(a.index, b.index)
	})
	return &result
}

func intCompare(a, b int) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}

func isNumber(char rune) bool {
	return unicode.IsDigit(char)
}
