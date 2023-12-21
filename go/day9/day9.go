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

func Solution1() {
	sequenceMaps := parse()
	sum := 0

	for _, sequenceMap := range sequenceMaps {
		nextValue := findNext(sequenceMap)
		sum += nextValue
	}

	fmt.Println(sum)
}

func Solution2() {
	sequenceMaps := parse()
	sum := 0

	for _, sequenceMap := range sequenceMaps {
		nextValue := findPrev(sequenceMap)
		sum += nextValue
	}

	fmt.Println(sum)
}

func findNext(sequenceMap SequenceMap) int {
	currentSequence := sequenceMap.end
	currentValue := 0
	for currentSequence != nil {
		currentValue += currentSequence.values[len(currentSequence.values)-1]
		currentSequence = currentSequence.parent
	}

	return currentValue
}

func findPrev(sequenceMap SequenceMap) int {
	currentSequence := sequenceMap.end
	currentValue := 0
	for currentSequence != nil {
		currentValue = currentSequence.values[0] - currentValue
		currentSequence = currentSequence.parent
	}

	return currentValue
}

func parse() []SequenceMap {
	sequenceMaps := []SequenceMap{}
	common.ReadFile("day9/day9.txt", func(line string) {
		sequenceValues := strings.Fields(line)

		sequence := Sequence{}
		for _, value := range sequenceValues {
			sequence.values = append(sequence.values, common.ParseInt(value))
		}

		currentSequence := &sequence
		lastSequence := currentSequence
		for slices.Max(currentSequence.values) != 0 || slices.Min(currentSequence.values) != 0 {
			currentSequence.subSequence = currentSequence.findSubSequence()
			currentSequence = currentSequence.subSequence
			lastSequence = currentSequence
		}

		currentSequence = &sequence

		sequenceMaps = append(sequenceMaps, SequenceMap{&sequence, lastSequence})
	})

	return sequenceMaps
}

type SequenceMap struct {
	start, end *Sequence
}

type Sequence struct {
	values      []int
	subSequence *Sequence
	parent      *Sequence
}

func (this *Sequence) findSubSequence() *Sequence {
	nextSeq := []int{}
	for index := 1; index < len(this.values); index++ {
		leftValue := this.values[index-1]
		rightValue := this.values[index]

		nextSequenceValue := rightValue - leftValue
		nextSeq = append(nextSeq, nextSequenceValue)
	}

	return &Sequence{nextSeq, nil, this}
}

func (this *Sequence) print() {
	currentSequence := this
	i := 0
	for {

		fmt.Print(strings.Repeat(" ", i*2))

		for _, value := range currentSequence.values {
			fmt.Print(value, " ")
		}

		fmt.Println()

		if currentSequence.subSequence == nil {
			break
		}

		currentSequence = currentSequence.subSequence
		i++
	}
}
