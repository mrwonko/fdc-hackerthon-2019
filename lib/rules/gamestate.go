package rules

type (
	PlayerID  int
	PlanetID  int
	FleetID   int
	ShipCount [3]int

	Gamestate struct {
		GameOver  bool      `json:"game_over"`
		Winner    *PlayerID `json:"winner"`
		Round     int       `json:"round"`
		MaxRounds int       `json:"max_rounds"` // 500, usually
		Fleets    []Fleet   `json:"fleets"`
		Players   []Player  `json:"players"`
		Planets   []Planet  `json:"planets"`
	}

	Fleet struct {
		ID     FleetID   `json:"ID"`
		Owner  PlayerID  `json:"owner_id"`
		Origin PlanetID  `json:"origin"`
		Target PlanetID  `json:"target"`
		Ships  ShipCount `json:"ships"`
		ETA    int       `json:"eta"` // in rounds
	}

	Player struct {
		ID    PlayerID `json:"id"`
		Name  string   `json:"name"`
		ItsMe bool     `json:"itsme"`
	}

	Planet struct {
		ID         PlanetID  `json:"id"`
		Owner      PlayerID  `json:"owner_id"`
		X          int       `json:"x"`
		Y          int       `json:"y"`
		Ships      ShipCount `json:"ships"`
		Production ShipCount `json:"production"`
	}
)

const (
	Neutral PlayerID = 0
)
