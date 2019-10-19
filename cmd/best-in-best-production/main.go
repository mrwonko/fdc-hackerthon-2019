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
	minInt  = -maxInt - 1
)

func main() {
	runner.Main(&runner.Player{
		Name:     "ws-petra",
		Password: "F792CF89-442A-4A6E-8F40-E5FCD407672B",
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
		projectedProductionDiff := scoreProduction(finalState.TotalProduction())

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
			src            int
			dst            int
			productionDiff int
			fleet          int
		}

		best := candidate{
			src:            potentialSources[0],
			dst:            potentialTargets[0],
			productionDiff: minInt,
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
							src:            src,
							dst:            dst,
							productionDiff: scoreProduction(moveResult.TotalProduction()),
							fleet:          moveResult.FleetSize()[rules.MyPlayer.ToIndex()].Total(),
						}
						if cur.productionDiff > best.productionDiff ||
							cur.productionDiff == best.productionDiff && cur.fleet > best.fleet {
							best = cur
						}
					}
				}
			}
		}()
		if best.productionDiff > projectedProductionDiff ||
			best.productionDiff == projectedProductionDiff && best.fleet > projectedFleet {
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

func scoreProduction(ppp rules.ShipsPerPlayer) int {
	myScore := scoreProduction2(ppp[rules.MyPlayer.ToIndex()], ppp[rules.EnemyPlayer.ToIndex()])
	enemyScore := scoreProduction2(ppp[rules.EnemyPlayer.ToIndex()], ppp[rules.MyPlayer.ToIndex()])
	return myScore - enemyScore
}

func scoreProduction2(a, b rules.ShipCount) int { // assymetric
	const l = len(a)
	res := 0
	for i, p := range a {
		// n is better than n+1
		gt := b[(i+1)%l]
		lt := b[(i+l-1)%l]
		// this is just a stab in the dark. Having units is generally good, but not if the enemy has the counter, but especially if they counter the enemy.
		res += p + gt - lt
	}
	return res
}
