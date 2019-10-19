package rules

import "github.com/mrwonko/fdc-hackerthon-2019/lib/gamemath"

const (
	numShipTypes = 3

	NeutralPlayer PlayerID = 0
	MyPlayer      PlayerID = 1
	EnemyPlayer   PlayerID = 2
)

func (p PlayerID) ToIndex() int { // excluding neutral player
	return int(p - MyPlayer)
}

type (
	PlayerID  int
	PlanetID  int
	FleetID   int
	ShipCount [numShipTypes]int

	Gamestate struct {
		TurnOrder []PlayerID
		// we don't care about game over / winner
		Round     int
		MaxRounds int // 500, usually
		Fleets    []Fleet
		Planets   []Planet
	}

	Fleet struct {
		ID          FleetID
		Owner       PlayerID
		OriginIndex int
		TargetIndex int
		Ships       ShipCount
		ETA         int // in rounds
	}

	Planet struct {
		ID         PlanetID
		Owner      PlayerID
		Coords     gamemath.Point
		Ships      ShipCount
		Production ShipCount
	}
)

func (sc ShipCount) Dead() bool {
	for _, c := range sc {
		if c > 0 {
			return false
		}
	}
	return true
}

func (sc ShipCount) Add(other ShipCount) ShipCount {
	var res ShipCount
	for i, c := range sc {
		res[i] = c + other[i]
	}
	return res
}

func (sc ShipCount) Mul(i int) ShipCount {
	var res ShipCount
	for i, c := range sc {
		res[i] = c * i
	}
	return res
}
