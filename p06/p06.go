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

func A(in io.Reader) {
	fishAtAge := make([]int, 9)

	for _, age := range Ints(strings.Split(ReadLines(in)[0], ",")) {
		fishAtAge[age]++
	}

	for i := 0; i < 80; i++ {
		fishAtAge = evolve(fishAtAge)
	}

	fmt.Println(fishAtAge)
	fmt.Println(Sum(fishAtAge...))
}

func B(in io.Reader) {
	fishAtAge := make([]int, 9)

	for _, age := range Ints(strings.Split(ReadLines(in)[0], ",")) {
		fishAtAge[age]++
	}

	for i := 0; i < 256; i++ {
		fishAtAge = evolve(fishAtAge)
	}

	fmt.Println(fishAtAge)
	fmt.Println(Sum(fishAtAge...))
}
