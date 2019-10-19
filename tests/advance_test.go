package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"testing"

	"github.com/mrwonko/fdc-hackerthon-2019/lib/rules"
	"github.com/mrwonko/fdc-hackerthon-2019/lib/runner"
)

func TestAdvance(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/game.ndjson")
	if err != nil {
		t.Fatal(err)
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	var cur, prev rules.Gamestate
	var curRaw, prevRaw runner.Gamestate
	first := true
	for {
		switch err = decoder.Decode(&curRaw); err {
		case io.EOF:
			return
		case nil:
		default:
			t.Fatal(err)
		}
		cur = *curRaw.Preprocess()
		if first {
			first = false
			prev, prevRaw = cur, curRaw
			continue
		}
		sends, err := detectSends(prevRaw, curRaw)
		if err != nil {
			t.Fatal(err)
		}
		withSend := prev
		for _, s := range sends {
			t.Logf("round %d: send: %v", prev.Round, s)
			withSend = *withSend.Send(s.src, s.dst, s.ships)
		}
		estimated, err := rules.Advance(context.Background(), &withSend, 1)
		if err != nil {
			t.Fatal(err)
		}
		checkGamestate(t, estimated, &cur)
		prev, prevRaw = cur, curRaw
	}
}

type Send struct {
	src, dst int
	ships    rules.ShipCount
}

func detectSends(prev, cur runner.Gamestate) ([]Send, error) {
	newFleets := detectNewFleets(prev, cur)
	if len(newFleets) > 2 {
		return nil, fmt.Errorf("more than 2 new fleets: %v", newFleets)
	}
	var res []Send
	for _, f := range newFleets {
		s := Send{
			ships: f.Ships,
		}
		for i, p := range cur.Planets {
			if p.ID == f.Target {
				s.dst = i
			}
			if p.ID == f.Origin {
				s.src = i
			}
		}
		res = append(res, s)
	}
	return res, nil
}

func detectNewFleets(prev, cur runner.Gamestate) []runner.Fleet {
	var ifNothing []runner.Fleet
	for _, f := range prev.Fleets {
		if f.ETA > 1 {
			f.ETA--
			ifNothing = append(ifNothing, f)
		}
	}
	var new []runner.Fleet
outer:
	for _, cf := range cur.Fleets {
		for _, pf := range ifNothing {
			if pf.ID == cf.ID {
				continue outer
			}
		}
		new = append(new, cf)
	}
	return new
}

func checkGamestate(t *testing.T, want, got *rules.Gamestate) {
	if want.Round != got.Round {
		t.Errorf("round=%d: want %d", got.Round, want.Round)
	}
	for i, wp := range want.Planets {
		gp := got.Planets[i]
		if wp != gp {
			t.Errorf("round=%d: Planets[%d]=%#v, want %#v", got.Round, i, gp, wp)
		}
	}
	sort.Sort(SortableFleets(want.Fleets))
	sort.Sort(SortableFleets(got.Fleets))
	for i, wf := range want.Fleets {
		gf := got.Fleets[i]
		wf.ID = 0
		gf.ID = 0
		if wf != gf {
			t.Errorf("round=%d: Fleets[%d]=%#v, want %#v", got.Round, i, gf, wf)
		}
	}
}

type SortableFleets []rules.Fleet

func (sf SortableFleets) Len() int {
	return len(sf)
}

func (sf SortableFleets) Less(i, j int) bool {
	l, r := sf[i], sf[j]
	if l.Owner < r.Owner {
		return true
	} else if l.Owner > r.Owner {
		return false
	}
	if l.OriginIndex < r.OriginIndex {
		return true
	} else if l.OriginIndex > r.OriginIndex {
		return false
	}
	if l.TargetIndex < r.TargetIndex {
		return true
	} else if l.TargetIndex > r.TargetIndex {
		return false
	}
	if l.ETA < r.ETA {
		return true
	} else if l.ETA > r.ETA {
		return false
	}
	if l.Ships[0] < r.Ships[0] {
		return true
	} else if l.Ships[0] > r.Ships[0] {
		return false
	}
	if l.Ships[1] < r.Ships[1] {
		return true
	} else if l.Ships[1] > r.Ships[1] {
		return false
	}
	if l.Ships[2] < r.Ships[2] {
		return true
	} else if l.Ships[2] > r.Ships[2] {
		return false
	}
	// this should be sufficient, as each player can only send one player each turn
	return false
}

func (sf SortableFleets) Swap(i, j int) {
	sf[j], sf[i] = sf[i], sf[j]
}
