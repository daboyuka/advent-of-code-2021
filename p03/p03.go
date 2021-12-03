package p03

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

// note: leastCommon = !mostCommon (ties invert to 0 as intended)
func mostCommon(lines []string, pos int) bool {
	count := 0
	for _, l := range lines {
		if l[pos] == '1' {
			count++
		}
	}
	return count >= len(lines)-count // tie goes to 1
}

func filterBitmasks(lines []string, pos int, useMostCommon bool) []string {
	if len(lines) < 2 {
		return lines
	}

	criterion := mostCommon(lines, pos) == useMostCommon // invert if !useMostCommon
	idx := 0
	for _, l := range lines {
		if (l[pos] == '1') == criterion {
			lines[idx] = l
			idx++
		}
	}
	return lines[:idx]
}

func A(in io.Reader) {
	lines := ReadLines(in)

	gamma, epsilon := 0, 0
	for pos := range lines[0] {
		gamma <<= 1
		epsilon <<= 1
		if mostCommon(lines, pos) {
			gamma |= 1
		} else {
			epsilon |= 1
		}
	}

	fmt.Println(gamma, epsilon)
	fmt.Println(gamma * epsilon)
}

func B(in io.Reader) {
	lines := ReadLines(in)

	oxyLines := append([]string(nil), lines...)
	co2Lines := append([]string(nil), lines...)

	for pos := range lines[0] { // bitpos
		oxyLines = filterBitmasks(oxyLines, pos, true)
		co2Lines = filterBitmasks(co2Lines, pos, false)
	}

	oxyBits, co2Bits := oxyLines[0], co2Lines[0]

	fmt.Println(oxyBits, co2Bits)

	oxyVal, co2Val := 0, 0
	for i := range oxyBits {
		oxyVal <<= 1
		co2Val <<= 1
		oxyVal |= int(oxyBits[i] - '0')
		co2Val |= int(co2Bits[i] - '0')
	}

	fmt.Println(oxyVal, co2Val)
	fmt.Println(oxyVal * co2Val)
}
