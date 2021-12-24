package p23

import (
	"fmt"
	"io"
	"sort"

	. "aoc2021/helpers"
)

const NumRooms = 4

type Amphipod struct {
	Kind   int  // 0–3 (A–D)
	Room   int  // coordinate (2, 4, 6 or 8)
	Depth  int  // depth in room (0–max)
	AtRest bool // already in final position
}

type Problem struct {
	Amphipods [][NumRooms]Amphipod // [depth][room]
	Targets   [NumRooms]int
}

type Move struct {
	Amphipod
	HallPos int
	Out     bool
}

func roomToHallPos(room int) (hpos int) {
	return 2*room + 2
}

var validHallPos = [...]int{0, 1, 3, 5, 7, 9, 10}

var kindEnergy = [...]int{
	0: 1,
	1: 10,
	2: 100,
	3: 1000,
}

func (m Move) ComesBefore(other Move, prob Problem) bool {
	log := func(v ...interface{}) {
		//if m.Kind == 0 && m.Room == 3 && m.Depth == 1 {
		//	fmt.Println(v...)
		//}
	}

	if m.Room == other.Room && m.Depth == other.Depth && m.Out && !other.Out {
		log(m, "comes before", other, "because leave-before-enter")
		return true // we must leave before entering
	}

	if m.Room == other.Room && m.Depth < other.Depth && m.Out && other.Out {
		log(m, "comes before", other, "because blocker-leave-before-leave")
		return true // blocker must leave before we do
	}

	//targetRoom := prob.Targets[m.Kind]
	//if targetRoom == other.Room && !m.Out && other.Out { // entering target room occupied by other
	//	if other.Depth == 0 {
	//		log(m, "comes before", other, "because other-top-leave-before-enter")
	//		return true // we're top amphi and must leave before other enters our room
	//	}
	//	if other.Depth == 1 && !other.AtRest {
	//		log(m, "comes before", other, "because other-bottom-non-atrest-leave-before-enter")
	//		return true // we're bottom non-at-rest amphi and must leave before other enters our room
	//	}
	//}

	otherTargetRoom := prob.Targets[other.Kind]
	//log("other entering our room", m, other, otherTargetRoom)
	if otherTargetRoom == m.Room && m.Out && !other.Out { // other entering our room
		if !m.AtRest {
			log(m, "comes before", other, "because non-atrest-leave-before-other-enter")
			return true // we're a non-at-rest amphi in original room and must leave before other enters our room
		}
	}

	return false
}

func findMoveInsertRange(ins Move, moves []Move, prob Problem) (minIdx, maxIdx int) {
	minIdx, maxIdx = 0, len(moves)
	for i, m := range moves {
		if ins.ComesBefore(m, prob) && maxIdx > i {
			maxIdx = i
		}
		if m.ComesBefore(ins, prob) {
			minIdx = i + 1
		}
	}
	return minIdx, maxIdx
}

func order(a, b int) (int, int) {
	if a > b {
		return b, a
	}
	return a, b
}

func findMinMoveRange(amph Amphipod, prob Problem) (min, max int) {
	return order(roomToHallPos(amph.Room), roomToHallPos(prob.Targets[amph.Kind]))
}

func listHallPoses(moveStart, moveEnd int, extraDist int) (out []int) {
	if extraDist == 0 {
		from, to := 0, 0
		for i, p := range validHallPos {
			if p < moveStart {
				from = i + 1
			} else if p > moveEnd {
				to = i
				break
			}
		}
		return validHallPos[from:to]
	}

	out = make([]int, 0, 2)
	for _, p := range validHallPos {
		if p < moveStart && moveStart-p == extraDist {
			out = append(out, p)
		} else if p > moveEnd && p-moveEnd == extraDist {
			out = append(out, p)
		}
	}
	return out
}

// N.B.: both bounds inclusive
func getMoveInterval(m Move, prob Problem) (int, int) {
	var from, to int
	if m.Out {
		from, to = roomToHallPos(m.Room), m.HallPos
	} else {
		from, to = m.HallPos, roomToHallPos(prob.Targets[m.Kind])
	}
	return order(from, to)
}

func checkBlocking(pos int, moves []Move, prob Problem) (ok bool) {
	for _, m := range moves {
		start, end := getMoveInterval(m, prob)
		if start <= pos && pos <= end {
			return false
		}
	}
	return true
}

func computeOccupied(moves []Move, prob Problem) (occupied int) {
	for _, m2 := range moves {
		if m2.Out {
			occupied |= 1 << m2.HallPos
		} else {
			occupied &^= 1 << m2.HallPos
		}
	}
	return occupied
}

