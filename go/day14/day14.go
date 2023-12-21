package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	// Solution1()
	Solution2()
	// rocks := parse()
	// printRocks(rocks)
	// fmt.Println()
	// printRocks(flipWest(flipWest(flipWest(rocks))))
}

func Solution1() {
	rocks := parse()

	tiltedRocks := tiltLever(rocks)

	totalLoad := calculateLoad(tiltedRocks)

	fmt.Println(totalLoad)
}

func Solution2() {
	steps := 1000000000
	rocks := parse()
	start, end, cache := findCycle(rocks, steps)

	lastStep := ((steps - start) % (end - start)) + start - 1

	if step, exists := cache[lastStep]; exists {
		fmt.Println(calculateLoad(step))
	}
}

func findCycle(rocks [][]string, steps int) (int, int, map[int][][]string) {
	cache := map[string]int{}
	stepCache := map[int][][]string{}

	for currentCycle := 0; currentCycle < steps; currentCycle++ {
		rocks = spinCycle(rocks)

		joinedRocks := joinRocks(rocks)
		if previous, exists := cache[joinedRocks]; exists {
			return previous, currentCycle, stepCache

		} else {
			cache[joinedRocks] = currentCycle
			stepCache[currentCycle] = rocks
		}
	}

	return -1, -1, stepCache
}

func spinCycle(rocks [][]string) [][]string {
	for i := 0; i < 4; i++ {
		rocks = tiltLever(rocks)
		rocks = flipEast(rocks)
	}

	return rocks
}

func createDirectionExtractor(directions [4]Direction) func(int) (int, Direction) {
	return func(current int) (int, Direction) {
		next := (current + 1) % 4
		return next, directions[next]
	}
}

func createInverterCreator() func(Direction) func([][]string) [][]string {
	currentDirection := NORTH
	westMap := map[Direction]Direction{
		NORTH: WEST,
		WEST:  SOUTH,
		SOUTH: EAST,
		EAST:  NORTH,
	}
	eastMap := map[Direction]Direction{
		NORTH: EAST,
		EAST:  SOUTH,
		SOUTH: WEST,
		WEST:  NORTH,
	}

	return func(direction Direction) func([][]string) [][]string {
		if direction == WEST {
			currentDirection = westMap[currentDirection]
		}

		if direction == EAST {
			currentDirection = eastMap[currentDirection]
		}

		if currentDirection == SOUTH {
			return func(s [][]string) [][]string {
				return flipEast(flipEast(s))
			}
		}

		if currentDirection == WEST {
			return flipEast
		}

		if currentDirection == EAST {
			return flipWest
		}

		return func(s [][]string) [][]string {
			return s
		}
	}
}

func calculateLoad(rocks [][]string) int {
	totalRockLevel := len(rocks)

	sum := 0
	for elevation, row := range rocks {
		for _, space := range row {
			if space == ROLLING_STONE {
				sum += totalRockLevel - elevation
			}
		}
	}

	return sum
}

func printRocks(rocks [][]string) {
	for _, row := range rocks {
		fmt.Println(strings.Join(row, ""))
	}
}

func joinRocks(rocks [][]string) string {
	joinedRocks := []string{}

	for _, row := range rocks {
		joinedRocks = append(joinedRocks, strings.Join(row, ""))
	}

	return strings.Join(joinedRocks, "\n")
}

func parse() [][]string {
	rocks := [][]string{}

	common.ReadFile("day14/day14.txt", func(line string) {
		rocks = append(rocks, strings.Split(line, ""))
	})

	return rocks
}

func splitRocks(s string) [][]string {
	rocks := [][]string{}

	for _, part := range strings.Split(s, "\n") {
		rocks = append(rocks, strings.Split(part, ""))
	}

	return rocks
}

func tiltLever(rocks [][]string) [][]string {
	newMap := [][]string{}
	emptySpaceMap := map[int]int{}

	for y := 0; y < len(rocks); y++ {
		row := rocks[y]
		newMap = append(newMap, make([]string, len(row)))

		for x := 0; x < len(row); x++ {
			space := row[x]

			if space == EMPTY_SPACE {
				if value, exists := emptySpaceMap[x]; (exists && y <= value) || !exists {
					emptySpaceMap[x] = y
				}
				newMap[y][x] = space
			} else if space == ROLLING_STONE {
				if value, exists := emptySpaceMap[x]; exists {
					emptySpaceMap[x] = value + 1

					newMap[value][x] = space
					newMap[y][x] = EMPTY_SPACE
				} else {
					newMap[y][x] = space
				}
			} else if space == CUBED_STONE {
				delete(emptySpaceMap, x)
				newMap[y][x] = space
			}
		}
	}

	return newMap
}

