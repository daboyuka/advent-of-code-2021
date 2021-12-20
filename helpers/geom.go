package helpers

import (
	"fmt"
	"math"
	"strings"
)

type Pos struct{ Row, Col int }

func (p Pos) Add(p2 Pos) Pos { return Pos{p.Row + p2.Row, p.Col + p2.Col} }
func (p Pos) Sub(p2 Pos) Pos { return Pos{p.Row - p2.Row, p.Col - p2.Col} }

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

type FixedGrid [][]rune

func ParseFixedGrid(lines []string) (g FixedGrid) {
	for _, l := range lines {
		g = append(g, []rune(l))
	}
	return g
}

func MakeFixedGrid(rows, cols int, fill rune) FixedGrid {
	g := make(FixedGrid, rows)
	for i := range g {
		g[i] = make([]rune, cols)
		for j := range g[i] {
			g[i][j] = fill
		}
	}
	return g
}

func (g FixedGrid) Copy() FixedGrid {
	g2 := make(FixedGrid, len(g))
	for i, r := range g {
		g2[i] = append([]rune(nil), r...)
	}
	return g2
}

func (g FixedGrid) InBounds(p Pos) bool {
	return p.Row >= 0 && p.Row < len(g) && p.Col >= 0 && p.Col < len(g[0])
}

func (g FixedGrid) At(p Pos) rune {
	if !g.InBounds(p) {
		return 0
	}
	return g[p.Row][p.Col]
}

func (g FixedGrid) Set(p Pos, v rune) {
	if !g.InBounds(p) {
		panic("out of bounds")
	}
	g[p.Row][p.Col] = v
}

func (g FixedGrid) Count(x rune) (count int) {
	for _, r := range g {
		for _, c := range r {
			if c == x {
				count++
			}
		}
	}
	return count
}

func (g FixedGrid) String() string {
	sb := strings.Builder{}
	for _, r := range g {
		sb.WriteString(string(r))
		sb.WriteByte('\n')
	}
	return sb.String()
}

type InfGrid map[Pos]rune

func ParseInfGrid(lines []string) InfGrid {
	g := make(InfGrid, len(lines)*len(lines[0]))
	for r, l := range lines {
		for c, v := range l {
			g[Pos{r, c}] = v
		}
	}
	return g
}

func (g InfGrid) Copy() InfGrid {
	g2 := make(InfGrid, len(g))
	for p, v := range g {
		g2[p] = v
	}
	return g2
}

func (g InfGrid) Bounds() (min, max Pos) {
	minR, maxR := math.MaxInt, math.MinInt
	minC, maxC := math.MaxInt, math.MinInt
	for p := range g {
		minR, maxR = Min(p.Row, minR), Max(p.Row, maxR)
		minC, maxC = Min(p.Col, minC), Max(p.Col, maxC)
	}

	return Pos{minR, minC}, Pos{maxR, maxC}
}

func (g InfGrid) ToFixedGrid(fill rune) FixedGrid {
	base, max := g.Bounds()
	g2 := MakeFixedGrid(max.Row-base.Row+1, max.Col-base.Col+1, fill)
	for p, v := range g {
		g2.Set(p.Sub(base), v)
	}

	return g2
}

func (g InfGrid) Count(x rune) (count int) {
	for _, c := range g {
		if c == x {
			count++
		}
	}
	return count
}

func (g InfGrid) String() string { return g.ToFixedGrid(' ').String() }
