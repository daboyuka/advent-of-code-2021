package p07

import (
	"fmt"
	"io"
	"strings"

	. "aoc2021/helpers"
)

func A(in io.Reader) {
	poss := Ints(strings.Split(ReadLines(in)[0], ","))

	cnts := map[int]int{}
	for _, pos := range poss {
		cnts[pos]++
	}

	sum := Sum(poss...)
	total := len(poss)
	behind := 0
	minSum := sum
	minSumAt := 0
	for i := 1; i <= Max(poss...); i++ {
		behind += cnts[i-1]

		sum += behind - (total - behind)
		if sum < minSum {
			minSum = sum
			minSumAt = i
		}
	}
	fmt.Println(minSumAt, minSum)
}

func B(in io.Reader) {
	poss := Ints(strings.Split(ReadLines(in)[0], ","))

	posCounts := map[int]int{}

	for _, pos := range poss {
		posCounts[pos]++
	}

	minSum := 1<<31 - 1
	minSumAt := 0
	for i := 0; i <= Max(poss...); i++ {
		sum := 0
		for p, c := range posCounts {
			d := Abs(p - i)
			sum += d * (d + 1) / 2 * c
		}
		if sum < minSum {
			minSum = sum
			minSumAt = i
		}
	}
	fmt.Println(minSumAt, minSum)
}
