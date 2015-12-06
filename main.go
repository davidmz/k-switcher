package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"
	"time"

	"github.com/AllenDang/w32"
	"github.com/davidmz/go-semaphore"
	"github.com/davidmz/k-switcher/consolewin"
	"github.com/davidmz/k-switcher/win32"
	"github.com/lxn/walk"
	"github.com/lxn/win"
)

const (
	MODE_SEL = 0
	MODE_ALL = 1
)

func main() {

	if ok := win32.RegisterHotKey(0, MODE_SEL, win32.MOD_SHIFT|win32.MOD_NOREPEAT, w32.VK_PAUSE); !ok {
		fmt.Println("Не удаётся захватить комбинацию Shift+Break. Возможно, она занята другим приложением")
		return
	}
	if ok := win32.RegisterHotKey(0, MODE_ALL, win32.MOD_SHIFT|win32.MOD_CONTROL|win32.MOD_NOREPEAT, w32.VK_PAUSE); !ok {
		fmt.Println("Не удаётся захватить клавишу Ctrl+Shift+Break. Возможно, она занята другим приложением")
		return
	}

	fmt.Println("K-Switcher запущен")

	time.AfterFunc(3*time.Second, func() {
		if consolewin.ShowState != w32.SW_HIDE {
			consolewin.Hide()
		}
	})

	var (
		trayIcon      *walk.NotifyIcon
		trayIconImage *walk.Icon
	)

	{
		img, _, err := image.Decode(bytes.NewReader(iconData))
		mustBeOk(err)
		trayIconImage, err = walk.NewIconFromImage(img)
		mustBeOk(err)
	}

	createTrayIcon := func() {
		var err error

		trayIcon, err = walk.NewNotifyIcon()
		mustBeOk(err)

		mustBeOk(trayIcon.SetToolTip("K-Switcher"))
		mustBeOk(trayIcon.SetIcon(trayIconImage))
		mustBeOk(trayIcon.SetVisible(true))

		trayIcon.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
			if button == walk.LeftButton {
				if consolewin.ShowState != w32.SW_HIDE {
					consolewin.Hide()
				} else {
					consolewin.Show()
				}
			}
		})
	}
	createTrayIcon()
	HandleExplorerCrash(createTrayIcon)

	consolewin.OnControlEvent(func(t int) bool {
		if t == win32.CTRL_C_EVENT || t == win32.CTRL_BREAK_EVENT || t == win32.CTRL_CLOSE_EVENT {
			mustBeOk(trayIcon.SetVisible(false))
			fmt.Println("Выходим…")
			os.Exit(0)
		}
		return false
	})

	consolewin.OnShowStateChange(func(showState uint32) {
		if showState == w32.SW_SHOWMINIMIZED {
			consolewin.Hide()
		}
	})

	consolewin.SetTitle("K-Switcher")

	myHKL := win32.GetKeyboardLayout(0)
	win32.ActivateKeyboardLayout(myHKL, 0)

	var lock = semaphore.Mutex()
	msg := new(w32.MSG)
	for w32.GetMessage(msg, 0, 0, 0) == 1 {
		if msg.Message == w32.WM_HOTKEY {
			if lock.Try() {
				processHotkey(msg.WParam)
				lock.Release()
			}
		}
	}
	fmt.Println("Bye!")
}

func processHotkey(mode uintptr) {
	bkup := backupClipboard()

	startClipNum := win32.GetClipboardSequenceNumber()
	clipNum := startClipNum

	if mode == MODE_LAST_WORD {
		w32.SendInput(CtrlShiftLeftSeq)
	} else if mode == MODE_SEL {
		w32.SendInput(ShiftUpSeq)
	}
	w32.SendInput(CtrlCSeq)

	timeout := time.NewTimer(500 * time.Millisecond)
	tick := time.NewTicker(20 * time.Millisecond)
	defer timeout.Stop()
	defer tick.Stop()

loop:
	for {
		select {
		case <-tick.C:
			if clipNum = win32.GetClipboardSequenceNumber(); clipNum > startClipNum {
				break loop
			}
		case <-timeout.C:
			break loop
		}
	}

	if clipNum > startClipNum {
		w32.OpenClipboard(0)
		txt := utf16zToString(readFromClipboard(w32.CF_UNICODETEXT))
		cnv, lastHKL := Transcode(txt)
		fmt.Printf("%q → %q\n", txt, cnv)
		w32.EmptyClipboard()
		sendToClipboard(w32.CF_UNICODETEXT, stringToUtf16z(cnv))
		w32.CloseClipboard()
		w32.SendInput(CtrlVSeq)
		if mode == MODE_SEL {
			w32.SendInput(ShiftDownSeq)
		}
		time.AfterFunc(250*time.Millisecond, func() { restoreClipboard(bkup) })

		// переключаем раскладку
		w32.PostMessage(win32.GetForegroundWindow(), w32.WM_INPUTLANGCHANGEREQUEST, 0, uintptr(lastHKL))
	}
}

var (
	ShiftDownSeq = []w32.INPUT{
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_SHIFT},
		},
	}

	ShiftUpSeq = []w32.INPUT{
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_SHIFT, DwFlags: win32.KEYEVENTF_KEYUP},
		},
	}

	CtrlCSeq = []w32.INPUT{
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_CONTROL},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: 0x43}, // C
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: 0x43, DwFlags: win32.KEYEVENTF_KEYUP},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_CONTROL, DwFlags: win32.KEYEVENTF_KEYUP},
		},
	}

	CtrlVSeq = []w32.INPUT{
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_CONTROL},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: 0x56}, // V
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: 0x56, DwFlags: win32.KEYEVENTF_KEYUP},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_CONTROL, DwFlags: win32.KEYEVENTF_KEYUP},
		},
	}

	CtrlShiftLeftSeq = []w32.INPUT{
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_SHIFT},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_CONTROL},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			// http://letcoderock.blogspot.de/2011/10/sendinput-with-shift-key-not-work.html
			Ki: w32.KEYBDINPUT{WVk: w32.VK_LEFT, DwFlags: win32.KEYEVENTF_EXTENDEDKEY},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_LEFT, DwFlags: win32.KEYEVENTF_KEYUP | win32.KEYEVENTF_EXTENDEDKEY},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_CONTROL, DwFlags: win32.KEYEVENTF_KEYUP},
		},
		w32.INPUT{
			Type: w32.INPUT_KEYBOARD,
			Ki:   w32.KEYBDINPUT{WVk: w32.VK_SHIFT, DwFlags: win32.KEYEVENTF_KEYUP},
		},
	}
)

func mustBeOk(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

type EWin struct {
	walk.WindowBase
	handler func()
}

const winClass = `HandleExplorerCrashWin`

var eCrashMessage uint32

func HandleExplorerCrash(foo func()) {
	walk.MustRegisterWindowClass(winClass)
	eCrashMessage = win32.RegisterWindowMessage("TaskbarCreated")
	w := new(EWin)
	w.handler = foo
	walk.InitWindow(w, nil, winClass, 0, 0)
}

func (e *EWin) WndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	if msg == eCrashMessage {
		log.Println("TaskbarCreated")
		time.AfterFunc(3*time.Second, e.handler)
	}
	return e.WindowBase.WndProc(hwnd, msg, wParam, lParam)
}
