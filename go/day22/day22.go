package main

import (
	"common"
	"fmt"
	"math"
	"slices"
	"strings"
)

func main() {
	Solution1()
	Solution2()
}

func Solution1() {
	bricks := parse()
	_, support, supportedBy, _ := dropBricks(bricks)

	removable := countRemovable(support, supportedBy)

	fmt.Println(removable)
}

func Solution2() {
	bricks := parse()
	_, support, supportedBy, supportsMap := dropBricks(bricks)

	count := countFallingBricks(support, supportedBy, supportsMap)

	fmt.Println(count)
}

func countRemovable(supports, supportedBy map[Brick][]Brick) int {
	removable := 0
	for _, row := range supports {
		if isRemovable(row, supportedBy) {
			removable++
		}
	}
	return removable
}

func countFallingBricks(supports, supportedBy map[Brick][]Brick, supportsMap map[Brick]SupportMap) int {
	fallCounter := createFallCounter(supports, supportedBy, supportsMap, false)
	count := 0

	for b, row := range supports {
		count += fallCounter(b, row)
	}

	return count
}

func dropBricks(bricks [][]Brick) ([][]Brick, map[Brick][]Brick, map[Brick][]Brick, map[Brick]SupportMap) {
	support := map[Brick][]Brick{}
	supportedBy := map[Brick][]Brick{}
	supportsMap := map[Brick]SupportMap{}

	collectSupports := createSupportCollector(&support, &supportedBy, &supportsMap)

	newBricks := make([][]Brick, len(bricks))
	newBricks[1] = bricks[1]
	for _, brick := range bricks[1] {
		collectSupports(brick, []Brick{})
		if brick.start.z != brick.end.z {
			newBricks[brick.end.z] = append(newBricks[brick.end.z], brick)
		}
	}

	for level := 2; level < len(bricks); level++ {
		for _, brick := range bricks[level] {
			if condition, supports := brick.fitsIn(newBricks[level-1]); condition && len(supports) <= 0 {
				nextBrick := brick
				for j := level - 1; j > 0 && condition; j-- {
					nextBrick = nextBrick.fall()
					condition, supports = nextBrick.fitsIn(newBricks[j-1])
				}
				newBricks[nextBrick.start.z] = append(newBricks[nextBrick.start.z], nextBrick)
				if nextBrick.start.z != nextBrick.end.z {
					newBricks[nextBrick.end.z] = append(newBricks[nextBrick.end.z], nextBrick)
				}
				collectSupports(nextBrick, supports)
			} else if condition {
				fallenBrick := brick.fall()
				newBricks[fallenBrick.start.z] = append(newBricks[fallenBrick.start.z], fallenBrick)
				if fallenBrick.start.z != fallenBrick.end.z {
					newBricks[fallenBrick.end.z] = append(newBricks[fallenBrick.end.z], fallenBrick)
				}
				collectSupports(fallenBrick, supports)
			} else {
				newBricks[level] = append(newBricks[level], brick)
				if brick.start.z != brick.end.z {
					newBricks[brick.end.z] = append(newBricks[brick.end.z], brick)
				}
				collectSupports(brick, supports)
			}
		}
	}

	return newBricks, support, supportedBy, supportsMap
}

func isRemovable(supports []Brick, supportedBy map[Brick][]Brick) bool {
	for _, s := range supports {
		if len(supportedBy[s]) <= 1 {
			return false
		}
	}

	return true
}

func createFallCounter(supportMap, supportedBy map[Brick][]Brick, supportedByMap map[Brick]SupportMap, print bool) func(current Brick, supports []Brick) int {

	return func(current Brick, supports []Brick) int {
		if isRemovable(supports, supportedBy) {
			return 0
		}

		var brick Brick
		bricksToCheck := []Brick{}
		bricksToCheck = append(bricksToCheck, supports...)

		supportedByMap = removeFrom(current, supportMap[current], supportedByMap)
		rollbacks := []func() map[Brick]SupportMap{
			func() map[Brick]SupportMap {
				return addTo(current, supportMap[current], supportedByMap)
			},
		}

		count := 0
		counted := map[Brick]struct{}{current: {}}
		for len(bricksToCheck) > 0 {
			brick, bricksToCheck = common.Shift(bricksToCheck)

			if _, exists := counted[brick]; exists {
				continue
			}
			counted[brick] = struct{}{}

			if fits, _ := brick.fitsInMap(supportedByMap[brick].supports); fits {
				supportedBricks := supportMap[brick]
				supportedByMap = removeFrom(brick, supportedBricks, supportedByMap)
				bricksToCheck = append(bricksToCheck, supportedBricks...)

				slices.SortFunc(bricksToCheck, func(a, b Brick) int {
					return common.IntCompare(a.start.z, b.start.z)
				})

				b := brick
				rollbacks = append(rollbacks, func() map[Brick]SupportMap {
					return addTo(b, supportedBricks, supportedByMap)
				})

				count++
			}
		}

		for _, r := range rollbacks {
			supportedByMap = r()
		}

		return count
	}
}

