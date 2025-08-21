package main

import (
	"common"
	"fmt"
	"math"
	"math/big"
	"strings"
)

var BigInt = big.NewInt(0)
var BigFloat = big.NewFloat(0)

func main() {
	Solution1()
	Solution2()
}

func Solution1() {
	lines := parse()

	// minX, minY, maxX, maxY := big.NewFloat(7), big.NewFloat(7), big.NewFloat(27), big.NewFloat(27)
	minX, minY, maxX, maxY := big.NewFloat(200000000000000), big.NewFloat(200000000000000), big.NewFloat(400000000000000), big.NewFloat(400000000000000)

	intersections := []Point{}

	for lv, nextLines := common.Shift(lines); len(nextLines) > 0; lv, nextLines = common.Shift(nextLines) {
		line := lv.line
		for _, lv1 := range nextLines {
			otherLine := lv1.line

			if intersectionPoint, intersects := lineIntersection(line, otherLine); intersects {
				nanoSecondA := getNano(lv, intersectionPoint)
				nanoSecondB := getNano(lv1, intersectionPoint)

				if intersectionPoint.x.Cmp(minX) >= 0 &&
					intersectionPoint.x.Cmp(maxX) <= 0 &&
					intersectionPoint.y.Cmp(minY) >= 0 &&
					intersectionPoint.y.Cmp(maxY) <= 0 &&
					nanoSecondA.Cmp(big.NewFloat(0)) >= 0 &&
					nanoSecondB.Cmp(big.NewFloat(0)) >= 0 {
					intersections = append(intersections, intersectionPoint)
				}
			}
		}
	}

	fmt.Println(len(intersections))
}

func Solution2() {
	lines := parse()

	for j := 1; j < len(lines)-3; j++ {
		A := [][]float64{}
		b := []float64{}

		for i := j; i < j+3; i++ {
			matrixPart1, matrixPart2 := buildCoefficientMatrixPart(lines[i-1], lines[i])
			resultPart := buildConstantTermMatrix(lines[i-1], lines[i])

			A = append(A, matrixPart1, matrixPart2)
			b = append(b, resultPart...)
		}

		solution := solveSystem(b, A)
		x, y, z := solution[0], solution[1], solution[2]

		if math.Mod(x, 1) <= 0.05 && math.Mod(y, 1) <= 0.05 && math.Mod(z, 1) <= 0.05 {
			fmt.Println(int64(math.Floor(x) + math.Floor(y) + math.Floor(z)))
			break
		}
	}
}

// 19 - 2 - 2 = 15
// x = 19 nx = 16 vx = -2  / s
// x + (s * vx) = nx
// s * vx = nx - x
// s = (nx - x) / vx

func getNano(l LineVelocity, p Point) *big.Float {
	return new(big.Float).Quo(new(big.Float).Sub(p.x, l.line.p1.x), l.velocity.x)
}

// func getNano(l LineVelocity, p Point) float64 {
// 	return (p.x - l.line.p1.x) / l.velocity.x
// }

func parse() []LineVelocity {
	lines := []LineVelocity{}

	common.ReadFile("day24/day24.txt", func(line string) {
		point, velocity := common.Split2(line, " @ ")
		pointCoords := strings.Split(point, ", ")
		velocityValues := strings.Split(velocity, ", ")

		// x, y, z := common.ParseInt64(pointCoords[0]), common.ParseInt64(pointCoords[1]), common.ParseInt64(pointCoords[2])
		// vx, vy, vz := common.ParseInt64(trimWhitespace(velocityValues[0])), common.ParseInt64(trimWhitespace(velocityValues[1])), common.ParseInt64(trimWhitespace(velocityValues[2]))

		x, y, z := parseBigFloat(pointCoords[0]), parseBigFloat(pointCoords[1]), parseBigFloat(pointCoords[2])

		vx := parseBigFloat(trimWhitespace(velocityValues[0]))
		vy := parseBigFloat(trimWhitespace(velocityValues[1]))
		vz := parseBigFloat(trimWhitespace(velocityValues[2]))

		px, py, pz := floatToInt(x), floatToInt(y), floatToInt(z)
		pxv, pyv, pzv := floatToInt(vx), floatToInt(vy), floatToInt(vz)

		base, exp := big.NewInt(10), big.NewInt(90)

		px2 := intToFloat(new(big.Int).Add(px, new(big.Int).Mul(pxv, new(big.Int).Exp(base, exp, nil))))
		py2 := intToFloat(new(big.Int).Add(py, new(big.Int).Mul(pyv, new(big.Int).Exp(base, exp, nil))))
		pz2 := intToFloat(new(big.Int).Add(pz, new(big.Int).Mul(pzv, new(big.Int).Exp(base, exp, nil))))

		p1 := Point{x, y, z}
		p2 := Point{px2, py2, pz2}

		lines = append(lines, LineVelocity{Line{p1, p2}, Point{vx, vy, vz}})
	})

	return lines
}

