package main

import (
	"github.com/mrwonko/fdc-hackerthon-2019/lib/rules"
	"github.com/mrwonko/fdc-hackerthon-2019/lib/runner"
)

func main() {
	runner.Main(&runner.Player{
		Name:     "ws-nop",
		Password: "F6B5108D-2446-44D4-A740-3D3B87D47AF9",
		Play:     playNop,
	})
}

func playNop(ticks <-chan *runner.Tick) {
	for tick := range ticks {
		tick.Move <- rules.Nop
	}
}
