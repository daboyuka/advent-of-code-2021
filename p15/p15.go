package p15

import (
	"container/heap"
	"fmt"
	"io"

	. "aoc2021/helpers"
)

type PosRisk struct {
	Pos
	From Pos
	Risk int
}
type posHeap []PosRisk

func (p posHeap) Len() int {
	return len(p)
}

func (p posHeap) Less(i, j int) bool {
	return p[i].Risk < p[j].Risk
}

func (p posHeap) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p *posHeap) Push(x interface{}) {
	*p = append(*p, x.(PosRisk))
}

func (p *posHeap) Pop() interface{} {
	last := len(*p) - 1
	x := (*p)[last]
	*p = (*p)[:last]
	return x
}

func A(in io.Reader) {
	g := ParseFixedGrid(ReadLines(in))

	riskPath := map[Pos]int{}
	toVisit := &posHeap{PosRisk{Pos{0, 0}, Pos{0, 0}, 0}}
	for len(*toVisit) > 0 {
		next := heap.Pop(toVisit).(PosRisk)
		fmt.Println(next)

		if r2, ok := riskPath[next.Pos]; ok {
			if r2 > next.Risk {
				panic(fmt.Errorf("%+v %d %d %+v", next.Pos, next.Risk, r2, *toVisit))
			}
			fmt.Println("skip", next)
			continue
		}

		r := next.Risk
		riskPath[next.Pos] = r
		fmt.Println("lock in", next)

		for _, p := range [...]Pos{next.Add(Pos{0, 1}), next.Add(Pos{0, -1}), next.Add(Pos{-1, 0}), next.Add(Pos{1, 0})} {
			if g.InBounds(p) {
				fmt.Println("push", PosRisk{p, next.Pos, r + int(g.At(p)-'0')})
				heap.Push(toVisit, PosRisk{p, next.Pos, r + int(g.At(p)-'0')})
			}
		}
	}

	fmt.Println(riskPath[Pos{len(g) - 1, len(g[0]) - 1}])
}

func boost(g FixedGrid, amt int) {
	for _, r := range g {
		for j, v := range r {
			r[j] = rune((int(v-'1')+amt)%9 + '1')
		}
	}
}

func B(in io.Reader) {
	lines := ReadLines(in)
	g := ParseFixedGrid(lines)

	// Horiz reps (cols)
	for rep := 1; rep <= 4; rep++ {
		g2 := ParseFixedGrid(lines)
		boost(g2, rep)

		for i := range g {
			g[i] = append(g[i], g2[i]...)
		}
	}

	rowTmpl := make(FixedGrid, len(g))
	for i := range g {
		rowTmpl[i] = append([]rune(nil), g[i]...)
	}
	for rep := 1; rep <= 4; rep++ {
		g2 := make(FixedGrid, len(rowTmpl))
		for i := range rowTmpl {
			g2[i] = append([]rune(nil), rowTmpl[i]...)
		}
		boost(g2, rep)
		g = append(g, g2...)
	}

	fmt.Println(g.String())

	riskPath := map[Pos]int{}
	toVisit := &posHeap{PosRisk{Pos{0, 0}, Pos{0, 0}, 0}}
	for len(*toVisit) > 0 {
		next := heap.Pop(toVisit).(PosRisk)

		if r2, ok := riskPath[next.Pos]; ok {
			if r2 > next.Risk {
				panic(fmt.Errorf("%+v %d %d %+v", next.Pos, next.Risk, r2, *toVisit))
			}
			continue
		}

		r := next.Risk
		riskPath[next.Pos] = r

		for _, p := range [...]Pos{next.Add(Pos{0, 1}), next.Add(Pos{0, -1}), next.Add(Pos{-1, 0}), next.Add(Pos{1, 0})} {
			if g.InBounds(p) {
				heap.Push(toVisit, PosRisk{p, next.Pos, r + int(g.At(p)-'0')})
			}
		}
	}

	fmt.Println(riskPath[Pos{len(g) - 1, len(g[0]) - 1}])
}
