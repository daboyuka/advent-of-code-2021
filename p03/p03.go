package p03

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

func mostCommon(lines []string, pos int, tie bool) bool {
	count := 0
	for _, l := range lines {
		if l[pos] == '1' {
			count++
		}
	}
	if count == len(lines)-count {
		return tie
	}
	return count > len(lines)-count
}

func leastCommon(lines []string, pos int, tie bool) bool {
	count := 0
	for _, l := range lines {
		if l[pos] == '1' {
			count++
		}
	}
	if count == len(lines)-count {
		return tie
	}
	return count < len(lines)-count
}

func A(in io.Reader) {
	lines := ReadLines(in)
	count := make([]int, len(lines[0]))
	for _, l := range lines {
		for i, c := range l {
			if c == '1' {
				count[i]++
			}
		}
	}

	mask := 0
	gamma := 0
	for _, c := range count {
		mask <<= 1
		mask |= 1
		gamma <<= 1
		if c > len(lines)/2 { // more 1 bits
			gamma |= 1
		}
	}

	epsilon := gamma ^ mask

	fmt.Println(gamma, epsilon)
	fmt.Println(gamma * epsilon)
}

func B(in io.Reader) {
	lines := ReadLines(in)

	oxyLines := append([]string(nil), lines...)
	co2Lines := append([]string(nil), lines...)

	for pos := range lines[0] { // bitpos
		oxyBit := mostCommon(oxyLines, pos, true)
		co2Bit := leastCommon(co2Lines, pos, false)

		fmt.Println(oxyBit, co2Bit)

		if len(oxyLines) > 1 {
			oxyPos := 0
			for _, l := range oxyLines {
				if (l[pos] == '1') == oxyBit {
					fmt.Println("oxy keep", l)
					oxyLines[oxyPos] = l
					oxyPos++
				} else {
					fmt.Println("oxy drop", l, l[pos], oxyBit)
				}
			}
			oxyLines = oxyLines[:oxyPos]
		}

		if len(co2Lines) > 1 {
			co2Pos := 0
			for _, l := range co2Lines {
				if (l[pos] == '1') == co2Bit {
					fmt.Println("co2 keep", l)
					co2Lines[co2Pos] = l
					co2Pos++
				} else {
					fmt.Println("co2 drop", l)
				}
			}
			co2Lines = co2Lines[:co2Pos]
		}

		fmt.Println()
	}

	fmt.Println(oxyLines[0], co2Lines[0])

	oxyVal, co2Val := 0, 0
	for i := range oxyLines[0] {
		oxyVal <<= 1
		co2Val <<= 1
		oxyVal |= int(oxyLines[0][i] - '0')
		co2Val |= int(co2Lines[0][i] - '0')
	}

	fmt.Println(oxyVal, co2Val)
	fmt.Println(oxyVal * co2Val)
}
