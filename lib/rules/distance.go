package rules

import (
	"math"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/gamemath"
)

func ETA(from, to gamemath.Point) int {
	return int(math.Ceil(to.Minus(from).Magnitude()))
}

func ReachablePlanets(planets []*Planet, from gamemath.Point, distance int) []int {
	res := make([]int, 0, len(planets))
	for i, p := range planets {
		squaredDistance := p.Coords.Minus(from).SquaredMagnitude()
		if squaredDistance <= distance*distance {
			res = append(res, i)
		}
	}
	return res
}
