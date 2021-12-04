package p01

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

func A(in io.Reader) {
	depths := Ints(ReadLines(in))
	chg := 0
	for i, d := range depths[1:] {
		if depths[i] < d {
			chg++
		}
	}
	fmt.Println(chg)
}

func B(in io.Reader) {
	depths := Ints(ReadLines(in))
	prev := Sum(depths[:3]...)
	chg := 0
	for i := range depths[3:] {
		next := prev - depths[i] + depths[i+3]
		if next > prev {
			chg++
		}
		prev = next
	}
	fmt.Println(chg)
}
