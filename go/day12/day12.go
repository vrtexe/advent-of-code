package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	Solution2()
}

func Solution1() {
	arrangements := parse()

	sum := 0

	for _, arrangement := range arrangements {
		sum += countSolutions(arrangement.value, arrangement.sequence)
	}

	fmt.Println(sum)
}

func Solution2() {
	arrangements := sliceMap(parse(), func(a Arrangement) Arrangement {
		return Arrangement{repeat(a.value, 5, " ", "?"), repeatSlice(a.sequence, 5, ",")}
	})

	sum := 0

	for _, arrangement := range arrangements {
		sum += countSolutions(arrangement.value, arrangement.sequence)
	}

	fmt.Println(sum)
}

func repeatSlice(sequence []int, times int, separator string) []int {
	return sliceMap(strings.Split(repeat(strings.Join(sliceMap(sequence, func(s int) string { return fmt.Sprint(s) }), separator), times, separator, separator), separator), common.ParseInt)
}

func repeat(value string, times int, separator string, join string) string {
	return strings.Join(strings.Split(strings.TrimSuffix(strings.Repeat(value+separator, 5), separator), separator), join)
}

func parse() []Arrangement {

	arrangements := []Arrangement{}

	common.ReadFile("day12/day12.txt", func(line string) {

		pattern, sequenceValues := common.Split2(line, " ")
		sequence := sliceMap(strings.Split(sequenceValues, ","), common.ParseInt)
		arrangements = append(arrangements, Arrangement{pattern, sequence})

	})

	return arrangements

}

type Arrangement struct {
	value    string
	sequence []int
}

type ArrangementSate struct {
	group, position int
	done            bool
}

func countSolutions(value string, sequence []int) int {
	row := []rune(value)

	currentStates := map[ArrangementSate]int{{0, 0, false}: 1}

	for _, spring := range row {
		if nextState, empty := computeNextState(rune(spring), sequence, currentStates); !empty {
			currentStates = nextState
		}
	}

	return sumSolutions(len(sequence), currentStates)

}

func sumSolutions(sequenceCount int, state map[ArrangementSate]int) int {
	sum := 0
	for state, count := range state {
		if state.group == sequenceCount {
			sum += count
		}
	}

	return sum
}

func computeNextState(spring rune, sequence []int, state map[ArrangementSate]int) (map[ArrangementSate]int, bool) {
	sequenceCount := len(sequence)
	nextState := map[ArrangementSate]int{}

	for state, count := range state {
		if (spring == '#' || spring == '?') && state.group < sequenceCount && !state.done {
			if spring == '?' && state.position == 0 {
				nextState[ArrangementSate{state.group, state.position, state.done}] += count
			}

			if nextStatePosition := state.position + 1; nextStatePosition == sequence[state.group] {
				nextState[ArrangementSate{state.group + 1, 0, true}] += count
			} else {
				nextState[ArrangementSate{state.group, nextStatePosition, state.done}] += count
			}

		} else if (spring == '.' || spring == '?') && state.position == 0 {
			nextState[ArrangementSate{state.group, state.position, false}] += count
		}
	}

	return nextState, len(nextState) == 0

}

func mapsClear[M ~map[K]V, K comparable, V any](m M) {
	for k := range m {
		delete(m, k)
	}
}

func pop[S ~[]E, E any](s S) (E, S) {
	lastIndex := len(s) - 1
	last := s[lastIndex]

	return last, s[:lastIndex]
}

func sliceMap[S ~[]E, E any, R any](s S, transform func(E) R) []R {
	newSlice := []R{}

	for _, e := range s {
		newSlice = append(newSlice, transform(e))
	}

	return newSlice
}

func calculateCombinationsRepeating(n, k int) int {
	return calculateCombinations(n+k-1, k)
}

func calculateCombinations(n, k int) int {
	return factorial(n) / factorial(k) * factorial(n-k)
}

func factorial(n int) int {
	fact := 1

	for i := 1; i <= n; i++ {
		fact *= i
	}

	return fact
}
