package main

import (
	"common"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

func Solution1() {
	mirrors := parse()
	sum := sumReflections(mirrors, 0)
	fmt.Println(sum)
}

func Solution2() {
	mirrors := parse()
	sum := sumReflections(mirrors, 1)
	fmt.Println(sum)
}

func sumReflections(mirrors [][]string, tolerance int) int {
	sum := 0

	for _, mirror := range mirrors {
		if horizontalLine, hasSmudge := findReflectionLine(mirror, tolerance); horizontalLine != -1 && hasSmudge {
			sum += (horizontalLine + 1) * 100
		} else if verticalLine, hasSmudge := findReflectionLine(flipSlice(mirror), tolerance); verticalLine != -1 && hasSmudge {
			sum += verticalLine + 1
		}
	}

	return sum
}

func parse() [][]string {
	mirrors := [][]string{{}}

	currentMirrorIndex := 0
	common.ReadFile("day13/day13.txt", func(line string) {
		if line == "" {
			currentMirrorIndex++
			mirrors = append(mirrors, []string{})
			return
		}

		mirrors[currentMirrorIndex] = append(mirrors[currentMirrorIndex], line)
	})

	return mirrors
}

func findReflectionLine(mirror []string, tolerance int) (int, bool) {
	mirrorSize := len(mirror)
	for i := 0; i < mirrorSize-1; i++ {
		if s, hasSmudge := isReflectionLine(i, i+1, mirror, tolerance); s {
			return i, hasSmudge
		}
	}

	return -1, false
}

func isReflectionLine(leftSide, rightSide int, mirror []string, tolerance int) (bool, bool) {
	mirrorSize := len(mirror)
	foundSmudge := false

	for left, right := leftSide, rightSide; left >= 0 && right < mirrorSize; left, right = left-1, right+1 {
		if s, hasSmudge := isReflecting(mirror[left], mirror[right], tolerance); !s {
			return false, foundSmudge
		} else if hasSmudge {
			foundSmudge = hasSmudge
		}
	}

	if !foundSmudge {
		return false, foundSmudge
	}

	return true, foundSmudge
}

func isReflecting(left, right string, allowedDifference int) (bool, bool) {
	leftValue, rightValue := convertToBinary(left), convertToBinary(right)
	difference := strconv.FormatInt(leftValue^rightValue, 2)
	differences := strings.Count(difference, "1")

	return differences <= allowedDifference, differences == allowedDifference
}

func abs(value int) int {
	return int(math.Abs(float64(value)))
}

func convertToBinary(s string) int64 {
	replacer := strings.NewReplacer("#", "1", ".", "0")
	value, _ := strconv.ParseInt(replacer.Replace(s), 2, 64)
	return value
}

func countMirrorObjects(value string) (int, int) {
	return strings.Count(value, "."), strings.Count(value, "#")
}

func flipSlice(s []string) []string {
	flippedSlice := make([][]string, len(s[0]))

	for y, row := range s {
		for x, col := range strings.Split(row, "") {
			if y >= len(flippedSlice[x]) {
				flippedSlice[x] = make([]string, len(s))
			}

			flippedSlice[x][y] = col
		}
	}

	collectedSlices := []string{}
	for _, slice := range flippedSlice {
		collectedSlices = append(collectedSlices, strings.Join(slice, ""))
	}

	return collectedSlices
}
