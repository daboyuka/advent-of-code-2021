package p19

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

type Pos3 [3]int

func parse(in io.Reader) (out []map[Pos3]bool) {
	lgs := ReadLinegroups(in)
	for _, lg := range lgs {
		scanner := map[Pos3]bool{}
		for _, l := range lg[1:] {
			xS, yS, zS := Split3(l, ",", ",")
			scanner[Pos3{Atoi(xS), Atoi(yS), Atoi(zS)}] = true
		}
		out = append(out, scanner)
	}
	return out
}

type Mat3 [3][3]int // [row][col]

func (m Mat3) Mul(p Pos3) (out Pos3) {
	for outI := range out {
		for i, l := range m[outI] {
			out[outI] += l * p[i]
		}
	}
	return out
}

func (m Mat3) Transpose() (out Mat3) {
	for i, r := range m {
		for j, v := range r {
			out[j][i] = v
		}
	}
	return out
}

func (p Pos3) Cross(p2 Pos3) Pos3 {
	return Pos3{
		p[1]*p2[2] - p[2]*p2[1],
		p[2]*p2[0] - p[0]*p2[2],
		p[0]*p2[1] - p[1]*p2[0],
	}
}

var Rots [24]Mat3

func init() {
	idx := 0
	for axis := 0; axis < 3; axis++ {
		for _, axisS := range [...]int{-1, 1} {
			for upAxis := 0; upAxis < 3; upAxis++ {
				if axis == upAxis {
					continue
				}
				for _, upAxisS := range [...]int{-1, 1} {
					var fwd, up, right Pos3
					fwd[axis] = axisS
					up[upAxis] = upAxisS
					right = up.Cross(fwd)
					Rots[idx] = Mat3{right, up, fwd}.Transpose()
					idx++
				}
			}
		}
	}
}

func (p Pos3) Add(p2 Pos3) Pos3 { return Pos3{p[0] + p2[0], p[1] + p2[1], p[2] + p2[2]} }
func (p Pos3) Sub(p2 Pos3) Pos3 { return Pos3{p[0] - p2[0], p[1] - p2[1], p[2] - p2[2]} }

const threshold = 12

func checkOffset(a, b map[Pos3]bool, offset Pos3, rot Mat3) bool {
	found := 0
	for bp := range b {
		ap := rot.Mul(bp).Add(offset)
		if a[ap] {
			found++
		}
	}
	return found >= threshold
}

func findOffset(a, b map[Pos3]bool) (offset Pos3, rot Mat3, found bool) {
	for _, rot = range Rots { // rot.Mul(bPt) -> aPt
		for ap := range a {
			for bp := range b {
				offset = ap.Sub(rot.Mul(bp)) // bPt.Add(offset) -> aPt
				if checkOffset(a, b, offset, rot) {
					return offset, rot, true
				}
			}
		}
	}
	return offset, rot, false
}

type Relative struct { // .Rot.Mul(rel).Add(.Off) -> abs
	Off Pos3
	Rot Mat3
}

func A(in io.Reader) {
	scanners := parse(in)

	known := map[int]Relative{0: {}}
	knownAbs := make([]map[Pos3]bool, len(scanners))
	knownAbs[0] = scanners[0]

	for len(known) < len(scanners) {
		for i, scanner := range scanners {
			if _, ok := known[i]; ok {
				continue
			}

			// i is unknown
			for knownIdx := range known {
				knownScanner := knownAbs[knownIdx]

				if off, rot, ok := findOffset(knownScanner, scanner); ok {
					fmt.Println("scanner", i, "overlaps scanner", knownIdx, "with", off, rot)

					known[i] = Relative{Off: off, Rot: rot}
					knownAbs[i] = make(map[Pos3]bool, len(scanner))
					for bp := range scanner {
						ap := rot.Mul(bp).Add(off)
						knownAbs[i][ap] = true
					}
					break
				}
			}
		}
	}

	total := map[Pos3]bool{}
	for _, knownScanner := range knownAbs {
		for p := range knownScanner {
			total[p] = true
		}
	}

	fmt.Println(total)
	fmt.Println(len(total))
}

func B(in io.Reader) {
	scanners := parse(in)

	known := map[int]Relative{0: {}}
	knownAbs := make([]map[Pos3]bool, len(scanners))
	knownAbs[0] = scanners[0]

	for len(known) < len(scanners) {
		for i, scanner := range scanners {
			if _, ok := known[i]; ok {
				continue
			}

			// i is unknown
			for knownIdx := range known {
				knownScanner := knownAbs[knownIdx]

				if off, rot, ok := findOffset(knownScanner, scanner); ok {
					fmt.Println("scanner", i, "overlaps scanner", knownIdx, "with", off, rot)

					known[i] = Relative{Off: off, Rot: rot}
					knownAbs[i] = make(map[Pos3]bool, len(scanner))
					for bp := range scanner {
						ap := rot.Mul(bp).Add(off)
						knownAbs[i][ap] = true
					}
					break
				}
			}
		}
	}

	maxDist := 0
	for _, a := range known {
		for _, b := range known {
			if dist := a.Off.Dist(b.Off); maxDist < dist {
				maxDist = dist
			}
		}
	}

	fmt.Println(maxDist)
}

func (p Pos3) Dist(p2 Pos3) int { return Abs(p[0]-p2[0]) + Abs(p[1]-p2[1]) + Abs(p[2]-p2[2]) }
