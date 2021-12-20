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

func window(center Pos, g InfGrid, bg rune) int {
	v := 0
	for roff := -1; roff <= 1; roff++ {
		for coff := -1; coff <= 1; coff++ {
			off := Pos{roff, coff}

			c, ok := g[center.Add(off)]
			if !ok {
				c = bg
			}

			v <<= 1
			if c == '#' {
				v++
			}
		}
	}
	return v
}

func evolve(e Enhancer, g InfGrid, bg rune) (g2 InfGrid, nextBG rune) {
	g2 = make(InfGrid)
	min, max := g.Bounds()

	for r := min.Row - 1; r <= max.Row+1; r++ {
		for c := min.Col - 1; c <= max.Col+1; c++ {
			p := Pos{r, c}
			v := window(p, g, bg)
			g2[p] = e[v]
		}
	}

	if bg == '#' {
		return g2, e[0b111111111]
	} else {
		return g2, e[0b000000000]
	}
}

func A(in io.Reader) {
	e, g := parse(in)
	for i, bg := 0, '.'; i < 2; i++ {
		g, bg = evolve(e, g, bg)
	}
	fmt.Println(g.Count('#'))
}

func B(in io.Reader) {
	e, g := parse(in)
	for i, bg := 0, '.'; i < 50; i++ {
		g, bg = evolve(e, g, bg)
	}
	fmt.Println(g.Count('#'))
}
