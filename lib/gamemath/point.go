package gamemath

import (
	"math"
)

type Point [2]int

func (p Point) Minus(other Point) Point {
	return Point{
		p[0] - other[0],
		p[1] - other[1],
	}
}

func (p Point) SquaredMagnitude() int {
	return p[0]*p[0] + p[1]*p[1]
}

func (p Point) Magnitude() float64 {
	return math.Sqrt(float64(p.SquaredMagnitude()))
}
