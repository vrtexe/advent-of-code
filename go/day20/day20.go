package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {

	// Solution1()
	Solution2()
}

func Solution1() {
	broadcast, config, _ := parse(nil)
	pressButton := createButton(&broadcast, config)

	for i := 0; i < 1000; i++ {
		pressButton()
	}

	fmt.Println(broadcast.counter.value[LOW] * broadcast.counter.value[HIGH])
}

func Solution2() {
	cycleSnapshot := CycleSnapshot{map[string]int{}, "", 0, 0}
	broadcast, config, inputMap := parse(&cycleSnapshot)

	destinationInput := inputMap["rx"][0]

	cycleSnapshot.output = destinationInput
	for input := range *extractInputs(destinationInput, config) {
		cycleSnapshot.lengths[input] = -1
	}

	for !cycleSnapshot.isDone(config) {
		cycleSnapshot.nextCycle()
		broadcast.send(LOW)

		for !broadcast.counter.done() {
			for _, module := range *config {
				module.receiveSignals()
			}
		}
	}

	buttonPresses := mLcm(cycleSnapshot.toLengthsList())
	fmt.Println(buttonPresses)
}

func createButton(broadcast *Broadcast, config *Config) func() {
	return func() {
		broadcast.send(LOW)

		for !broadcast.counter.done() {
			for _, module := range *config {
				module.receiveSignals()
			}
		}
	}
}

func printConfig(config *Config) {
	for _, c := range *config {
		switch v := c.(type) {
		case *Conjunction:
			fmt.Println(len(v.channel), v)
		case *FlipFlop:
			fmt.Println(len(v.channel), v)
		}
	}
}

func parse(cycleSnapshot *CycleSnapshot) (Broadcast, *Config, map[string][]string) {
	broadcastValue := ""
	config := Config{}
	inputsMapping := map[string][]string{}
	pulseCounter := PulseCounter{map[Pulse]int{}, 0, 0}

	mapInputs := createInputMapper(&inputsMapping)

	common.ReadFile("day20/day20.txt", func(line string) {
		if strings.HasPrefix(line, FLIP_FLOP) {
			value, connections := parseModule(strings.TrimPrefix(line, FLIP_FLOP))
			config[value] = &FlipFlop{value, false, connections, make(chan PulseDto, 10), make(chan PulseDto, 10), &config, &pulseCounter, cycleSnapshot, LOW}
			mapInputs(value, connections)
		} else if strings.HasPrefix(line, CONJUNCTION) {
			value, connections := parseModule(strings.TrimPrefix(line, CONJUNCTION))
			config[value] = &Conjunction{value, map[string]Pulse{}, connections, make(chan PulseDto, 10), make(chan PulseDto, 10), &config, &pulseCounter, cycleSnapshot, HIGH}
			mapInputs(value, connections)
		} else {
			broadcastValue = line
		}
	})

	for module, inputs := range inputsMapping {
		switch v := config[module].(type) {
		case *Conjunction:
			v.inputs = setupInputs(inputs)
		}
	}

	broadcast := parseBroadcast(broadcastValue, &config, &pulseCounter)

	return broadcast, &config, inputsMapping
}

func setupInputs(inputs []string) map[string]Pulse {
	inputPulses := map[string]Pulse{}

	for _, input := range inputs {
		inputPulses[input] = LOW
	}

	return inputPulses
}

func createInputMapper(inputMap *map[string][]string) func(string, []string) {
	return func(value string, outputs []string) {
		for _, output := range outputs {
			if input, exists := (*inputMap)[output]; !exists {
				(*inputMap)[output] = []string{value}
			} else {
				(*inputMap)[output] = append(input, value)
			}
		}

	}
}

func parseModule(line string) (string, []string) {
	value, connections := common.Split2(line, "->")
	return value, strings.Split(connections, ", ")
}

func parseBroadcast(line string, config *Config, pulseCounter *PulseCounter) Broadcast {
	_, modules := common.Split2(line, "->")

	outputs := []Module{}

	for _, module := range strings.Split(modules, ", ") {
		outputs = append(outputs, (*config)[module])
	}

	return Broadcast{outputs, pulseCounter}
}

type Config map[string]Module

type Broadcast struct {
	modules []Module
	counter *PulseCounter
}

func (this Broadcast) send(pulse Pulse) {
	this.counter.send(pulse)
	this.counter.receive(pulse)
	for _, module := range this.modules {
		module.send(PulseDto{pulse, "broadcast"})
	}
}

type Module interface {
	receiveSignals()
	transfer(PulseDto)
	send(PulseDto)
	getPulse() Pulse
}

type FlipFlop struct {
	value         string
	state         bool
	connections   []string
	channel       chan PulseDto
	next          chan PulseDto
	config        *Config
	counter       *PulseCounter
	cycleSnapshot *CycleSnapshot
	pulse         Pulse
}

func (this FlipFlop) getPulse() Pulse {
	return this.nextPulse()
}

