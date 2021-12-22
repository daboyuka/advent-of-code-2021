package p22

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

type Pos3 [3]int

type Cube struct {
	Min, Max Pos3
	On       bool
}

func parse(in io.Reader) (steps []Cube) {
	for _, l := range ReadLines(in) {
		step := Cube{}
		state, bounds := Split2(l, " ")
		step.On = state == "on"

		xb, yb, zb := Split3(bounds, ",", ",")
		x1, x2 := Split2(xb[2:], "..")
		y1, y2 := Split2(yb[2:], "..")
		z1, z2 := Split2(zb[2:], "..")

		step.Min = Pos3{Atoi(x1), Atoi(y1), Atoi(z1)}
		step.Max = Pos3{Atoi(x2), Atoi(y2), Atoi(z2)}
		steps = append(steps, step)
	}
	return steps
}

type InfGrid3 map[Pos3]bool

func doStep(step Cube, g InfGrid3, min, max Pos3) {
	for x := Max(step.Min[0], min[0]); x <= Min(step.Max[0], max[0]); x++ {
		for y := Max(step.Min[1], min[1]); y <= Min(step.Max[1], max[1]); y++ {
			for z := Max(step.Min[2], min[2]); z <= Min(step.Max[2], max[2]); z++ {
				g[Pos3{x, y, z}] = step.On
			}
		}
	}
}

func A(in io.Reader) {
	g := make(InfGrid3, 50*50*50)
	steps := parse(in)

	for i, step := range steps {
		doStep(step, g, Pos3{-50, -50, -50}, Pos3{50, 50, 50})
		fmt.Println(i)
	}

	count := 0
	for _, v := range g {
		if v {
			count++
		}
	}
	fmt.Println(count)
}

type Square = Cube
type Line = Cube

type Point struct {
	X     int
	Start bool
	On    bool
}

//func computeLines(lines []Line) (area int) {
//	pts := make([]Point, 0, 2*len(lines))
//	for _, l := range lines {
//		pts = append(pts, Point{l.Min.X, l.On, true})
//		pts = append(pts, Point{l.Max.X, l.On, false})
//	}
//	sort.SliceStable(pts, func(i, j int) bool { return pts[i].X < pts[j].X })
//
//	onC, offC := 0, 0
//	on := false
//	lastPos := 0
//	for i, pt := range pts {
//		if pt.Start {
//			if pt.On {
//				onC++
//				on = true
//			} else {
//				offC++
//				on = false
//			}
//		} else {
//
//		}
//	}
//}

// Split - high keeps at
func (c Cube) Split(dim int, at int) (low, high Cube) {
	if at < c.Min[dim] || at > c.Max[dim] {
		panic(fmt.Errorf("oboe %+v %+v %d %d", c.Min, c.Max, dim, at))
	}

	low, high = c, c
	low.Max[dim] = at
	high.Min[dim] = at
	return low, high
}

func (c Cube) IsDegenerate() bool {
	for i, v := range c.Min {
		if v == c.Max[i] {
			return true
		}
	}
	return false
}

func (c Cube) Intersects(other Cube) bool {
	if c.IsDegenerate() || other.IsDegenerate() {
		return false
	}
	for dim := range c.Min {
		if c.Max[dim] <= other.Min[dim] || c.Min[dim] >= other.Max[dim] {
			return false
		}
	}
	return true
}

func doSplits(a, b Cube, dim int) (arem, brem Cube, splits []Cube) { // return pieces of a - b
	// Intersect on Z
	if a.Min[dim] < b.Min[dim] { // a extends lower; lop off that part as its own cube
		var low Cube
		low, a = a.Split(dim, b.Min[dim])
		splits = append(splits, low)
	} else { // b extends lower; lop off that part of b
		_, b = b.Split(dim, a.Min[dim])
	}
	if a.Max[dim] > b.Max[dim] { // a extends higher; lop off that part as its own cube
		var high Cube
		a, high = a.Split(dim, b.Max[dim])
		splits = append(splits, high)
	} else { // b extends higher; lop off that part of b
		b, _ = b.Split(dim, a.Max[dim])
	}

	return a, b, splits
}

func Subtract(a, b Cube) (trimmed Cube, splits []Cube) { // return pieces of a - b
	if !a.Intersects(b) {
		return a, nil
	}

	fmt.Println("Subtract start", a, b)
	a, b, splitsX := doSplits(a, b, 0)
	fmt.Println("Subtract X", a, b)
	a, b, splitsY := doSplits(a, b, 1)
	fmt.Println("Subtract Y", a, b)
	a, b, splitsZ := doSplits(a, b, 2)
	fmt.Println("Subtract Z", a, b)

	splits = append(splits, splitsX...)
	splits = append(splits, splitsY...)
	splits = append(splits, splitsZ...)

	if b.IsDegenerate() { // b as completed eliminated, so a is valid; otherwise a == b and a is covered
		return a, splits
	} else {
		return Cube{}, splits
	}
}

func B(in io.Reader) {
	steps := parse(in)
	for i, step := range steps { // make upper bound exclusive
		for dim := range step.Max {
			step.Max[dim]++
		}
		steps[i] = step
	}

	var cur []Cube

	var stack []Cube
	for i, step := range steps {
		fmt.Println("step", i, "cur len", len(cur))
		if step.On {
			stack = append(stack[:0], step)
			for len(stack) > 0 {
				l := len(stack)
				rem := stack[l-1]
				stack = stack[:l-1]
				if rem.IsDegenerate() {
					continue
				}

				for _, other := range cur {
					var splits []Cube
					rem, splits = Subtract(rem, other)
					//fmt.Printf("adding %d new splits: %+v\n", len(splits), splits)
					stack = append(stack, splits...)
					if rem.IsDegenerate() {
						break
					}
				}

				if !rem.IsDegenerate() {
					cur = append(cur, rem)
				}
			}
		} else {
			stack = stack[:0]
			for _, other := range cur {
				rem, remSplits := Subtract(other, step)
				stack = append(stack, remSplits...)
				if !rem.IsDegenerate() {
					stack = append(stack, rem)
				}
			}

			cur = append(cur[:0], stack...)
			stack = stack[:0]
		}
	}

	on := 0
	for _, c := range cur {
		if c.IsDegenerate() {
			continue
		}
		size := 1
		for i, min := range c.Min {
			size *= c.Max[i] - min
		}
		on += size
	}
	fmt.Println(len(cur))
	fmt.Println(on)
}
