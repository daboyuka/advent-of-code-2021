package main

import "testing"

func TestParseArg(t *testing.T) {
	ck := func(arg string, expProb, expPart int) {
		prob, part := parseArg(arg)
		if prob != expProb {
			t.Fatalf("%s: expected problem %d, got %d", arg, expProb, prob)
		} else if part != expPart {
			t.Fatalf("%s: expected part %d, got %d", arg, expPart, part)
		}
	}

	ck("1a", 0, 0)
	ck("1b", 0, 1)
	ck("1B", 0, 1)
	ck("2a", 1, 0)
	ck("2b", 1, 1)
	ck("10b", 9, 1)
	ck("25b", 24, 1)
}
