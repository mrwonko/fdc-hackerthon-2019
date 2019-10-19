package main

import (
	"log"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/rules"
	"github.com/mrwonko/fdc-hackerthon-2019/lib/runner"
)

func main() {
	runner.Main(&runner.Player{
		Name:     "ws-banshee",
		Password: "D7740F2A-9112-4349-9E4D-9A80DFA152C5",
		Play:     play,
	})
}

func play(ticks <-chan *runner.Tick) {
	for tick := range ticks {
		remainingRounds := tick.Gamestate.MaxRounds - tick.Gamestate.Round
		finalState, err := rules.Advance(tick.Ctx, tick.Gamestate, remainingRounds)
		if err != nil {
			log.Printf("round %d initial projection: %s", tick.Gamestate.Round, err)
			tick.Move <- rules.Nop
			continue
		}
		projectedFleet := finalState.FleetSize()[rules.MyPlayer.ToIndex()].Total()

		// consider all of our planets as ship sources
		potentialSources := make([]int, 0, len(tick.Gamestate.Planets))
		for i := range tick.Gamestate.Planets {
			p := &tick.Gamestate.Planets[i]
			if p.Owner == rules.MyPlayer {
				potentialSources = append(potentialSources, i)
			}
		}
		if len(potentialSources) == 0 { // oh no
			tick.Move <- rules.Nop
			continue
		}

		// try to send them to a planet we wouldn't otherwise get
		potentialTargets := make([]int, 0, len(finalState.Planets))
		projectedPlanets := 0
		for i := range finalState.Planets {
			p := &finalState.Planets[i]
			if p.Owner != rules.MyPlayer {
				potentialTargets = append(potentialTargets, i)
			} else {
				projectedPlanets++
			}
		}
		if len(potentialTargets) == 0 {
			tick.Move <- rules.Nop
			continue
		}
		type candidate struct {
			src     int
			dst     int
			planets int
			fleet   int
		}

		best := candidate{
			src: potentialSources[0],
			dst: potentialTargets[0],
		}

		func() {
			i := 0
			for _, src := range potentialSources {
				for _, dst := range potentialTargets {
					i++
					move := tick.Gamestate.Send(src, dst, tick.Gamestate.Planets[src].Ships)
					moveResult, err := rules.Advance(tick.Ctx, move, remainingRounds)
					if err != nil {
						log.Printf("round %d: predicting %d/%d: %s", tick.Gamestate.Round, i, len(potentialSources)*len(potentialTargets), err)
						return
					}
					cur := candidate{
						src:     src,
						dst:     dst,
						planets: moveResult.NumPlanets(rules.MyPlayer),
						fleet:   moveResult.FleetSize()[rules.MyPlayer.ToIndex()].Total(),
					}
					if cur.fleet > best.fleet ||
						cur.fleet == best.fleet && cur.planets > best.planets {
						best = cur
					}
				}
			}
		}()
		if best.fleet > projectedFleet ||
			best.fleet == projectedFleet && best.planets > projectedPlanets {
			tick.Move <- &rules.Send{
				Src:   tick.Gamestate.Planets[best.src].ID,
				Dst:   tick.Gamestate.Planets[best.dst].ID,
				Ships: tick.Gamestate.Planets[best.src].Ships,
			}
		} else {
			tick.Move <- rules.Nop
		}
	}
}
