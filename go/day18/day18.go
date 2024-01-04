package main

import (
	"common"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	Solution2()
}

func Solution2() {

	directionMap := createDirectionMapping()

	moves := parse()
	points, boundaryPoints := findPositions(moves, func(move Move) (int, Direction) {
		lastIndex := len(move.color) - 1
		value, _ := strconv.ParseInt(string(move.color[1:lastIndex]), 16, 32)
		directionValue := common.ParseInt(string(move.color[lastIndex:]))
		direction := directionMap[directionValue]

		return int(value), direction
	})
	area := polygonArea(points)
	innerPoints := calculateInnerPoints(int(area), boundaryPoints)

	fmt.Println(boundaryPoints + innerPoints)
}

func Solution1() {
	moves := parse()
	terrain := digTerrain(moves)
	terrainMap := mapOutTerrain(terrain)
	_, transformed := filterOutside(Position{0, 0}, &terrainMap)
	terrainSize := len(terrainMap) * len(terrainMap[0])
	fmt.Println(terrainSize - transformed)
}

func filterOutside(start Position, terrain *[][]string) (*[][]string, int) {
	var cube Position
	cubes := []Position{start}
	transformed := 0

	for len(cubes) > 0 {
		cube, cubes = common.Pop(cubes)
		if (*terrain)[cube.y][cube.x] == HIGH {
			continue
		}

		transformed++
		(*terrain)[cube.y][cube.x] = HIGH

		cubes = append(cubes, nextGroundMoves(cube, *terrain)...)
	}

	return terrain, transformed
}

func nextGroundMoves(position Position, terrain [][]string) []Position {
	nextPositions := []Position{}
	for _, direction := range []Direction{UP, DOWN, LEFT, RIGHT} {
		if next := position.move(direction); next.isValid(terrain) && terrain[next.y][next.x] == GROUND {
			nextPositions = append(nextPositions, next)
		}
	}

	return nextPositions
}

func (this Position) isValid(terrain [][]string) bool {
	return this.y >= 0 && this.y < len(terrain) && this.x >= 0 && this.x < len(terrain[0])
}

func digTerrain(moves []Move) Terrain {
	terrain := Terrain{Position{0, 0}, Rows{}, TerrainSize{Count{}, Count{}}}

	for _, move := range moves {
		terrain.dig(move)
	}

	return terrain
}

func findPositions(moves []Move, extract func(Move) (int, Direction)) ([]Position, int) {
	start := Position{}
	positions := []Position{}
	pointCount := 0

	for _, move := range moves {
		distance, direction := extract(move)
		start = start.moveBy(distance, direction)

		pointCount += distance
		positions = append(positions, start)
	}

	return positions, pointCount
}

func mapOutTerrain(terrain Terrain) [][]string {
	rows, cols := terrain.size.rows, terrain.size.cols
	terrainMapping := [][]string{}

	for y := rows.start - 1; y < rows.end+2; y++ {
		nextRow := []string{}
		for x := cols.start - 1; x < cols.end+2; x++ {
			if value, exists := terrain.trench[y][x]; exists {
				nextRow = append(nextRow, value)
			} else {
				nextRow = append(nextRow, GROUND)
			}
		}

		terrainMapping = append(terrainMapping, nextRow)
	}

	return terrainMapping
}

func parse() []Move {
	moves := []Move{}
	common.ReadFile("day18/day18.txt", func(line string) {
		fields := strings.Fields(line)
		direction, value, color := fields[0], fields[1], fields[2]
		moves = append(moves, Move{
			value:     common.ParseInt(value),
			direction: Direction(direction),
			color:     Color(strings.Trim(color, "(|)"))},
		)
	})

	return moves
}

func polygonArea(positions []Position) float64 {
	area := 0
	for i, j := 0, len(positions)-1; i < len(positions); j, i = i, i+1 {
		area += (positions[j].y + positions[i].y) * (positions[j].x - positions[i].x)
	}

	return math.Abs(float64(area) / 2)
}

// A = i + b/2 - 1 | solve i
// i = A - b/2 + 1
func calculateInnerPoints(area int, boundaryPoints int) int {
	return area - (boundaryPoints / 2) + 1
}

type Terrain struct {
	position Position
	trench   Rows
	size     TerrainSize
}

type TerrainSize struct {
	rows Count
	cols Count
}

type Count struct {
	start, end int
}

type Rows map[int]Columns
type Columns map[int]string

func (this *Terrain) dig(dig Move) {
	for i := 0; i < dig.value; i++ {
		this.position = this.nextPosition(dig)
		this.recalculateSize()

		if _, exists := this.trench[this.position.y]; !exists {
			this.trench[this.position.y] = Columns{}
		}

		this.trench[this.position.y][this.position.x] = string(TRENCH)
	}
}

func (this Terrain) nextPosition(dig Move) Position {
	return this.position.move(dig.direction)
}

func (this Position) move(direction Direction) Position {
	return this.moveBy(1, direction)
}

func (this Position) moveBy(value int, direction Direction) Position {
	switch direction {
	case UP:
		return Position{this.x, this.y - value}
	case DOWN:
		return Position{this.x, this.y + value}
	case LEFT:
		return Position{this.x - value, this.y}
	case RIGHT:
		return Position{this.x + value, this.y}
	default:
		return this
	}
}

func (this *Terrain) recalculateSize() {
	this.size.rows = Count{
		start: min(this.size.rows.start, this.position.y),
		end:   max(this.size.rows.end, this.position.y),
	}

	this.size.cols = Count{
		start: min(this.size.cols.start, this.position.x),
		end:   max(this.size.cols.end, this.position.x),
	}
}

type Position struct {
	x, y int
}

type Color string

type Move struct {
	value     int
	direction Direction
	color     Color
}

type Direction string

const (
	UP    Direction = "U"
	DOWN            = "D"
	RIGHT           = "R"
	LEFT            = "L"
)

type Level string

const (
	TRENCH Level = "#"
	GROUND       = "."
	HIGH         = "+"
)

func createDirectionMapping() map[int]Direction {
	return map[int]Direction{
		0: RIGHT,
		1: DOWN,
		2: LEFT,
		3: UP,
	}
}
