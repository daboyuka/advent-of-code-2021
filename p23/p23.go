package p23

import (
	"fmt"
	"io"
	"math"

	. "aoc2021/helpers"
)

const NumRooms = 4

type Loc struct {
	Room  int
	Depth int
}

type Move struct {
	Loc
	Hall int
	Park []HallAmphipod // hall slots
}

type HallAmphipod struct {
	Kind int
	Hall int
}

type Amphipod struct {
	Loc
	Kind   int  // 0–3 (A–D)
	AtRest bool // already in final position
}

type Problem [][NumRooms]Amphipod // [depth][room]

var hallPoses = [...]int{0, 1, 3, 5, 7, 9, 10}

var kindEnergy = [...]int{
	0: 1,
	1: 10,
	2: 100,
	3: 1000,
}

// IsAtRest returns true iff amphipod at room/depth is in its final location and need not move (in target room and not
// blocking another amphipod).
func (prob Problem) IsAtRest(room int, depth int) bool {
	a := prob[depth][room]
	if a.Room == a.Kind { // in target room
		if a.Depth == len(prob)-1 { // at bottom of target room
			return true
		} else if prob.IsAtRest(room, depth+1) { // amphipods below is at rest, so can be are
			return true
		}
	}
	return false
}

func (prob Problem) Depth() int { return len(prob) }

// MinimumEnergy returns minimum energy to shuffle the amphipods, ignoring all blocking.
// Actual energy will be this plus extra distance/energy needed to avoid blockers.
func (prob Problem) MinimumEnergy() (energy int) {
	nonAtRestByKind := [4]int{0, 0, 0, 0}
	for depth, rooms := range prob {
		for room, amph := range rooms {
			if amph.AtRest {
				continue
			}
			nonAtRestByKind[amph.Kind]++
			sideDist := Abs(roomToHallPos(room) - roomToHallPos(amph.Kind))
			// out + over (distance into target slot counted below)
			dist := depth + 1 + sideDist
			energy += dist * kindEnergy[amph.Kind]
		}
	}
	for kind, num := range nonAtRestByKind {
		// extra distance for amphipods moving into target slot; first takes 0 extra, second takes 1 extra
		energy += num * (num + 1) / 2 * kindEnergy[kind]
	}
	return energy
}

type PartialSoln struct {
	Moves           []Move
	ExtraSideEnergy int
	ExtraDist       [NumRooms]int

	Hall      [7]int // correspond to hallPoses
	NextDepth [NumRooms]int
	RoomClean [NumRooms]bool
}

func roomToHallSlot(room int) (hslot int) { return room + 1 }

func (soln PartialSoln) PossibleHallSlots(leaveRoom int) (minHSlot, maxHSlot int) { // both inclusive
	minHSlot, maxHSlot = roomToHallSlot(leaveRoom)+1, roomToHallSlot(leaveRoom) // left and right backwards = no interval
	for i := minHSlot - 1; i >= 0 && soln.Hall[i] == -1; i-- {
		minHSlot = i
	}
	for i := maxHSlot + 1; i < len(soln.Hall) && soln.Hall[i] == -1; i++ {
		maxHSlot = i
	}
	return minHSlot, maxHSlot
}

func (soln *PartialSoln) ParkSingle(hslot int, afterMove *Move) bool {
	kind := soln.Hall[hslot]
	if kind == -1 {
		return true
	} else if !soln.RoomClean[kind] {
		return false
	}

	target := roomToHallSlot(kind) // left side slot
	if target < hslot {
		target++ // right side slot
	}

	cur := hslot
	for cur != target {
		if cur < target {
			cur++
		} else {
			cur--
		}
		if soln.Hall[cur] != -1 {
			return false
		}
	}

	soln.Hall[hslot] = -1
	afterMove.Park = append(afterMove.Park, HallAmphipod{Kind: kind, Hall: hslot})
	return true
}

func (soln *PartialSoln) Park(leftHSlot, rightHSlot int, afterMove *Move) {
	for {
		if leftHSlot >= 0 && soln.ParkSingle(leftHSlot, afterMove) {
			leftHSlot--
		} else if rightHSlot < len(hallPoses) && soln.ParkSingle(rightHSlot, afterMove) {
			rightHSlot++
		} else {
			break // no progress on either end
		}
	}
}

