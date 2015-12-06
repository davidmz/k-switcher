package clipboard

import (
	"unicode/utf16"
	"unsafe"

	"github.com/davidmz/k-switcher.v2/win32"
)

func Get() string {
	win32.OpenClipboard()
	defer win32.CloseClipboard()

	h := win32.HGLOBAL(win32.GetClipboardData())
	if h == 0 {
		return ""
	}
	sz := win32.GlobalSize(h)
	if sz == 0 {
		return ""
	}
	p := win32.GlobalLock(h)
	defer win32.GlobalUnlock(h)
	return utf16zToString(readMemory(p, sz))
}

func Put(text string) {
	win32.OpenClipboard()
	defer win32.CloseClipboard()

	data := stringToUtf16z(text)
	h := win32.GlobalAlloc(win32.GHND, win32.SIZE_T(len(data)))
	p := win32.GlobalLock(h)
	writeMemory(p, data)
	win32.GlobalUnlock(h)
	win32.SetClipboardData(win32.HANDLE(h))
}

func GetSeqNumber() uint32 {
	return uint32(win32.GetClipboardSequenceNumber())
}

func readMemory(ptr unsafe.Pointer, size uint64) (result []byte) {
	const bufSize = 512
	var (
		start uint64 = uint64(uintptr(ptr))
		p     uint64
	)
	for p = start; p < start+size; p += bufSize {
		sz := start + size - p
		if sz > bufSize {
			sz = bufSize
		}
		src := (*[bufSize]byte)(unsafe.Pointer(uintptr(p)))[:sz]
		result = append(result, src...)
	}
	return result[:size:size]
}

func writeMemory(ptr unsafe.Pointer, data []byte) {
	const bufSize = 512
	var (
		start uint64 = uint64(uintptr(ptr))
		p     uint64
	)
	for p = start; p < start+uint64(len(data)); p += bufSize {
		src := data[p-start:]
		dst := (*[bufSize]byte)(unsafe.Pointer(uintptr(p)))[:]
		copy(dst, src)
	}
}

func utf16zToString(in []byte) string {
	out := make([]uint16, 0, len(in)/2)
	x := uint16(0)
	for i, b := range in {
		if i%2 == 0 {
			x = uint16(b)
		} else {
			x += uint16(b) << 8
			if x == 0 {
				break
			}
			out = append(out, x)
		}
	}
	return string(utf16.Decode(out))
}

func stringToUtf16z(s string) (out []byte) {
	us := utf16.Encode([]rune(s))
	out = make([]byte, 2*(len(us)+1))
	for i, u := range us {
		out[2*i] = byte(u & 0xff)
		out[2*i+1] = byte((u >> 8) & 0xff)
	}
	return
}
