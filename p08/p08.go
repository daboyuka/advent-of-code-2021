package p08

import (
	"fmt"
	"io"
	"sort"
	"strings"

	. "aoc2021/helpers"
)

func parseProblem(line string) (samples, actual []string) {
	words := Words(line)
	barIdx := 0
	for i, w := range words {
		if w == "|" {
			barIdx = i
			break
		}
	}
	fmt.Println(words)
	return words[:barIdx], words[barIdx+1:]
}

func A(in io.Reader) {
	lines := ReadLines(in)
	uniq := 0

	for _, line := range lines {
		_, actual := parseProblem(line)

		// 1 -> 2seg, 4 -> 4seg, 7 = 3seg, 8 = 7seg
		for _, dig := range actual {
			fmt.Println(dig, len(dig))
			if l := len(dig); l == 2 || l == 3 || l == 4 || l == 7 {
				uniq++
			}
		}
	}

	fmt.Println(uniq)
}

func B(in io.Reader) {
	sum := 0
	for _, line := range ReadLines(in) {
		samples, actual := parseProblem(line)

		segMayBeWire := map[rune]map[rune]bool{} // segWirePoss[seg] = set of wires possible

		for w := 'a'; w <= 'g'; w++ {
			segMayBeWire[w] = make(map[rune]bool)
			for w2 := 'a'; w2 <= 'g'; w2++ {
				segMayBeWire[w][w2] = true
			}
		}

		record := func(wires, segsOn string, segsOff string) {
			for _, seg := range segsOn {
				for w := 'a'; w <= 'g'; w++ {
					if !strings.ContainsRune(wires, w) { // clear possibility for other inactive wires
						delete(segMayBeWire[seg], w)
					}
				}
			}
			for _, seg := range segsOff {
				for _, w := range wires {
					delete(segMayBeWire[seg], w) // clear possibility for active wires
				}
			}
		}

		for _, sample := range samples { // sample -> wires on
			switch len(sample) {
			case 2: // 1
				record(sample, "cf", "abdeg")
			case 3: // 7
				record(sample, "acf", "bdeg")
			case 4: // 4
				record(sample, "bcdf", "aeg")
			case 5: // 2, 3, 5
				record(sample, "adg", "")
			case 6: // 6, 9, 0
				record(sample, "abfg", "")
			case 7: // 8
				record(sample, "abcdefg", "")
			}
		}

		wireToSeg := map[rune]rune{}
		for len(wireToSeg) < 7 {
			for seg, mayBeWire := range segMayBeWire {
				if len(mayBeWire) == 1 {
					var wire rune
					for wire = range mayBeWire {
					}
					wireToSeg[wire] = seg

					for _, poss := range segMayBeWire {
						delete(poss, wire)
					}
					delete(segMayBeWire, seg)
				}
			}
		}

		digs := ""
		for _, wires := range actual {
			segmentsRunes := []rune(strings.Map(func(r rune) rune { return wireToSeg[r] }, wires))
			sort.Slice(segmentsRunes, func(i, j int) bool { return segmentsRunes[i] < segmentsRunes[j] })
			segments := string(segmentsRunes)

			switch segments {
			case "abcefg":
				digs += "0"
			case "cf":
				digs += "1"
			case "acdeg":
				digs += "2"
			case "acdfg":
				digs += "3"
			case "bcdf":
				digs += "4"
			case "abdfg":
				digs += "5"
			case "abdefg":
				digs += "6"
			case "acf":
				digs += "7"
			case "abcdefg":
				digs += "8"
			case "abcdfg":
				digs += "9"
			default:
				panic(fmt.Errorf("bad segments '%s'", segments))
			}
		}
		fmt.Println(digs)
		sum += Atoi(digs)
	}
	fmt.Println(sum)
}
