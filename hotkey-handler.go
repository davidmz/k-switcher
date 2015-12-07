package main

import (
	"time"

	"github.com/davidmz/k-switcher/clipboard"
	"github.com/davidmz/k-switcher/win32"
)

func HandleHotkey(trans *Transcoder) {
	clipData := clipboard.Get()
	debug.Printf("clipData: %q", clipData)

	clipNum := clipboard.GetSeqNumber()
	SendCtrlC()
	if !WaitForClipboard(clipNum, 500*time.Millisecond) {
		debug.Println("No changes in clipboard")
		return
	}
	clipText := clipboard.Get()
	debug.Printf("clipText: %q", clipText)
	newText, lt := trans.Transcode(clipboard.Get())
	if lt == nil {
		debug.Printf("%q not converted", clipText)
		clipboard.Put(clipData)
		return
	}
	debug.Printf("%q → %q", clipText, newText)
	clipboard.Empty()
	clipboard.Put(newText)
	SendCtrlV()

	time.Sleep(250 * time.Millisecond)
	clipboard.Put(clipData)

	// переключаем раскладку
	win32.PostMessage(win32.GetForegroundWindow(), win32.WM_INPUTLANGCHANGEREQUEST, 0, win32.LPARAM(lt.Key))
}

func WaitForClipboard(initNum uint32, timeout time.Duration) bool {
	timeStart := time.Now()
	for clipboard.GetSeqNumber() == initNum {
		time.Sleep(20 * time.Millisecond)
		if time.Since(timeStart) >= timeout {
			return false
		}
	}
	return true
}

func SendCtrlC() {
	keys := []win32.INPUT{}
	shift := win32.GetAsyncKeyState(win32.VK_SHIFT)
	if shift {
		keys = append(keys, kbdInputUp(win32.VK_SHIFT))
	}
	keys = append(keys,
		kbdInputDown(win32.VK_CONTROL),
		kbdInputDown(0x43), // C
		kbdInputUp(0x43),   // C
		kbdInputUp(win32.VK_CONTROL),
	)
	if shift {
		keys = append(keys, kbdInputDown(win32.VK_SHIFT))
	}
	win32.SendInput(keys)
}

func SendCtrlV() {
	keys := []win32.INPUT{}
	shift := win32.GetAsyncKeyState(win32.VK_SHIFT)
	if shift {
		keys = append(keys, kbdInputUp(win32.VK_SHIFT))
	}
	keys = append(keys,
		kbdInputDown(win32.VK_CONTROL),
		kbdInputDown(0x56), // C
		kbdInputUp(0x56),   // C
		kbdInputUp(win32.VK_CONTROL),
	)
	if shift {
		keys = append(keys, kbdInputDown(win32.VK_SHIFT))
	}
	win32.SendInput(keys)
}

func kbdInputUp(vKey uint16) win32.INPUT   { return kbdInput(vKey, false) }
func kbdInputDown(vKey uint16) win32.INPUT { return kbdInput(vKey, true) }

func kbdInput(vKey uint16, down bool) win32.INPUT {
	var dwFlags uint32
	if !down {
		dwFlags = dwFlags | win32.KEYEVENTF_KEYUP
	}
	return win32.INPUT{
		Type: win32.INPUT_KEYBOARD,
		Ki: win32.KEYBDINPUT{
			WVk:     vKey,
			DwFlags: dwFlags,
		},
	}
}
