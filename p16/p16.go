package p16

import (
	"fmt"
	"io"
	"math"

	. "aoc2021/helpers"
)

type packet struct {
	Version int
	TypeID  int

	Literal uint
	Sub     []packet
}

func (p packet) Eval() (v int) {
	switch p.TypeID {
	case 0:
		for _, p := range p.Sub {
			v += p.Eval()
		}
		return v
	case 1:
		v = 1
		for _, p := range p.Sub {
			v *= p.Eval()
		}
		return v
	case 2:
		v = math.MaxInt
		for _, p := range p.Sub {
			v = Min(v, p.Eval())
		}
		return v
	case 3:
		v = math.MinInt
		for _, p := range p.Sub {
			v = Max(v, p.Eval())
		}
		return v
	case 4:
		return int(p.Literal)
	case 5:
		if p.Sub[0].Eval() > p.Sub[1].Eval() {
			return 1
		}
		return 0
	case 6:
		if p.Sub[0].Eval() < p.Sub[1].Eval() {
			return 1
		}
		return 0
	case 7:
		if p.Sub[0].Eval() == p.Sub[1].Eval() {
			return 1
		}
		return 0
	}
	panic(p)
}

func parsePacketsByBits(b *bitstream, nbits int) (ps []packet) {
	start := b.Pos()
	fmt.Println("parsing", nbits, "bits")
	for b.Pos() < start+nbits {
		fmt.Println("parsing at", b.Pos())
		ps = append(ps, parsePacket(b))
	}
	return ps
}

func parsePacketsByCount(b *bitstream, npack int) (ps []packet) {
	fmt.Println("parsing", npack, "packets")
	for i := 0; i < npack; i++ {
		ps = append(ps, parsePacket(b))
	}
	return ps
}

func parsePacket(b *bitstream) (p packet) {
	p.Version = int(b.Next(3))
	p.TypeID = int(b.Next(3))
	fmt.Println("packet", p.Version, p.TypeID)
	switch p.TypeID {
	case 4: // literal
		for {
			p.Literal <<= 4
			word := b.Next(5)
			p.Literal |= word & 0xF
			if word&0x10 == 0 {
				break
			}
		}
	default:
		var ps []packet
		if mode := b.Next(1); mode == 0 {
			nbits := b.Next(15)
			ps = parsePacketsByBits(b, int(nbits))
		} else {
			npack := b.Next(11)
			ps = parsePacketsByCount(b, int(npack))
		}

		p.Sub = ps
	}

	return p
}

type bitstream struct {
	S    string
	bIdx int
	pos  int
}

func (b *bitstream) charToHex(x byte) byte {
	if x >= '0' && x <= '9' {
		return x - '0'
	} else if x >= 'A' && x <= 'F' {
		return x - 'A' + 10
	} else {
		panic(x)
	}
}

func (b *bitstream) Pos() int { return b.pos }

func (b *bitstream) EOF() bool { return len(b.S) == 0 }

func (b *bitstream) Next(n int) (out uint) { // first bit is placed at 2^(n-1) bit, down to 2^0
	for i := 0; i < n; i++ {
		out <<= 1
		v := b.charToHex(b.S[0])
		bit := uint((v >> (3 - b.bIdx)) & 1)
		out |= bit
		fmt.Println("bit", v, b.bIdx, bit)
		b.bIdx++
		b.pos++

		if b.bIdx == 4 {
			b.bIdx, b.S = 0, b.S[1:]
		}
	}

	fmt.Printf("Next %d: %0x\n", n, out)
	return out
}

func sumVersions(p packet) (sum int) {
	sum = p.Version
	for _, p2 := range p.Sub {
		sum += sumVersions(p2)
	}
	return sum
}

func A(in io.Reader) {
	line := ReadLines(in)[0]
	b := &bitstream{S: line}

	p := parsePacket(b)
	fmt.Println(sumVersions(p))
}

func B(in io.Reader) {
	line := ReadLines(in)[0]
	b := &bitstream{S: line}

	p := parsePacket(b)
	fmt.Println(p.Eval())
}
