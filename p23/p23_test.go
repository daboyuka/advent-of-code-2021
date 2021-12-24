package p23

import "testing"

func TestCheckBlocking(t *testing.T) {
	ck := func(moves []Move, prob Problem, poses []int, expects []bool) {
		for i, pos := range poses {
			expect := expects[i]
			obtained := checkBlocking(pos, moves, prob)
			if obtained != expect {
				t.Fatalf("failed %d %+v %+v: obtained %t, expected %t", pos, moves, prob.Targets, obtained, expect)
			}
		}
	}

	prob := Problem{Targets: [4]int{0, 1, 2, 3}}

	ck([]Move{
		{Amphipod: Amphipod{Kind: 0, Room: 1}, HallPos: 1},
		{Amphipod: Amphipod{Kind: 1, Room: 2}, HallPos: 7},
	}, prob, []int{0, 1, 3, 5, 7, 9, 10}, []bool{true, false, true, false, false, true, true})
}

func TestCheckBlocked(t *testing.T) {
	ck := func(moves []Move, prob Problem, expectOccupied int, checkMoves []Move, expects []bool) {
		if obtained := computeOccupied(moves, prob); obtained != expectOccupied {
			t.Fatalf("failed computeOccupied %+v %+v: obtained %011b, expected %011b", moves, prob.Targets, obtained, expectOccupied)
		}
		for i, checkMove := range checkMoves {
			expect := expects[i]
			obtained := checkBlocked(checkMove, moves, prob)
			if obtained != expect {
				t.Errorf("failed %+v: obtained %t, expected %t", checkMove, obtained, expect)
			}
		}
	}

	prob := Problem{Targets: [4]int{0, 1, 2, 3}}

	ck(
		[]Move{ // expected outcome:
			{Amphipod: Amphipod{Kind: 1, Room: 2}, HallPos: 3, Out: true},
			{Amphipod: Amphipod{Kind: 0, Room: 1}, HallPos: 0, Out: true},
			{Amphipod: Amphipod{Kind: 1, Room: 2}, HallPos: 3, Out: false},
			{Amphipod: Amphipod{Kind: 0, Room: 3}, HallPos: 5, Out: true},
		}, prob,
		0b100001,
		[]Move{
			{Amphipod: Amphipod{Kind: 1, Room: 0}, HallPos: 0, Out: true},
			{Amphipod: Amphipod{Kind: 1, Room: 0}, HallPos: 1, Out: true},
			{Amphipod: Amphipod{Kind: 1, Room: 0}, HallPos: 3, Out: true},
			{Amphipod: Amphipod{Kind: 1, Room: 0}, HallPos: 5, Out: true},
			{Amphipod: Amphipod{Kind: 1, Room: 0}, HallPos: 7, Out: true},
			{Amphipod: Amphipod{Kind: 1, Room: 2}, HallPos: 5, Out: true},
			{Amphipod: Amphipod{Kind: 1, Room: 2}, HallPos: 7, Out: true},
		},
		[]bool{false, true, true, false, false, false, true},
	)
}