func parseBigFloat(value string) *big.Float {
	v, _, _ := new(big.Float).Parse(value, 10)
	return v
}

func floatToInt(bigFloat *big.Float) *big.Int {
	result := new(big.Int)
	bigFloat.Int(result)
	return result
}

func intToFloat(bigInt *big.Int) *big.Float {
	return big.NewFloat(0).SetInt(bigInt)
}

func trimWhitespace(s string) string {
	return strings.Trim(s, " ")
}

func lineIntersection(l1, l2 Line) (Point, bool) {
	a1 := new(big.Float).Sub(l1.p1.y, l1.p2.y)
	b1 := new(big.Float).Sub(l1.p1.x, l1.p2.x)
	c1 := new(big.Float).Sub(new(big.Float).Mul(l1.p1.x, l1.p2.y), new(big.Float).Mul(l1.p2.x, l1.p1.y))

	a2 := new(big.Float).Sub(l2.p1.y, l2.p2.y)
	b2 := new(big.Float).Sub(l2.p1.x, l2.p2.x)
	c2 := new(big.Float).Sub(new(big.Float).Mul(l2.p1.x, l2.p2.y), new(big.Float).Mul(l2.p2.x, l2.p1.y))

	px_n := new(big.Float).Sub(new(big.Float).Mul(c1, b2), new(big.Float).Mul(b1, c2))
	py_n := new(big.Float).Sub(new(big.Float).Mul(c1, a2), new(big.Float).Mul(a1, c2))
	denominator := new(big.Float).Sub(new(big.Float).Mul(b1, a2), new(big.Float).Mul(a1, b2))

	if denominator.Cmp(big.NewFloat(0)) == 0 {
		return Point{}, false
	}

	return Point{new(big.Float).Quo(px_n, denominator), new(big.Float).Quo(py_n, denominator), big.NewFloat(0)}, true
}

func solveSystem(results []float64, matrix [][]float64) []float64 {
	nominators := generateCramerNominators(results, matrix)
	denominator := common.Det(matrix)
	solutions := []float64{}

	for _, nominator := range nominators {
		solutions = append(solutions, common.Det(nominator)/denominator)
	}

	return solutions
}

func generateCramerNominators(column []float64, matrix [][]float64) [][][]float64 {
	result := [][][]float64{}

	for i := 0; i < len(matrix); i++ {
		result = append(result, replaceMatrixColumn(i, column, matrix))
	}

	return result
}

func replaceMatrixColumn(replaceIndex int, column []float64, matrix [][]float64) [][]float64 {
	result := [][]float64{}
	for rowIndex := 0; rowIndex < len(matrix); rowIndex++ {
		row := []float64{}
		for colIndex := 0; colIndex < len(matrix[rowIndex]); colIndex++ {
			if colIndex == replaceIndex {
				row = append(row, column[rowIndex])
			} else {
				row = append(row, matrix[rowIndex][colIndex])
			}
		}
		result = append(result, row)
	}

	return result
}

