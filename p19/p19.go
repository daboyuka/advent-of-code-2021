package p19

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

type Pos3 [3]int

func (p Pos3) Add(p2 Pos3) Pos3 { return Pos3{p[0] + p2[0], p[1] + p2[1], p[2] + p2[2]} }
func (p Pos3) Sub(p2 Pos3) Pos3 { return Pos3{p[0] - p2[0], p[1] - p2[1], p[2] - p2[2]} }
func (p Pos3) Dist(p2 Pos3) int { return Abs(p[0]-p2[0]) + Abs(p[1]-p2[1]) + Abs(p[2]-p2[2]) }

// Rot3 models a set of 90-deg rotations. [local axis] -> absolute axis, or -(absolute axis)-1 if negative dir
// Example: input Pos3{10,20,30}, rot [2, 0, 1], output Pos3{20,30,10}
// Example: input Pos3{10,20,30}, rot [-1, -2, 1], output Pos3{-10,-20,30}
type Rot3 [3]struct {
	Dim int
	Neg bool
}

func AxesToRot(rightAxis int, rightNeg bool, fwdAxis int, fwdNeg bool) (out Rot3) {
	out = Rot3{
		{rightAxis, rightNeg}, // x axis
		{fwdAxis, fwdNeg},     // y axis
		{},                    // z axis (filled in below)
	}

	// Z axis is positive to start, then flipped for each of rightNeg, upNeg, and if the
	// rightAxis is "1 axis cycled downward" from upAxis (e.g., right/up = 0/1, 1/2, 2/0)
	//
	// This trick works due to how 3D cross product breaks down: Z = X cross Y; X and Y
	// have exactly one non-zero dimension each; the non-zero dimension D of Z will be
	// the dimension that both X and Y have as 0; and finally, its value is:
	//
	//   X[D+1]*Y[D+2] - X[D+2]*Y[D+1]  (all dimension calcs are mod 3)
	//
	// This is +1 when an even number of (X non-zero dim is negative),
	// (Y non-zero dim is negative), and (non-zero dim for X is one "higher" than Y),
	//
	// Another trick: != is XOR for bools in Go (whyyyyy can't we just use ^)
	// Final trick: Go's modulus operator is not like C: -1 % 3 == -1, not 2; we work
	//              around this here with a +3
	//
	out[2].Dim = (3 - rightAxis - fwdAxis) % 3                        // "the unused axis"
	out[2].Neg = rightNeg != fwdNeg != ((rightAxis-fwdAxis+3)%3 == 1) // gather negative signs
	return out
}

func (r Rot3) Rot(p Pos3) (out Pos3) {
	for i, v := range p {
		j := r[i].Dim
		if r[i].Neg {
			v = -v
		}
		out[j] = v
	}
	return out
}

var NewRots [24]Rot3

func init() {
	idx := 0
	for rightAxis := 0; rightAxis < 3; rightAxis++ {
		for _, rightNeg := range [...]bool{true, false} {
			for fwdAxis := 0; fwdAxis < 3; fwdAxis++ {
				if rightAxis == fwdAxis {
					continue
				}
				for _, fwdNeg := range [...]bool{true, false} {
					NewRots[idx] = AxesToRot(rightAxis, rightNeg, fwdAxis, fwdNeg)
					idx++
				}
			}
		}
	}
}

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

type Translation struct { // .Rot.Rot(rel).Add(.Off) -> abs
	Off Pos3
	Rot Rot3
}

func (t Translation) RelToAbs(p Pos3) Pos3 {
	return t.Rot.Rot(p).Add(t.Off)
}

const threshold = 12

func checkOffset(a, b map[Pos3]bool, trans Translation) bool {
	found := 0
	for bp := range b {
		ap := trans.RelToAbs(bp)
		if a[ap] {
			found++
		}
		if found >= threshold {
			return true
		}
	}
	return false
}

func findOffset(a, b map[Pos3]bool) (Translation, bool) {
	for _, rot := range NewRots { // rot.Mul(bPt) -> aPt
		for ap := range a {
			for bp := range b {
				offset := ap.Sub(rot.Rot(bp)) // bPt.Add(offset) -> aPt
				trans := Translation{Off: offset, Rot: rot}
				if checkOffset(a, b, trans) {
					return trans, true
				}
			}
		}
	}
	return Translation{}, false
}

func solveScanners(in []map[Pos3]bool) (abs []map[Pos3]bool, absTrans []Translation) {
	abs = make([]map[Pos3]bool, len(in))
	absTrans = make([]Translation, len(in))
	abs[0], absTrans[0] = in[0], Translation{}

	tried := map[[2]int]bool{}

	numKnown := 1
	for numKnown < len(in) {
		for relIdx, relScanner := range in {
			if abs[relIdx] != nil {
				continue
			}
			for absIdx, absScanner := range abs {
				if absScanner == nil {
					continue
				}

				// Don't try a pair more than once
				if try := [2]int{relIdx, absIdx}; tried[try] {
					continue
				} else {
					tried[try] = true
				}

				if trans, ok := findOffset(absScanner, relScanner); ok {
					fmt.Println("scanner", relIdx, "overlaps absolute scanner", absIdx, "with translation", trans)
					numKnown++
					absTrans[relIdx] = trans
					abs[relIdx] = make(map[Pos3]bool, len(relScanner))
					for bp := range relScanner {
						ap := trans.RelToAbs(bp)
						abs[relIdx][ap] = true
					}
					break
				}
			}
		}
	}

	return abs, absTrans
}

func A(in io.Reader) {
	scanners := parse(in)
	abs, _ := solveScanners(scanners)

	total := map[Pos3]bool{}
	for _, scanner := range abs {
		for p := range scanner {
			total[p] = true
		}
	}

	fmt.Println(len(total))
}

func B(in io.Reader) {
	scanners := parse(in)
	_, absTrans := solveScanners(scanners)

	maxDist := 0
	for _, a := range absTrans {
		for _, b := range absTrans {
			if dist := a.Off.Dist(b.Off); maxDist < dist {
				maxDist = dist
			}
		}
	}

	fmt.Println(maxDist)
}
