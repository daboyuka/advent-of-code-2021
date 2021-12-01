package p01

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	. "aoc2021/helpers"
)

func A(in io.Reader) {
	depths := IntLines(ReadLines(in))
	chg := 0
	for i, d := range depths[1:] {
		if depths[i] < d {
			chg++
		}
	}
	fmt.Println(chg)
}

func B(in io.Reader) {
	depths := IntLines(ReadLines(in))
	prev := depths[0] + depths[1] + depths[2]
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

//

var _, _, _ = fmt.Sprintf, strings.Split, strconv.ParseFloat // use something from imports so we can always have them in our template