func (this *FlipFlop) receiveSignals() {
	// this.cycleSnapshot.calculateCycle(this.value, this.getPulse(), nil)

	if len(this.channel) > 0 {
		this.transfer(<-this.channel)
	}

	this.channel, this.next = this.next, this.channel
}

func (this *FlipFlop) transfer(pulse PulseDto) {
	this.counter.receive(pulse.value)

	if pulse.value == HIGH {
		return
	} else if pulse.value == LOW {
		this.state = !this.state
		sendPulse := this.nextPulse()
		this.flow(sendPulse)
	}
}

func (this *FlipFlop) flow(pulse Pulse) {
	for _, connection := range this.connections {
		(*this.config)[connection].send(PulseDto{pulse, this.value})
	}
}

func (this *FlipFlop) send(pulse PulseDto) {
	this.counter.send(pulse.value)
	this.next <- pulse

	// fmt.Println(pulse.in, fmt.Sprintf("-%s->", pulse.value), this.value)
}

func (this FlipFlop) nextPulse() Pulse {
	if this.state {
		this.pulse = HIGH
		return HIGH
	} else {
		this.pulse = LOW
		return LOW
	}
}

type Conjunction struct {
	value         string
	inputs        map[string]Pulse
	connections   []string
	channel       chan PulseDto
	next          chan PulseDto
	config        *Config
	counter       *PulseCounter
	cycleSnapshot *CycleSnapshot
	pulse         Pulse
}

func (this Conjunction) getPulse() Pulse {
	return this.pulse
}

func (this *Conjunction) receiveSignals() {
	if len(this.channel) > 0 {
		this.transfer(<-this.channel)
	}

	this.channel, this.next = this.next, this.channel
}

func (this *Conjunction) transfer(pulse PulseDto) {
	this.counter.receive(pulse.value)

	this.cycleSnapshot.calculateCycle(this.value, pulse.in, pulse.value)

	this.inputs[pulse.in] = pulse.value
	for _, pulse := range this.inputs {
		if pulse == LOW {
			this.pulse = HIGH
			this.flow(HIGH)
			return
		}
	}

	this.pulse = LOW
	this.flow(LOW)
}

func (this *Conjunction) flow(pulse Pulse) {
	for _, connection := range this.connections {
		if module, exists := (*this.config)[connection]; exists {
			module.send(PulseDto{pulse, this.value})
		} else {
			// fmt.Println(this.value, fmt.Sprintf("-%s->", pulse), connection)
			this.counter.send(pulse)
			this.counter.receive(pulse)
		}
	}
}

func (this *Conjunction) send(pulse PulseDto) {
	this.counter.send(pulse.value)
	this.next <- pulse

	// fmt.Println(pulse.in, fmt.Sprintf("-%s->", pulse.value), this.value)
}

type PulseCounterCondition func(value string, pulse Pulse) bool
type PulseCounter struct {
	value    map[Pulse]int
	sent     int
	received int
}

type CycleSnapshot struct {
	lengths map[string]int
	output  string
	current int
	inserts int
}

func (this *CycleSnapshot) isDone(config *Config) bool {
	return this.inserts == len(this.lengths)
}

func extractInputs(moduleValue string, config *Config) *map[string]Pulse {
	switch module := (*config)[moduleValue].(type) {
	case *Conjunction:
		return &module.inputs
	default:
		return nil
	}
}

func (this *CycleSnapshot) calculateCycle(output, input string, pulse Pulse) {
	if cycleLength, exists := this.lengths[input]; exists && this.output == output && cycleLength == -1 && pulse == HIGH {
		this.lengths[input] = this.current
		this.inserts++
	}
}

func (this *CycleSnapshot) toLengthsList() []int {
	list := []int{}
	for _, length := range this.lengths {
		list = append(list, length)
	}
	return list
}

func (this *CycleSnapshot) validateInputs(inputs *map[string]Pulse) ([]int, bool) {
	values := []int{}
	if inputs == nil {
		return values, false
	}

	for input := range *inputs {
		if value, exists := this.lengths[input]; !exists {
			return values, false
		} else {
			values = append(values, value)
		}
	}

	return values, true
}

func (this *CycleSnapshot) nextCycle() {
	this.current++
}

func (this *PulseCounter) send(pulse Pulse) {
	this.value[pulse]++
	this.sent++
}

func (this *PulseCounter) receive(pulse Pulse) {
	this.received++
}

func (this *PulseCounter) done() bool {
	return this.received == this.sent
}

type Pulse string

type PulseDto struct {
	value Pulse
	in    string
}

func (this Pulse) flip() Pulse {
	if this == LOW {
		return HIGH
	} else if this == HIGH {
		return LOW
	}

	return LOW
}

const (
	LOW  Pulse = "LOW"
	HIGH       = "HIGH"
)

const (
	FLIP_FLOP   = "%"
	CONJUNCTION = "&"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return (a * b) / gcd(a, b)
}

func mLcm(args []int) int {
	currentLcm := 1

	for _, arg := range args {
		currentLcm = lcm(currentLcm, arg)
	}

	return currentLcm
}
