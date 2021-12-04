package p04

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	. "aoc2021/helpers"
)

type bingo [][]int

func (b bingo) Mark(x int) bool {
	for i, row := range b {
		for j, v := range row {
			if v == x {
				b[i][j] = -1
				return true
			}
		}
	}
	return false
}

func (b bingo) Sum() (s int) {
	for _, row := range b {
		for _, v := range row {
			if v != -1 {
				s += v
			}
		}
	}
	return s
}

func (b bingo) Won() bool {
rowcheck:
	for i, row := range b {
		for _, v := range row {
			if v != -1 {
				continue rowcheck
			}
		}
		fmt.Println("won on row ", i)
		return true
	}

colcheck:
	for j := range b[0] {
		for _, row := range b {
			if row[j] != -1 {
				continue colcheck
			}
		}
		fmt.Println("won on col ", j)
		return true
	}

	return false
}

var word = regexp.MustCompile(` +`)

func toBingo(lines []string) (b bingo) {
	for _, line := range lines {
		b = append(b, IntLines(word.Split(strings.TrimSpace(line), -1)))
	}
	return b
}

func A(in io.Reader) {
	lg := ReadLinegroups(in)
	numbers := IntLines(strings.Split(lg[0][0], ","))

	var bingos []bingo
	for _, bingoLines := range lg[1:] {
		bingos = append(bingos, toBingo(bingoLines))
	}

	for _, n := range numbers {
		for bi, b := range bingos {
			b.Mark(n)
			if b.Won() {
				fmt.Println(bi, "won on", n, "with sum", b.Sum())
				fmt.Println(b.Sum(), n)
				fmt.Println(b.Sum() * n)
				return
			}
		}
	}
}

func B(in io.Reader) {
	lg := ReadLinegroups(in)
	numbers := IntLines(strings.Split(lg[0][0], ","))

	var bingos []bingo
	for _, bingoLines := range lg[1:] {
		bingos = append(bingos, toBingo(bingoLines))
	}

	for _, n := range numbers {
		rem := 0
		for bi, b := range bingos {
			b.Mark(n)
			if b.Won() {
				fmt.Println(bi, "won on", n, "with sum", b.Sum())
				if len(bingos) == 1 {
					fmt.Println(b.Sum(), n)
					fmt.Println(b.Sum() * n)
					return
				}
			} else {
				bingos[rem] = b
				rem++
			}
		}
		bingos = bingos[:rem]
	}
}
