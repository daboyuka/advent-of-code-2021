package p14

import (
	"fmt"
	"io"
	"math"
	"strings"

	. "aoc2021/helpers"
)

type rule struct {
	A, B rune
	To   rune
}

func Apply(in string, rules []rule) string {
	sb := strings.Builder{}

	prev := rune(in[0])
	for _, c := range in[1:] {
		insert := rune(0)
		for _, r := range rules {
			if prev == r.A && c == r.B {
				insert = r.To
				break
			}
		}

		sb.WriteRune(prev)
		if insert != 0 {
			sb.WriteRune(insert)
		}

		prev = c
	}

	sb.WriteByte(in[len(in)-1])
	return sb.String()
}

func A(in io.Reader) {
	lines := ReadLines(in)

	start := lines[0]
	var rules []rule
	for _, r := range lines[2:] {
		l, r := Split2(r, " -> ")
		rules = append(rules, rule{rune(l[0]), rune(l[1]), rune(r[0])})
	}

	out := start
	for i := 0; i < 10; i++ {
		out = Apply(out, rules)
		fmt.Println(i, len(out))
	}

	counts := [256]int{}
	for _, c := range out {
		counts[c]++
	}

	min, max, minC, maxC := len(out), 0, byte(0), byte(0)
	for c, count := range counts {
		if count > 0 && min > count {
			min = count
			minC = byte(c)
		}
		if max < count {
			max = count
			maxC = byte(c)
		}
	}

	fmt.Println(min, max, minC, maxC)
	fmt.Println(max - min)
}

func B(in io.Reader) {
	lines := ReadLines(in)

	start := lines[0]
	var rules []rule
	for _, r := range lines[2:] {
		l, r := Split2(r, " -> ")
		rules = append(rules, rule{rune(l[0]), rune(l[1]), rune(r[0])})
	}

	pairCounts := map[string]int{}
	for i, c := range start[1:] {
		pairCounts[string(rune(start[i]))+string(c)]++
	}

	for i := 0; i < 40; i++ {
		newPairCounts := map[string]int{}
		for pair, count := range pairCounts {
			a, b := rune(pair[0]), rune(pair[1])
			var theRule rule
			for _, r := range rules {
				if r.A == a && r.B == b {
					theRule = r
					break
				}
			}

			if theRule.A == 0 {
				continue
			}

			newPairCounts[string(a)+string(theRule.To)] += count
			newPairCounts[string(theRule.To)+string(b)] += count
		}

		pairCounts = newPairCounts
		fmt.Println(pairCounts)
	}

	counts := [256]int{}
	for pair, count := range pairCounts {
		counts[pair[0]] += count
		counts[pair[1]] += count
	}

	// This double-counted everything, EXCEPT the endpoints
	counts[start[0]]++
	counts[start[len(start)-1]]++
	for i := range counts {
		counts[i] /= 2
	}

	min, max, minC, maxC := math.MaxInt, 0, byte(0), byte(0)
	for c, count := range counts {
		if count > 0 && min > count {
			min = count
			minC = byte(c)
		}
		if max < count {
			max = count
			maxC = byte(c)
		}
	}

	fmt.Println(min, max, minC, maxC)
	fmt.Println(max - min)
}
