package p10

import (
	"fmt"
	"io"
	"sort"

	. "aoc2021/helpers"
)

var braces = map[rune]rune{
	'[': ']',
	'(': ')',
	'{': '}',
	'<': '>',
}

func check(line string) (corruptSym rune, missingTail string) {
	var stack []rune
	for _, c := range line {
		if _, ok := braces[c]; ok { // opening symbol (push)
			stack = append(stack, c)
		} else if len(stack) == 0 { // closing symbol but empty stack (behavior not specified by problem)
			panic("wut")
		} else if top := stack[len(stack)-1]; braces[top] != c { // closing symbol but does not match opening symbol (corrupt)
			return c, ""
		} else { // closing symbol matches opening symbol (pop)
			stack = stack[:len(stack)-1]
		}
	}

	// If we get this far, no corruption; return missing tail (may be empty if well-formed line)
	l := len(stack)
	for i := range stack[:l/2] { // reverse stack
		j := l - i - 1
		stack[i], stack[j] = stack[j], stack[i]
	}
	for i, c := range stack { // map to closing braces
		stack[i] = braces[c]
	}
	return 0, string(stack)
}

var score = map[rune]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

func A(in io.Reader) {
	total := 0
	for i, l := range ReadLines(in) {
		if corruptSym, _ := check(l); corruptSym != 0 {
			fmt.Println(i, string(corruptSym))
			total += score[corruptSym]
		}
	}
	fmt.Println(total)
}

var score2 = map[rune]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

func B(in io.Reader) {
	var scores []int
	for _, l := range ReadLines(in) {
		if corruptSym, missingTail := check(l); corruptSym != 0 {
			continue // skip corrupt lines
		} else if missingTail != "" {
			s := 0
			for _, c := range missingTail {
				s = 5*s + score2[c]
			}
			scores = append(scores, s)
		}
	}
	sort.Ints(scores)
	fmt.Println(scores)
	fmt.Println(scores[len(scores)/2])
}
