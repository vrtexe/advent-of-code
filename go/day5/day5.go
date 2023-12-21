package main

import (
	"common"
	"fmt"
	"slices"
	"strings"
)

func main() {
	Solution2()
}

type ParserState struct {
	text string
	step int
	next func(state ParserStates) *ParserState
}

func (this ParserState) isBefore(that ParserState) bool {
	return this.step <= that.step
}

type ParserStates struct {
	Seeds,
	SeedToSoil,
	SoilToFertilizer,
	FertilizerToWater,
	WaterToLight,
	LightToTemperature,
	TemperatureToHumidity,
	HumidityToLocation *ParserState
}

type FactorMap struct {
	from, to, count uint64
}

type FactorMaps struct {
	seedToSoil,
	soilToFertilizer,
	fertilizerToWater,
	waterToLight,
	lightToTemperature,
	temperatureToHumidity,
	humidityToLocation []FactorMap
}

func (this *FactorMaps) Set(state *ParserState, factorMap *FactorMap) {
	switch state {
	case States.SeedToSoil:
		this.seedToSoil = appendFactor(this.seedToSoil, *factorMap)
	case States.SoilToFertilizer:
		this.soilToFertilizer = appendFactor(this.soilToFertilizer, *factorMap)
	case States.FertilizerToWater:
		this.fertilizerToWater = appendFactor(this.fertilizerToWater, *factorMap)
	case States.WaterToLight:
		this.waterToLight = appendFactor(this.waterToLight, *factorMap)
	case States.LightToTemperature:
		this.lightToTemperature = appendFactor(this.lightToTemperature, *factorMap)
	case States.TemperatureToHumidity:
		this.temperatureToHumidity = appendFactor(this.temperatureToHumidity, *factorMap)
	case States.HumidityToLocation:
		this.humidityToLocation = appendFactor(this.humidityToLocation, *factorMap)
	}
}

var States = ParserStates{
	Seeds:                 &ParserState{"seeds:", 1, func(state ParserStates) *ParserState { return state.SeedToSoil }},
	SeedToSoil:            &ParserState{"seed-to-soil map:", 2, func(state ParserStates) *ParserState { return state.SoilToFertilizer }},
	SoilToFertilizer:      &ParserState{"soil-to-fertilizer map:", 3, func(state ParserStates) *ParserState { return state.FertilizerToWater }},
	FertilizerToWater:     &ParserState{"fertilizer-to-water map:", 4, func(state ParserStates) *ParserState { return state.WaterToLight }},
	WaterToLight:          &ParserState{"water-to-light map:", 5, func(state ParserStates) *ParserState { return state.LightToTemperature }},
	LightToTemperature:    &ParserState{"light-to-temperature map:", 6, func(state ParserStates) *ParserState { return state.TemperatureToHumidity }},
	TemperatureToHumidity: &ParserState{"temperature-to-humidity map:", 7, func(state ParserStates) *ParserState { return state.HumidityToLocation }},
	HumidityToLocation:    &ParserState{"humidity-to-location map:", 8, func(state ParserStates) *ParserState { return nil }},
}

type FactorInterval struct {
	start, count uint64
}

func (this FactorInterval) end() uint64 {
	return (this.start + this.count) - 1
}

func (this FactorMap) getMapping(value uint64) (uint64, bool) {
	if value < this.from || value > this.from+this.count-1 {
		return 0, false
	}

	return this.to + (value - this.from), true
}

func (this FactorMap) getIntervalMapping(value FactorInterval) (FactorInterval, FactorInterval, bool) {
	if start, end := max(this.from, value.start), min(this.from+this.count-1, value.end()); start <= end {
		mappedStart, _ := this.getMapping(start)
		mappedEnd, _ := this.getMapping(end)
		return FactorInterval{mappedStart, (mappedEnd - mappedStart) + 1}, FactorInterval{start, (end - start) + 1}, true
	} else {
		return FactorInterval{}, FactorInterval{}, false
	}
}

func (this FactorMaps) getOrderedMaps() [][]FactorMap {
	return [][]FactorMap{
		this.seedToSoil,
		this.soilToFertilizer,
		this.fertilizerToWater,
		this.waterToLight,
		this.lightToTemperature,
		this.temperatureToHumidity,
		this.humidityToLocation,
	}
}

func Solution2() {
	seeds, factorMaps := parse()

	locations := []FactorInterval{}

	for i := 0; i < len(seeds); i += 2 {
		interval := FactorInterval{common.ParseUint64(seeds[i]), common.ParseUint64(seeds[i+1])}
		locationIntervals := findSeedLocationInterval(interval, factorMaps)
		locations = append(locations, locationIntervals...)

	}

	minValue := slices.MinFunc(locations, func(a, b FactorInterval) int {
		return common.Uint64Compare(a.start, b.start)
	})

	fmt.Println(minValue.start)
}