// func tiltSouth(rocks [][]string) [][]string {
// 	t := tilt(rocks)

// 	return tiltSouthDir(len(rocks), len(rocks[0]), t)
// }

// func tilt(rocks [][]string) func(x, y int) [][]string {
// 	newMap := [][]string{}
// 	emptySpaceMap := map[int]int{}

// 	return func(x, y int) [][]string {
// 		if y == 0 {
// 			newMap = append(newMap, make([]string, len(rocks[0])))
// 		}
// 		space := rocks[y][x]

// 		if space == EMPTY_SPACE {
// 			if value, exists := emptySpaceMap[x]; (exists && y <= value) || !exists {
// 				emptySpaceMap[x] = y
// 			}
// 			newMap[y][x] = space
// 		} else if space == ROLLING_STONE {
// 			if value, exists := emptySpaceMap[x]; exists {
// 				emptySpaceMap[x] = value + 1

// 				newMap[value][x] = space
// 				newMap[y][x] = EMPTY_SPACE
// 			} else {
// 				newMap[y][x] = space
// 			}
// 		} else if space == CUBED_STONE {
// 			delete(emptySpaceMap, x)
// 			newMap[y][x] = space
// 		}

// 		return newMap
// 	}

// }

// func tiltNorth(rocks [][]string) [][]string {
// 	t := tilt(rocks)

// 	return tiltNorthDir(len(rocks), len(rocks[0]), t)
// }

// func tiltRight(rows int, cols int, operate func(x, y int)) {
// 	for y := 0; y < rows; y++ {
// 		for x := cols - 1; x > 0; x++ {
// 			operate(x, y)
// 		}
// 	}
// }

// func tiltNorthDir(rows int, cols int, operate func(x, y int) [][]string) [][]string {
// 	newMap := [][]string{}

// 	for y := 0; y < rows; y++ {
// 		for x := 0; x < cols; x++ {
// 			newMap = operate(x, y)
// 		}
// 	}

// 	return newMap
// }

// func tiltSouthDir(rows int, cols int, operate func(x, y int) [][]string) [][]string {
// 	newMap := [][]string{}

// 	for y := rows; y > 0; y-- {
// 		for x := 0; x < cols; x++ {
// 			newMap = operate(x, y)
// 		}
// 	}

// 	return newMap
// }

func flipUpsideDown(rocks [][]string) [][]string {
	newRocks := [][]string{}
	lastRowIndex := len(rocks) - 1
	for y, row := lastRowIndex, rocks[lastRowIndex]; y > 0; y, row = y-1, rocks[y-1] {
		newRocks = append(newRocks, row)
	}

	return newRocks
}

func flipWest(s [][]string) [][]string {
	flippedSlice := make([][]string, len(s[0]))

	for y, row := range s {
		for x, col := range row {
			currentRow := (len(row) - 1) - x
			currentColumn := y
			if y >= len(flippedSlice[currentRow]) {
				flippedSlice[currentRow] = make([]string, len(s))
			}

			flippedSlice[currentRow][currentColumn] = col
		}
	}

	return flippedSlice
}

func flipEast(s [][]string) [][]string {
	flippedSlice := make([][]string, len(s[0]))

	for y, row := range s {
		for x, col := range row {
			currentRow := x
			currentColumn := (len(s) - 1) - y
			if y >= len(flippedSlice[currentRow]) {
				flippedSlice[currentRow] = make([]string, len(s))
			}

			flippedSlice[currentRow][currentColumn] = col
		}
	}

	return flippedSlice
}

func collectColumn(column int, rocks [][]string) []string {
	value := []string{}

	for y := 0; y < len(rocks); y++ {
		value = append(value, rocks[y][column])
	}

	return value
}

const (
	ROLLING_STONE = "O"
	CUBED_STONE   = "#"
	EMPTY_SPACE   = "."
)

type Direction string

const (
	NORTH Direction = "NORTH"
	WEST            = "WEST"
	SOUTH           = "SOUTH"
	EAST            = "EAST"
)
