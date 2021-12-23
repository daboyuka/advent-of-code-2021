package p23

import (
	"fmt"
	"io"
	"sort"
	"strings"

	. "aoc2021/helpers"
)

type State struct {
	Hall  [11]rune
	Rooms [4][2]rune // [4] left to right, [2] top to bottom
}

func (s State) String() string {
	sb := strings.Builder{}
	sb.WriteString("#############\n")

	sb.WriteRune('#')
	for _, c := range s.Hall {
		if c == 0 {
			sb.WriteRune('.')
		} else {
			sb.WriteRune(c)
		}
	}
	sb.WriteString("#\n")

	sb.WriteString("###")
	for _, room := range s.Rooms {
		if room[0] == 0 {
			sb.WriteRune('.')
		} else {
			sb.WriteRune(room[0])
		}
		sb.WriteRune('#')
	}
	sb.WriteString("##\n")

	sb.WriteString("  #")
	for _, room := range s.Rooms {
		if room[1] == 0 {
			sb.WriteRune('.')
		} else {
			sb.WriteRune(room[1])
		}
		sb.WriteRune('#')
	}
	sb.WriteString("  \n")

	sb.WriteString("  #########  \n")
	return sb.String()
}

func hallPosForRoom(room int) (hpos int) {
	return 2*room + 2
}

func tryMoveHallToRoom(state State, kind rune, fromPos int, targets [4]int) (dist int, outState State, ok bool) {
	room := targets[kind-'A']
	if roomTop := state.Rooms[room][0]; roomTop != 0 {
		return 0, State{}, false // room full
	} else if bottomKind := state.Rooms[room][1]; bottomKind != 0 && bottomKind != kind && bottomKind != kind-'A'+'a' {
		return 0, State{}, false // room has wrong kind of amphi inside
	}

	bottomOpen := state.Rooms[room][1] == 0

	entryPos := hallPosForRoom(room)
	if fromPos < entryPos {
		for _, c := range state.Hall[fromPos+1 : entryPos] {
			if c != 0 {
				return 0, State{}, false
			}
		}
		dist = 1 + (entryPos - fromPos)
	} else { // hpos > leftHall
		for _, c := range state.Hall[entryPos+1 : fromPos] {
			if c != 0 {
				return 0, State{}, false
			}
		}
		dist = 1 + (fromPos - entryPos)
	}

	newKind := kind - 'A' + 'a'
	state.Hall[fromPos] = 0
	if bottomOpen {
		state.Rooms[room][1] = newKind
		dist++
	} else {
		state.Rooms[room][0] = newKind
	}
	return dist, state, true
}

func tryMoveRoomToHall(state State, kind rune, fromRoom int, fromBottom bool, hpos int) (dist int, outState State, ok bool) {
	if fromBottom && state.Rooms[fromRoom][0] != 0 {
		return 0, State{}, false // blocked
	}
	entryPos := hallPosForRoom(fromRoom)
	if hpos < entryPos {
		for _, c := range state.Hall[hpos:entryPos] {
			if c != 0 {
				return 0, State{}, false
			}
		}
		dist = 1 + (entryPos - hpos)
	} else { // hpos > entryPos
		for _, c := range state.Hall[entryPos+1 : hpos+1] {
			if c != 0 {
				return 0, State{}, false
			}
		}
		dist = 1 + (hpos - entryPos)
	}

	state.Hall[hpos] = kind
	if fromBottom {
		state.Rooms[fromRoom][1] = 0
		dist++
	} else {
		state.Rooms[fromRoom][0] = 0
	}
	return dist, state, true
}

func tryMoveRoomToRoom(state State, kind rune, fromRoom int, fromBottom bool, targets [4]int) (dist int, outState State, ok bool) {
	toRoom := targets[kind-'A']
	if fromRoom == toRoom {
		return 0, State{}, false // already home
	} else if fromBottom && state.Rooms[fromRoom][0] != 0 {
		return 0, State{}, false // blocked in
	} else if state.Rooms[toRoom][0] != 0 {
		return 0, State{}, false // blocked out
	} else if bottomKind := state.Rooms[toRoom][1]; bottomKind != 0 && bottomKind != kind && bottomKind != kind-'A'+'a' {
		return 0, State{}, false // room mixed
	}

	fromEntry, toEntry := hallPosForRoom(fromRoom), hallPosForRoom(toRoom)
	if fromEntry > toEntry {
		fromEntry, toEntry = toEntry, fromEntry
	}

	for _, c := range state.Hall[fromEntry+1 : toEntry] {
		if c != 0 {
			return 0, State{}, false
		}
	}

	newKind := kind - 'A' + 'a'
	dist = 2 + toEntry - fromEntry // move out 1, over, in 1
	if fromBottom {
		dist++ // move from bottom
		state.Rooms[fromRoom][1] = 0
	} else {
		state.Rooms[fromRoom][0] = 0
	}
	if state.Rooms[toRoom][1] == 0 {
		dist++ // move to bottom
		state.Rooms[toRoom][1] = newKind
	} else {
		state.Rooms[toRoom][0] = newKind
	}

	return dist, state, true
}

var validHallPos = [...]int{0, 1, 3, 5, 7, 9, 10}

type SolnSet struct {
	States   []State
	Energies []int
}

