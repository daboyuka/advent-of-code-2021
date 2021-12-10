package p10

import (
	"fmt"
	"io"
	"sort"

	. "aoc2021/helpers"
)

var braces = map[byte]byte{
	'[': ']',
	'(': ')',
	'{': '}',
	'<': '>',
}

func checkCorrupted(line string) byte {
	stack := []byte{}
	for _, c := range line {
		if _, ok := braces[byte(c)]; ok {
			stack = append(stack, byte(c))
		} else if len(stack) == 0 {
			panic("wut")
		} else if top := stack[len(stack)-1]; braces[top] != byte(c) {
			return byte(c)
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	return 0
}

func findRemaining(line string) (end string) {
	stack := []byte{}
	for _, c := range line {
		if _, ok := braces[byte(c)]; ok {
			stack = append(stack, byte(c))
		} else if len(stack) == 0 {
			panic("wut")
		} else if top := stack[len(stack)-1]; braces[top] != byte(c) {
			panic("corrupt")
		} else {
			stack = stack[:len(stack)-1]
		}
	}

	for _, c := range stack {
		end = string(rune(braces[c])) + end
	}
	return end
}

var score = map[byte]int{
	')': 3,
	']': 57,
	'}': 1197,
	'>': 25137,
}

func A(in io.Reader) {
	lines := ReadLines(in)
	total := 0
	for i, l := range lines {
		if badBrace := checkCorrupted(l); badBrace != 0 {
			fmt.Println(i, string(rune(badBrace)))
			total += score[badBrace]
		}
	}
	fmt.Println(total)
}

var score2 = map[byte]int{
	')': 1,
	']': 2,
	'}': 3,
	'>': 4,
}

func B(in io.Reader) {
	lines := ReadLines(in)
	keep := 0
	for _, l := range lines {
		if checkCorrupted(l) == 0 {
			lines[keep] = l
			keep++
		}
	}
	lines = lines[:keep]

	scores := []int{}
	for _, l := range lines {
		if rem := findRemaining(l); rem != "" {
			s := 0
			for _, c := range rem {
				s = 5*s + score2[byte(c)]
			}
			scores = append(scores, s)
		}
	}
	sort.Ints(scores)
	fmt.Println(scores)
	fmt.Println(scores[len(scores)/2])
}