func checkBlocked(m Move, moves []Move, prob Problem) (ok bool) {
	start, end := getMoveInterval(m, prob)
	occupied := computeOccupied(moves, prob)
	return (occupied&(1<<(end+1)-1))>>start == 0
}

func iterateMoveOrders(indent string, moves []Move, order []Amphipod, extraDistByKind [4]int, prob Problem) (winner []Move) {
	if len(order) == 0 {
		return moves // all moves added successfully
	}

	log := func(...interface{}) {} //fmt.Println

	amph := order[0]
	order = order[1:]
	if amph.AtRest {
		return iterateMoveOrders(indent, moves, order, extraDistByKind, prob)
	}

	log(indent, string(rune(amph.Kind+'A')), amph.Room, amph.Depth)

	out, in := Move{Amphipod: amph, Out: true}, Move{Amphipod: amph, Out: false}
	outMinIdx, outMaxIdx := findMoveInsertRange(out, moves, prob)
	inMinIdx, inMaxIdx := findMoveInsertRange(in, moves, prob)
	//fmt.Println(amph.Room, amph.Depth, "step", len(moves)/2, ":", outMinIdx, outMaxIdx, inMinIdx, inMaxIdx)
	if outMaxIdx < outMinIdx || inMaxIdx < inMinIdx {
		log(indent+"  ", out, in, "no valid insertion point", moves)
		return // no valid insertion point
	}

	moveStart, moveEnd := findMinMoveRange(amph, prob)
	//fmt.Println("  move", moveStart, moveEnd)

	extraDistAllowed := extraDistByKind[amph.Kind]
	for extraDist := 0; extraDist <= extraDistAllowed; extraDist++ {
		extraDistByKind[amph.Kind] = extraDistAllowed - extraDist

		hallPoses := listHallPoses(moveStart, moveEnd, extraDist)

		if len(hallPoses) == 0 {
			log(indent+"  ", out, in, "insufficient energy with", extraDist, "of", extraDistAllowed, moves)
			continue
		}

		for outIdx := outMinIdx; outIdx <= outMaxIdx; outIdx++ {
			for inIdx := inMinIdx; inIdx <= inMaxIdx; inIdx++ {
				if inIdx < outIdx {
					continue
				}

				for _, hpos := range hallPoses {
					log(indent+"  ", string(rune(amph.Kind+'A')), amph.Room, amph.Depth, hpos)

					in.HallPos, out.HallPos = hpos, hpos
					if !checkBlocking(hpos, moves[outIdx:inIdx], prob) ||
						!checkBlocked(out, moves[:outIdx], prob) ||
						!checkBlocked(in, moves[:inIdx], prob) {
						log(indent+"    ", out, in, hpos, "fully blocked", moves)
						continue
					}

					nextMoves := make([]Move, 0, len(moves)+2)
					nextMoves = append(nextMoves, moves[:outIdx]...)
					nextMoves = append(nextMoves, out)
					nextMoves = append(nextMoves, moves[outIdx:inIdx]...)
					nextMoves = append(nextMoves, in)
					nextMoves = append(nextMoves, moves[inIdx:]...)

					if winner = iterateMoveOrders(indent+"  ", nextMoves, order, extraDistByKind, prob); winner != nil {
						return winner
					}
				}
			}
		}
	}
	return nil
}

// IsAtRest returns true iff amphipod at room/depth is in its final location and need not move (in target room and not
// blocking another amphipod).
func (prob Problem) IsAtRest(room int, depth int) bool {
	a := prob.Amphipods[depth][room]
	if a.Room == prob.Targets[a.Kind] { // in target room
		if a.Depth == len(prob.Amphipods)-1 { // at bottom of target room
			return true
		} else if prob.IsAtRest(room, depth+1) { // amphipods below is at rest, so can be are
			return true
		}
	}
	return false
}

// MinimumEnergy returns minimum energy to shuffle the amphipods, ignoring all blocking.
// Actual energy will be this plus extra distance/energy needed to avoid blockers.
func (prob Problem) MinimumEnergy() (energy int) {
	nonAtRestByKind := [4]int{0, 0, 0, 0}
	for depth, rooms := range prob.Amphipods {
		for room, amph := range rooms {
			if amph.AtRest {
				continue
			}
			nonAtRestByKind[amph.Kind]++
			sideDist := Abs(roomToHallPos(room) - roomToHallPos(prob.Targets[amph.Kind]))
			// out + over + in to top slot (dist for 1/2 amphipod to reach bottom slot is added below)
			dist := (depth + 1) + (sideDist) + 1
			fmt.Println("extra energy", room, depth, prob.Targets[amph.Kind], sideDist, dist)
			energy += dist * kindEnergy[amph.Kind]
		}
	}
	for kind, num := range nonAtRestByKind {
		// extra distance for amphipods moving into target slot; first takes 0 extra, second takes 1 extra
		energy += (num - 1) * num / 2 * kindEnergy[kind]
	}
	return energy
}

