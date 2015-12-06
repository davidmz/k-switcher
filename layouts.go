package main

import "strings"

const (
	stdEnLayout = "`1234567890-=\\qwertyuiop[]asdfghjkl;'zxcvbnm,./~!@#$%^&*()_+|QWERTYUIOP{}ASDFGHJKL:\"ZXCVBNM<>?"
	stdRuLayout = "ё1234567890-=\\йцукенгшщзхъфывапролджэячсмитьбю.Ё!\"№;%:?*()_+/ЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ,"
)

var Layouts = LayoutsMap{
	// Стандартные раскладки US и Russian
	strings.ToLower("00000409"): stdEnLayout,
	strings.ToLower("00000419"): stdRuLayout,
	// Раскладки Бирмана
	strings.ToLower("A0000409"): stdEnLayout,
	strings.ToLower("A0000419"): stdRuLayout,
}

type LayoutsMap map[string]string

func (l LayoutsMap) Get(name string) (s string, ok bool) { s, ok = l[strings.ToLower(name)]; return }
