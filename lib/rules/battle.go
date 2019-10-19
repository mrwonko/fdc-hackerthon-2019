package rules

func Battle(fleet1, fleet2 ShipCount) (ShipCount, ShipCount) {
	fFleet1, fFleet2 := newFShipCount(fleet1), newFShipCount(fleet2)
	for fFleet1.Total() > 0 && fFleet2.Total() > 0 {
		fFleet1, fFleet2 = battleRound(fFleet2, fFleet1), battleRound(fFleet1, fFleet2)
	}
	return fFleet1.ShipCount(), fFleet2.ShipCount()
}

// battleRound does one asymetric round. this needs to be called twice.
func battleRound(attacker, defender fShipCount) fShipCount {
	for curDef := range defender {
		for curAtt := range attacker {
			if attacker[curAtt] <= 0 {
				continue
			}
			multiplier, minimum := attackValues(curDef, curAtt)
			destroyed := attacker[curAtt] * multiplier
			if destroyed < minimum {
				destroyed = minimum
			}
			defender[curDef] -= destroyed
		}
		if defender[curDef] < 0 {
			defender[curDef] = 0
		}
	}
	return defender
}

func attackValues(defType, attType int) (multiplier float64, minimum float64) {
	strongAgainst := (attType + 1) % numShipTypes
	switch {
	case defType == attType:
		return 0.1, 1
	case defType == strongAgainst:
		return 0.25, 2
	default:
		return 0.01, 1
	}
}

type fShipCount [numShipTypes]float64

func newFShipCount(sc ShipCount) fShipCount {
	var res fShipCount
	for i, c := range sc {
		res[i] = float64(c)
	}
	return res
}

func (fsc fShipCount) ShipCount() ShipCount {
	var res ShipCount
	for i, c := range fsc {
		res[i] = int(c)
	}
	return res
}

func (fsc fShipCount) Total() float64 {
	var sum float64
	for _, c := range fsc {
		if c > 0 { // some weirdos send negative fleets
			sum += c
		}
	}
	return sum
}