func Solution1() {
	seeds, factorMaps := parse()

	locations := []uint64{}

	for _, seed := range seeds {
		locations = append(locations, findSeedLocation(seed, factorMaps))
	}

	minValue := slices.Min(locations)

	fmt.Println(minValue)
}

func parse() ([]string, FactorMaps) {
	var seeds []string

	factorMaps := FactorMaps{}

	stateCompute := createStateComputer(States.Seeds)

	common.ReadFile("../day5/day5.txt", func(line string) {
		if strings.Trim(line, " ") == "" {
			return
		}

		currentState := stateCompute(line)

		if currentState == States.Seeds {
			seeds = parseSeeds(line)
		} else if factorMap := parseFactorMap(line, *currentState); factorMap != nil {
			factorMaps.Set(currentState, factorMap)
		}
	})

	return seeds, factorMaps
}

func findSeedLocationInterval(seedInterval FactorInterval, factorMaps FactorMaps) []FactorInterval {
	currentState := []FactorInterval{seedInterval}

	for _, stateMap := range factorMaps.getOrderedMaps() {
		nextState := []FactorInterval{}

		for _, state := range currentState {
			mappedIntervals := findIntervalMapping(state, stateMap)
			nextState = append(nextState, mappedIntervals...)
		}

		currentState = nextState
	}

	return currentState

}

func findSeedLocation(seed string, factorMaps FactorMaps) uint64 {
	currentState := common.ParseUint64(seed)

	for _, stateMap := range factorMaps.getOrderedMaps() {
		currentState = findValueMapping(currentState, stateMap)
	}

	return currentState
}

func findIntervalMapping(seed FactorInterval, factorMaps []FactorMap) []FactorInterval {
	mappedIntervals := []FactorInterval{}
	mappedParts := []FactorInterval{}

	for _, factorMap := range factorMaps {
		if mapping, mapped, ok := factorMap.getIntervalMapping(seed); ok {
			mappedParts = append(mappedParts, mapped)
			mappedIntervals = append(mappedIntervals, mapping)
		}
	}

	if len(mappedIntervals) == 0 {
		mappedIntervals = append(mappedIntervals, seed)
	} else {
		mappedIntervals = append(mappedIntervals, extractUnmapped(seed, mappedParts)...)
	}

	return mappedIntervals
}

func extractUnmapped(seed FactorInterval, mapped []FactorInterval) []FactorInterval {
	unmappedIntervals := []FactorInterval{}

	minStartFactor := slices.MinFunc(mapped, func(a, b FactorInterval) int {
		return common.Uint64Compare(a.start, b.start)
	})

	maxEndFactor := slices.MaxFunc(mapped, func(a, b FactorInterval) int {
		return common.Uint64Compare(a.end(), b.end())
	})

	if seed.start < minStartFactor.start {
		unmappedIntervals = append(unmappedIntervals, FactorInterval{seed.start, minStartFactor.start - seed.start})
	}

	if seed.end() > maxEndFactor.end() {
		unmappedIntervals = append(unmappedIntervals, FactorInterval{maxEndFactor.end() + 1, (seed.end() - (maxEndFactor.end() + 1)) + 1})
	}

	return unmappedIntervals
}

func findValueMapping(seed uint64, factorMaps []FactorMap) uint64 {
	for _, factorMap := range factorMaps {
		if mapping, ok := factorMap.getMapping(seed); ok {
			return mapping
		}
	}

	return seed
}

func parseSeeds(line string) []string {
	_, seedText := common.Split2(line, ":")
	return strings.Fields(seedText)
}

func parseFactorMap(line string, state ParserState) *FactorMap {
	if strings.HasPrefix(line, state.text) {
		return nil
	}

	fields := strings.Fields(line)

	return &FactorMap{
		from:  common.ParseUint64(fields[1]),
		to:    common.ParseUint64(fields[0]),
		count: common.ParseUint64(fields[2]),
	}
}

func appendFactor(factorMap []FactorMap, item FactorMap) []FactorMap {
	if factorMap == nil {
		return []FactorMap{item}
	} else {
		return append(factorMap, item)
	}
}

func createStateComputer(startingState *ParserState) func(line string) *ParserState {
	currentState := startingState

	return func(line string) *ParserState {
		nextState := currentState.next(States)

		if nextState == nil {
			return currentState
		}

		if strings.HasPrefix(line, nextState.text) {
			currentState = nextState
		}

		return currentState
	}
}
