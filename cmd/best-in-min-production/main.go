package main

import (
	"log"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/rules"
	"github.com/mrwonko/fdc-hackerthon-2019/lib/runner"
)

const (
	buckets = 5 // 0, 25, 50, 75, 100
	maxUint = ^uint(0)
	maxInt  = int(maxUint >> 1)
)

func main() {
	runner.Main(&runner.Player{
		Name:     "ws-eris",
		Password: "23B2A1D7-A77F-47CE-8045-C1BA0502A3CD",
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
		projectedFleet := finalState.FleetSize()[rules.EnemyPlayer.ToIndex()].Total()
		projectedProduction := finalState.TotalProduction()[rules.EnemyPlayer.ToIndex()].Total()

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
		for i := range finalState.Planets {
			p := &finalState.Planets[i]
			if p.Owner != rules.MyPlayer {
				potentialTargets = append(potentialTargets, i)
			}
		}
		if len(potentialTargets) == 0 {
			tick.Move <- rules.Nop
			continue
		}
		type candidate struct {
			src        int
			dst        int
			production int
			fleet      int
		}

		best := candidate{
			src:        potentialSources[0],
			dst:        potentialTargets[0],
			production: maxInt,
			fleet:      maxInt,
		}

		func() {
			i := 0
			for _, src := range potentialSources {
				for _, dst := range potentialTargets {
					i++
					maxShips := tick.Gamestate.Planets[src].Ships
					for ships := range rules.GenerateShipSubsets(tick.Ctx, maxShips, buckets) {
						move := tick.Gamestate.Send(src, dst, ships)
						moveResult, err := rules.Advance(tick.Ctx, move, remainingRounds)
						if err != nil {
							log.Printf("round %d: predicting %d/%d: %s", tick.Gamestate.Round, i, len(potentialSources)*len(potentialTargets), err)
							return
						}
						cur := candidate{
							src:        src,
							dst:        dst,
							production: moveResult.TotalProduction()[rules.EnemyPlayer.ToIndex()].Total(),
							fleet:      moveResult.FleetSize()[rules.EnemyPlayer.ToIndex()].Total(),
						}
						if cur.production < best.production ||
							cur.production == best.production && cur.fleet < best.fleet {
							best = cur
						}
					}
				}
			}
		}()
		if best.production < projectedProduction ||
			best.production == projectedProduction && best.fleet < projectedFleet {
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