func buildCoefficientMatrixPart(
	left LineVelocity,
	right LineVelocity,
) ([]float64, []float64) {
	l1, v1 := left.line, left.velocity
	x1, y1, z1 := l1.p1.x, l1.p1.y, l1.p1.z
	vx1, vy1, vz1 := v1.x, v1.y, v1.z

	l2, v2 := right.line, right.velocity
	x2, y2, z2 := l2.p1.x, l2.p1.y, l2.p1.z
	vx2, vy2, vz2 := v2.x, v2.y, v2.z

	return []float64{
			toFloat(new(big.Float).Sub(vy1, vy2)),
			toFloat(new(big.Float).Sub(vx2, vx1)),
			0,
			toFloat(new(big.Float).Sub(y2, y1)), toFloat(new(big.Float).Sub(x1, x2)),
			0,
		},
		[]float64{
			toFloat(new(big.Float).Sub(vz1, vz2)),
			0,
			toFloat(new(big.Float).Sub(vx2, vx1)),
			toFloat(new(big.Float).Sub(z2, z1)),
			0,
			toFloat(new(big.Float).Sub(x1, x2)),
		}
}

func toFloat(bigFloat *big.Float) float64 {
	float, _ := bigFloat.Float64()
	return float
}

// func buildCoefficientMatrixPart(
// 	left LineVelocity,
// 	right LineVelocity,
// ) ([]*big.Float, []*big.Float) {
// 	l1, v1 := left.line, left.velocity
// 	x1, y1, z1 := l1.p1.x, l1.p1.y, l1.p1.z
// 	vx1, vy1, vz1 := v1.x, v1.y, v1.z

// 	l2, v2 := right.line, right.velocity
// 	x2, y2, z2 := l2.p1.x, l2.p1.y, l2.p1.z
// 	vx2, vy2, vz2 := v2.x, v2.y, v2.z

// 	return []*big.Float{new(big.Float).Sub(vy1, vy2), new(big.Float).Sub(vx2, vx1), big.NewFloat(0), new(big.Float).Sub(y2, y1), new(big.Float).Sub(x1, x2), big.NewFloat(0)},
// 		[]*big.Float{new(big.Float).Sub(vz1, vz2), big.NewFloat(0), new(big.Float).Sub(vx2, vx1), new(big.Float).Sub(z2, z1), big.NewFloat(0), new(big.Float).Sub(x1, x2)}
// }

// func buildCoefficientMatrixPart(
// 	left LineVelocity,
// 	right LineVelocity,
// ) ([]*big.Float, []*big.Float) {
// 	l1, v1 := left.line, left.velocity
// 	x1, y1, z1 := l1.p1.x, l1.p1.y, l1.p1.z
// 	vx1, vy1, vz1 := v1.x, v1.y, v1.z

// 	l2, v2 := right.line, right.velocity
// 	x2, y2, z2 := l2.p1.x, l2.p1.y, l2.p1.z
// 	vx2, vy2, vz2 := v2.x, v2.y, v2.z

// 	return []float64{vy1 - vy2, vx2 - vx1, 0, y2 - y1, x1 - x2, 0}, []float64{vz1 - vz2, 0, vx2 - vx1, z2 - z1, 0, x1 - x2}
// }

func buildConstantTermMatrix(
	left LineVelocity,
	right LineVelocity,
) []float64 {
	l1, v1 := left.line, left.velocity
	x1, y1, z1 := l1.p1.x, l1.p1.y, l1.p1.z
	vx1, vy1, vz1 := v1.x, v1.y, v1.z

	l2, v2 := right.line, right.velocity
	x2, y2, z2 := l2.p1.x, l2.p1.y, l2.p1.z
	vx2, vy2, vz2 := v2.x, v2.y, v2.z

	return []float64{
		toFloat(new(big.Float).Add(new(big.Float).Sub(new(big.Float).Sub(new(big.Float).Mul(x1, vy1), new(big.Float).Mul(y1, vx1)), new(big.Float).Mul(x2, vy2)), new(big.Float).Mul(y2, vx2))),
		toFloat(new(big.Float).Add(new(big.Float).Sub(new(big.Float).Sub(new(big.Float).Mul(x1, vz1), new(big.Float).Mul(z1, vx1)), new(big.Float).Mul(x2, vz2)), new(big.Float).Mul(z2, vx2))),
	}
}

