package rules

import (
	"context"
	"sort"
)

func Advance(ctx context.Context, cur *Gamestate, rounds int) (*Gamestate, error) {
	state := *cur
	// copy slices we modify
	state.Planets = make([]Planet, len(cur.Planets))
	copy(state.Planets, cur.Planets)
	state.Fleets = make([]Fleet, len(cur.Fleets))
	copy(state.Fleets, cur.Fleets)

	// instead of calculating every single step, find
	// TODO: benchmark if/when this is actually faster
	sort.Sort(FleetByETA(state.Fleets))
	t := state.Round
	steps := make([]int, 0, len(state.Fleets))
	for i := range state.Fleets {
		f := &state.Fleets[i]
		if f.ETA > state.Round+rounds {
			break // we don't care about anything that far into the future
		}
		if f.ETA == t {
			continue // multiple fleets arriving simultaneously
		}
		steps = append(steps, f.ETA-t)
		t = f.ETA
	}
	if t < state.Round+rounds {
		steps = append(steps, (state.Round+rounds)-t)
		t = state.Round + rounds
	}

	sppByPlanet := make([]ShipsPerPlayer, len(state.Planets))
	for _, step := range steps {
		// consider timeout
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		// produce
		for i := range state.Planets {
			p := &state.Planets[i]
			if p.Owner != NeutralPlayer {
				p.Ships = p.Ships.Add(p.Production.Mul(step))
			}
		}

		// find all fleets arriving this tick
		for i := range sppByPlanet {
			sppByPlanet[i].Reset()
		}
		for i := 0; i < len(state.Fleets); {
			f := &state.Fleets[i]
			if f.ETA <= state.Round {
				sppByPlanet[f.TargetIndex].Add(f)
				// move the last fleet into this one to delete it
				*f = state.Fleets[len(state.Fleets)-1]
				state.Fleets = state.Fleets[:len(state.Fleets)-1]
				continue
			}
			i++
		}

		// land
		for i := range sppByPlanet {
			for _, id := range state.TurnOrder {
				ships := sppByPlanet[i][id.ToIndex()]
				if ships.Dead() {
					continue // nothing landed here
				}
				p := &state.Planets[i]
				if id == p.Owner { // same owner reinforces
					p.Ships = p.Ships.Add(ships)
				} else { // different owner battles
					p.Ships, ships = Battle(p.Ships, ships)
					if p.Ships.Dead() {
						p.Ships = ships
						p.Owner = id
					}
				}
			}
		}
		state.Round += step
	}
	return &state, nil
}

type FleetByETA []Fleet

func (f FleetByETA) Len() int {
	return len(f)
}

func (f FleetByETA) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f FleetByETA) Less(i, j int) bool {
	return f[i].ETA < f[j].ETA
}
