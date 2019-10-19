package rules

import (
	"fmt"
	"io"
)

type Move interface {
	isMove() // only this package may implement this sum type
	io.WriterTo
}

const Nop NopMove = ""

type NopMove string // just any type we can store in a constant

var _ Move = NopMove("")

func (NopMove) isMove() {}

var nopMessage = []byte("nop\n")

func (NopMove) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(nopMessage)
	return int64(n), err
}

type Send struct {
	Src   PlanetID
	Dst   PlanetID
	Ships ShipCount
}

var _ Move = (*Send)(nil)

func (s *Send) isMove() {}

func (s *Send) WriteTo(w io.Writer) (int64, error) {
	n, err := fmt.Fprintf(w, "send %d %d %d %d %d\n", s.Src, s.Dst, s.Ships[0], s.Ships[1], s.Ships[2])
	return int64(n), err
}
