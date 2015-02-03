package main

import (
	"unicode/utf16"
	"unsafe"

	"github.com/AllenDang/w32"
	"github.com/davidmz/k-switcher/win32"
)

type cbRecord struct {
	Type uint
	Data []byte
}
type CBRecords []*cbRecord

func backupClipboard() (records CBRecords) {
	w32.OpenClipboard(0)
	defer w32.CloseClipboard()

	// сохраняем только текст, потому что остальные форматы
	// в общем случае сохранить невозможно
	if d := readFromClipboard(w32.CF_UNICODETEXT); d != nil {
		records = append(records, &cbRecord{Type: w32.CF_UNICODETEXT, Data: d})
	}
	return
}

func restoreClipboard(records CBRecords) {
	w32.OpenClipboard(0)
	defer w32.CloseClipboard()

	w32.EmptyClipboard()
	for _, rec := range records {
		sendToClipboard(rec.Type, rec.Data)
	}
}

func readFromClipboard(f uint) []byte {
	h := w32.HGLOBAL(w32.GetClipboardData(f))
	if h == 0 {
		return nil
	}
	sz := win32.GlobalSize(h)
	if sz == 0 {
		return nil
	}
	p := w32.GlobalLock(h)
	defer w32.GlobalUnlock(h)
	return readMemory(p, sz)
}

func sendToClipboard(typ uint, data []byte) {
	if len(data) == 0 {
		return
	}
	h := w32.GlobalAlloc(w32.GHND, uint32(len(data)))
	p := w32.GlobalLock(h)
	writeMemory(p, data)
	w32.GlobalUnlock(h)
	w32.SetClipboardData(typ, w32.HANDLE(h))
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
