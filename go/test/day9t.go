package main

import "fmt"

type Point struct {
	x, y float64
}

func doSegmentsIntersect(p1, p2, p3, p4 Point) bool {
	denom := (p4.y-p3.y)*(p2.x-p1.x) - (p4.x-p3.x)*(p2.y-p1.y)
	numA := (p4.x-p3.x)*(p1.y-p3.y) - (p4.y-p3.y)*(p1.x-p3.x)
	numB := (p2.x-p1.x)*(p1.y-p3.y) - (p2.y-p1.y)*(p1.x-p3.x)

	if denom == 0 {
		return numA == 0 && numB == 0 &&
			((p1.x >= p3.x && p1.x <= p4.x) || (p2.x >= p3.x && p2.x <= p4.x) ||
				(p3.x >= p1.x && p3.x <= p2.x) || (p4.x >= p1.x && p4.x <= p2.x))
	}

	uA := numA / denom
	uB := numB / denom

	return uA >= 0 && uA <= 1 && uB >= 0 && uB <= 1
}

func main() {
	p1 := Point{1, 0}
	p2 := Point{1, 2}
	p3 := Point{0, 2}
	p4 := Point{2, 2}

	intersects := doSegmentsIntersect(p1, p2, p3, p4)

	if intersects {
		fmt.Println("Lines intersect within the segments.")
	} else {
		fmt.Println("Lines do not intersect within the segments.")
	}
}