// func buildConstantTermMatrix(
// 	left LineVelocity,
// 	right LineVelocity,
// ) []float64 {
// 	l1, v1 := left.line, left.velocity
// 	x1, y1, z1 := l1.p1.x, l1.p1.y, l1.p1.z
// 	vx1, vy1, vz1 := v1.x, v1.y, v1.z

// 	l2, v2 := right.line, right.velocity
// 	x2, y2, z2 := l2.p1.x, l2.p1.y, l2.p1.z
// 	vx2, vy2, vz2 := v2.x, v2.y, v2.z

// 	return []float64{
// 		x1*vy1 - y1*vx1 - x2*vy2 + y2*vx2,
// 		x1*vz1 - z1*vx1 - x2*vz2 + z2*vx2,
// 	}
// }

// func lineIntersection(l1, l2 Line) (Point, bool) {
// 	a1 := l1.p1.y - l1.p2.y
// 	b1 := l1.p1.x - l1.p2.x
// 	c1 := l1.p1.x*l1.p2.y - l1.p2.x*l1.p1.y

// 	a2 := l2.p1.y - l2.p2.y
// 	b2 := l2.p1.x - l2.p2.x
// 	c2 := l2.p1.x*l2.p2.y - l2.p2.x*l2.p1.y

// 	px_n := float64(c1*b2 - b1*c2)
// 	py_n := float64(c1*a2 - a1*c2)
// 	denominator := float64(b1*a2 - a1*b2)

// 	if denominator == 0 {
// 		return Point{}, false
// 	}

// 	return Point{px_n / denominator, py_n / denominator, 0}, true
// }

// func segmentIntersection(l1, l2 Segment) (Point, bool) {
// 	P1, P2, P3, P4 := l1.start, l1.end, l2.start, l2.end
// 	Ax, Ay := P2.x-P1.x, P2.y-P1.y
// 	Bx, By := P3.x-P4.x, P3.y-P4.y
// 	Cx, Cy := P1.x-P3.x, P1.y-P3.y

// 	alphaNumerator := By*Cx - Bx*Cy
// 	betaNumerator := Ax*Cy - Ay*Cx
// 	denominator := Ay*Bx - Ax*By

// 	if denominator == 0 {
// 		return Point{}, alphaNumerator == 0 && betaNumerator == 0 && (onLineX(P1, P2, P3, P4) && onLineY(P1, P2, P3, P4))
// 	}

// 	px, py := float64(alphaNumerator)/float64(denominator), float64(betaNumerator)/float64(denominator)

// 	return Point{px, py, 0}, between(px, 0, 1) && between(py, 0, 1)
// }

// func between(n float64, start, end float64) bool {
// 	return n >= start && n <= end
// }

// func onLineX(P1, P2, P3, P4 Point) bool {
// 	return ((P1.x >= P3.x && P1.x <= P4.x) || (P2.x >= P3.x && P2.x <= P4.x) ||
// 		(P3.x >= P1.x && P3.x <= P2.x) || (P4.x >= P1.x && P4.x <= P2.x))
// }

// func onLineY(P1, P2, P3, P4 Point) bool {
// 	return ((P1.y >= P3.y && P1.y <= P4.y) || (P2.y >= P3.y && P2.y <= P4.y) ||
// 		(P3.y >= P1.y && P3.y <= P2.y) || (P4.y >= P1.y && P4.y <= P2.y))
// }

type Point struct {
	x, y, z *big.Float
}

// type Point struct {
// 	x, y, z float64
// }

type Line struct {
	p1, p2 Point
}

type LineVelocity struct {
	line     Line
	velocity Point
}

type Segment struct {
	start, end Point
}