func createFallCounterLeveled(supportMap, supportedBy map[Brick][]Brick, supportedByMap map[Brick]SupportMap, print bool) func(current Brick, supports []Brick) int {

	return func(current Brick, supports []Brick) int {
		if isRemovable(supports, supportedBy) {
			return 0
		}

		var brick Brick
		bricksToCheck := map[int][]Brick{}
		bricksToCheck[current.end.z] = append(bricksToCheck[current.end.z], supports...)

		supportedByMap = removeFrom(current, supportMap[current], supportedByMap)
		rollbacks := []func() map[Brick]SupportMap{
			func() map[Brick]SupportMap {
				return addTo(current, supportMap[current], supportedByMap)
			},
		}

		count := 0
		counted := map[Brick]struct{}{current: {}}
		currentLevel := current.start.z
		for len(bricksToCheck) > 0 {
			if len(bricksToCheck[currentLevel]) == 0 {
				delete(bricksToCheck, currentLevel)
				currentLevel++
				continue
			}

			brick, bricksToCheck[currentLevel] = common.Shift(bricksToCheck[currentLevel])

			if _, exists := counted[brick]; exists {
				continue
			}
			counted[brick] = struct{}{}

			if fits, _ := brick.fitsInMap(supportedByMap[brick].supports); fits {
				supportedBricks := supportMap[brick]
				supportedByMap = removeFrom(brick, supportedBricks, supportedByMap)
				bricksToCheck[brick.end.z] = append(bricksToCheck[brick.end.z], supportedBricks...)

				b := brick
				rollbacks = append(rollbacks, func() map[Brick]SupportMap {
					return addTo(b, supportedBricks, supportedByMap)
				})

				count++
			}

		}

		for _, r := range rollbacks {
			supportedByMap = r()
		}

		return count
	}
}

func removeFrom(brick Brick, supports []Brick, supportsMap map[Brick]SupportMap) map[Brick]SupportMap {
	for _, s := range supports {
		delete(supportsMap[s].supports, brick)
	}

	return supportsMap
}

func addTo(brick Brick, supports []Brick, supportsMap map[Brick]SupportMap) map[Brick]SupportMap {
	for _, s := range supports {
		supportsMap[s].supports[brick] = struct{}{}
	}

	return supportsMap
}

func createSupportCollector(supports, supportedBy *map[Brick][]Brick, supportsMap *map[Brick]SupportMap) func(brick Brick, s []Brick) {
	return func(brick Brick, bricks []Brick) {
		if _, exists := (*supportedBy)[brick]; !exists {
			(*supportedBy)[brick] = []Brick{}
		}

		if _, exists := (*supports)[brick]; !exists {
			(*supports)[brick] = []Brick{}
		}

		if _, exists := (*supportsMap)[brick]; !exists {
			(*supportsMap)[brick] = SupportMap{brick, map[Brick]struct{}{}}
		}

		for _, brickSupport := range bricks {
			if _, exists := (*supports)[brickSupport]; !exists {
				(*supports)[brickSupport] = []Brick{}
			}

			(*supports)[brickSupport] = append((*supports)[brickSupport], brick)
			(*supportedBy)[brick] = append((*supportedBy)[brick], brickSupport)
			(*supportsMap)[brick].supports[brickSupport] = struct{}{}
		}
	}
}

func parse() [][]Brick {
	bricks := [][]Brick{}
	common.ReadFile("day22/day22.txt", func(line string) {

		brick := parseBrick(line)
		if brickCount, maxHeight := len(bricks), max(brick.start.z, brick.end.z); maxHeight >= brickCount {
			bricks = append(bricks, make([][]Brick, (maxHeight-brickCount)+1)...)
		}

		bricks[brick.start.z] = append(bricks[brick.start.z], brick)
	})

	return bricks
}

func overlap(a, b Brick) bool {
	return rangeOverlap(a.getRangeX(), b.getRangeX()) && rangeOverlap(a.getRangeY(), b.getRangeY())
}

func canFall(a, b Brick) (bool, bool) {
	intersect := intersects(a, b)
	return (intersect && rangeDistance(a.getRangeZ(), b.getRangeZ()) > 1) || !intersect, intersect
}

func rangeOverlap(left, right Interval) bool {
	return rangeDistance(left, right) > 0
}

func rangeDistance(left, right Interval) int {
	return int(math.Abs(float64(min(left.start, left.end) - max(right.start, right.end))))
}

