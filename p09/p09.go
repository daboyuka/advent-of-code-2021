package p09

import (
	"fmt"
	"io"
	"sort"

	. "aoc2021/helpers"
)

type hm []string

func (h hm) At(p Pos) int {
	if p.Row < 0 || p.Row >= len(h) {
		return -1
	}
	r := h[p.Row]
	if p.Col < 0 || p.Col >= len(r) {
		return -1
	}
	return int(r[p.Col] - '0')
}

func A(in io.Reader) {
	h := hm(ReadLines(in))

	risk := 0
	for r := 0; r < len(h); r++ {
		rs := h[r]
		for c := range rs {
			p := Pos{r, c}

			center := h.At(p)
			lowerNeigh := 0
			for _, neigh := range [...]Pos{p.Add(Pos{0, 1}), p.Add(Pos{0, -1}), p.Add(Pos{1, 0}), p.Add(Pos{-1, 0})} {
				v := h.At(neigh)
				fmt.Println(p, neigh, v, center)
				if v != -1 && v <= center {
					lowerNeigh++
				}
			}

			if lowerNeigh == 0 {
				fmt.Println(center)
				risk += center + 1
			}
		}
	}

	fmt.Println(risk)
}

func basinSize(h hm, at Pos, visited map[Pos]bool) (size int) {
	if visited[at] || h.At(at) == -1 || h.At(at) == 9 {
		return 0
	}
	visited[at] = true
	size = 1
	for _, neigh := range [...]Pos{at.Add(Pos{0, 1}), at.Add(Pos{0, -1}), at.Add(Pos{1, 0}), at.Add(Pos{-1, 0})} {
		size += basinSize(h, neigh, visited)
	}
	return size
}

func B(in io.Reader) {
	h := hm(ReadLines(in))

	basins := []Pos{}
	for r := 0; r < len(h); r++ {
		rs := h[r]
		for c := range rs {
			p := Pos{r, c}

			center := h.At(p)
			lowerNeigh := 0
			for _, neigh := range [...]Pos{p.Add(Pos{0, 1}), p.Add(Pos{0, -1}), p.Add(Pos{1, 0}), p.Add(Pos{-1, 0})} {
				v := h.At(neigh)
				fmt.Println(p, neigh, v, center)
				if v != -1 && v <= center {
					lowerNeigh++
				}
			}

			if lowerNeigh == 0 {
				basins = append(basins, p)
			}
		}
	}

	sizes := []int{}
	v := make(map[Pos]bool)
	for _, basin := range basins {
		size := basinSize(h, basin, v)
		sizes = append(sizes, size)
	}
	sort.Ints(sizes)
	top3 := sizes[len(sizes)-3:]
	fmt.Println(top3)
	fmt.Println(top3[0] * top3[1] * top3[2])
}
