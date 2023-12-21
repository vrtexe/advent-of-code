package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Game struct {
	id    int
	draws []Draw
}

type Draw struct {
	red, green, blue int
}

type Color string

const (
	Red   Color = "red"
	Green       = "green"
	Blue        = "blue"
)

type Ball struct {
	color Color
	count int
}

type Summable interface {
	Sum() int
}

func (draw Draw) Sum() int {
	return draw.blue + draw.red + draw.green
}

func (box Box) Sum() int {
	return box.blue + box.red + box.green
}

func (box Box) Power() int {
	return box.blue * box.red * box.green
}

func main() {
	Solution1()
	Solution2()
}

func Solution2() {
	games := parse("day2/day2_2.txt")

	powerSum := 0
	for _, game := range games {
		box := findMinBox(game)
		powerSum += box.Power()
	}

	fmt.Println(powerSum)
}

func Solution1() {
	games := parse("day2/day2.txt")
	gameIdsSum := sumValidGameIds(games)
	fmt.Println(gameIdsSum)
}

type Box struct {
	red, green, blue int
}

type Checker func(Draw) bool

func sumValidGameIds(games []Game) int {
	sum := 0
	handlePossibleGames(games, func(g Game) {
		sum += g.id
	})

	return sum
}

func findMinBox(game Game) Box {
	red := findMaxColor(game.draws, func(d Draw) int { return d.red })
	green := findMaxColor(game.draws, func(d Draw) int { return d.green })
	blue := findMaxColor(game.draws, func(d Draw) int { return d.blue })

	return Box{red, green, blue}
}

func findMaxColor(draws []Draw, color func(Draw) int) int {
	colorFinder := createMaxFinder(color)
	return color(slices.MaxFunc(draws, colorFinder))
}

func createMaxFinder(color func(Draw) int) func(a, b Draw) int {
	return func(a, b Draw) int {
		return intCompare(color(a), color(b))
	}
}

func handlePossibleGames(games []Game, handle func(Game)) {
	checker := createChecker(Box{12, 13, 14})

	for _, game := range games {
		if isGamePossible(game, checker) {
			handle(game)
		}
	}
}

func isGamePossible(game Game, isPossible Checker) bool {
	for _, draw := range game.draws {
		if !isPossible(draw) {
			return false
		}
	}

	return true
}

func createChecker(box Box) Checker {
	return func(draw Draw) bool {

		if draw.blue > box.blue {
			return false
		}

		if draw.green > box.green {
			return false
		}

		if draw.red > box.red {
			return false
		}

		if draw.Sum() == box.Sum() {

		}

		return true
	}
}

func parse(path string) []Game {
	games := []Game{}
	readFile(path, func(line string) {
		games = append(games, parseGame(line))
	})

	return games
}

func parseGame(line string) Game {
	game, draws := split2(line, ":")
	_, id := split2(game, " ")

	parsedDraws := parseDraws(draws)

	return Game{parseInt(id), parsedDraws}
}

func parseDraws(line string) []Draw {
	draws := strings.Split(line, ";")

	result := []Draw{}

	for _, value := range draws {
		red, green, blue := parseBalls(value)
		result = append(result, Draw{red, green, blue})
	}

	return result
}

func parseBalls(line string) (r, g, b int) {
	trimmedLine := strings.Trim(line, " ")
	balls := strings.Split(trimmedLine, ",")

	red := 0
	green := 0
	blue := 0

	for _, value := range balls {
		count, color := split2(strings.Trim(value, " "), " ")
		countInt := parseInt(count)
		switch Color(color) {
		case Red:
			red = countInt
		case Blue:
			blue = countInt
		case Green:
			green = countInt
		}
	}

	return red, green, blue
}

func split2(s string, sep string) (s1, s2 string) {
	slices := strings.Split(s, sep)
	return strings.Trim(slices[0], " "), strings.Trim(slices[1], " ")
}

func parseInt(number string) int {
	if i, err := strconv.Atoi(number); err == nil {
		return i
	} else {
		panic(err)
	}
}

func readFile(path string, handle func(string)) {
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

func intCompare(a, b int) int {
	if a == b {
		return 0
	} else if a > b {
		return 1
	} else {
		return -1
	}
}