func parseBrick(line string) Brick {
	startValue, endValue := common.Split2(line, "~")
	start, end := parsePosition(startValue), parsePosition(endValue)
	return Brick{start, end}
}

func parsePosition(values string) Position {
	coordinates := strings.Split(values, ",")
	x, y, z := common.ParseInt(coordinates[0]), common.ParseInt(coordinates[1]), common.ParseInt(coordinates[2])

	return Position{x, y, z}
}

type Brick struct {
	start, end Position
}

type VisitBrick struct {
	brick         Brick
	disintegrated bool
}

type Position struct {
	x, y, z int
}

type Interval struct {
	start, end int
}

type cluster struct {
	Interval
}

func (this Brick) fitsIn(pile []Brick) (bool, []Brick) {
	supportedBy := []Brick{}

	state := true

	if min(this.start.z, this.end.z) == 1 {
		return false, supportedBy
	}

	for _, brick := range pile {
		if fall, overlap := canFall(this, brick); !fall && overlap {
			supportedBy = append(supportedBy, brick)
			state = false
		} else if !fall {
			state = false
		}
	}

	return state, supportedBy
}

func (this Brick) fitsInMap(pile map[Brick]struct{}) (bool, []Brick) {
	supportedBy := []Brick{}

	state := true

	if min(this.start.z, this.end.z) == 1 {
		return false, supportedBy
	}

	for brick := range pile {
		if fall, overlap := canFall(this, brick); !fall && overlap {
			supportedBy = append(supportedBy, brick)
			state = false
		} else if !fall {
			state = false
		}
	}

	return state, supportedBy
}

func (this *Brick) fall() Brick {
	this.start.z--
	this.end.z--
	return *this
}

func (this Brick) getRangeY() Interval {
	return Interval{this.start.y, this.end.y}
}

func (this Brick) getRangeX() Interval {
	return Interval{this.start.x, this.end.x}
}

func (this Brick) getRangeZ() Interval {
	return Interval{this.start.z, this.end.z}
}

func intersectsPrint(l1, l2 Brick, print bool) bool {
	P1, P2, P3, P4 := l1.start, l1.end, l2.start, l2.end
	Ax, Ay := P2.x-P1.x, P2.y-P1.y
	Bx, By := P3.x-P4.x, P3.y-P4.y
	Cx, Cy := P1.x-P3.x, P1.y-P3.y

	alphaNumerator := By*Cx - Bx*Cy
	betaNumerator := Ax*Cy - Ay*Cx
	denominator := Ay*Bx - Ax*By

	if denominator == 0 {
		if print {
			fmt.Println(alphaNumerator, betaNumerator, denominator, "h:", onLineX(P1, P2, P3, P4), onLineY(P1, P2, P3, P4))
		}
		return alphaNumerator == 0 && betaNumerator == 0 && (onLineX(P1, P2, P3, P4) && onLineY(P1, P2, P3, P4))
	}

	return between(float64(alphaNumerator)/float64(denominator), 0, 1) && between(float64(betaNumerator)/float64(denominator), 0, 1)
}

func intersects(l1, l2 Brick) bool {
	return intersectsPrint(l1, l2, false)
}

func onLineX(P1, P2, P3, P4 Position) bool {
	return ((P1.x >= P3.x && P1.x <= P4.x) || (P2.x >= P3.x && P2.x <= P4.x) ||
		(P3.x >= P1.x && P3.x <= P2.x) || (P4.x >= P1.x && P4.x <= P2.x))
}

func onLineY(P1, P2, P3, P4 Position) bool {
	return ((P1.y >= P3.y && P1.y <= P4.y) || (P2.y >= P3.y && P2.y <= P4.y) ||
		(P3.y >= P1.y && P3.y <= P2.y) || (P4.y >= P1.y && P4.y <= P2.y))
}

func between(n float64, start, end float64) bool {
	return n >= start && n <= end
}

func areTouching(l1, l2 Brick) bool {
	return isTouchingPoints(l1, l2) || isTouchingPoints(l2, l1)
}

func isTouchingPoints(left, right Brick) bool {
	return isOnLine(left, right.start) || isOnLine(left, right.end)
}

func isOnLine(line Brick, point Position) bool {
	return (point.x <= line.end.x && point.x >= line.start.x) ||
		(point.y <= line.end.y && point.y >= line.start.y)
}

func inInterval(numerator, denominator int) bool {
	if denominator > 0 {
		if numerator < 0 || numerator > denominator {
			return false
		}
	} else {
		if numerator > 0 || numerator < denominator {
			return false
		}
	}

	return true
}

type SupportMap struct {
	brick    Brick
	supports map[Brick]struct{}
}

func (this *SupportMap) setSupports(m map[Brick]struct{}) {
	this.supports = m
}
