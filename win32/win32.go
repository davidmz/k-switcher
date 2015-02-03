package win32

import (
	"syscall"
	"unicode/utf16"
	"unsafe"

	"github.com/AllenDang/w32"
)

const (
	MOD_ALT      = 0x0001
	MOD_CONTROL  = 0x0002
	MOD_SHIFT    = 0x0004
	MOD_WIN      = 0x0008
	MOD_NOREPEAT = 0x4000

	KEYEVENTF_EXTENDEDKEY = 0x0001
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_UNICODE     = 0x0004
	KEYEVENTF_SCANCODE    = 0x0008

	KL_NAMELENGTH = 9

	CTRL_C_EVENT        = 0
	CTRL_BREAK_EVENT    = 1
	CTRL_CLOSE_EVENT    = 2
	CTRL_LOGOFF_EVENT   = 5
	CTRL_SHUTDOWN_EVENT = 6
)

type HandlerRoutine func(int) int

func RegisterHotKey(h w32.HWND, id, mod, vk int) bool {
	r, _, _ := __RegisterHotKey(uintptr(h), uintptr(id), uintptr(mod), uintptr(vk))
	return r != 0
}

func GlobalSize(h w32.HGLOBAL) uint64 {
	r, _, _ := __GlobalSize(uintptr(h))
	return uint64(r)
}

func GetClipboardSequenceNumber() uint32 {
	r, _, _ := __GetClipboardSequenceNumber()
	return uint32(r)
}

func GetClipboardData(f uint) w32.HANDLE {
	r, _, _ := __GetClipboardData(uintptr(f))
	return w32.HANDLE(r)
}

func GetKeyboardLayoutList() []w32.HKL {
	r, _, _ := __GetKeyboardLayoutList(0, 0)
	list := make([]w32.HKL, int(r))
	__GetKeyboardLayoutList(r, uintptr(unsafe.Pointer(&list[0])))
	return list
}

func GetKeyboardLayout(id uint32) w32.HKL {
	r, _, _ := __GetKeyboardLayout(uintptr(id))
	return w32.HKL(r)
}

func GetKeyboardLayoutName() string {
	buf := make([]byte, 2*KL_NAMELENGTH)
	__GetKeyboardLayoutNameW(uintptr(unsafe.Pointer(&buf[0])))
	return Utf16zToString(buf)
}

func ActivateKeyboardLayout(hkl w32.HKL, flags uint32) w32.HKL {
	r, _, _ := __ActivateKeyboardLayout(uintptr(hkl), uintptr(flags))
	return w32.HKL(r)
}

func LoadKeyboardLayout(id string, flags uint32) w32.HKL {
	b := StringToUtf16z(id)
	r, _, _ := __LoadKeyboardLayoutW(uintptr(unsafe.Pointer(&b[0])), uintptr(flags))
	return w32.HKL(r)
}

func UnloadKeyboardLayout(hkl w32.HKL) bool {
	r, _, _ := __UnloadKeyboardLayout(uintptr(hkl))
	return r != 0
}

func GetForegroundWindow() w32.HWND {
	r, _, _ := __GetForegroundWindow()
	return w32.HWND(r)
}

func SetForegroundWindow(h w32.HWND) bool {
	r, _, _ := __SetForegroundWindow(uintptr(h))
	return r != 0
}

func SetConsoleCtrlHandler(h HandlerRoutine, add bool) bool {
	a := 0
	if add {
		a = 1
	}
	r, _, _ := __SetConsoleCtrlHandler(syscall.NewCallback(h), uintptr(a))
	return r != 0
}

type tagPoint struct{ x, y int32 }

type tagRect struct{ left, top, right, bottom int32 }

type PlacementInfo struct {
	Length           uint32
	Flags            uint32
	ShowCmd          uint32
	PtMinPosition    tagPoint
	PtMaxPosition    tagPoint
	RcNormalPosition tagRect
}

func GetWindowPlacement(h w32.HWND) *PlacementInfo {
	pl := &PlacementInfo{}
	pl.Length = uint32(unsafe.Sizeof(*pl))
	r, _, _ := __GetWindowPlacement(uintptr(h), uintptr(unsafe.Pointer(pl)))
	if r != 0 {
		return pl
	}
	return nil
}

func SetConsoleTitle(title string) bool {
	p, _ := syscall.UTF16PtrFromString(title)
	r, _, _ := __SetConsoleTitleW(uintptr(unsafe.Pointer(p)))
	return r != 0
}

func RegisterWindowMessage(msg string) uint32 {
	p, _ := syscall.UTF16PtrFromString(msg)
	r, _, _ := __RegisterWindowMessage(uintptr(unsafe.Pointer(p)))
	return uint32(r)
}

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

////////////////

var (
	__RegisterHotKey             = winAPI("user32.dll", "RegisterHotKey")
	__GetClipboardSequenceNumber = winAPI("user32.dll", "GetClipboardSequenceNumber")
	__GetClipboardData           = winAPI("user32.dll", "GetClipboardData")
	__GetKeyboardLayoutList      = winAPI("user32.dll", "GetKeyboardLayoutList")
	__GetKeyboardLayout          = winAPI("user32.dll", "GetKeyboardLayout")
	__GetKeyboardLayoutNameW     = winAPI("user32.dll", "GetKeyboardLayoutNameW")
	__ActivateKeyboardLayout     = winAPI("user32.dll", "ActivateKeyboardLayout")
	__LoadKeyboardLayoutW        = winAPI("user32.dll", "LoadKeyboardLayoutW")
	__UnloadKeyboardLayout       = winAPI("user32.dll", "UnloadKeyboardLayout")
	__GetForegroundWindow        = winAPI("user32.dll", "GetForegroundWindow")
	__SetForegroundWindow        = winAPI("user32.dll", "SetForegroundWindow")
	__GetWindowPlacement         = winAPI("user32.dll", "GetWindowPlacement")
	__RegisterWindowMessage      = winAPI("user32.dll", "RegisterWindowMessageW")

	__GlobalSize            = winAPI("kernel32.dll", "GlobalSize")
	__SetConsoleCtrlHandler = winAPI("kernel32.dll", "SetConsoleCtrlHandler")
	__SetConsoleTitleW      = winAPI("kernel32.dll", "SetConsoleTitleW")
)

func winAPI(dllName, funcName string) func(...uintptr) (uintptr, uintptr, error) {
	proc := syscall.MustLoadDLL(dllName).MustFindProc(funcName)
	return func(a ...uintptr) (uintptr, uintptr, error) { return proc.Call(a...) }
}
