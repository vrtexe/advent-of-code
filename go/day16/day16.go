package main

import (
	"common"
	"fmt"
	"slices"
)

func main() {
	Solution1()
	Solution2()
}

func Solution2() {
	grid := parse()

	startingTiles := extractStartingTiles(grid)
	energyMaps := []EnergizedMap{}

	for _, start := range startingTiles {
		energizedTiles := energizeTiles(start, grid)
		energyMaps = append(energyMaps, EnergizedMap{start, energizedTiles, len(energizedTiles)})
	}

	mostEnergizedMap := slices.MaxFunc(energyMaps, func(a, b EnergizedMap) int {
		return common.IntCompare(a.size, b.size)
	})

	fmt.Println(mostEnergizedMap.size)
}

func extractStartingTiles(grid []string) []Movement {
	moves := []Movement{}

	gridRows := len(grid)
	gridCols := len(grid[0])

	lastRow := gridRows - 1
	lastCol := gridCols - 1

	for i := 0; i < gridCols; i++ {
		moves = append(moves, Movement{Position{i, 0}, DOWN})
		moves = append(moves, Movement{Position{i, lastRow}, UP})
	}

	for i := 0; i < gridRows; i++ {
		moves = append(moves, Movement{Position{0, i}, RIGHT})
		moves = append(moves, Movement{Position{lastCol, i}, LEFT})
	}

	return moves
}

func Solution1() {
	grid := parse()
	start := Movement{Position{0, 0}, RIGHT}

	energized := energizeTiles(start, grid)

	// for _, g := range grid {
	// 	fmt.Println(g)
	// }

	// fmt.Println()
	// for p := range cached {
	// 	r := []rune(grid[p.position.y])
	// 	if r[p.position.x] == []rune(string(EMPTY))[0] {
	// 		r[p.position.x] = []rune(string(p.direction))[0]
	// 	}
	// 	r[p.position.x] = '#'
	// 	grid[p.position.y] = string(r)
	// }

	// for _, g := range grid {
	// 	fmt.Println(g)
	// }

	fmt.Println(len(energized))
}

func energizeTiles(start Movement, grid []string) map[Position]struct{} {
	current := Movement{}
	positions := []Movement{start}
	cached := map[Movement]struct{}{}
	energized := map[Position]struct{}{}

	for len(positions) > 0 {
		current, positions = common.Pop(positions)

		if _, exists := cached[current]; exists {
			continue
		}

		cached[current] = struct{}{}
		energized[current.position] = struct{}{}
		positions = append(positions, nextMovement(current, grid)...)
	}

	return energized
}

func nextMovement(movement Movement, grid []string) []Movement {
	value := string(grid[movement.position.y][movement.position.x])

	switch GridTile(value) {
	case EMPTY:
		return handleEmptyMovement(movement, grid)
	case VERTICAL_SPLITTER:
		movementMap := verticalSplitterMovement()
		return handleSplitterMovement(movement, movementMap, grid)
	case HORIZONTAL_SPLITTER:
		movementMap := horizontalSplitterMovement()
		return handleSplitterMovement(movement, movementMap, grid)
	case RIGHT_MIRROR:
		rightMirrorMovements := rightMirrorMovement()
		return handleMirrorMovement(movement, rightMirrorMovements, grid)
	case LEFT_MIRROR:
		leftMirrorMovements := leftMirrorMovement()
		return handleMirrorMovement(movement, leftMirrorMovements, grid)
	}

	return []Movement{}
}

func leftMirrorMovement() map[Direction]Direction {
	return map[Direction]Direction{
		DOWN:  RIGHT,
		UP:    LEFT,
		LEFT:  UP,
		RIGHT: DOWN,
	}
}

func rightMirrorMovement() map[Direction]Direction {
	return map[Direction]Direction{
		DOWN:  LEFT,
		UP:    RIGHT,
		LEFT:  DOWN,
		RIGHT: UP,
	}
}

func handleEmptyMovement(movement Movement, grid []string) []Movement {
	if next, valid := nextPosition(movement.direction, movement.position, grid); valid {
		return []Movement{{next, movement.direction}}
	}

	return []Movement{}
}

func handleSplitterMovement(movement Movement, movementMap map[Direction][]Direction, grid []string) []Movement {
	moves := []Movement{}
	for _, direction := range movementMap[movement.direction] {
		if next, valid := nextPosition(direction, movement.position, grid); valid {
			moves = append(moves, Movement{next, direction})
		}
	}

	return moves
}

func handleMirrorMovement(movement Movement, movementMap map[Direction]Direction, grid []string) []Movement {
	nextDirection := movementMap[movement.direction]
	if next, valid := nextPosition(nextDirection, movement.position, grid); valid {
		return []Movement{{next, nextDirection}}
	}

	return []Movement{}
}

func horizontalSplitterMovement() map[Direction][]Direction {
	return map[Direction][]Direction{
		DOWN:  {RIGHT, LEFT},
		UP:    {RIGHT, LEFT},
		LEFT:  {LEFT},
		RIGHT: {RIGHT},
	}
}

func verticalSplitterMovement() map[Direction][]Direction {
	return map[Direction][]Direction{
		DOWN:  {DOWN},
		UP:    {UP},
		LEFT:  {UP, DOWN},
		RIGHT: {UP, DOWN},
	}
}

func nextPosition(direction Direction, position Position, grid []string) (Position, bool) {
	switch direction {
	case UP:
		return Position{position.x, position.y - 1}, position.y-1 >= 0
	case DOWN:
		return Position{position.x, position.y + 1}, position.y+1 < len(grid)
	case LEFT:
		return Position{position.x - 1, position.y}, position.x-1 >= 0
	case RIGHT:
		return Position{position.x + 1, position.y}, position.x+1 < len(grid[0])
	default:
		panic("No direction provided")
	}
}

func parse() []string {
	grid := []string{}

	common.ReadFile("day16/day16.txt", func(line string) {
		grid = append(grid, line)
	})

	return grid
}

type EnergizedMap struct {
	start     Movement
	energized map[Position]struct{}
	size      int
}

type Movement struct {
	position  Position
	direction Direction
}

type Position struct {
	x, y int
}

type GridTile string

const (
	RIGHT_MIRROR        GridTile = "/"
	LEFT_MIRROR         GridTile = "\\"
	VERTICAL_SPLITTER   GridTile = "|"
	HORIZONTAL_SPLITTER GridTile = "-"
	EMPTY               GridTile = "."
)

type Direction string

const (
	UP    Direction = "^"
	DOWN  Direction = "v"
	LEFT  Direction = "<"
	RIGHT Direction = ">"
)
