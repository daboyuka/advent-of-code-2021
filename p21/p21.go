package p21

import (
	"fmt"
	"io"
	"strings"

	. "aoc2021/helpers"
)

func parse(in io.Reader) (p1, p2 int) {
	lines := ReadLines(in)
	p1 = Atoi(strings.TrimPrefix(lines[0], "Player 1 starting position: "))
	p2 = Atoi(strings.TrimPrefix(lines[1], "Player 2 starting position: "))
	return p1, p2
}

func play(p1, p2 int) (p1score, p2score, rolls int) {
	nextDie := 1
	p1Turn := true
	for p1score < 1000 && p2score < 1000 {
		d1 := nextDie
		nextDie = nextDie%100 + 1
		d2 := nextDie
		nextDie = nextDie%100 + 1
		d3 := nextDie
		nextDie = nextDie%100 + 1
		rolls += 3
		fmt.Println(d1, d2, d3, p1, p2)
		if p1Turn {
			p1 = (p1+d1+d2+d3-1)%10 + 1
			p1score += p1
			fmt.Println(d1, d2, d3, p1, p2)
		} else {
			p2 = (p2+d1+d2+d3-1)%10 + 1
			p2score += p2
			fmt.Println(d1, d2, d3, p1, p2)
		}
		p1Turn = !p1Turn
	}
	return
}

func A(in io.Reader) {
	p1, p2 := parse(in)
	p1score, p2score, rolls := play(p1, p2)

	lowerScore := p1score
	if p1score >= 1000 {
		lowerScore = p2score
	}
	fmt.Println(p1score, p2score, rolls)
	fmt.Println(lowerScore * rolls)
}

type GameState struct {
	P1, P2           int
	P1Score, P2Score int
}

var diracRolls = [...]int{ // [roll] = instances
	3: 1,
	4: 3,
	5: 6,
	6: 7,
	7: 6,
	8: 3,
	9: 1,
}

// for 0-based
func addpos(p int, delta int) int {
	return (p + delta + 10) % 10
}

func playDirac(firstP1, firstP2 int) (p1wins, p2wins int) {
	const maxScore = 30
	states := [maxScore + 1][maxScore + 1][2][10][10]int{} // [p1s][p2s][0 == next turn p1][p1][p2]
	states[0][0][0][firstP1-1][firstP2-1] = 1

	// p1p and p2p are 0-based
	inner := func(p1s, p2s, turn int) {
		for p1p := range states[0][0][0] {
			for p2p := range states[0][0][0][0] {
				for roll, combos := range diracRolls {
					if combos == 0 {
						continue
					}
					p1sPrev, p2sPrev := p1s, p2s
					p1pPrev, p2pPrev := p1p, p2p
					if turn == 0 {
						p2sPrev -= p2p + 1 // lose score based on current position
						p2pPrev = addpos(p2p, -roll)
					} else {
						p1sPrev -= p1p + 1 // lose score based on current position
						p1pPrev = addpos(p1p, -roll)
					}

					if p1sPrev < 0 || p2sPrev < 0 || p1sPrev >= 21 || p2sPrev >= 21 {
						continue
					}

					states[p1s][p2s][turn][p1p][p2p] += combos *
						states[p1sPrev][p2sPrev][1-turn][p1pPrev][p2pPrev]

					if false && states[p1sPrev][p2sPrev][1-turn][p1pPrev][p2pPrev] != 0 {
						fmt.Println(
							"from",
							p1sPrev, p2sPrev, 1-turn, p1pPrev, p2pPrev, "=",
							states[p1sPrev][p2sPrev][1-turn][p1pPrev][p2pPrev],
							"roll", roll, "=", combos,
							"to",
							p1s, p2s, turn, p1p, p2p, "=",
							states[p1s][p2s][turn][p1p][p2p],
						)
					}
				}
			}
		}
	}

	for diag := 1; diag < len(states)*2-1; diag++ {
		for idx := 0; idx < len(states); idx++ {
			p1s, p2s := diag-idx, idx
			if p1s < 0 || p2s < 0 || p1s >= len(states) || p2s >= len(states) {
				continue
			}
			fmt.Println(p1s, p2s)

			inner(p1s, p2s, 0)
			inner(p1s, p2s, 1)
		}
	}

	for p1s, x := range states {
		for p2s, x := range x {
			for _, x := range x {
				for _, x := range x {
					for _, combos := range x {
						if p1s >= 21 && p2s < 21 {
							p1wins += combos
						} else if p1s < 21 && p2s >= 21 {
							p2wins += combos
						}
					}
				}
			}
		}
	}
	return
}

func B(in io.Reader) {
	p1, p2 := parse(in)
	p1wins, p2wins := playDirac(p1, p2)
	fmt.Println(p1wins, p2wins)
	fmt.Println(Max(p1wins, p2wins))
}
