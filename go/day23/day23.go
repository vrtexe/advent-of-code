package main

import (
	"common"
	"fmt"
	"slices"
	"strings"
)

type Direction = string

const (
	North Direction = "^"
	South           = "v"
	West            = "<"
	East            = ">"
)

var DirectionMap = map[Direction]Position{
	North: {x: 0, y: -1},
	South: {x: 0, y: 1},
	West:  {x: -1, y: 0},
	East:  {x: 1, y: 0},
}

var DirectionValues = []Direction{North, South, West, East}

func main() {
	// Solution1()
	Solution2()
}

func Solution1() {
	// bricks := parse()
	// _, support, supportedBy, _ := dropBricks(bricks)

	// removable := countRemovable(support, supportedBy)

	// fmt.Println(removable)
}

func Solution2() {
	data := parse()

	startPosition := findStartPosition(data)
	endPosition := findEndPosition(data)

	// positionsMap := generatePositions(data)
	// stringsMap := generatePositionStrings(positionsMap)

	distances := findLongestPath(EfficientPoint{
		position: startPosition,
		path:     positionToString(startPosition),
		weight:   1,
	}, func(point EfficientPoint) []EfficientPoint {
		nextPositions := efficientlyGetNextTiles(point.position, point.path, data)
		nextPoints := []EfficientPoint{}

		for _, nextPosition := range nextPositions {
			nextPoints = append(nextPoints, EfficientPoint{
				position: nextPosition,
				path:     fmt.Sprintf("%s-%s", point.path, positionToString(nextPosition)),
				weight:   1,
			})
		}

		return nextPoints
	}, func(v EfficientPoint) string {
		// return positionToString(v.position)
		return fmt.Sprintf("%s<-[%s]", positionToString(v.position), v.path)
	})

	endPositionString := positionToString(endPosition)

	destinations := []Distance[EfficientPoint]{}

	for key, distance := range distances {
		if strings.HasPrefix(key, endPositionString) {
			destinations = append(destinations, distance)
		}
	}

	maxDistance := slices.MaxFunc(destinations, func(left, right Distance[EfficientPoint]) int {
		return left.weight - right.weight
	})

	// currentDistance := maxDistance
	// c := 0
	// for currentDistance.through != nil {
	// 	fmt.Println(currentDistance.weight, currentDistance.vertex.position, currentDistance.through.position)
	// 	// fmt.Println(distances[positionToString(currentDistance.through.position)].vertex.position)
	// 	currentDistance = distances[positionToString(currentDistance.through.position)]
	// 	// break
	// }
	// fmt.Println(c)
	fmt.Println(maxDistance.weight)
}

func efficientlyGetNextTiles(
	position Position,
	invalid string,
	data [][]string,
) []Position {
	result := []Position{}

	for _, direction := range DirectionValues {
		nextDirection := DirectionMap[direction]
		nextPosition := addPosition(position, nextDirection)

		if !isValidPosition(nextPosition, data) || strings.Contains(invalid, positionToString(nextPosition)) || data[nextPosition.y][nextPosition.x] == "#" {
			continue
		}

		result = append(result, nextPosition)
	}

	return result
}

func positionToString(position Position) string {
	return fmt.Sprintf("(%d,%d)", position.y, position.x)
}

func extractAddedPosition(
	left Position,
	right Position,
	positionGrid [][]Position,
) *Position {
	newY := left.y + right.y
	if newY < 0 || newY >= len(positionGrid) {
		return nil
	}

	line := positionGrid[newY]
	newX := left.x + right.x

	if newX < 0 || newX >= len(line) {
		return nil
	}

	return &line[newX]
}

func isValidPosition(position Position, grid [][]string) bool {
	if position.y < 0 || position.y >= len(grid) {
		return false
	}

	if position.x < 0 || position.x >= len(grid[position.y]) {
		return false
	}

	return true
}

func addPosition(left Position, right Position) Position {
	return Position{x: left.x + right.x, y: left.y + right.y}
}

func generatePositions(data Data) [][]Position {
	positionMap := [][]Position{}

	for y, line := range data {
		mappedLine := make([]Position, len(line))
		for x := range line {
			mappedLine = append(mappedLine, Position{x: x, y: y})
		}

		positionMap = append(positionMap, mappedLine)
	}

	return positionMap
}

func generatePositionStrings(data [][]Position) [][]string {
	positionMap := [][]string{}

	for _, line := range data {
		mappedLine := make([]string, len(line))
		for _, value := range line {
			mappedLine = append(mappedLine, positionToString(value))
		}

		positionMap = append(positionMap, mappedLine)
	}

	return positionMap
}

func findStartPosition(data Data) Position {
	return Position{y: 0, x: slices.Index(data[0], ".")}
}

func findEndPosition(data Data) Position {
	return Position{y: len(data) - 1, x: slices.Index(data[len(data)-1], ".")}
}

type EfficientPoint struct {
	position Position
	path     string
	weight   int
}

func (this EfficientPoint) getWeight() int {
	return this.weight
}

type Distance[V any] struct {
	vertex  V
	through *V
	weight  int
}

type Pair[T, U any] struct {
	First  T
	Second U
}

type Vertex interface {
	getWeight() int
}

func findLongestPath[V Vertex](
	vertex V,
	getAdjacent func(v V) []V,
	keyOf func(v V) string,
) map[string]Distance[V] {
	distances := map[string]Distance[V]{keyOf(vertex): {vertex: vertex, weight: 0}}

	queue := []Pair[int, V]{{0, vertex}}

	var current Pair[int, V]
	for len(queue) > 0 {
		current, queue = common.Shift(queue)
		currentDistance := distances[keyOf(current.Second)]

		nextVertices := []Pair[int, V]{}

		for _, next := range getAdjacent(current.Second) {
			nextKey := keyOf(next)
			nextDistance, nextDistanceExists := distances[nextKey]
			nextWeight := currentDistance.weight + next.getWeight()

			if !nextDistanceExists || nextWeight > nextDistance.weight {
				secondCopy := current.Second
				distances[nextKey] = Distance[V]{
					vertex:  next,
					through: &secondCopy,
					weight:  nextWeight,
				}

				nextVertices = append(nextVertices, Pair[int, V]{nextWeight, next})
			}
		}

		queue = append(queue, nextVertices...)
	}

	return distances
}

func parse() Data {
	data := [][]string{}
	common.ReadFile("day23/day23_ex.txt", func(line string) {
		data = append(data, strings.Split(line, ""))
	})

	return data
}

type Data = [][]string
type Position struct {
	x, y int
}
