package flowhash

import (
	"fmt"
	"strings"
	"unsafe"
)

const IndirectTableItemSize = unsafe.Sizeof(uint32(0))

func IndirectTableSize(n uint32) uintptr {
	return IndirectTableItemSize * uintptr(n)
}

type IndirectTable []uint32

func NewIndirectTable(n int) IndirectTable {
	return make(IndirectTable, n)
}

func UnsafeRawIndirectTable(ptr unsafe.Pointer, len int) IndirectTable {
	return (*[1 << 24]uint32)(unsafe.Pointer(ptr))[:len]
}

func (t IndirectTable) Size() uintptr {
	return uintptr(len(t)) * IndirectTableItemSize
}

func (t IndirectTable) Clone() IndirectTable {
	n := make(IndirectTable, len(t))
	copy(n, t)
	return n
}

func (t IndirectTable) String() string {
	var b strings.Builder

	for i, n := range t {
		if i%8 == 0 {
			fmt.Fprintf(&b, "%5d: ", i)
		}

		fmt.Fprintf(&b, " %5d", n)

		if i%8 == 7 || i == len(t)-1 {
			fmt.Fprintln(&b, "")
		}
	}

	return b.String()
}
