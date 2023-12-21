package main

import (
	"common"
	"fmt"
	"math"
	"slices"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

type Position struct {
	x, y int
}

type Path struct {
	position *Position
	previous *Position
	length   int
	parent   *Path
}

func Solution2() {
	start, maze := parse()

	loop, _ := getMainLoop(start, &maze)

	loopMap := map[string]struct{}{fmt.Sprint(start.y) + fmt.Sprint(start.x): {}}

	for _, point := range loop {
		loopMap[fmt.Sprint(point.y)+fmt.Sprint(point.x)] = struct{}{}
	}

	for y, row := range maze.tiles {
		for x, tile := range strings.Split(row, "") {
			(&maze).passed[y][x] = false
			if _, exists := loopMap[fmt.Sprint(y)+fmt.Sprint(x)]; !exists && tile != GROUND {
				maze.tiles[y] = replaceCharacterAt(maze.tiles[y], x, rune(GROUND[0]))
			}
		}
	}

	expandedMaze := Field{expandMaze(start, maze), []map[int]bool{}}

	nextGroundMoves := expandedMaze.createNextGroundMover()

	currentMove := ValuePosition{}
	currentMoves := nextGroundMoves(&Position{0, 0})

	for len(currentMoves) > 0 {
		currentMove, currentMoves = pop(currentMoves)

		if currentMove.value == GROUND {
			expandedMaze.tiles[currentMove.position.y] = replaceCharacterAt(expandedMaze.tiles[currentMove.position.y], currentMove.position.x, OUTSIDE_LOOP_MARK)
		}

		currentMoves = append(currentMoves, nextGroundMoves(currentMove.position)...)
	}

	sum := 0
	for _, row := range expandedMaze.tiles {
		sum += strings.Count(row, GROUND)
	}

	// fmt.Println(strings.Join(expandedMaze.tiles, "\n"))
	fmt.Println(sum)
}

func expandMaze(start *Position, maze Field) []string {
	directionMap := DirectionMap()
	startValue := findStartValue(start, maze, &directionMap)
	maze.tiles[start.y] = replaceCharacterAt(maze.tiles[start.y], start.x, rune(startValue[0]))

	newTiles := []string{}
	horizontalEdge := EXPANSION_VERTICAL_EDGE + strings.Repeat(EXPANSION_HORIZONTAL_EDGE, (len(maze.tiles[0])*2)-1) + EXPANSION_VERTICAL_EDGE
	newTiles = append(newTiles, horizontalEdge)

	for y, row := range maze.tiles {
		nextRow := expandRowHorizontally(row, func(x int) string { return maze.getAt(&Position{x, y}) })
		if len(newTiles)-1 == 0 {
			newTiles = append(newTiles, nextRow)
			continue
		}

		middleRow := generateRowBetween(newTiles[len(newTiles)-1], nextRow)
		newTiles = append(newTiles, middleRow, nextRow)
	}

	return append(newTiles, horizontalEdge)
}

func generateRowBetween(prevRow, nextRow string) string {
	middleRow := ""

	prevRowTiles := strings.Split(prevRow, "")
	nextRowTiles := strings.Split(nextRow, "")

	for i := 0; i < len(nextRowTiles); i++ {
		prevTile := prevRowTiles[i]
		tile := nextRowTiles[i]

		if prevTile == tile && tile == EXPANSION_VERTICAL_EDGE {
			middleRow += EXPANSION_VERTICAL_EDGE
			continue
		}

		if isConnectedDown(prevTile, tile) {
			middleRow += EXPANSION_VERTICAL_CONNECTOR
		} else {
			middleRow += EXPANSION_EMPTY_SPACE
		}
	}

	return middleRow
}

func expandRowHorizontally(row string, nextTile func(x int) string) string {
	nextRow := ""
	rowTiles := strings.Split(row, "")

	for x, tile := range rowTiles {
		if x+1 >= len(rowTiles) {
			nextRow += tile + EXPANSION_VERTICAL_EDGE
			continue
		}

		if x-1 < 0 {
			nextRow += EXPANSION_VERTICAL_EDGE
		}

		rightTile := nextTile(x + 1)

		if isConnectedRight(tile, rightTile) {
			nextRow += tile + EXPANSION_HORIZONTAL_CONNECTOR
		} else {
			nextRow += tile + EXPANSION_EMPTY_SPACE
		}
	}

	return nextRow
}

func findStartValue(startPosition *Position, maze Field, directionMap *map[Direction]string) string {
	positions := maze.Next(startPosition, 0, nil)
	p1, p2 := positions[0], positions[1]
	d1, d2 := calculatePositionDirection(p1.position, startPosition), calculatePositionDirection(p2.position, startPosition)

	return (*directionMap)[d1+d2]
}

func calculatePositionDirection(from, to *Position) Direction {
	return getDirection(to.x-from.x, to.y-from.y)
}

func isConnectedRight(from, to string) bool {
	leftPossibilities := []string{HORIZONTAL_PIPE, UP_LEFT_EDGE, DOWN_LEFT_EDGE}
	rightPossibilities := []string{HORIZONTAL_PIPE, DOWN_RIGHT_EDGE, UP_RIGHT_EDGE}

	return slices.Contains(leftPossibilities, from) && slices.Contains(rightPossibilities, to)
}

func isConnectedDown(from, to string) bool {
	leftPossibilities := []string{VERTICAL_PIPE, UP_LEFT_EDGE, UP_RIGHT_EDGE}
	rightPossibilities := []string{VERTICAL_PIPE, DOWN_LEFT_EDGE, DOWN_RIGHT_EDGE}

	return slices.Contains(leftPossibilities, from) && slices.Contains(rightPossibilities, to)
}

func replaceCharacterAt(s string, i int, new rune) string {
	newValue := []rune(s)
	newValue[i] = new
	return string(newValue)
}

func Solution1() {
	start, maze := parse()

	_, count := getMainLoop(start, &maze)

	fmt.Println(count / 2)
}

func getMainLoop(start *Position, maze *Field) ([]Position, int) {
	path := map[int]map[int]int{}
	nextPositions := maze.Next(start, 0, nil)
	currentPosition := Path{}

	for {
		if len(nextPositions) <= 0 {
			break
		}

		currentPosition, nextPositions = pop(nextPositions)

		if maze.canMoveTo(currentPosition.position, start) && currentPosition.previous != start {
			break
		}

		connections := maze.Next(currentPosition.position, currentPosition.length, &Path{
			currentPosition.position, currentPosition.previous, currentPosition.length, currentPosition.parent,
		})
		nextPositions = append(nextPositions, connections...)

		if _, exists := path[currentPosition.position.y]; !exists {
			path[currentPosition.position.y] = map[int]int{}
		}
	}

	return extractLoop(&currentPosition), currentPosition.length + 1
}

func extractLoop(path *Path) []Position {
	positions := []Position{}
	for currentPosition := path; currentPosition != nil; currentPosition = currentPosition.parent {
		positions = append(positions, *currentPosition.position)
	}

	return positions
}

func parse() (*Position, Field) {
	var startingPosition *Position = nil

	maze := Field{[]string{}, []map[int]bool{}}

	common.ReadFileLines("day10/day10.txt", func(line string, index int) {
		if startingPosition == nil {
			if x := strings.IndexRune(line, START_PIPE); x != -1 {
				startingPosition = &Position{x, index}
			}
		}

		maze.tiles = append(maze.tiles, line)
		maze.passed = append(maze.passed, map[int]bool{})
	})

	return startingPosition, maze
}

func pop[S ~[]E, E any](s S) (E, S) {
	lastIndex := len(s) - 1
	last := s[lastIndex]

	return last, s[:lastIndex]
}

type Field struct {
	tiles  []string
	passed []map[int]bool
}

type ValuePosition struct {
	value    string
	position *Position
}

func (this *Field) createNextGroundMover() func(position *Position) []ValuePosition {
	passedMap := make([]map[int]bool, len(this.tiles))
	movableGroundTargets := []string{GROUND, EXPANSION_EMPTY_SPACE, EXPANSION_HORIZONTAL_EDGE, EXPANSION_VERTICAL_EDGE}

	return func(position *Position) []ValuePosition {
		flag(position, &passedMap)
		groundMoves := []ValuePosition{}

		for _, direction := range [4]Direction{UP, DOWN, LEFT, RIGHT} {
			if next := position.move(this, direction); next != nil {
				if passedMap[next.y][next.x] {
					continue
				}

				if nextValue := this.getAt(next); slices.Contains(movableGroundTargets, nextValue) {
					groundMoves = append(groundMoves, ValuePosition{nextValue, next})
				}
			}
		}

		return groundMoves
	}
}

func flag(position *Position, m *[]map[int]bool) {
	if (*m)[position.y] == nil {
		(*m)[position.y] = map[int]bool{}
	}

	(*m)[position.y][position.x] = true
}

func (this Position) move(maze *Field, direction Direction) *Position {
	switch direction {
	case UP:
		return this.moveUp(maze)
	case DOWN:
		return this.moveDown(maze)
	case LEFT:
		return this.moveLeft(maze)
	case RIGHT:
		return this.moveRight(maze)
	}

	return nil
}

func (this Position) moveUp(maze *Field) *Position {
	if this.y-1 < 0 {
		return nil
	}

	return &Position{this.x, this.y - 1}
}

func (this Position) moveDown(maze *Field) *Position {
	if this.y+1 >= len(maze.tiles) {
		return nil
	}

	return &Position{this.x, this.y + 1}
}

func (this Position) moveLeft(maze *Field) *Position {
	if this.x-1 < 0 {
		return nil
	}
	return &Position{this.x - 1, this.y}
}

func (this Position) moveRight(maze *Field) *Position {
	if this.x+1 >= len(maze.tiles[0]) {
		return nil
	}

	return &Position{this.x + 1, this.y}
}

func (this Field) canMoveTo(from, position *Position) bool {
	if xDiff, yDiff := math.Abs(float64(from.x-position.x)), math.Abs(float64(from.y-position.y)); (xDiff <= 1 && yDiff == 0) || (yDiff <= 1 && xDiff == 0) {
		byX, byY := position.x-from.x, position.y-from.y

		direction := getDirection(byX, byY)

		return canMoveInDirection(this.getAt(from), direction)
	}

	return false
}

func canMoveInDirection(pipe string, direction Direction) bool {
	switch direction {
	case UP:
		return pipe == VERTICAL_PIPE || pipe == DOWN_LEFT_EDGE || pipe == DOWN_RIGHT_EDGE
	case DOWN:
		return pipe == VERTICAL_PIPE || pipe == UP_LEFT_EDGE || pipe == UP_RIGHT_EDGE
	case RIGHT:
		return pipe == HORIZONTAL_PIPE || pipe == UP_LEFT_EDGE || pipe == DOWN_LEFT_EDGE
	case LEFT:
		return pipe == HORIZONTAL_PIPE || pipe == UP_RIGHT_EDGE || pipe == DOWN_RIGHT_EDGE
	}

	return false
}

func getDirection(toX, toY int) Direction {
	if toX < 0 {
		return LEFT
	} else if toX > 0 {
		return RIGHT
	} else if toY < 0 {
		return UP
	} else if toY > 0 {
		return DOWN
	}

	return DOWN
}

func (this *Field) Next(position *Position, count int, path *Path) []Path {
	this.passed[position.y][position.x] = true

	nextPositions := []Path{}

	for _, direction := range [4]Direction{UP, DOWN, LEFT, RIGHT} {
		if next, canMove := this.move(position, direction); canMove {
			nextPositions = append(nextPositions, Path{next, position, count + 1, path})
		}
	}

	return nextPositions
}

func (this Field) move(from *Position, direction Direction) (*Position, bool) {
	switch direction {
	case UP:
		return this.canMoveUp(from)
	case DOWN:
		return this.canMoveDown(from)
	case LEFT:
		return this.canMoveLeft(from)
	case RIGHT:
		return this.canMoveRight(from)
	}

	return nil, false
}

func (this Field) canMoveUp(from *Position) (*Position, bool) {
	if from.y-1 < 0 || this.passed[from.y-1][from.x] {
		return nil, false
	}

	position := &Position{from.x, from.y - 1}
	tile := this.getAt(position)

	return position, tile == VERTICAL_PIPE || tile == UP_RIGHT_EDGE || tile == UP_LEFT_EDGE
}

func (this Field) canMoveDown(from *Position) (*Position, bool) {
	if from.y+1 >= len(this.tiles) || this.passed[from.y+1][from.x] {
		return nil, false
	}

	position := &Position{from.x, from.y + 1}
	tile := this.getAt(position)

	return position, tile == VERTICAL_PIPE || tile == DOWN_RIGHT_EDGE || tile == DOWN_LEFT_EDGE
}

func (this Field) canMoveRight(from *Position) (*Position, bool) {
	if from.x+1 >= len(this.tiles[0]) || this.passed[from.y][from.x+1] {
		return nil, false
	}

	position := &Position{from.x + 1, from.y}
	tile := this.getAt(position)

	return position, tile == HORIZONTAL_PIPE || tile == DOWN_RIGHT_EDGE || tile == UP_RIGHT_EDGE
}

func (this Field) canMoveLeft(from *Position) (*Position, bool) {
	if from.x-1 < 0 || this.passed[from.y][from.x-1] {
		return nil, false
	}

	position := &Position{from.x - 1, from.y}
	tile := this.getAt(position)

	return position, tile == HORIZONTAL_PIPE || tile == DOWN_LEFT_EDGE || tile == UP_LEFT_EDGE
}

func (this Field) getAt(position *Position) string {
	return string(this.tiles[position.y][position.x])
}

type Direction string

const (
	UP    Direction = "UP"
	DOWN            = "DOWN"
	LEFT            = "LEFT"
	RIGHT           = "RIGHT"
)

const (
	START_PIPE = 'S'
)

const (
	VERTICAL_PIPE   = "|"
	HORIZONTAL_PIPE = "-"
	DOWN_LEFT_EDGE  = "L"
	DOWN_RIGHT_EDGE = "J"
	UP_RIGHT_EDGE   = "7"
	UP_LEFT_EDGE    = "F"
	GROUND          = "."
)

func DirectionMap() map[Direction]string {
	return map[Direction]string{
		UP + LEFT: UP_LEFT_EDGE,
		LEFT + UP: UP_LEFT_EDGE,

		UP + RIGHT: UP_RIGHT_EDGE,
		RIGHT + UP: UP_RIGHT_EDGE,

		DOWN + LEFT: DOWN_LEFT_EDGE,
		LEFT + DOWN: DOWN_LEFT_EDGE,

		DOWN + RIGHT: DOWN_RIGHT_EDGE,
		RIGHT + DOWN: DOWN_RIGHT_EDGE,

		UP + DOWN: VERTICAL_PIPE,
		DOWN + UP: VERTICAL_PIPE,

		LEFT + RIGHT: HORIZONTAL_PIPE,
		RIGHT + LEFT: HORIZONTAL_PIPE,
	}
}

const (
	EXPANSION_VERTICAL_CONNECTOR   = "#"
	EXPANSION_HORIZONTAL_CONNECTOR = "="
	EXPANSION_EMPTY_SPACE          = " "
	EXPANSION_VERTICAL_EDGE        = "!"
	EXPANSION_HORIZONTAL_EDGE      = "~"
)

const (
	OUTSIDE_LOOP_MARK = 'O'
	INSIDE_LOOP_MARK  = 'I'
)
