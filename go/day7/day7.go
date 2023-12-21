package main

import (
	"common"
	"fmt"
	"slices"
	"strings"
)

func main() {
	Solution2()
}

func Solution2() {
	game := Game{[]Hand{}, &CardWeight{}, &Joker{"J", 1}}
	game.rankings.game = &game

	common.ReadFile("day7/day7.txt", func(line string) {
		cards, bid := parseHand(line)
		game.Add(cards, bid)
	})

	slices.SortFunc(game.hands, func(a, b Hand) int {
		return common.IntCompare(a.rank, b.rank)
	})

	points := 0
	for index, hand := range game.hands {
		points += (index + 1) * hand.bid
	}

	fmt.Println(points)
}

func parse() {
	game := Game{[]Hand{}, &CardWeight{}, nil}

	common.ReadFile("day7/day7.txt", func(line string) {
		cards, bid := parseHand(line)
		game.Add(cards, bid)
	})

	slices.SortFunc(game.hands, func(a, b Hand) int {
		return common.IntCompare(a.rank, b.rank)
	})

	points := 0
	for index, hand := range game.hands {
		points += (index + 1) * hand.bid
	}

	fmt.Println(points)
}

func parseHand(line string) ([]string, int) {
	hand, bid := common.Split2(line, " ")
	return strings.Split(hand, ""), common.ParseInt(bid)
}

type Game struct {
	hands    []Hand
	rankings *CardWeight
	joker    *Joker
}

type Joker struct {
	value string
	rank  int
}

func (this *Game) Add(cards []string, bid int) {
	hand := Hand{cards, bid, 0, this}
	hand.Rank()

	this.hands = append(this.hands, hand)
}

type Hand struct {
	cards []string
	bid   int
	rank  int
	game  *Game
}

func (this Hand) getWeights() *CardWeight {
	return this.game.rankings
}

func (this *Hand) Rank() int {
	combo := this.getHighestCombo()
	comboRank := this.getWeights().GetComboRankings()[combo]

	cardsValue := this.findCardsValue()

	this.rank = comboRank + cardsValue
	return this.rank
}

func (this Hand) findCardsValue() int {
	cardCount := len(this.cards)
	sum := 0

	for index, card := range this.cards {
		cardWeight := this.getWeights().GetCardRankings()[cardCount-index]
		sum += this.getWeights().GetRanking()[card] * cardWeight
	}

	return sum
}

func (this Hand) getHighestCombo() string {
	countMap := map[string]int{}
	calculateCombo := createComboCalculator(this.game)
	for _, card := range this.cards {
		countMap[card]++
		calculateCombo(card)
	}

	maxValue, secondMaxValue, jokers := calculateCombo("")

	return this.getWeights().GetComboMap()[maxValue+jokers][secondMaxValue]
}

func createComboCalculator(game *Game) func(string) (int, int, int) {
	countMap := map[string]int{}

	maxCard := Card{}
	secondMaxCard := Card{}

	return func(card string) (int, int, int) {
		jokerCount := 0
		if game.joker != nil {
			jokerCount = countMap[game.joker.value]
		}

		if card == "" {
			return maxCard.count, secondMaxCard.count, jokerCount
		}

		countMap[card]++

		if game.joker != nil && card == game.joker.value {
			return -1, -1, jokerCount
		}

		if card == maxCard.value { // Card is already set as max
			maxCard.count = countMap[card]
		} else if secondMaxCard.value == card && countMap[card] > maxCard.count { // Card is set as second max, but should be max
			secondMaxCard.count = countMap[card]
			maxCard, secondMaxCard = secondMaxCard, maxCard
		} else if countMap[card] > maxCard.count { // Card not set to anything but should be max, and max drops to second max
			secondMaxCard.value = maxCard.value
			secondMaxCard.count = maxCard.count
			maxCard.count = countMap[card]
			maxCard.value = card
		} else if secondMaxCard.value == card { // Card is same as second max (update count)
			secondMaxCard.count = countMap[card]
		} else if countMap[card] > secondMaxCard.count { // Card can not fit in to max, but is greater that second max
			secondMaxCard.count = countMap[card]
			secondMaxCard.value = card
		}

		return -1, -1, jokerCount
	}
}

type Card struct {
	value string
	count int
}

type CardWeight struct {
	comboMap      *[6][6]string
	rankings      *map[string]int
	comboRankings *map[string]int
	cardRankings  *map[int]int
	game          *Game
}

func (this *CardWeight) GetComboMap() [6][6]string {
	if this.comboMap != nil {
		return *this.comboMap
	}

	this.comboMap = &[6][6]string{
		{"##", "##", "##", "##", "##", "##"}, // 0
		{"##", "HC", "1P", "TK", "4K", "##"}, // 1
		{"##", "1P", "2P", "FH", "##", "##"}, // 2
		{"##", "TK", "FH", "##", "##", "##"}, // 3
		{"##", "4K", "##", "##", "##", "##"}, // 4
		{"5K", "##", "##", "##", "##", "##"}, // 5
	}

	return *this.comboMap
}

func (this *CardWeight) GetRanking() map[string]int {
	if this.rankings != nil {
		return *this.rankings
	}

	jValue := 11
	if this.game != nil && this.game.joker != nil {
		jValue = this.game.joker.rank
	}

	this.rankings = &map[string]int{
		"A": 14,
		"K": 13,
		"Q": 12,
		"J": jValue,
		"T": 10,
		"9": 9,
		"8": 8,
		"7": 7,
		"6": 6,
		"5": 5,
		"4": 4,
		"3": 3,
		"2": 2,
	}

	return *this.rankings
}

func (this *CardWeight) GetComboRankings() map[string]int {
	if this.comboRankings != nil {
		return *this.comboRankings
	}

	this.comboRankings = &map[string]int{
		"5K": 100_000_000,
		"4K": 75_000_000,
		"FH": 60_000_000,
		"TK": 45_000_000,
		"2P": 30_000_000,
		"1P": 15_000_000,
		"HC": 0,
	}

	return *this.comboRankings
}

func (this *CardWeight) GetCardRankings() map[int]int {
	if this.cardRankings != nil {
		return *this.cardRankings
	}

	this.cardRankings = &map[int]int{
		5: 1_000_000,
		4: 50_000,
		3: 1_000,
		2: 20,
		1: 1,
	}

	return *this.cardRankings
}
