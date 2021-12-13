package p13

import (
	"fmt"
	"io"
	"strings"

	. "aoc2021/helpers"
)

type fold struct {
	Axis  rune
	Coord int
}

func parseFold(line string) fold {
	dimS, coordS := Split2(strings.TrimPrefix(line, "fold along "), "=")
	return fold{rune(dimS[0]), Atoi(coordS)}
}

func parsePos(line string) Pos {
	xS, yS := Split2(line, ",")
	x, y := Atoi(xS), Atoi(yS)
	return Pos{y, x} // r, c
}

func parseInput(in io.Reader) (g InfGrid, folds []fold) {
	g = InfGrid{}
	lg := ReadLinegroups(in)

	for _, line := range lg[0] {
		g[parsePos(line)] = '#'
	}
	for _, line := range lg[1] {
		folds = append(folds, parseFold(line))
	}
	return g, folds
}

func (f fold) Fold(g InfGrid) {
	for p, v := range g {
		switch f.Axis {
		case 'x':
			if p.Col > f.Coord {
				delete(g, p)
				g[Pos{p.Row, 2*f.Coord - p.Col}] = v
			}
		case 'y':
			if p.Row > f.Coord {
				delete(g, p)
				g[Pos{2*f.Coord - p.Row, p.Col}] = v
			}
		default:
			panic(fmt.Errorf("unknwon fold %c", f.Axis))
		}
	}
}

func A(in io.Reader) {
	g, folds := parseInput(in)

	folds[0].Fold(g)

	fmt.Println(len(g))
}

func B(in io.Reader) {
	g, folds := parseInput(in)

	for _, f := range folds {
		f.Fold(g)
	}

	fmt.Println(g.ToFixedGrid(' '))
}
