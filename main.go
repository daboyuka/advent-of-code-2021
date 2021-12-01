package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	. "aoc2021/helpers"
	"aoc2021/p01"
	"aoc2021/p02"
	"aoc2021/p03"
	"aoc2021/p04"
	"aoc2021/p05"
	"aoc2021/p06"
	"aoc2021/p07"
	"aoc2021/p08"
	"aoc2021/p09"
	"aoc2021/p10"
	"aoc2021/p11"
	"aoc2021/p12"
	"aoc2021/p13"
	"aoc2021/p14"
	"aoc2021/p15"
	"aoc2021/p16"
	"aoc2021/p17"
	"aoc2021/p18"
	"aoc2021/p19"
	"aoc2021/p20"
	"aoc2021/p21"
	"aoc2021/p22"
	"aoc2021/p23"
	"aoc2021/p24"
	"aoc2021/p25"
)

var problems = [...][2]func(io.Reader){
	{p01.A, p01.B},
	{p02.A, p02.B},
	{p03.A, p03.B},
	{p04.A, p04.B},
	{p05.A, p05.B},
	{p06.A, p06.B},
	{p07.A, p07.B},
	{p08.A, p08.B},
	{p09.A, p09.B},
	{p10.A, p10.B},
	{p11.A, p11.B},
	{p12.A, p12.B},
	{p13.A, p13.B},
	{p14.A, p14.B},
	{p15.A, p15.B},
	{p16.A, p16.B},
	{p17.A, p17.B},
	{p18.A, p18.B},
	{p19.A, p19.B},
	{p20.A, p20.B},
	{p21.A, p21.B},
	{p22.A, p22.B},
	{p23.A, p23.B},
	{p24.A, p24.B},
	{p25.A, p25.B},
}

func parseArg(s string) (problem int, part byte) {
	pre, suf := s[:len(s)-1], s[len(s)-1]
	return Atoi(pre), suf
}

func main() {
	prob, part := parseArg(strings.ToLower(os.Args[1]))
	if prob < 1 || prob > len(problems) {
		panic(fmt.Errorf("bad problem number %d", prob))
	} else if part != 'a' && part != 'b' {
		panic(fmt.Errorf("bad part %c", part))
	}
	problems[prob-1][part-'a'](os.Stdin)
}
