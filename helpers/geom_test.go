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
