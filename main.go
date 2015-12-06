package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/davidmz/k-switcher/win32"
)

var debug *log.Logger

func main() {
	flags := &struct {
		ShowHelp  bool
		DebugMode bool
	}{}

	flag.Usage = func() {
		fmt.Println("Usage: " + filepath.Base(os.Args[0]) + " [options] [LAYOUT_NAME_1 LAYOUT_NAME_2]")
		fmt.Println("Where options are:")
		flag.PrintDefaults()
		Pause()
	}

	flag.BoolVar(&flags.ShowHelp, "h", false, "Show options")
	flag.BoolVar(&flags.DebugMode, "debug", false, "Turn debug log on")
	flag.Parse()

	if flags.ShowHelp {
		flag.Usage()
		return
	}

	if flags.DebugMode {
		debug = log.New(os.Stderr, "DEBUG ", log.LstdFlags)
	} else {
		debug = log.New(ioutil.Discard, "DEBUG ", log.LstdFlags)
	}

	kList, err := GetSystemLayouts()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error fetching system layouts:", err)
		Pause()
		os.Exit(1)
	}
	fmt.Printf("Found %d layout(s) in your system:\n", len(kList))
	ourLayouts := []*KLayout{}
	for _, k := range kList {
		fmt.Printf("%s\t%s\n", k.Name, k.Title)
		if l, ok := Layouts.Get(k.Name); ok {
			k.Layout = l
			ourLayouts = append(ourLayouts, k)
		}
	}
	fmt.Println("")

	switch len(ourLayouts) {
	case 0:
		fmt.Fprintf(os.Stderr, "No layouts is supported, sorry.\n")
		Pause()
		os.Exit(1)
	case 1:
		fmt.Fprintf(os.Stderr, "Only one layout is supported (%s), k-switcher is useless.\n", ourLayouts[0].Name)
		Pause()
		os.Exit(1)
	case 2:
		fmt.Printf("OK, will switch between %s and %s.\n", ourLayouts[0].Name, ourLayouts[1].Name)
	default:
		lNames := make([]string, len(ourLayouts))
		for i, l := range ourLayouts {
			lNames[i] = l.Name
		}
		fmt.Fprintf(os.Stderr, "%d layouts is supported (%s) but we can switch only between two.\n", len(ourLayouts), strings.Join(lNames, ", "))
		fmt.Fprintf(os.Stderr, "Please specify layout names in program arguments.\n")
		Pause()
		os.Exit(1)
	}

	trans := NewTranscoder(ourLayouts[0], ourLayouts[1])

	if !win32.RegisterHotKey(0, 0, win32.MOD_SHIFT|win32.MOD_NOREPEAT, win32.VK_PAUSE) {
		fmt.Fprintln(os.Stderr, "Can not register Shift+Break hotkey.")
		Pause()
		os.Exit(1)
	}

	msg := new(win32.MSG)
	for win32.GetMessage(msg, 0, 0, 0) {
		if msg.Message == win32.WM_HOTKEY {
			debug.Println("===========", time.Now())
			debug.Println("Start handler")
			HandleHotkey(trans)
			debug.Println("End handler")
		}
	}
}

func Pause() {
	fmt.Fprintf(os.Stderr, "Press Enter to exit...\n")
	bufio.NewReader(os.Stdin).ReadLine()
}
