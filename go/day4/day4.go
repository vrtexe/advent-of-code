package main

import (
	"common"
	"fmt"
	"math"
	"strings"
)

func main() {
	parseInput()
}

type Card struct {
	id             int
	winningNumbers map[string]struct{}
	numbers        []string
	copies         int
}

func (this Card) Points() int {
	count := this.WinningNumbersCount()

	if count-1 < 0 {
		return 0
	}

	return int(math.Pow(2, float64(count-1)))
}

func (this Card) WinningNumbersCount() int {
	count := 0

	for _, number := range this.numbers {
		if _, containsKey := this.winningNumbers[number]; containsKey {
			count += 1
		}
	}

	return count
}

func parseInput() {
	Solution1()
	Solution2()
}

func Solution2() {
	cards := []*Card{}

	calculateCopies := createCopyCalculator()
	sum := 0

	common.ReadFile("../day4/day4.txt", func(line string) {
		card := parseCard(line)
		cards = append(cards, card)

		card.copies = calculateCopies(card)
		sum += card.copies + 1
	})

	fmt.Println(sum)
}

func createCopyCalculator() func(card *Card) int {
	copies := map[int]int{}

	return func(card *Card) int {
		cardsCount := card.WinningNumbersCount()
		currentCopies := extractOrDefault(copies, card.id, 0)

		for i := 1; i <= cardsCount; i++ {
			previousValue := extractOrDefault(copies, card.id+i, 0)
			copies[card.id+i] = previousValue + (currentCopies + 1)
		}

		return currentCopies
	}
}

func extractOrDefault(copies map[int]int, key int, defaultValue int) int {
	if value, present := copies[key]; !present {
		return defaultValue
	} else {
		return value
	}
}

func Solution1() {
	cards := []*Card{}

	sum := 0
	common.ReadFile("../day4/day4.txt", func(line string) {
		card := parseCard(line)
		cards = append(cards, card)

		sum += card.Points()
	})

	fmt.Println(sum)
}

func parseCard(line string) *Card {
	cardTitle, numbers := common.Split2(line, ":")

	_, cardNumber := common.Split2(cardTitle, " ")
	winningNumbersString, myNumbers := common.Split2(numbers, "|")

	return &Card{
		id:             common.ParseInt(cardNumber),
		winningNumbers: toMap(strings.Split(winningNumbersString, " ")),
		numbers:        strings.Fields(myNumbers),
		copies:         0,
	}
}

func toMap(values []string) map[string]struct{} {
	mappings := map[string]struct{}{}

	for _, value := range values {
		mappings[value] = struct{}{}
	}

	return mappings
}

func Map[S interface{}, T any](slice []S, mapper func(S) T) []T {
	newValue := []T{}

	for _, value := range slice {
		newValue = append(newValue, mapper(value))
	}

	return newValue
}

func Filter[S any](slice []S, filter func(item S) bool) []S {
	newValue := []S{}

	for _, value := range slice {
		if filter(value) {
			newValue = append(newValue, value)
		}
	}

	return newValue
}
