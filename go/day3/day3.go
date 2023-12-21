package main

import (
	"common"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"unicode"
)

type Engine struct {
	parts []Part
}

type Part struct {
	number   int
	position Position
	symbol   *Symbol
}

type Symbol struct {
	value    string
	position Position
	parts    []*Part
}

type Position struct {
	row, start, end int
}

func (this Position) Adjacent(that Position) bool {
	if math.Abs(float64(this.row-that.row)) > 1 {
		return false
	}

	return max(this.start, that.start) <= min(this.end, that.end)
}

func (this Symbol) IsGear() bool {
	return this.value == "*" && len(this.parts) == 2
}

func (this Symbol) GearRatio() int {
	n1, n2 := this.parts[0], this.parts[1]
	return n1.number * n2.number
}

func main() {

	regex, _ := regexp.Compile("(\\d+)|([^\\.])")
	wd, _ := os.Getwd()
	fmt.Println("here:", wd)
	parts := [][]*Part{}
	symbols := [][]*Symbol{}

	common.ReadFile("../day3/day3_2.txt", func(line string) {
		indexMatches := regex.FindAllStringIndex(line, -1)

		partsRow := []*Part{}
		symbolsRow := []*Symbol{}
		for _, i := range indexMatches {
			start, end := i[0], i[1]
			m := line[start:end]
			if !unicode.IsDigit(rune(m[0])) {
				symbolsRow = append(symbolsRow, &Symbol{m, Position{len(parts), start, end}, []*Part{}})
			} else {
				partsRow = append(partsRow, &Part{common.ParseInt(m), Position{len(parts), start, end}, nil})
			}
		}

		prevSymbols := getPreviousRow(symbols)
		prevParts := getPreviousRow(parts)

		symbols = append(symbols, symbolsRow)
		parts = append(parts, partsRow)

		cp := append(prevParts, partsRow...)
		cs := append(prevSymbols, symbolsRow...)

		for _, part := range cp {
			for _, symbol := range cs {
				if part.position.Adjacent(symbol.position) {
					part.symbol = symbol
					if !slices.Contains(symbol.parts, part) {
						symbol.parts = append(symbol.parts, part)
					}
				}
			}
		}

	})

	sumParts(parts)
	sumGearPower(symbols)

}

func sumParts(parts [][]*Part) {
	sum := 0

	for _, row := range parts {
		for _, part := range row {
			if part.symbol != nil {
				sum += part.number
			}
		}
	}

	fmt.Println(sum)
}

func sumGearPower(symbols [][]*Symbol) int {
	sum := 0

	for _, row := range symbols {
		for _, symbol := range row {

			if symbol.IsGear() {
				sum += symbol.GearRatio()
			}

		}
	}

	fmt.Println(sum)

	return sum
}

func getPreviousRow[T any](slice [][]*T) []*T {
	if len(slice) <= 0 {
		return []*T{}
	}

	return slice[len(slice)-1]
}
