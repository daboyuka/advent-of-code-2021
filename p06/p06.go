package p06

import (
	"fmt"
	"io"
	"strings"

	. "aoc2021/helpers"
)

func evolve(fishAtAge []int) []int {
	newFish := fishAtAge[0]

	newFishAtAge := fishAtAge[1:] // all ages shift down
	newFishAtAge[6] += newFish
	newFishAtAge = append(newFishAtAge, newFish)
	return newFishAtAge
}

func run(input string, n int) []int {
	fishAtAge := make([]int, 9)

	for _, age := range Ints(strings.Split(input, ",")) {
		fishAtAge[age]++
	}

	for i := 0; i < n; i++ {
		fishAtAge = evolve(fishAtAge)
	}

	return fishAtAge
}

func A(in io.Reader) {
	fishAtAge := run(ReadLines(in)[0], 80)
	fmt.Println(fishAtAge)
	fmt.Println(Sum(fishAtAge...))
}

func B(in io.Reader) {
	fishAtAge := run(ReadLines(in)[0], 256)
	fmt.Println(fishAtAge)
	fmt.Println(Sum(fishAtAge...))
}
