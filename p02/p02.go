package p02

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

func A(in io.Reader) {
	lines := ReadLines(in)
	x, d := 0, 0
	for _, l := range lines {
		dir, amt := Split2(l, " ")
		amtV := Atoi(amt)
		switch dir {
		case "forward":
			x += amtV
		case "down":
			d += amtV
		case "up":
			d -= amtV
		}
	}
	fmt.Println(x, d)
	fmt.Println(x * d)
}

func B(in io.Reader) {
	lines := ReadLines(in)
	x, d, a := 0, 0, 0
	for _, l := range lines {
		dir, amt := Split2(l, " ")
		amtV := Atoi(amt)
		switch dir {
		case "forward":
			x += amtV
			d += amtV * a
		case "down":
			a += amtV
		case "up":
			a -= amtV
		}
	}
	fmt.Println(x, d)
	fmt.Println(x * d)
}
