package p17

import (
	"fmt"
	"io"

	. "aoc2021/helpers"
)

const (
	above = 0
	in    = 1
	below = 2
)

func getState(y, minY, maxY int) (state int) {
	if y > maxY {
		return above
	} else if y >= minY {
		return in
	} else {
		return below
	}
}

func simulateY(yvel, minY, maxY int) (peak int, hit bool) {
	y := 0
	for {
		switch getState(y, minY, maxY) {
		case above:
		case in:
			return peak, true
		case below:
			return peak, false
		}

		y = y + yvel
		yvel--

		if peak < y {
			peak = y
		}
	}
}

func simulateX(yvel, xvel, minX, maxX, minY, maxY int) (hit bool) {
	x, y := 0, 0
	for {
		switch getState(y, minY, maxY) {
		case above:
		case below:
			return false
		case in:
			switch getState(x, minX, maxX) {
			case below:
			case in:
				return true
			case above:
				return false
			}
		}

		y += yvel
		x += xvel
		yvel--
		if xvel > 0 {
			xvel--
		} else if xvel < 0 {
			xvel++
		}
	}
}

func A(in io.Reader) {
	minX, maxX, minY, maxY := 0, 0, 0, 0

	l := ReadLines(in)[0]
	fmt.Println(l)
	fmt.Sscanf(l, "target area: x=%d..%d, y=%d..%d", &minX, &maxX, &minY, &maxY)

	bestXVel, bestYVel, bestPeak := 0, 0, 0
	for yvel := 0; yvel <= Abs(minY); yvel++ {
		peak, hit := simulateY(yvel, minY, maxY)
		fmt.Println(yvel, hit, peak)
		if !hit || peak < bestPeak {
			continue
		}

		for xvel := 0; xvel < Abs(maxX); xvel++ {
			if simulateX(yvel, xvel, minX, maxX, minY, maxY) {
				bestXVel, bestYVel, bestPeak = xvel, yvel, peak
				break
			}
		}
	}

	fmt.Println(bestXVel, bestYVel, bestPeak)
}

func B(in io.Reader) {
	minX, maxX, minY, maxY := 0, 0, 0, 0

	l := ReadLines(in)[0]
	fmt.Println(l)
	fmt.Sscanf(l, "target area: x=%d..%d, y=%d..%d", &minX, &maxX, &minY, &maxY)

	combos := 0
	for yvel := -Abs(minY); yvel <= Abs(minY); yvel++ {
		_, hit := simulateY(yvel, minY, maxY)
		if !hit {
			continue
		}

		for xvel := 0; xvel <= Abs(maxX); xvel++ {
			if simulateX(yvel, xvel, minX, maxX, minY, maxY) {
				//fmt.Printf("%d,%d\n", xvel, yvel)
				combos++
			}
		}
	}

	fmt.Println(combos)
}
