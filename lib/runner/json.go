package runner

import (
	"sort"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/gamemath"
	"github.com/mrwonko/fdc-hackerthon-2019/lib/rules"
)

type (
	playerID int

	Gamestate struct {
		GameOver  bool      `json:"game_over"`
		Winner    *playerID `json:"winner"`
		Round     int       `json:"round"`
		MaxRounds int       `json:"max_rounds"` // 500, usually
		Fleets    []Fleet   `json:"fleets"`
		Players   []player  `json:"players"`
		Planets   []planet  `json:"planets"`
	}

	Fleet struct {
		ID     rules.FleetID   `json:"ID"`
		Owner  playerID        `json:"owner_id"`
		Origin rules.PlanetID  `json:"origin"`
		Target rules.PlanetID  `json:"target"`
		Ships  rules.ShipCount `json:"ships"`
		ETA    int             `json:"eta"` // in rounds
	}

	player struct {
		ID    playerID `json:"id"`
		Name  string   `json:"name"`
		ItsMe bool     `json:"itsme"`
	}

	planet struct {
		ID         rules.PlanetID  `json:"id"`
		Owner      playerID        `json:"owner_id"`
		X          int             `json:"x"`
		Y          int             `json:"y"`
		Ships      rules.ShipCount `json:"ships"`
		Production rules.ShipCount `json:"production"`
	}
)

func (gs Gamestate) Preprocess() *rules.Gamestate {
	// FIXME: this assumes 2 players
	playerIDLookup := map[playerID]rules.PlayerID{
		0: rules.NeutralPlayer,
	}
	playerIDs := []int{} // as int for easier sorting
	for i := range gs.Players {
		p := &gs.Players[i]
		playerIDs = append(playerIDs, int(p.ID))
		if p.ItsMe {
			playerIDLookup[p.ID] = rules.MyPlayer
		} else {
			playerIDLookup[p.ID] = rules.EnemyPlayer
		}
	}
	sort.Ints(playerIDs)
	turnOrder := make([]rules.PlayerID, len(playerIDs))
	for i, id := range playerIDs {
		turnOrder[i] = playerIDLookup[playerID(id)]
	}

	planetIDLookup := make(map[rules.PlanetID]int, len(gs.Planets))
	planets := make([]rules.Planet, len(gs.Planets))
	for i := range gs.Planets {
		p := &gs.Planets[i]
		planetIDLookup[p.ID] = i
		planets[i] = rules.Planet{
			ID:         p.ID,
			Owner:      playerIDLookup[p.Owner],
			Coords:     gamemath.Point{p.X, p.Y},
			Ships:      p.Ships,
			Production: p.Production,
		}
	}
	fleets := make([]rules.Fleet, len(gs.Fleets))
	for i := range gs.Fleets {
		f := &gs.Fleets[i]
		fleets[i] = rules.Fleet{
			ID:          f.ID,
			Owner:       playerIDLookup[f.Owner],
			OriginIndex: planetIDLookup[f.Origin],
			TargetIndex: planetIDLookup[f.Target],
			Ships:       f.Ships,
			ETA:         f.ETA,
		}
	}
	return &rules.Gamestate{
		TurnOrder: turnOrder,
		Round:     gs.Round,
		MaxRounds: gs.MaxRounds,
		Planets:   planets,
		Fleets:    fleets,
	}
}
