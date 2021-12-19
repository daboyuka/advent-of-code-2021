package p18

import (
	"fmt"
	"io"
	"strconv"

	. "aoc2021/helpers"
)

type Value struct {
	V int
	P *Pair
}

type Pair struct {
	L, R Value
}

func (s *Value) String() string {
	if s.IsNum() {
		return strconv.Itoa(s.V)
	} else {
		return s.P.String()
	}
}

func (p *Pair) String() string {
	return fmt.Sprintf("[%s,%s]", p.L.String(), p.R.String())
}

func (s *Value) IsNum() bool { return s.P == nil }

func (s *Value) Split() bool {
	if s.V >= 10 {
		*s = Value{P: &Pair{Value{V: s.V / 2}, Value{V: (s.V + 1) / 2}}}
		return true
	} else if s.P != nil {
		return s.P.Split()
	} else {
		return false
	}
}

func (s *Value) AddRegular(v int, left bool) {
	if v == 0 {
		return
	} else if s.IsNum() {
		s.V += v
	} else {
		s.P.AddRegular(v, left)
	}
}

func (p *Pair) AddRegular(v int, left bool) {
	if left {
		p.L.AddRegular(v, left)
	} else {
		p.R.AddRegular(v, left)
	}
}

func (p *Pair) Split() bool {
	return p.L.Split() || p.R.Split()
}

func (p *Pair) Explode(depth int) (l, r int, exploded, immedexploded bool) {
	if p.L.IsNum() && p.R.IsNum() && depth >= 4 {
		return p.L.V, p.R.V, true, true
	}

	if !p.L.IsNum() {
		l, r, exploded, immedexploded = p.L.P.Explode(depth + 1)
		if immedexploded {
			p.L = Value{V: 0}
		}
		if exploded {
			p.R.AddRegular(r, true)
			return l, 0, true, false
		}
	}
	if !p.R.IsNum() {
		l, r, exploded, immedexploded = p.R.P.Explode(depth + 1)
		if immedexploded {
			p.R = Value{V: 0}
		}
		if exploded {
			p.L.AddRegular(l, false)
			return 0, r, true, false
		}
	}

	return 0, 0, false, false
}

func (v *Value) Apply(depth int) {
	for {
		for {
			if !v.IsNum() {
				if _, _, ok, _ := v.P.Explode(depth); !ok {
					break
				}
			}
			fmt.Println(v.String())
		}
		if !v.Split() {
			break
		}
		fmt.Println(v.String())
	}
}

func (v Value) Magnitude() int {
	if v.IsNum() {
		return v.V
	} else {
		return 3*v.P.L.Magnitude() + 2*v.P.R.Magnitude()
	}
}

func (v Value) Add(other Value) Value {
	return Value{P: &Pair{L: v, R: other}}
}

func (p *Pair) Copy() *Pair {
	return &Pair{L: p.L.Copy(), R: p.R.Copy()}
}

func (v Value) Copy() Value {
	if v.IsNum() {
		return v
	} else {
		return Value{P: v.P.Copy()}
	}
}

func parse(s string) (v Value, rem string) {
	switch s[0] {
	case '[':
		var l, r Value
		l, s = parse(s[1:])
		r, s = parse(s[1:])
		s = s[1:]
		return Value{P: &Pair{L: l, R: r}}, s
	default:
		return Value{V: Atoi(s[:1])}, s[1:]
	}
}

func A(in io.Reader) {
	var v Value
	for _, l := range ReadLines(in) {
		v2, _ := parse(l)
		if v == (Value{}) {
			v = v2
			fmt.Println(v.String())
		} else {
			fmt.Println("+", v2.String())
			v = v.Add(v2)
		}

		v.Apply(0)
		fmt.Println("=", v.String())
	}

	fmt.Println(v.String())
	fmt.Println(v.Magnitude())
}

func B(in io.Reader) {
	var vs []Value
	for _, l := range ReadLines(in) {
		v, _ := parse(l)
		vs = append(vs, v)
	}

	maxMag := 0
	for i, v1 := range vs {
		for j, v2 := range vs {
			if i == j {
				continue
			}

			vOut := v1.Copy().Add(v2.Copy())
			vOut.Apply(0)
			if m := vOut.Magnitude(); maxMag < m {
				maxMag = m
			}
		}
	}

	fmt.Println(maxMag)
}
