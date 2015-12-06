package win32

import "unsafe"

type (
	ATOM            uint16
	BOOL            int32
	COLORREF        uint32
	DWM_FRAME_COUNT uint64
	DWORD           uint32
	HACCEL          HANDLE
	HANDLE          uintptr
	HBITMAP         HANDLE
	HBRUSH          HANDLE
	HCURSOR         HANDLE
	HDC             HANDLE
	HDROP           HANDLE
	HDWP            HANDLE
	HENHMETAFILE    HANDLE
	HFONT           HANDLE
	HGDIOBJ         HANDLE
	HGLOBAL         HANDLE
	HGLRC           HANDLE
	HHOOK           HANDLE
	HICON           HANDLE
	HIMAGELIST      HANDLE
	HINSTANCE       HANDLE
	HKEY            HANDLE
	HKL             HANDLE
	HMENU           HANDLE
	HMODULE         HANDLE
	HMONITOR        HANDLE
	HPEN            HANDLE
	HRESULT         int32
	HRGN            HANDLE
	HRSRC           HANDLE
	HTHUMBNAIL      HANDLE
	HWND            HANDLE
	LPARAM          uintptr
	LPCVOID         unsafe.Pointer
	LRESULT         uintptr
	PVOID           unsafe.Pointer
	QPC_TIME        uint64
	ULONG_PTR       uintptr
	WPARAM          uintptr
	UINT            uint32
	LONG            int32
	SHORT           int16
	SIZE_T          ULONG_PTR
)

type MSG struct {
	Hwnd    HWND
	Message UINT
	WParam  WPARAM
	LParam  LPARAM
	Time    DWORD
	Pt      POINT
}

type POINT struct {
	X, Y LONG
}

type LPGUITHREADINFO struct {
	CbSize        DWORD
	Flags         DWORD
	HwndActive    HWND
	HwndFocus     HWND
	HwndCapture   HWND
	HwndMenuOwner HWND
	HwndMoveSize  HWND
	HwndCaret     HWND
	RcCaret       RECT
}

type RECT struct {
	Left   LONG
	Top    LONG
	Right  LONG
	Bottom LONG
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms646270(v=vs.85).aspx
type INPUT struct {
	Type DWORD
	Ki   KEYBDINPUT
	Mi   MOUSEINPUT
	Hi   HARDWAREINPUT
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms646273(v=vs.85).aspx
type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms646271(v=vs.85).aspx
type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/ms646269(v=vs.85).aspx
type HARDWAREINPUT struct {
	UMsg    uint32
	WParamL uint16
	WParamH uint16
}

type _INPUT_KEYBDINPUT struct {
	Type DWORD
	Ki   KEYBDINPUT
}

type _INPUT_MOUSEINPUT struct {
	Type DWORD
	Mi   MOUSEINPUT
}

type _INPUT_HARDWAREINPUT struct {
	Type DWORD
	Hi   HARDWAREINPUT
}

var inputStructSize = 0

func (i *INPUT) toWinStruct() []byte {
	b := make([]byte, inputStructSize)
	pb := unsafe.Pointer(&b[0])
	switch i.Type {
	case INPUT_KEYBOARD:
		x := (*_INPUT_KEYBDINPUT)(pb)
		x.Type = i.Type
		x.Ki = i.Ki
	case INPUT_MOUSE:
		x := (*_INPUT_MOUSEINPUT)(pb)
		x.Type = i.Type
		x.Mi = i.Mi
	case INPUT_HARDWARE:
		x := (*_INPUT_HARDWAREINPUT)(pb)
		x.Type = i.Type
		x.Hi = i.Hi
	}
	return b
}

func init() {
	var a int

	a = int(unsafe.Sizeof(_INPUT_KEYBDINPUT{}))
	if inputStructSize < a {
		inputStructSize = a
	}
	a = int(unsafe.Sizeof(_INPUT_MOUSEINPUT{}))
	if inputStructSize < a {
		inputStructSize = a
	}
	a = int(unsafe.Sizeof(_INPUT_HARDWAREINPUT{}))
	if inputStructSize < a {
		inputStructSize = a
	}
}
