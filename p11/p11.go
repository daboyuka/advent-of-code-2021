package p11

import (
	"fmt"
	"io"
	"strings"

	. "aoc2021/helpers"
)

type grid [][]int

func (g grid) InBounds(i, j int) bool {
	return i >= 0 && i < len(g) && j >= 0 && j < len(g[0])
}

func (g grid) String() string {
	s := ""
	for _, r := range g {
		for _, v := range r {
			s += string(rune(v) + '0')
		}
		s += "\n"
	}
	return s
}

func simulate(g grid) (flashes int) {
	for _, r := range g {
		for j := range r {
			r[j]++
		}
	}

	for {
		flashed := 0
		for i, r := range g {
			for j, v := range r {
				if v > 9 {
					g[i][j] = 0
					flashed++
					for dx := -1; dx <= 1; dx++ {
						for dy := -1; dy <= 1; dy++ {
							if dx == 0 && dy == 0 {
								continue
							} else if i2, j2 := i+dx, j+dy; !g.InBounds(i2, j2) {
								continue
							} else if g[i2][j2] != 0 {
								g[i2][j2]++
							}
						}
					}
				}
			}
		}
		if flashed == 0 {
			break
		} else {
			flashes += flashed
		}
	}
	return flashes
}

func A(in io.Reader) {
	lines := ReadLines(in)

	g := grid{}
	for _, l := range lines {
		r := Ints(strings.Split(l, ""))
		g = append(g, r)
	}

	fmt.Println(g.String())

	flashes := 0
	for i := 0; i < 100; i++ {
		flashes += simulate(g)
		fmt.Println(g.String())
	}
	fmt.Println(flashes)
}

func B(in io.Reader) {
	lines := ReadLines(in)

	g := grid{}
	for _, l := range lines {
		r := Ints(strings.Split(l, ""))
		g = append(g, r)
	}

	for i := 1; ; i++ {
		flashes := simulate(g)
		if flashes == len(g)*len(g[0]) {
			fmt.Println(i)
			return
		}
	}
}
