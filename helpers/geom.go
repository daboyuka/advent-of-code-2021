package helpers

import "fmt"

type Pos struct{ Row, Col int }

func (p Pos) Add(p2 Pos) Pos { return Pos{p.Row + p2.Row, p.Col + p2.Col} }

type Dir int

const (
	South = Dir(iota)
	East
	North
	West
)

func (d Dir) Left() Dir {
	switch d {
	case South:
		return East
	case East:
		return North
	case North:
	case West:
	}
	panic(fmt.Errorf("bad dir %d", d))
}
func (p Pos) Move(d Dir, amt int) Pos {
	switch d {
	case South:
		return Pos{p.Row + amt, p.Col}
	case East:
		return Pos{p.Row, p.Col + amt}
	case North:
		return Pos{p.Row - amt, p.Col}
	case West:
		return Pos{p.Row, p.Col - amt}
	}
	panic(fmt.Errorf("bad dir %d", d))
}
