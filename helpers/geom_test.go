package helpers

import (
	"reflect"
	"testing"
)

func TestInfGrid_ToFixedGrid(t *testing.T) {
	g := InfGrid{
		{-1, -1}: 'A',
		{-2, 1}:  'B',
		{1, -2}:  'C',
	}

	g2 := g.ToFixedGrid(' ')
	expect := FixedGrid{
		[]rune("   B"),
		[]rune(" A  "),
		[]rune("    "),
		[]rune("C   "),
	}
	if !reflect.DeepEqual(g2, expect) {
		t.Fatalf("grids do not match:\nexpect:%s\nobtained:\n%s", expect, g2)
	}
}

func TestFixedGrid_Copy(t *testing.T) {
	g := FixedGrid{
		{'A', 'B', 'C'},
		{'D', 'E', 'F'},
	}

	g2 := g.Copy()
	if !reflect.DeepEqual(g, g2) {
		t.Fatalf("grids do not match:\nexpect:%s\nobtained:\n%s", g, g2)
	}
}

func TestInfGrid_Copy(t *testing.T) {
	g := InfGrid{
		{0, 1}: 'A',
		{2, 3}: 'B',
		{4, 5}: 'C',
	}

	g2 := g.Copy()
	if !reflect.DeepEqual(g, g2) {
		t.Fatalf("grids do not match:\nexpect:%s\nobtained:\n%s", g, g2)
	}
}
