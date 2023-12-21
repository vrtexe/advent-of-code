package main

import (
	"common"
	"fmt"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

type Graph struct {
	nodes map[string]*Node
}

type Node struct {
	value string
	left  *Node
	right *Node

	next map[string]*Node
}

func Solution1() {
	graph, sequence, _ := parse()
	findPath := createPathCountFinder(&graph, sequence)
	count := findPath("AAA", func(s string) bool { return s == "ZZZ" })

	fmt.Println(count)
}

func Solution2() {
	graph, sequence, startingPoints := parse()

	findPathCount := createPathCountFinder(&graph, sequence)
	pointCounts := []int{}
	condition := func(s string) bool {
		return strings.HasSuffix(s, "Z")
	}

	for _, point := range startingPoints {
		pointCounts = append(pointCounts, findPathCount(point, condition))
	}

	count := mLcm(pointCounts)

	fmt.Println(count)
}

func createPathCountFinder(graph *Graph, sequence []string) func(from string, condition func(string) bool) int {
	return func(from string, condition func(string) bool) int {
		currentState := graph.nodes[from]
		destination := false
		count := 0

		for index := 0; !destination; index = (index + 1) % len(sequence) {
			nextStep := sequence[index]
			currentState = currentState.Move(nextStep)
			destination = condition(currentState.value)
			count++
		}

		return count
	}
}

func findSinglePathCount(graph *Graph, sequence []string) {
	destination := false

	currentState := graph.nodes["AAA"]
	count := 0
	for index := 0; !destination; index = (index + 1) % len(sequence) {
		nextStep := sequence[index]
		currentState = currentState.Move(nextStep)
		destination = currentState.value == "ZZZ"
		count++
	}

	fmt.Println(count)
}

func parse() (Graph, []string, []string) {
	sequence := []string{}
	graph := Graph{map[string]*Node{}}
	startingPoints := []string{}

	parseNode := createNodeParser(&graph)

	common.ReadFileLines("day8/day8.txt", func(line string, index int) {
		if line == "" {
			return
		}

		if index == 0 {
			sequence = strings.Split(line, "")
			return
		}

		node := parseNode(line)
		graph.Add(node)

		if strings.HasSuffix(node.value, "A") {
			startingPoints = append(startingPoints, node.value)
		}
	})

	return graph, sequence, startingPoints
}

func createNodeParser(graph *Graph) func(string) *Node {
	replacer := strings.NewReplacer("(", "", ")", "")

	return func(line string) *Node {
		value, neighbors := common.Split2(line, "=")
		leftValue, rightValue := common.Split2(replacer.Replace(neighbors), ", ")

		currentNode := findNodeOrCreate(value, graph)
		if leftValue != rightValue {
			leftNode := findNodeOrCreate(leftValue, graph)
			rightNode := findNodeOrCreate(rightValue, graph)

			currentNode.SetLeft(leftNode)
			currentNode.SetRight(rightNode)
		} else {
			node := findNodeOrCreate(leftValue, graph)
			currentNode.SetLeft(node)
			currentNode.SetRight(node)
		}

		return currentNode
	}
}

func findNodeOrCreate(value string, graph *Graph) *Node {
	if node, exists := graph.nodes[value]; exists {
		return node
	} else {
		return &Node{value, nil, nil, map[string]*Node{}}
	}
}

func (this Graph) Add(node *Node) {
	if _, exists := this.nodes[node.left.value]; !exists {
		this.nodes[node.left.value] = node.left
	}

	if _, exists := this.nodes[node.right.value]; !exists {
		this.nodes[node.right.value] = node.right
	}

	this.nodes[node.value] = node
}

func (this *Node) SetLeft(node *Node) {
	this.left = node
	this.next["L"] = node
}

func (this *Node) SetRight(node *Node) {
	this.right = node
	this.next["R"] = node
}

func (this *Node) Move(direction string) *Node {
	return this.next[direction]
}

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
