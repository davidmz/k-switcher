package win32

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func MAKELCID(LangId, SortId uint16) uint32 {
	return uint32(SortId)<<16 | uint32(LangId)
}

func ActivateKeyboardLayout(hkl HKL, flags uint32) HKL {
	r, _, _ := winAPI("user32.dll", "ActivateKeyboardLayout", uintptr(hkl), uintptr(flags))
	return HKL(r)
}

func GetKeyboardLayoutName() string {
	buf := make([]byte, 2*KL_NAMELENGTH)
	winAPI("user32.dll", "GetKeyboardLayoutNameW", uintptr(unsafe.Pointer(&buf[0])))
	return Utf16zToString(buf)
}

func GetKeyboardLayoutList() []HKL {
	r, _, _ := winAPI("user32.dll", "GetKeyboardLayoutList", 0, 0)
	list := make([]HKL, int(r))
	winAPI("user32.dll", "GetKeyboardLayoutList", r, uintptr(unsafe.Pointer(&list[0])))
	return list
}

func GetLocaleInfo(Locale, LCType int64) string {
	r, _, _ := winAPI("kernel32.dll", "GetLocaleInfoW", uintptr(Locale), uintptr(LCType), 0, 0)
	buf := make([]byte, int(r)*2)
	winAPI("kernel32.dll", "GetLocaleInfoW", uintptr(Locale), uintptr(LCType), uintptr(unsafe.Pointer(&buf[0])), r)
	return Utf16zToString(buf)
}

func OpenClipboard()  { winAPI("user32.dll", "OpenClipboard", 0) }
func CloseClipboard() { winAPI("user32.dll", "CloseClipboard") }

func GetClipboardData() HANDLE {
	r, _, _ := winAPI("user32.dll", "GetClipboardData", CF_UNICODETEXT)
	return HANDLE(r)
}

func SetClipboardData(h HANDLE) {
	winAPI("user32.dll", "SetClipboardData", CF_UNICODETEXT, uintptr(h))
}

func EmptyClipboard() {
	winAPI("user32.dll", "EmptyClipboard")
}

func GlobalSize(h HGLOBAL) uint64 {
	r, _, _ := winAPI("kernel32.dll", "GlobalSize", uintptr(h))
	return uint64(r)
}

func GlobalLock(h HGLOBAL) unsafe.Pointer {
	r, _, _ := winAPI("kernel32.dll", "GlobalLock", uintptr(h))
	return unsafe.Pointer(r)
}

func GlobalUnlock(h HGLOBAL) {
	winAPI("kernel32.dll", "GlobalUnlock", uintptr(h))
}

func GlobalAlloc(uFlags UINT, dwBytes SIZE_T) HGLOBAL {
	r, _, _ := winAPI("kernel32.dll", "GlobalAlloc", uintptr(uFlags), uintptr(dwBytes))
	return HGLOBAL(r)
}

func RegisterHotKey(h HWND, id, mod, vk int) bool {
	r, _, _ := winAPI("user32.dll", "RegisterHotKey", uintptr(h), uintptr(id), uintptr(mod), uintptr(vk))
	return r != 0
}

func GetMessage(msg *MSG, hWnd HWND, wMsgFilterMin, wMsgFilterMax int) bool {
	r, _, _ := winAPI("user32.dll", "GetMessageW", uintptr(unsafe.Pointer(msg)), uintptr(hWnd), uintptr(wMsgFilterMin), uintptr(wMsgFilterMax))
	return r != 0
}

func GetClipboardSequenceNumber() DWORD {
	r, _, _ := winAPI("user32.dll", "GetClipboardSequenceNumber")
	return DWORD(r)
}

func SendMessage(hWnd HWND, Msg UINT, wParam WPARAM, lParam LPARAM) LRESULT {
	r, _, _ := winAPI("user32.dll", "SendMessageW", uintptr(hWnd), uintptr(Msg), uintptr(wParam), uintptr(lParam))
	return LRESULT(r)
}

func GetActiveWindow() HWND {
	r, _, _ := winAPI("user32.dll", "GetActiveWindow")
	return HWND(r)
}

func GetForegroundWindow() HWND {
	r, _, _ := winAPI("user32.dll", "GetForegroundWindow")
	return HWND(r)
}

func PostMessage(hWnd HWND, Msg UINT, wParam WPARAM, lParam LPARAM) bool {
	r, _, _ := winAPI("user32.dll", "PostMessageW", uintptr(hWnd), uintptr(Msg), uintptr(wParam), uintptr(lParam))
	return r != 0
}

func GetGUIThreadInfo(idThread DWORD) *LPGUITHREADINFO {
	lpgui := new(LPGUITHREADINFO)
	lpgui.CbSize = DWORD(unsafe.Sizeof(*lpgui))
	winAPI("user32.dll", "GetGUIThreadInfo", uintptr(idThread), uintptr(unsafe.Pointer(lpgui)))
	return lpgui
}

func GetAsyncKeyState(vKey int) bool {
	r, _, _ := winAPI("user32.dll", "GetAsyncKeyState", uintptr(vKey))
	return SHORT(r) < 0
}

func SendInput(inputs []INPUT) UINT {
	var xInputs []byte
	for _, inp := range inputs {
		xInputs = append(xInputs, inp.toWinStruct()...)
	}
	r, _, _ := winAPI("user32.dll", "SendInput",
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&xInputs[0])),
		uintptr(inputStructSize),
	)
	return UINT(r)
}

func SendInputs(inputs ...INPUT) UINT { return SendInput(inputs) }

func GetLastError() DWORD {
	r, _, _ := winAPI("kernel32.dll", "GetLastError")
	return DWORD(r)
}

///////////////////////

func Utf16zToString(in []byte) string {
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

func StringToUtf16z(s string) (out []byte) {
	us := utf16.Encode([]rune(s))
	out = make([]byte, 2*(len(us)+1))
	for i, u := range us {
		out[2*i] = byte(u & 0xff)
		out[2*i+1] = byte((u >> 8) & 0xff)
	}
	return
}

func winAPI(dllName, funcName string, args ...uintptr) (uintptr, uintptr, error) {
	proc := syscall.MustLoadDLL(dllName).MustFindProc(funcName)
	return proc.Call(args...)
}
