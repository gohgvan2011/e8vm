package arch8

import (
	"testing"
)

func TestPage(t *testing.T) {
	e := func(cond bool, s string, args ...interface{}) {
		if cond {
			t.Fatalf(s, args...)
		}
	}

	e(PageSize < 4, "page size too small")
	p := newPage()
	e(((PageSize-1)&PageSize) != 0, "page size not exponential of 2")

	for i := uint32(0); i < PageSize; i++ {
		b := p.ReadByte(i)
		e(b != 0, "byte %d not zero on new page", i)
	}

	for i := uint32(0); i < PageSize/4; i++ {
		b := p.ReadWord(i * 4)
		e(b != 0, "word %d not zero on new page", i)
	}

	off := uint32(56)
	p.WriteByte(off+0, 0x37)
	p.WriteByte(off+1, 0x21)
	p.WriteByte(off+2, 0x5a)
	p.WriteByte(off+3, 0x70)

	exp := uint32(0x705a2137)
	w := p.ReadWord(off)
	e(w != exp, "expect 0x%08x got 0x%08x", exp, w)
	w = p.ReadWord(off + 3)
	e(w != exp, "expect 0x%08x got 0x%08x", exp, w)
	w = p.ReadWord(off + 3 + 2*PageSize)
	e(w != exp, "expect 0x%08x got 0x%08x", exp, w)

	b := p.ReadByte(off + 2)
	e(b != 0x5a, "got incorrect byte 0x%02x", b)
	b = p.ReadByte(off + 2 + 2*PageSize)
	e(b != 0x5a, "got incorrect byte 0x%02x", b)

	p.WriteWord(off+3+2*PageSize, exp)
	b = p.ReadByte(off + 2)
	e(b != 0x5a, "got incorrect byte 0x%02x", b)
}
