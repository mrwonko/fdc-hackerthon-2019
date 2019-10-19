package rules

import (
	"math"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/gamemath"
)

func ETA(from, to gamemath.Point) int {
	return int(math.Ceil(to.Minus(from).Magnitude()))
}

func ReachablePlanets(planets []*Planet, from gamemath.Point, distance int) []PlanetID {
	res := make([]PlanetID, 0, len(planets))
	for _, p := range planets {
		squaredDistance := gamemath.Point{p.X, p.Y}.Minus(from).SquaredMagnitude()
		if squaredDistance <= distance*distance {
			res = append(res, p.ID)
		}
	}
	return res
}
