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

func Solution1() {
	data := parse()
	startingPosition := findStart(data)

	result := walkThroughMap(startingPosition, []int{64}, data)

	fmt.Println("Result:", result)
}

func Solution2() {

	data := parse()
	startingPosition := findStart(data)

	boardSize := len(data)
	stepsTakenToFillSingleDiamond := 65
	steps := 26501365
	target := (steps - stepsTakenToFillSingleDiamond) / boardSize

	depths := generateDepths(3, stepsTakenToFillSingleDiamond, boardSize)
	depthValues := extractRightValues(depths)

	results := walkThroughMap(startingPosition, depthValues, data)

	// fmt.Println(results) // [3877, 34674, 96159]

	points := []Pair[int, int]{}
	for index, value := range extractLeftValues(depths) {
		points = append(points, Pair[int, int]{value, results[index]})
	}

	f := findPolynomialCramer(points)

	result := f(target)

	fmt.Println("Result:", result)
}

func parse() [][]string {
	data := [][]string{}

	common.ReadFile("day21/day21.txt", func(line string) {
		data = append(data, strings.Split(line, ""))
	})

	return data
}

func findPolynomialCramer(points []Pair[int, int]) func(x int) int {
	A := [][]float64{}

	for _, point := range points {
		A = append(A, []float64{math.Pow(float64(point.left), 2), float64(point.left), 1})
	}

	B := []float64{}

	for _, value := range extractRightValues(points) {
		B = append(B, float64(value))
	}

	values := []int{}
	for _, v := range solveSystem(B, A) {
		values = append(values, int(math.Round(v)))
	}

	a, b, c := values[0], values[1], values[2]

	return func(x int) int {
		return a*common.IntPow(x, 2) + b*x + c
	}
}

func solveSystem(results []float64, matrix [][]float64) []float64 {
	nominators := generateCramerNominators(results, matrix)
	denominator := common.Det(matrix)
	solutions := []float64{}

	for _, nominator := range nominators {
		solutions = append(solutions, common.Det(nominator)/denominator)
	}

	return solutions
}

func generateCramerNominators(column []float64, matrix [][]float64) [][][]float64 {
	result := [][][]float64{}

	for i := 0; i < len(matrix); i++ {
		result = append(result, replaceMatrixColumn(i, column, matrix))
	}

	return result
}

func replaceMatrixColumn(replaceIndex int, column []float64, matrix [][]float64) [][]float64 {
	result := [][]float64{}
	for rowIndex := 0; rowIndex < len(matrix); rowIndex++ {
		row := []float64{}
		for colIndex := 0; colIndex < len(matrix[rowIndex]); colIndex++ {
			if colIndex == replaceIndex {
				row = append(row, column[rowIndex])
			} else {
				row = append(row, matrix[rowIndex][colIndex])
			}
		}
		result = append(result, row)
	}

	return result
}

type Pair[L, R any] struct {
	left  L
	right R
}

func generateDepths(count int, singleFill int, boardSize int) []Pair[int, int] {
	result := []Pair[int, int]{}

	for depth := 0; depth < count; depth++ {
		result = append(result, Pair[int, int]{depth, singleFill + boardSize*depth})
	}

	return result
}

func extractLeftValues[L, R any](pairs []Pair[L, R]) []L {
	result := []L{}

	for _, v := range pairs {
		result = append(result, v.left)
	}

	return result
}

func extractRightValues[L, R any](pairs []Pair[L, R]) []R {
	result := []R{}

	for _, v := range pairs {
		result = append(result, v.right)
	}

	return result
}

func walkThroughMap(start Position, depths []int, data [][]string) []int {
	positions := map[string]struct{}{start.toString(): {}}
	results := []int{}

	sortedDepths := append([]int{}, depths...)
	slices.Sort(sortedDepths)

	maxDepth := sortedDepths[len(sortedDepths)-1]
	sortedDepthsSize := len(sortedDepths)

	for i := 0; i < maxDepth; i++ {
		if sortedDepthsSize > 0 && i == sortedDepths[0] {
			results = append(results, len(positions))
			sortedDepths = sortedDepths[1:sortedDepthsSize]
			sortedDepthsSize = len(sortedDepths)
		}

		nextPositions := []Position{}

		for position := range positions {
			for _, p := range getNextPositions(parsePosition(position), data) {
				nextPositions = append(nextPositions, p)
			}
		}

		positions = map[string]struct{}{}
		for _, position := range nextPositions {
			positions[position.toString()] = struct{}{}
		}
	}

	results = append(results, len(positions))

	return results
}

func parsePosition(value string) Position {
	x, y := common.Split2(value, ",")
	return Position{x: common.ParseInt(x), y: common.ParseInt(y)}
}

func getNextPositions(position Position, data [][]string) []Position {
	result := []Position{}
	rows, cols := len(data), len(data[0])

	for _, direction := range Directions {
		nextPosition := moveInDirection(position, direction)
		yb, xb := adjustToBoundary(nextPosition.y, rows), adjustToBoundary(nextPosition.x, cols)

		tile := data[yb][xb]

		if tile == string(Rock) {
			continue
		}

		result = append(result, nextPosition)
	}

	return result
}

func adjustToBoundary(value int, max int) int {
	adjustedValue := value % max
	if adjustedValue < 0 {
		return max + adjustedValue
	}

	return adjustedValue
}

func moveInDirection(position Position, direction Direction) Position {
	switch direction {
	case North:
		return Position{x: position.x, y: position.y - 1}
	case South:
		return Position{x: position.x, y: position.y + 1}
	case West:
		return Position{x: position.x - 1, y: position.y}
	case East:
		return Position{x: position.x + 1, y: position.y}
	}

	panic("Invalid direction")
}

func sortAscending(left, right int) int {
	return left - right
}

func findStart(data [][]string) Position {
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			if data[i][j] == Start {
				return Position{x: j, y: i}
			}
		}
	}

	panic("Could not find starting position")
}

type Position struct {
	x, y int
}

func (it Position) toString() string {
	return fmt.Sprintf("%d,%d", it.x, it.y)
}

var Directions = []Direction{North, South, East, West}

type Direction string

const (
	North Direction = "N"
	South           = "S"
	East            = "E"
	West            = "W"
)

type Tile string

const (
	Rock  Tile = "#"
	Plot       = "."
	Start      = "S"
)
