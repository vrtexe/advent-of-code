package main

import (
	"common"
	"fmt"
	"math"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

func Solution1() {
	lines := []string{}
	common.ReadFile("day6/day6.txt", func(line string) {
		lines = append(lines, line)
	})

	_, timeData := common.Split2(lines[0], ":")
	times := strings.Fields(timeData)

	_, distanceData := common.Split2(lines[1], ":")
	distances := strings.Fields(distanceData)

	prod := 1
	for i := 0; i < len(times); i++ {
		race := Race{
			time:     common.ParseInt(times[i]),
			distance: common.ParseInt(distances[i]),
		}

		minValue, maxValue := findPoints(race)
		possibilities := (maxValue - minValue) + 1

		prod *= possibilities
	}

	fmt.Println(prod)
}

func Solution2() {
	lines := []string{}
	common.ReadFile("day6/day6.txt", func(line string) {
		lines = append(lines, line)
	})

	_, timeData := common.Split2(lines[0], ":")
	time := strings.Join(strings.Fields(timeData), "")

	_, distanceData := common.Split2(lines[1], ":")
	distance := strings.Join(strings.Fields(distanceData), "")

	race := Race{
		time:     common.ParseInt(time),
		distance: common.ParseInt(distance),
	}

	minValue, maxValue := findPoints(race)
	possibilities := (maxValue - minValue) + 1

	fmt.Println(possibilities)
}

func findPoints(race Race) (start, end int) {
	startingValue := calculateNextValue(0, race.time)

	currentMinValue, currentMaxValue := createDistance(startingValue, &race), createDistance(startingValue, &race)

	minValue, maxValue := 0, race.time
	currentMinLowerLimit, currentMinUpperLimit := minValue, maxValue
	currentMaxLowerLimit, currentMaxUpperLimit := minValue, maxValue

	for minFound, maxFound := false, false; !minFound || !maxFound; minFound, maxFound = areWinnablePointsFound(currentMinValue, currentMaxValue, race) {
		if !minFound {
			if race.isWinnable(currentMinValue.distance) {
				currentMinUpperLimit = currentMinValue.value
			}

			if race.isNotWinnable(currentMinValue.distance) {
				currentMinLowerLimit = currentMinValue.value
			}

			currentMinValue.updateValue(calculateNextValue(currentMinLowerLimit, currentMinUpperLimit))
		}

		if !maxFound {
			if race.isWinnable(currentMaxValue.distance) {
				currentMaxLowerLimit = currentMaxValue.value
			}

			if race.isNotWinnable(currentMaxValue.distance) {
				currentMaxUpperLimit = currentMaxValue.value
			}

			currentMaxValue.updateValue(calculateNextValue(currentMaxLowerLimit, currentMaxUpperLimit))
		}

	}

	fmt.Println(currentMinValue, currentMaxValue)

	return currentMinValue.value, currentMaxValue.value
}

func calculateNextValue(from, to int) int {
	return int(math.Round(float64(to-from)/2)) + from
}

func areWinnablePointsFound(minValue, maxValue Distance, race Race) (minFound, maxFound bool) {
	return isMinValue(minValue, race), isMaxValue(maxValue, race)
}

func isMinValue(value Distance, race Race) bool {
	return race.isWinnable(value.distance) && race.isNotWinnable(value.prevDistance())
}

func isMaxValue(value Distance, race Race) bool {
	return race.isWinnable(value.distance) && race.isNotWinnable(value.nextDistance())
}

type Race struct {
	time, distance int
}

func (this Race) isWinnable(value int) bool {
	return value > this.distance
}

func (this Race) isNotWinnable(value int) bool {
	return !this.isWinnable(value)
}

func (this Race) measureDistance(holdTime int) int {
	return (this.time - holdTime) * holdTime
}

type Distance struct {
	value    int
	distance int
	race     *Race
}

func (this *Distance) updateValue(newValue int) {
	this.value = newValue
	this.recalculateDistance()
}

func (this *Distance) recalculateDistance() {
	this.distance = this.currentDistance()
}

func (this Distance) currentDistance() int {
	return this.race.measureDistance(this.value)
}

func (this Distance) prevDistance() int {
	return this.race.measureDistance(this.value - 1)
}

func (this Distance) nextDistance() int {
	return this.race.measureDistance(this.value + 1)
}

func createDistance(value int, race *Race) Distance {
	return Distance{value, race.measureDistance(value), race}
}
