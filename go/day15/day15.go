package main

import (
	"common"
	"fmt"
	"slices"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

func Solution1() {
	sequenceSteps := parse()
	hashes := calculateHashes(sequenceSteps)
	hashSum := sumHashes(hashes)

	fmt.Println(hashSum)
}

func Solution2() {
	sequenceSteps := parse()
	hashMap := computeHashMap(sequenceSteps)
	totalFocusPower := sumLensFocusPower(hashMap)

	fmt.Println(totalFocusPower)
}

func sumLensFocusPower(hashMap HashMap) int {
	sum := 0
	for boxNumber, box := range hashMap.boxes {
		for slot, lens := range box {
			sum += calculateFocusPower(boxNumber, slot, lens.focalLength)
		}
	}

	return sum
}

func calculateFocusPower(boxNumber int, slot int, focalLength int) int {
	return (boxNumber + 1) * (slot + 1) * focalLength
}

func computeHashMap(sequenceSteps []string) HashMap {
	hashMap := HashMap{map[int][]Lens{}}
	executeSequenceStep := createSequenceStepExecutor(&hashMap)

	for _, step := range sequenceSteps {
		executeSequenceStep(step)
	}

	return hashMap
}

func createSequenceStepExecutor(hashMap *HashMap) func(step string) HashMap {
	return func(step string) HashMap {
		strings.Split(step, "=")

		if strings.Contains(step, "=") {
			label, focalLength := common.Split2(step, "=")
			hashMap.put(label, common.ParseInt(focalLength))
		}

		if strings.Contains(step, "-") {
			label := strings.TrimSuffix(step, "-")
			hashMap.remove(label)
		}

		return *hashMap
	}
}

func parse() []string {
	sequenceSteps := []string{}

	common.ReadFile("day15/day15.txt", func(line string) {
		sequenceSteps = append(sequenceSteps, strings.Split(line, ",")...)
	})

	return sequenceSteps
}

func calculateHashes(sequenceSteps []string) []int {
	hashes := []int{}

	for _, step := range sequenceSteps {
		hashes = append(hashes, hash(step))
	}

	return hashes
}

func hash(value string) int {
	currentValue := 0

	for _, char := range []rune(value) {
		currentValue += int(char)
		currentValue *= 17
		currentValue %= 256
	}

	return currentValue
}

func sumHashes(hashes []int) int {
	sum := 0

	for _, hash := range hashes {
		sum += hash
	}

	return sum
}

type HashMap struct {
	boxes map[int][]Lens
}

type Lens struct {
	label       string
	focalLength int
}

func (this *HashMap) put(label string, focalLength int) {
	hashValue := hash(label)
	if _, exists := this.boxes[hashValue]; !exists {
		this.boxes[hashValue] = []Lens{}
	}

	if index := slices.IndexFunc(this.boxes[hashValue], func(b Lens) bool { return b.label == label }); index != -1 {
		this.boxes[hashValue][index].focalLength = focalLength
		return
	}

	this.boxes[hashValue] = append(this.boxes[hashValue], Lens{label, focalLength})
}

func (this *HashMap) remove(label string) {
	hashValue := hash(label)

	if !slices.ContainsFunc(this.boxes[hashValue], func(b Lens) bool { return b.label == label }) {
		return
	}

	this.boxes[hashValue] = slices.DeleteFunc(this.boxes[hashValue], func(b Lens) bool { return b.label == label })
}
