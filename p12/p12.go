package p12

import (
	"fmt"
	"io"
	"sort"
	"strings"

	. "aoc2021/helpers"
)

type PathSeg struct{ A, B string }

func traverse(at string, nPaths int, path []string, edges map[string][]string, visited map[string]bool) (nToEnd int) {
	for _, next := range edges[at] {
		if next == "start" {
			continue
		} else if next == "end" {
			nToEnd += nPaths
			//fmt.Println(nPaths, "x", path)
		} else if strings.ToLower(next) == next {
			if visited[next] {
				continue
			}

			visited[next] = true
			nToEnd += traverse(next, nPaths, append(path, next), edges, visited)
			delete(visited, next)
		} else {
			nToEnd += traverse(next, nPaths, append(path, next), edges, visited)
		}
	}
	return nToEnd
}

func traverse2(at string, nPaths int, path []string, edges map[string][]string, visited map[string]bool, doubleVisited bool) (nToEnd int) {
	for _, next := range edges[at] {
		if next == "start" {
			continue
		} else if next == "end" {
			nToEnd += nPaths
			//fmt.Println(nPaths, "x", append(path, next))
		} else if strings.ToLower(next) == next {
			isV := visited[next]
			if isV {
				if doubleVisited {
					continue
				}
				nToEnd += traverse2(next, nPaths, append(path, next), edges, visited, true)
			} else {
				visited[next] = true
				nToEnd += traverse2(next, nPaths, append(path, next), edges, visited, doubleVisited)
				delete(visited, next)
			}
		} else {
			nToEnd += traverse2(next, nPaths, append(path, next), edges, visited, doubleVisited)
		}
	}
	return nToEnd
}

func A(in io.Reader) {
	lines := ReadLines(in)
	segs := make([]PathSeg, len(lines))
	for i, l := range lines {
		segs[i].A, segs[i].B = Split2(l, "-")
	}

	edges := make(map[string][]string)
	for _, seg := range segs {
		edges[seg.A] = append(edges[seg.A], seg.B)
		edges[seg.B] = append(edges[seg.B], seg.A)
	}

	nToEnd := traverse("start", 1, []string{"start"}, edges, make(map[string]bool, len(edges)))
	fmt.Println(nToEnd)
}

func B(in io.Reader) {
	lines := ReadLines(in)
	segs := make([]PathSeg, len(lines))
	for i, l := range lines {
		segs[i].A, segs[i].B = Split2(l, "-")
	}

	edges := make(map[string][]string)
	for _, seg := range segs {
		edges[seg.A] = append(edges[seg.A], seg.B)
		edges[seg.B] = append(edges[seg.B], seg.A)
	}

	for _, out := range edges {
		sort.Strings(out)
	}

	nToEnd := traverse2("start", 1, []string{"start"}, edges, make(map[string]bool, len(edges)), false)
	fmt.Println(nToEnd)
}
