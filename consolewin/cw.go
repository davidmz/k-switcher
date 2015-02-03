package consolewin

import (
	"time"

	"github.com/AllenDang/w32"
	"github.com/davidmz/k-switcher/win32"
)

var (
	win                  = w32.GetConsoleWindow()
	isHidden             bool
	controlEventHandlers []ControlEventHandler
	showStateHandlers    []ShowStateHandler

	ShowState uint32 = w32.SW_HIDE
)

type (
	ControlEventHandler func(eventType int) bool
	ShowStateHandler    func(showState uint32)
)

func Show() {
	w32.ShowWindow(win, w32.SW_NORMAL)
	win32.SetForegroundWindow(win)
	isHidden = false
}

func Hide() {
	w32.ShowWindow(win, w32.SW_HIDE)
	isHidden = true
}

func SetTitle(title string) {
	win32.SetConsoleTitle(title)
}

func OnControlEvent(h ControlEventHandler) {
	controlEventHandlers = append(controlEventHandlers, h)
}

func OnShowStateChange(h ShowStateHandler) {
	showStateHandlers = append(showStateHandlers, h)
}

func init() {
	win32.SetConsoleCtrlHandler(func(t int) int {
		for _, h := range controlEventHandlers {
			if h(t) {
				return 1
			}
		}
		return 0
	}, true)

	go func() {
		tick := time.Tick(500 * time.Millisecond)
		for _ = range tick {
			if len(showStateHandlers) > 0 {
				var st uint32 = w32.SW_HIDE
				if !isHidden {
					if pl := win32.GetWindowPlacement(win); pl != nil {
						st = pl.ShowCmd
					}
				}
				if st != ShowState {
					ShowState = st
					for _, h := range showStateHandlers {
						h(st)
					}
				}
			}
		}
	}()
}