func (soln PartialSoln) Win() bool {
	for _, k := range soln.Hall { // hall must be clear
		if k != -1 {
			return false
		}
	}
	for _, c := range soln.RoomClean { // all rooms must be clean
		if !c {
			return false
		}
	}
	return true
}

func roomToHallPos(room int) int { return 2*room + 2 }

func getExcessSideDist(amph Amphipod, hslot int) int {
	startPos, endPos := roomToHallPos(amph.Room), roomToHallPos(amph.Kind)
	hallPos := hallPoses[hslot]
	minDist := Abs(startPos - endPos)
	actualDist := Abs(startPos-hallPos) + Abs(hallPos-endPos)
	return actualDist - minDist
}

// for debugging: expected steps of sample in part 2
func onGoldenPath(soln PartialSoln) bool {
	for i, m := range soln.Moves {
		ok := false
		switch i {
		case 0:
			ok = m.Room == 3 && m.Depth == 0 && m.Hall == 6
		case 1:
			ok = m.Room == 3 && m.Depth == 1 && m.Hall == 0
		case 2:
			ok = m.Room == 2 && m.Depth == 0 && m.Hall == 5
		case 3:
			ok = m.Room == 2 && m.Depth == 1 && m.Hall == 4
		case 4:
			ok = m.Room == 2 && m.Depth == 2 && m.Hall == 1
		case 5:
			ok = m.Room == 1 && m.Depth == 0 && m.Hall == 3
		case 6:
			ok = m.Room == 1 && m.Depth == 1 && m.Hall == 3
		case 7:
			ok = m.Room == 1 && m.Depth == 2 && m.Hall == 3
		case 8:
			ok = m.Room == 1 && m.Depth == 3 && m.Hall == 2
		case 9:
			ok = m.Room == 3 && m.Depth == 2 && m.Hall == 4
		case 10:
			ok = m.Room == 3 && m.Depth == 3 && m.Hall == 5
		case 11:
			ok = m.Room == 0 && m.Depth == 0 && m.Hall == 2
		case 12:
			ok = m.Room == 0 && m.Depth == 1 && m.Hall == 2
		case 13:
			ok = m.Room == 0 && m.Depth == 2 && m.Hall == 2
		}
		if !ok {
			return false
		}
	}
	return true
}

func iterateMoveOrders(prob Problem, soln PartialSoln, bestSoln *PartialSoln) (foundWin bool) {
	if soln.ExtraSideEnergy >= bestSoln.ExtraSideEnergy {
		return false
	} else if soln.Win() {
		if onGoldenPath(soln) {
			fmt.Println("GOLDEN WIN")
		}
		fmt.Println("win", soln.ExtraSideEnergy)
		*bestSoln = soln
		return true
	}

	for room, depth := range soln.NextDepth {
		if soln.RoomClean[room] {
			continue
		}
		amph := prob[depth][room]

		startHSlot := roomToHallSlot(room)
		minHSlot, maxHSlot := soln.PossibleHallSlots(room)
		for combo := 0; combo < 12; combo++ {
			// Zigzag outward, keeping min dist
			leftOfRoom := combo%2 == 0
			var hslot int
			if leftOfRoom {
				hslot = startHSlot - combo/2 // leftward
			} else {
				hslot = startHSlot + 1 + combo/2 // rightward
			}
			if hslot < minHSlot || hslot > maxHSlot {
				continue // out of bounds
			}

			nextMove := Move{Loc: Loc{room, depth}, Hall: hslot}

			soln2 := soln

			soln2.NextDepth[room] = depth + 1
			if depth+1 == prob.Depth() || prob[depth+1][room].AtRest {
				soln2.RoomClean[room] = true
			}

			soln2.Hall[hslot] = amph.Kind
			soln2.Moves = append(make([]Move, 0, len(soln.Moves)+1), soln.Moves...)
			soln2.Moves = append(soln2.Moves, nextMove)

			excessDist := getExcessSideDist(amph, hslot)
			soln2.ExtraSideEnergy += excessDist * kindEnergy[amph.Kind]
			soln2.ExtraDist[amph.Kind] += excessDist / 2

			lastMove := &soln2.Moves[len(soln2.Moves)-1]
			if leftOfRoom {
				soln2.Park(hslot, maxHSlot+1, lastMove)
			} else {
				soln2.Park(minHSlot-1, hslot, lastMove)
			}

			if onGoldenPath(soln2) {
				fmt.Println("golden path", len(soln2.Moves), ":", lastMove.Loc, "->", lastMove.Hall, ";", leftOfRoom, minHSlot, maxHSlot)
				fmt.Println("golden path", len(soln2.Moves), soln2)
			}

			if iterateMoveOrders(prob, soln2, bestSoln) {
				foundWin = true
			}
		}
	}

	return foundWin
}

