package p20

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

type Enhancer [512]rune

func parse(in io.Reader) (e Enhancer, g InfGrid) {
	lines := ReadLines(in)
	if len(lines[0]) != 512 {
		panic("oh noes")
	}
	for i, c := range lines[0] {
		e[i] = c
	}

	return e, ParseInfGrid(lines[2:])
}

func window(p Pos, g InfGrid, background rune) int {
	v := 0
	for roff := -1; roff <= 1; roff++ {
		for coff := -1; coff <= 1; coff++ {
			off := Pos{roff, coff}
			v <<= 1
			if c, ok := g[p.Add(off)]; (ok && c == '#') || (!ok && background == '#') {
				v++
			}
		}
	}
	return v
}

func evolve(e Enhancer, g InfGrid, background rune) InfGrid {
	g2 := make(InfGrid)
	min, max := g.Bounds()
	fmt.Println(min, max)

	for r := min.Row - 2; r <= max.Row+2; r++ {
		for c := min.Col - 2; c <= max.Col+2; c++ {
			p := Pos{r, c}
			v := window(p, g, background)
			g2[p] = e[v]
		}
	}

	return g2
}

func A(in io.Reader) {
	e, g := parse(in)

	g = evolve(e, g, '.')
	g = evolve(e, g, e[0])
	count := 0
	for _, c := range g {
		if c == '#' {
			count++
		}
	}
	fmt.Println(count)
}

func B(in io.Reader) {
	e, g := parse(in)

	bg := '.'
	for i := 0; i < 50; i++ {
		g = evolve(e, g, bg)

		if bg == '.' {
			bg = e[0]
		} else {
			bg = e[511]
		}
	}

	count := 0
	for _, c := range g {
		if c == '#' {
			count++
		}
	}
	fmt.Println(count)
}
