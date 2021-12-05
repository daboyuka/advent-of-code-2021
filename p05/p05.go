package p05

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

type pos struct{ X, Y int }

type vent struct{ A, B pos }

func parseLine(l string) vent {
	as, _, bs := Split3(l, " ", " ")
	ax, ay := Split2(as, ",")
	bx, by := Split2(bs, ",")
	a := pos{Atoi(ax), Atoi(ay)}
	b := pos{Atoi(bx), Atoi(by)}
	return vent{a, b}
}

func (v vent) draw(out map[pos]int) {
	dx := v.A.X - v.B.X
	dy := v.A.Y - v.B.Y
	dmax := Max(Abs(dx), Abs(dy))
	for i := 0; i <= dmax; i++ {
		x := v.B.X + i*dx/dmax
		y := v.B.Y + i*dy/dmax
		fmt.Println(v, ":", i, x, y)
		out[pos{x, y}]++
	}
}

func (v vent) isDiag() bool {
	return v.A.X != v.B.X && v.A.Y != v.B.Y
}

func A(in io.Reader) {
	var vents []vent
	for _, l := range ReadLines(in) {
		vents = append(vents, parseLine(l))
	}

	m := make(map[pos]int)
	for _, v := range vents {
		if !v.isDiag() {
			v.draw(m)
		}
	}

	bad := 0
	for _, x := range m {
		if x >= 2 {
			bad++
		}
	}

	fmt.Println(bad)
}

func B(in io.Reader) {
	var vents []vent
	for _, l := range ReadLines(in) {
		vents = append(vents, parseLine(l))
	}

	m := make(map[pos]int)
	for _, v := range vents {
		v.draw(m)
	}

	bad := 0
	for _, x := range m {
		if x >= 2 {
			bad++
		}
	}

	fmt.Println(bad)
}