func (s SolnSet) Len() int { return len(s.States) }

func (s SolnSet) Less(i, j int) bool { return s.Energies[i] < s.Energies[j] }

func (s SolnSet) Swap(i, j int) {
	s.States[i], s.States[j] = s.States[j], s.States[i]
	s.Energies[i], s.Energies[j] = s.Energies[j], s.Energies[i]
}

func possibleMoves(state State, targets [4]int) (solns SolnSet) {
	var r2r SolnSet
	for roomIdx, room := range state.Rooms {
		for depth, c := range room {
			if c != 0 && c >= 'A' && c <= 'D' {
				if dist, outState, ok := tryMoveRoomToRoom(state, c, roomIdx, depth != 0, targets); ok {
					r2r.States, r2r.Energies = append(r2r.States, outState), append(r2r.Energies, dist*kindEnergy[c])
				}
			}
		}
	}
	var h2r SolnSet
	for i, c := range state.Hall {
		if c != 0 && c >= 'A' && c <= 'D' {
			if dist, outState, ok := tryMoveHallToRoom(state, c, i, targets); ok {
				h2r.States, h2r.Energies = append(h2r.States, outState), append(h2r.Energies, dist*kindEnergy[c])
			}
		}
	}
	var r2h SolnSet
	for roomIdx, room := range state.Rooms {
		for depth, c := range room {
			if c != 0 && c >= 'A' && c <= 'D' {
				for _, hpos := range validHallPos {
					if dist, outState, ok := tryMoveRoomToHall(state, c, roomIdx, depth != 0, hpos); ok {
						r2h.States, r2h.Energies = append(r2h.States, outState), append(r2h.Energies, dist*kindEnergy[c])
					}
				}
			}
		}
	}

	sort.Sort(r2r)
	sort.Sort(h2r)
	sort.Sort(r2h)

	solns = r2r
	solns.States, solns.Energies = append(solns.States, h2r.States...), append(solns.Energies, h2r.Energies...)
	solns.States, solns.Energies = append(solns.States, r2h.States...), append(solns.Energies, r2h.Energies...)
	return solns
}

func minRemainingEnergy(state State, targets [4]int) (minEnergy int) {
	for fromRoom, room := range state.Rooms {
		for depth, c := range room {
			if c < 'A' || c > 'D' {
				continue
			}

			toRoom := targets[c-'A']
			if fromRoom == toRoom {
				continue
			}

			minDist := 2 + 2*Abs(fromRoom-toRoom) + depth
			minEnergy += minDist * kindEnergy[c]
		}
	}
	return minEnergy
}

func (s State) Won() bool {
	for _, room := range s.Rooms {
		a, b := room[0], room[1]
		if a >= 'a' {
			a -= 'a' - 'A'
		}
		if b >= 'b' {
			b -= 'b' - 'B'
		}
		if a == 0 || a != b {
			return false
		}
	}
	return true
}

var kindEnergy = [...]int{
	'A': 1,
	'B': 10,
	'C': 100,
	'D': 1000,
}

type Solution struct {
	Path   []State
	Energy int
}

func solve(cur State, energy int, path []State, bestSoln *Solution, targets [4]int) {
	if cur.Won() {
		if bestSoln.Energy == -1 || bestSoln.Energy > energy {
			fmt.Println("new winner", energy)
			*bestSoln = Solution{Path: path, Energy: energy}
		}
		return
	} else if bestSoln.Energy != -1 && energy >= bestSoln.Energy {
		return // already too much energy
	} else if bestSoln.Energy != -1 && energy+minRemainingEnergy(cur, targets) >= bestSoln.Energy {
		return // already too much energy
	}

	next := possibleMoves(cur, targets)
	for i, nextState := range next.States {
		energyToNext := energy + next.Energies[i]
		solve(nextState, energyToNext, append(path, nextState), bestSoln, targets)
	}
}

func A(in io.Reader) {
	//start := State{
	//	Rooms: [4][2]rune{
	//		{'B', 'A'},
	//		{'C', 'D'},
	//		{'B', 'C'},
	//		{'D', 'A'},
	//	},
	//}
	start := State{
		Rooms: [4][2]rune{
			{'B', 'C'},
			{'B', 'A'},
			{'D', 'A'},
			{'D', 'C'},
		},
	}

	fmt.Println(start)

	move := func(which int) {
		next := possibleMoves(start, [4]int{0, 1, 2, 3})
		fmt.Println(which, next.Energies[which])
		fmt.Println(next.States[which])

		start = next.States[which]
	}
	_ = move

	//move(16)
	//move(0)
	//move(2)
	//move(0)
	//move(0)
	//move(3)
	//move(3)
	//move(0)
	//move(0)
	//move(0)

	//states, energies := possibleMoves(start, [4]int{0, 1, 2, 3})
	//for i := range states {
	//	fmt.Println(i, energies[i])
	//	fmt.Println(states[i])
	//}

	//return

	bestSoln := Solution{Energy: -1}
	solve(start, 0, []State{start}, &bestSoln, [4]int{0, 1, 2, 3})
	for _, s := range bestSoln.Path {
		fmt.Println(s)
	}
	fmt.Println(bestSoln.Energy)
}

func B(in io.Reader) {
	lines := ReadLines(in)
	_ = lines
}
