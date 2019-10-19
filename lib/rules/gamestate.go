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
	ShipsPerPlayer [2]ShipCount
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

func (sc ShipCount) Sub(other ShipCount) ShipCount {
	var res ShipCount
	for i, c := range sc {
		res[i] = c - other[i]
	}
	return res
}

func (sc ShipCount) Mul(f int) ShipCount {
	var res ShipCount
	for i, c := range sc {
		res[i] = c * f
	}
	return res
}

func (sc ShipCount) Mulf(f float64) ShipCount {
	var res ShipCount
	for i, c := range sc {
		res[i] = int(float64(c) * f)
	}
	return res
}

func (sc ShipCount) Total() int {
	res := 0
	for _, c := range sc {
		res += c
	}
	return res
}

func (ifs *ShipsPerPlayer) Reset() {
	*ifs = ShipsPerPlayer{}
}

func (ifs *ShipsPerPlayer) Add(f *Fleet) {
	i := f.Owner.ToIndex()
	ifs[i] = ifs[i].Add(f.Ships)
}

func (gs *Gamestate) TotalProduction() ShipsPerPlayer {
	res := ShipsPerPlayer{}
	for i := range gs.Planets {
		p := &gs.Planets[i]
		if p.Owner != NeutralPlayer {
			idx := p.Owner.ToIndex()
			res[idx] = res[idx].Add(p.Production)
		}
	}
	return res
}

func (gs *Gamestate) FleetSize() ShipsPerPlayer {
	var res ShipsPerPlayer
	for i := range gs.Planets {
		p := &gs.Planets[i]
		if p.Owner == NeutralPlayer {
			continue
		}
		idx := p.Owner.ToIndex()
		res[idx] = res[idx].Add(p.Ships)
	}
	for i := range gs.Fleets {
		res.Add(&gs.Fleets[i])
	}
	return res
}

func (gs *Gamestate) Send(src, dst int, ships ShipCount) *Gamestate {
	res := *gs
	res.Fleets = make([]Fleet, len(gs.Fleets), len(gs.Fleets)+1)
	copy(res.Fleets, gs.Fleets)
	res.Planets = make([]Planet, len(gs.Planets))
	copy(res.Planets, gs.Planets)
	p := &res.Planets[src]
	p.Ships = p.Ships.Sub(ships)
	var lastID FleetID
	for i := range res.Fleets {
		f := &res.Fleets[i]
		if f.ID > lastID {
			lastID = f.ID
		}
	}
	res.Fleets = append(res.Fleets, Fleet{
		ID:          lastID + 1,
		Owner:       p.Owner,
		ETA:         res.Round + ETA(p.Coords, res.Planets[dst].Coords),
		OriginIndex: src,
		TargetIndex: dst,
		Ships:       ships,
	})
	return &res
}

func (gs *Gamestate) NumPlanets(owner PlayerID) int {
	c := 0
	for i := range gs.Planets {
		p := &gs.Planets[i]
		if p.Owner == owner {
			c++
		}
	}
	return c
}
