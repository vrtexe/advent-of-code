package main

import (
	"common"
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

type Data struct {
	cuts   int
	edges  []string
	graphs *[][]string
}

func main() {
	data := parse()
	graph := buildGraph(data)

	results := []Data{}

	for {
		count, eg, g := kargerMinCut(graph)
		results = append(results, Data{cuts: count, edges: eg, graphs: g})
		if count == 3 {
			break
		}
	}

	minEdges := slices.MinFunc(results, func(left, right Data) int {
		return left.cuts - right.cuts
	})

	counts := []int{}
	countProduct := 1
	for _, g := range *minEdges.graphs {
		counts = append(counts, len(g))
		countProduct *= len(g)
	}

	fmt.Println(countProduct)
}

func kargerMinCut(graph Graph[string]) (int, []string, *[][]string) {
	disjointSet := NewDisjointedSet[string]()
	vertices := []string{}
	for v := range graph.vertexes {
		disjointSet.MakeSet(v)
		vertices = append(vertices, v)
	}

	verticesLeft := len(vertices)

	for verticesLeft > 2 {
		vertex := vertices[rand.Intn(len(vertices))]
		if edges, exists := graph.edges[vertex]; !exists || len(edges) <= 0 {
			continue
		}
		edge := graph.edges[vertex][rand.Intn(len(graph.edges[vertex]))]

		v := disjointSet.Find(edge.left)
		u := disjointSet.Find(edge.right)

		if v == u {
			continue
		}

		disjointSet.Union(edge.left, edge.right)
		verticesLeft--
	}

	cuts := 0
	cutEdges := []string{}
	bySubset := map[*Subset[string]][]string{}
	for _, e := range graph.edges {
		for _, c := range e {
			v := disjointSet.Find(c.left)
			u := disjointSet.Find(c.right)

			if v != u {
				cutEdges = append(cutEdges, fmt.Sprintf("%s-%s", c.left, c.right))
				cuts++
			}
		}
	}

	for _, v := range vertices {
		subset := disjointSet.Find(v)
		if _, exists := bySubset[subset]; exists {
			bySubset[subset] = append(bySubset[subset], v)
		} else {
			bySubset[subset] = []string{v}
		}
	}

	lists := [][]string{}
	for _, v := range bySubset {
		lists = append(lists, v)
	}

	return cuts, cutEdges, &lists
}

func buildGraph(data map[string][]string) Graph[string] {
	vertexes := map[string]struct{}{}
	edges := map[string][]Edge[string]{}

	for v, e := range data {
		vertexes[v] = struct{}{}
		if _, exists := edges[v]; !exists {
			edges[v] = []Edge[string]{}
		}

		for _, u := range e {
			vertexes[u] = struct{}{}
			edges[v] = append(edges[v], Edge[string]{left: v, right: u})
		}
	}

	return Graph[string]{vertexes: vertexes, edges: edges}
}

func parse() map[string][]string {
	data := map[string][]string{}

	common.ReadFile("day25/day25.txt", func(line string) {
		vertex, connections := common.Split2(line, ": ")
		data[vertex] = strings.Split(connections, " ")
	})

	return data
}

type Graph[V comparable] struct {
	vertexes map[V]struct{}
	edges    map[V][]Edge[V]
}

type Edge[V any] struct {
	left, right V
}

type Subset[T any] struct {
	parent *Subset[T]
	rank   int
	size   int
}

type DisjointedSetForest[T comparable] struct {
	subset map[T]*Subset[T]
}

type DisjointedSet[V any] interface {
	MakeSet(v V)
	Find(v V) *Subset[V]
	Union(x V, y V)
}

func NewDisjointedSet[T comparable]() DisjointedSet[T] {
	return DisjointedSetForest[T]{subset: map[T]*Subset[T]{}}
}

func (this DisjointedSetForest[V]) MakeSet(v V) {
	if _, exists := this.subset[v]; exists {
		return
	}

	node := Subset[V]{parent: nil, rank: 0, size: 1}
	node.parent = &node
	this.subset[v] = &node
}

func (this DisjointedSetForest[V]) Find(v V) *Subset[V] {
	if _, exists := this.subset[v]; !exists {
		return nil
	}

	n := this.subset[v]

	for n.parent != n {
		n.parent = n.parent.parent
		n = n.parent
	}

	return n
}

func (this DisjointedSetForest[T]) Union(x T, y T) {
	xRoot, yRoot := this.Find(x), this.Find(y)

	if xRoot == yRoot {
		return
	}

	if xRoot.rank < yRoot.rank {
		xRoot.parent = yRoot
	} else if yRoot.rank < xRoot.rank {
		yRoot.parent = xRoot
	} else {
		yRoot.parent = xRoot
		xRoot.rank++
	}
}
