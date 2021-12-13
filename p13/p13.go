package p13

import (
	"fmt"
	"io"
	"strings"

	. "aoc2021/helpers"
)

type Pt struct{ X, Y int }

func A(in io.Reader) {
	g := map[Pt]bool{}
	lg := ReadLinegroups(in)

	for _, line := range lg[0] {
		xS, yS := Split2(line, ",")
		x, y := Atoi(xS), Atoi(yS)
		g[Pt{x, y}] = true
	}

	folds := lg[1]
	for _, f := range folds[:1] {
		f := strings.TrimPrefix(f, "fold along ")
		dim, coordS := Split2(f, "=")
		coord := Atoi(coordS)

		for p := range g {
			switch dim {
			case "x":
				if p.X > coord {
					delete(g, p)
					g[Pt{2*coord - p.X, p.Y}] = true
				}
			case "y":
				if p.Y > coord {
					delete(g, p)
					g[Pt{p.X, 2*coord - p.Y}] = true
				}
			}
		}
	}

	fmt.Println(len(g))
}

func B(in io.Reader) {
	g := map[Pt]bool{}
	lg := ReadLinegroups(in)

	for _, line := range lg[0] {
		xS, yS := Split2(line, ",")
		x, y := Atoi(xS), Atoi(yS)
		g[Pt{x, y}] = true
	}

	folds := lg[1]
	for _, f := range folds {
		f := strings.TrimPrefix(f, "fold along ")
		dim, coordS := Split2(f, "=")
		coord := Atoi(coordS)

		for p := range g {
			switch dim {
			case "x":
				if p.X > coord {
					delete(g, p)
					g[Pt{2*coord - p.X, p.Y}] = true
				}
			case "y":
				if p.Y > coord {
					delete(g, p)
					g[Pt{p.X, 2*coord - p.Y}] = true
				}
			}
		}
	}

	maxX, maxY := 0, 0
	for p := range g {
		if maxX < p.X {
			maxX = p.X
		}
		if maxY < p.Y {
			maxY = p.Y
		}
	}

	g2 := make(FixedGrid, maxY+1)
	for i := range g2 {
		g2[i] = make([]rune, maxX+1)
		for j := range g2[i] {
			g2[i][j] = ' '
		}
	}
	for p := range g {
		g2[p.Y][p.X] = '#'
	}

	for _, r := range g2 {
		fmt.Println(string(r))
	}
}
