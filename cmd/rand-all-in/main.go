package main

import (
	"math/rand"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/rules"
	"github.com/mrwonko/fdc-hackerthon-2019/lib/runner"
)

func main() {
	runner.Main(&runner.Player{
		Name:     "ws-rand-all-in",
		Password: "94A6BC3B-66AE-4EB6-BE1D-EEF0D48E06AB",
		Play:     playRandAllIn,
	})
}

func playRandAllIn(ticks <-chan *runner.Tick) {
	for tick := range ticks {
		myPlanets := make([]rules.PlanetID, 0, len(tick.Gamestate.Planets))
		otherPlanets := make([]rules.PlanetID, 0, len(tick.Gamestate.Planets))
		for i := range tick.Gamestate.Planets {
			id := tick.Gamestate.Planets[i].ID
			if tick.Gamestate.Planets[i].Owner == rules.MyPlayer {
				myPlanets = append(myPlanets, id)
			} else {
				otherPlanets = append(otherPlanets, id)
			}
		}
		if len(otherPlanets) == 0 || len(myPlanets) == 0 {
			tick.Move <- rules.Nop
			continue
		}
		src := myPlanets[rand.Intn(len(myPlanets))]
		tick.Move <- &rules.Send{
			Src:   src,
			Dst:   otherPlanets[rand.Intn(len(otherPlanets))],
			Ships: tick.Gamestate.Planets[src].Ships,
		}
	}
}
