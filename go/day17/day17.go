package main

import (
	"common"
	"container/heap"
	"fmt"
	"slices"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

func Solution1() {
	data := parse()
	constraint := Constraint{min: 1, max: 3}

	minDistance := getShortestPathFor(constraint, data)

	fmt.Println("Result: ", minDistance.weight)
}

func Solution2() {
	data := parse()
	constraint := Constraint{min: 4, max: 10}

	minDistance := getShortestPathFor(constraint, data)

	fmt.Println("Result: ", minDistance.weight)
}

func getShortestPathFor(constraint Constraint, data [][]string) Distance[Block] {
	rows, columns := len(data), len(data[0])

	graph, neighborExtractor := buildGraph(data, buildNeighborExtractor(constraint))

	start := graph[Horizontal][0][0]

	distances := findShortestPath(Block{
		position:  start.position,
		direction: start.direction,
		weight:    0,
	}, vertexKey, neighborExtractor)

	ends := []Block{
		graph[Horizontal][rows-1][columns-1],
		graph[Vertical][rows-1][columns-1],
	}

	results := []Distance[Block]{}

	for _, end := range ends {
		if v, exists := distances[vertexKey(end)]; exists {
			results = append(results, v)
		}
	}

	minDistance := slices.MinFunc(results, func(left, right Distance[Block]) int {
		return left.weight - right.weight
	})

	return minDistance
}

func parse() [][]string {
	grid := [][]string{}

	common.ReadFile("day17/day17.txt", func(line string) {
		grid = append(grid, strings.Split(line, ""))
	})

	return grid
}

type Constraint struct {
	min int
	max int
}

func buildGraph(
	city [][]string,
	neighborExtractor NeighborExtractor,
) (map[Direction][][]Block, func(b Block) []Block) {
	// (Graph<Block, Edge<Block>>, Record<Direction, Block[][]>)
	horizontalBlocks := buildCityBlocks(Horizontal, city)
	verticalBlocks := buildCityBlocks(Vertical, city)

	blocks := map[Direction][][]Block{
		Horizontal: horizontalBlocks,
		Vertical:   verticalBlocks,
	}

	return blocks, func(v Block) []Block {
		return neighborExtractor(v, blocks)
	}
}

type NeighborExtractor = func(block Block, blocks map[Direction][][]Block) []Block

func vertexKey(vertex Block) string {
	return fmt.Sprintf("(%d,%d,%s)", vertex.position.x, vertex.position.y, *vertex.direction)
}

func buildCityBlocks(direction Direction, city [][]string) [][]Block {
	result := [][]Block{}

	for y, line := range city {
		tempBlocks := []Block{}

		for x, block := range line {
			tempBlocks = append(tempBlocks, Block{
				position:  Position{x: x, y: y},
				weight:    common.ParseInt(block),
				direction: &direction,
			})
		}

		result = append(result, tempBlocks)
	}

	return result
}

func buildNeighborExtractor(constraint Constraint) func(Block, map[Direction][][]Block) []Block {
	return func(block Block, blocks map[Direction][][]Block) []Block {
		return extractNeighboringBlocks(block, blocks, constraint)
	}
}

func extractNeighboringBlocks(
	block Block,
	blocks map[Direction][][]Block,
	constraint Constraint,
) []Block {
	position := block.position

	return extractNextNeighbors(
		block,
		constraint,
		createNeighborResolver(position, Direction(*block.direction), blocks),
	)
}

func createNeighborResolver(
	position Position,
	direction Direction,
	blocks map[Direction][][]Block,
) func(increment int) (Block, bool) {
	flippedDirection := nextDirection(direction)

	switch direction {
	case Horizontal:
		return func(increment int) (Block, bool) {
			if value, exists := blocks[flippedDirection]; !exists {
				return Block{}, false
			} else {
				if position.y < 0 || position.y >= len(value) {
					return Block{}, false
				}

				nextX := position.x + increment
				if nextX < 0 || nextX >= len(value[position.y]) {
					return Block{}, false
				}

				return value[position.y][nextX], true
			}
		}
	case Vertical:
		{
			return func(increment int) (Block, bool) {
				if value, exists := blocks[flippedDirection]; !exists {
					return Block{}, false
				} else {
					nextY := position.y + increment
					if nextY < 0 || nextY >= len(value) {
						return Block{}, false
					}

					if position.x < 0 || position.x >= len(value[nextY]) {
						return Block{}, false
					}

					return value[nextY][position.x], true
				}
			}
		}
	}

	panic("Invalid direction")
}

func extractNextNeighbors(
	block Block,
	constraint Constraint,
	next func(increment int) (Block, bool),
) []Block {
	direction := nextDirection(*block.direction)
	min, max := constraint.min, constraint.max

	result := []Block{}
	result = append(result, calculateNegativeWeights(-min, -max, direction, next)...)
	result = append(result, calculatePositiveWeights(min, max, direction, next)...)

	return result
}

func calculateNegativeWeights(
	start int,
	end int,
	direction Direction,
	next func(increment int) (Block, bool),
) []Block {
	weight := 0
	result := []Block{}

	for i := -1; i > start; i-- {
		nextBlock, exists := next(i)
		if !exists {
			break
		}

		weight += nextBlock.weight
	}

	for i := start; i >= end; i-- {
		nextBlock, exists := next(i)
		if !exists {
			break
		}

		weight += nextBlock.weight
		result = append(result, Block{
			position:  nextBlock.position,
			direction: &direction,
			weight:    weight,
		})
	}

	return result
}

func calculatePositiveWeights(
	start int,
	end int,
	direction Direction,
	next func(increment int) (Block, bool),
) []Block {
	weight := 0
	result := []Block{}

	for i := 1; i < start; i++ {
		nextBlock, nextBlockExists := next(i)
		if !nextBlockExists {
			break
		}

		weight += nextBlock.weight
	}

	for i := start; i <= end; i++ {
		nextBlock, nextBlockExists := next(i)
		if !nextBlockExists {
			break
		}
		weight += nextBlock.weight
		result = append(result, Block{
			position:  nextBlock.position,
			direction: &direction,
			weight:    weight,
		})
	}

	return result
}

func nextDirection(direction Direction) Direction {
	if direction == Horizontal {
		return Vertical
	}

	if direction == Vertical {
		return Horizontal
	}

	panic("Invalid direction")
}

type Direction string

const (
	Horizontal Direction = "H"
	Vertical   Direction = "V"
)

type Block struct {
	weight    int
	position  Position
	direction *Direction
}

func (this Block) getWeight() int {
	return this.weight
}

type Position struct {
	x, y int
}

func findShortestPath[V Weighted](vertex V, keyOf func(vertex V) string, getAdjacent func(vertex V) []V) map[string]Distance[V] {
	distances := map[string]Distance[V]{
		keyOf(vertex): {vertex: vertex, weight: 0},
	}

	queue := PriorityQueue[V]{{value: vertex, weight: 0, index: 0}}
	heap.Init(&queue)

	for queue.Len() > 0 {
		current := heap.Pop(&queue).(*Item[V])
		currentDistance := distances[keyOf(current.value)]

		for _, next := range getAdjacent(current.value) {
			nextKey := keyOf(next)
			if nextDistance, exists := distances[nextKey]; !exists || nextDistance.weight > currentDistance.weight+next.getWeight() {
				nextWeight := currentDistance.weight + next.getWeight()
				distances[nextKey] = Distance[V]{
					vertex:  next,
					through: &current.value,
					weight:  nextWeight,
				}

				heap.Push(&queue, &Item[V]{weight: nextWeight, value: next})
			}
		}
	}

	return distances
}

type Distance[V any] struct {
	vertex  V
	through *V
	weight  int
	index   int
}

type Item[V any] struct {
	value  V
	weight int
	index  int
}

type Weighted interface {
	getWeight() int
}

type PriorityQueue[V any] []*Item[V]

func (pq *PriorityQueue[V]) Push(x any) {
	pq.Poll(x.(*Item[V]))
}

func (pq *PriorityQueue[V]) Poll(x *Item[V]) {
	n := len(*pq)
	item := x
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[V]) Shift() *Item[V] {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1

	old[n-1] = nil
	// *pq = old[1:length]
	*pq = old[:n-1]

	return item
}

func (pq *PriorityQueue[V]) Pop() any {
	return pq.Shift()
}

func (pq *PriorityQueue[V]) update(item *Item[V], value V, weight int) {
	item.value = value
	item.weight = weight
	heap.Fix(pq, item.index)
}

func (pq PriorityQueue[V]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq PriorityQueue[V]) Len() int {
	return len(pq)
}

func (pq PriorityQueue[V]) Less(left, right int) bool {
	return pq[left].weight <= pq[right].weight
}
