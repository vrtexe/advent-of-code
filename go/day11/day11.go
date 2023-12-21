package main

import (
	"common"
	"fmt"
	"math"
	"slices"
	"strings"
)

func main() {
	Solution2()
}

func Solution1() {
	universeImage := parse()

	universeImage = expandUniverseHorizontally(universeImage)
	universeImage = expandUniverseVertically(universeImage)
	galaxies := findGalaxies(universeImage)
	galaxyDistances := calculateGalaxyDistances(galaxies)
	completeDistance := sumGalaxyDistances(galaxyDistances)

	fmt.Println(completeDistance)
}

func Solution2() {
	universeImage := parse()

	galaxies := findGalaxies(universeImage)
	expandingHorizontalLines := findExpandingHorizontalLines(universeImage)
	expandingVerticalLines := findExpandingVerticalLines(universeImage)

	galaxyDistances := calculateGalaxyDistancesUnexpanded(galaxies, &ExpandingLines{expandingHorizontalLines, expandingVerticalLines}, 1000000)

	completeGalaxyDistance := sumGalaxyDistances(galaxyDistances)

	fmt.Println(completeGalaxyDistance)
}

func parse() [][]string {
	universeImage := [][]string{}

	common.ReadFileLines("day11/day11.txt", func(line string, y int) {
		universeLine := strings.Split(line, "")
		universeImage = append(universeImage, universeLine)
	})

	return universeImage
}

func sumGalaxyDistances(galaxyDistances map[PositionPair]int) int {
	sum := 0
	for _, galaxyDistance := range galaxyDistances {
		sum += galaxyDistance
	}
	return sum
}

func calculateGalaxyDistances(galaxies []Position) map[PositionPair]int {
	galaxyDistances := map[PositionPair]int{}

	for galaxy, galaxies := pop(galaxies); len(galaxies) > 0; galaxy, galaxies = pop(galaxies) {
		for _, otherGalaxy := range galaxies {
			distance := galaxy.calculateDistance(otherGalaxy)
			galaxyDistances[PositionPair{galaxy, otherGalaxy}] = distance
		}
	}

	return galaxyDistances
}

func calculateGalaxyDistancesUnexpanded(galaxies []Position, lines *ExpandingLines, expansionLevel int) map[PositionPair]int {
	galaxyDistances := map[PositionPair]int{}

	for galaxy, galaxies := pop(galaxies); len(galaxies) > 0; galaxy, galaxies = pop(galaxies) {
		for _, otherGalaxy := range galaxies {
			distance := galaxy.calculateDistance(otherGalaxy)
			horizontalLinesCrossed, verticalLinesCrossed := galaxy.getCrossedExpandingLines(otherGalaxy, lines)

			galaxyDistances[PositionPair{galaxy, otherGalaxy}] = (distance - (horizontalLinesCrossed + verticalLinesCrossed)) + ((horizontalLinesCrossed + verticalLinesCrossed) * expansionLevel)
		}
	}

	return galaxyDistances
}

func findGalaxies(universeImage [][]string) []Position {
	nodes := []Position{}

	for y, row := range universeImage {
		for x, space := range row {
			if space == GALAXY {
				nodes = append(nodes, Position{x, y})
			}
		}
	}

	return nodes
}

func expandUniverseHorizontally(universeImage [][]string) [][]string {
	newUniverseImage := [][]string{}

	for _, value := range universeImage {
		if !strings.Contains(strings.Join(value, ""), GALAXY) {
			newUniverseImage = append(newUniverseImage, value, value)
		} else {
			newUniverseImage = append(newUniverseImage, value)
		}
	}

	return newUniverseImage
}

func expandUniverseVertically(universeImage [][]string) [][]string {
	newUniverseImage := [][]string{}

	for _, value := range flipSlice(universeImage) {
		if !strings.Contains(strings.Join(value, ""), GALAXY) {
			newUniverseImage = append(newUniverseImage, value, value)
		} else {
			newUniverseImage = append(newUniverseImage, value)
		}
	}

	return flipSlice(newUniverseImage)
}

func findExpandingHorizontalLines(universeImage [][]string) []int {
	expandingLines := []int{}
	for row, value := range universeImage {
		if !strings.Contains(strings.Join(value, ""), GALAXY) {
			expandingLines = append(expandingLines, row)
		}
	}

	return expandingLines
}

func findExpandingVerticalLines(universeImage [][]string) []int {
	expandingLines := []int{}
	for col, value := range flipSlice(universeImage) {
		if !strings.Contains(strings.Join(value, ""), GALAXY) {
			expandingLines = append(expandingLines, col)
		}
	}

	return expandingLines
}

func flipSlice(s [][]string) [][]string {
	flippedSlice := make([][]string, len(s[0]))

	for y, row := range s {
		for x, col := range row {
			if y >= len(flippedSlice[x]) {
				flippedSlice[x] = make([]string, len(s))
			}

			flippedSlice[x][y] = col
		}
	}

	for x, slice := range flippedSlice {
		flippedSlice[x] = slices.Clip(slice)
	}

	return slices.Clip(flippedSlice)
}

func pop[S ~[]E, E any](s S) (E, S) {
	lastIndex := len(s) - 1
	last := s[lastIndex]

	return last, s[:lastIndex]
}

const (
	GALAXY = "#"
)

type Position struct {
	x, y int
}

func (this Position) calculateDistance(position Position) int {
	return int(math.Abs(float64(position.y)-float64(this.y)) + math.Abs(float64(position.x)-float64(this.x)))
}

func (this Position) getCrossedExpandingLines(position Position, lines *ExpandingLines) (horizontalLinesCrossed, verticalLinesCrossed int) {
	crossedHorizontalLines := []int{}
	for _, line := range lines.horizontal {

		if line > min(this.y, position.y) && line < max(this.y, position.y) {
			crossedHorizontalLines = append(crossedHorizontalLines, line)
		}
	}

	crossedVerticalLines := []int{}
	for _, line := range lines.vertical {
		if line > min(this.x, position.x) && line < max(this.x, position.x) {
			crossedVerticalLines = append(crossedVerticalLines, line)
		}
	}

	return len(crossedHorizontalLines), len(crossedVerticalLines)
}

//  (4,0) -> (9,10)

type ExpandingLines struct {
	horizontal, vertical []int
}

type PositionPair struct {
	left, right Position
}