func A(in io.Reader) {
	//start := [2][4]rune{ // sample
	//  {'B','C','B','D'},
	//  {'A','D','C','A'},
	//}
	start := [2][4]rune{ // real
		{'B', 'B', 'D', 'D'},
		{'C', 'A', 'A', 'C'},
	}

	const NumDepths = 2

	prob := Problem{Amphipods: make([][NumRooms]Amphipod, NumDepths), Targets: [4]int{0, 1, 2, 3}}
	for depth, rooms := range start {
		for room, c := range rooms {
			amph := Amphipod{
				Kind:  int(c - 'A'),
				Room:  room,
				Depth: depth,
			}
			prob.Amphipods[depth][room] = amph
		}
	}

	var sortedAmph []Amphipod
	for depth, rooms := range prob.Amphipods {
		for room, amph := range rooms {
			amph.AtRest = prob.IsAtRest(room, depth)
			prob.Amphipods[depth][room] = amph
			sortedAmph = append(sortedAmph, amph)
		}
	}
	sort.Slice(sortedAmph, func(i, j int) bool { return sortedAmph[i].Kind > sortedAmph[j].Kind })

	fmt.Println("problem", prob)
	fmt.Println("minimum energy", prob.MinimumEnergy())

	var extraDists [NumRooms]int
	for {
		win := iterateMoveOrders("", nil, sortedAmph, extraDists, prob)
		if win != nil {
			fmt.Println("win", extraDists)
			printMoves(win, prob)
			break
		}

		for i, v := range extraDists {
			if v == 8 {
				extraDists[i] = 0
			} else {
				extraDists[i]++
				break
			}
		}
		if extraDists == ([NumRooms]int{}) {
			break
		}
	}

	energy := prob.MinimumEnergy()
	for kind, extraDist := range extraDists {
		energy += 2 * extraDist * kindEnergy[kind]
	}
	fmt.Println(energy)
}

func printMoves(moves []Move, prob Problem) {
	for _, m := range moves {
		kindStr := string(rune(m.Kind + 'A'))
		if m.Out {
			fmt.Println(kindStr, "at", m.Room, ",", m.Depth, "moves to hall", m.HallPos)
		} else {
			fmt.Println(kindStr, "at hall", m.HallPos, "moves to room", prob.Targets[m.Kind])
		}
	}
}

func B(in io.Reader) {
	start := [4][4]rune{ // sample
		{'B', 'C', 'B', 'D'},
		{'D', 'C', 'B', 'A'},
		{'D', 'B', 'A', 'C'},
		{'A', 'D', 'C', 'A'},
	}
	//start := [4][4]rune{ // real
	//	{'B', 'B', 'D', 'D'},
	//	{'D', 'C', 'B', 'A'},
	//	{'D', 'B', 'A', 'C'},
	//	{'C', 'A', 'A', 'C'},
	//}

	const NumDepths = 4

	prob := Problem{Amphipods: make([][NumRooms]Amphipod, NumDepths), Targets: [4]int{0, 1, 2, 3}}
	for depth, rooms := range start {
		for room, c := range rooms {
			amph := Amphipod{
				Kind:  int(c - 'A'),
				Room:  room,
				Depth: depth,
			}
			prob.Amphipods[depth][room] = amph
		}
	}

	var sortedAmph []Amphipod
	for depth, rooms := range prob.Amphipods {
		for room, amph := range rooms {
			amph.AtRest = prob.IsAtRest(room, depth)
			prob.Amphipods[depth][room] = amph
			sortedAmph = append(sortedAmph, amph)
		}
	}
	sort.Slice(sortedAmph, func(i, j int) bool {
		return sortedAmph[i].Depth < sortedAmph[j].Depth ||
			sortedAmph[i].Kind > sortedAmph[j].Kind
	})

	fmt.Println("problem", prob)
	fmt.Println("minimum energy", prob.MinimumEnergy())

	var extraDists [NumRooms]int
	for {
		fmt.Println(extraDists)
		win := iterateMoveOrders("", nil, sortedAmph, extraDists, prob)
		if win != nil {
			fmt.Println("win", extraDists)
			printMoves(win, prob)
			break
		}

		for i, v := range extraDists {
			if v == 16 {
				extraDists[i] = 0
			} else {
				extraDists[i]++
				break
			}
		}
		if extraDists == ([NumRooms]int{}) {
			break
		}
	}

	energy := prob.MinimumEnergy()
	for kind, extraDist := range extraDists {
		energy += 2 * extraDist * kindEnergy[kind]
	}
	fmt.Println(energy)
}
