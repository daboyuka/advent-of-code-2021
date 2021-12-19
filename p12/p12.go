package p12

import (
	"fmt"
	"io"
	"unicode"

	. "aoc2021/helpers"
)

type Caves map[string][]string // [from] -> tos

func parse(in io.Reader) Caves {
	edges := make(Caves)
	for _, l := range ReadLines(in) {
		a, b := Split2(l, "-")
		edges[a] = append(edges[a], b)
		edges[b] = append(edges[b], a)
	}
	return edges
}

func traverse(at string, edges Caves, visited map[string]bool, allowDoubleVisit bool) (nPaths int) {
	for _, next := range edges[at] {
		if next == "start" {
			continue
		} else if next == "end" {
			nPaths++
		} else if unicode.IsLower(rune(next[0])) {
			if alreadyVisited := visited[next]; alreadyVisited {
				if !allowDoubleVisit {
					continue
				}
				nPaths += traverse(next, edges, visited, false)
			} else {
				visited[next] = true
				nPaths += traverse(next, edges, visited, allowDoubleVisit)
				delete(visited, next)
			}
		} else {
			nPaths += traverse(next, edges, visited, allowDoubleVisit)
		}
	}
	return nPaths
}

func A(in io.Reader) {
	caves := parse(in)
	nToEnd := traverse("start", caves, make(map[string]bool, len(caves)), false)
	fmt.Println(nToEnd)
}

func B(in io.Reader) {
	caves := parse(in)
	nToEnd := traverse("start", caves, make(map[string]bool, len(caves)), true)
	fmt.Println(nToEnd)
}
