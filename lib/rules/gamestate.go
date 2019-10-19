package rules

import "github.com/mrwonko/fdc-hackerthon-2019/lib/gamemath"

const (
	numShipTypes = 3

	NeutralPlayer PlayerID = -1
	MyPlayer      PlayerID = 0
	EnemyPlayer   PlayerID = 1
)

type (
	PlayerID  int
	PlanetID  int
	FleetID   int
	ShipCount [numShipTypes]int

	Gamestate struct {
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
