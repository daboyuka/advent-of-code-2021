package p07

import (
	"fmt"
	"io"
	"strings"

	. "aoc2021/helpers"
)

func run(in io.Reader, distMetric func(int) int) (minSumAt int, minSum int) {
	poss := Ints(strings.Split(ReadLines(in)[0], ","))

	posCounts := map[int]int{}
	for _, pos := range poss {
		posCounts[pos]++
	}

	minSum = -1
	for i := 0; i <= Max(poss...); i++ {
		sum := 0
		for p, c := range posCounts {
			d := Abs(p - i)
			sum += distMetric(d) * c
		}
		if minSum == -1 || sum < minSum {
			minSumAt, minSum = i, sum
		}
	}

	return minSumAt, minSum
}

func A(in io.Reader) {
	at, sum := run(in, func(d int) int { return d })
	fmt.Println(at)
	fmt.Println(sum)
}

func B(in io.Reader) {
	at, sum := run(in, func(d int) int { return d * (d + 1) / 2 })
	fmt.Println(at)
	fmt.Println(sum)
}