func A(in io.Reader) {
	//start := [2][4]rune{ // sample
	//	{'B', 'C', 'B', 'D'},
	//	{'A', 'D', 'C', 'A'},
	//}
	start := [2][4]rune{ // real
		{'B', 'B', 'D', 'D'},
		{'C', 'A', 'A', 'C'},
	}

	const NumDepths = 2

	prob := make(Problem, NumDepths)
	for depth, rooms := range start {
		for room, c := range rooms {
			prob[depth][room] = Amphipod{
				Kind: int(c - 'A'),
				Loc:  Loc{Room: room, Depth: depth},
			}
		}
	}
	for depth, rooms := range prob {
		for room := range rooms {
			prob[depth][room].AtRest = prob.IsAtRest(room, depth)
		}
	}

	fmt.Println("problem", prob)

	initSoln := PartialSoln{
		Hall: [7]int{-1, -1, -1, -1, -1, -1, -1},
	}

	bestSoln := PartialSoln{ExtraSideEnergy: math.MaxInt}
	iterateMoveOrders(prob, initSoln, &bestSoln)

	for _, m := range bestSoln.Moves {
		fmt.Println("move", string(rune(prob[m.Depth][m.Room].Kind+'A')), "from", m.Room, "/", m.Depth, "to", m.Hall)
		for _, p := range m.Park {
			fmt.Println("return", string(rune(p.Kind+'A')), "from", p.Hall)
		}
	}

	fmt.Println(prob.MinimumEnergy(), bestSoln)
	fmt.Println(prob.MinimumEnergy() + bestSoln.ExtraSideEnergy)
}

//|||||||||||||
//|D..D.A.....|
//||| |B|C| |||
//  |A|B|C| |
//  |A|B|C|D|
//  |A|B|C|D|
//  |||||||||

// guess 55071 -> too low
func B(in io.Reader) {
	//start := [4][4]rune{ // sample
	//	{'B', 'C', 'B', 'D'},
	//	{'D', 'C', 'B', 'A'},
	//	{'D', 'B', 'A', 'C'},
	//	{'A', 'D', 'C', 'A'},
	//}
	start := [4][4]rune{ // real
		{'B', 'B', 'D', 'D'},
		{'D', 'C', 'B', 'A'},
		{'D', 'B', 'A', 'C'},
		{'C', 'A', 'A', 'C'},
	}

	const NumDepths = 4

	prob := make(Problem, NumDepths)
	for depth, rooms := range start {
		for room, c := range rooms {
			prob[depth][room] = Amphipod{
				Kind: int(c - 'A'),
				Loc:  Loc{Room: room, Depth: depth},
			}
		}
	}
	for depth, rooms := range prob {
		for room := range rooms {
			prob[depth][room].AtRest = prob.IsAtRest(room, depth)
		}
	}

	fmt.Println("problem", prob)

	initSoln := PartialSoln{
		Hall: [7]int{-1, -1, -1, -1, -1, -1, -1},
	}

	bestSoln := PartialSoln{ExtraSideEnergy: math.MaxInt}
	iterateMoveOrders(prob, initSoln, &bestSoln)

	for _, m := range bestSoln.Moves {
		fmt.Println("move", string(rune(prob[m.Depth][m.Room].Kind+'A')), "from", m.Room, "/", m.Depth, "to", m.Hall)
		for _, p := range m.Park {
			fmt.Println("return", string(rune(p.Kind+'A')), "from", p.Hall)
		}
	}

	fmt.Println(prob.MinimumEnergy(), bestSoln)
	fmt.Println(prob.MinimumEnergy() + bestSoln.ExtraSideEnergy)
}
