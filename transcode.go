package main

import (
	"regexp"
	"strings"

	"github.com/AllenDang/w32"
	"github.com/davidmz/k-switcher/win32"
)

type Layout struct {
	AllChars string
	Chars    string
	Suffix   string
	Handle   w32.HKL
}

var (
	LayoutA = "`1234567890-=\\qwertyuiop[]asdfghjkl;'zxcvbnm,./" +
		"~!@#$%^&*()_+|QWERTYUIOP{}ASDFGHJKL:\"ZXCVBNM<>?"

	LayoutB = "ё1234567890-=\\йцукенгшщзхъфывапролджэячсмитьбю." +
		"Ё!\"№;%:?*()_+/ЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ,"

	LayoutASuffix = "0409"
	LayoutBSuffix = "0419"

	wordRe = regexp.MustCompile(`\S+(\s+|$)`)

	MapAtoB, MapBtoA       map[rune]rune
	CharsA, CharsB         string
	LayoutAHKL, LayoutBHKL w32.HKL
)

const (
	LtUndef = 0
	LtA     = 1
	LtB     = 2
)

func Transcode(str string) (outs string, lastHKL w32.HKL) {
	var word, out []rune
	var lt = LtUndef

	for _, r := range str {
		rLt := LtUndef
		if strings.ContainsRune(CharsA, r) {
			rLt = LtA
		} else if strings.ContainsRune(CharsB, r) {
			rLt = LtB
		}
		if lt == LtUndef || rLt == LtUndef || lt == rLt { // продолжаем то же слово
			if lt == LtUndef {
				lt = rLt
			}
			word = append(word, r)
		} else { // новое слово, с другой раскладкой
			// конвертируем и сохраняем имеющееся слово
			if lt == LtA {
				out = append(out, trns(word, MapAtoB)...)
			} else if lt == LtB {
				out = append(out, trns(word, MapBtoA)...)
			} else { // unreachable
				out = append(out, word...)
			}
			// стартуем новое слово
			word = []rune{r}
			lt = rLt
		}
	}
	// конвертируем и сохраняем имеющееся слово
	if lt == LtA {
		out = append(out, trns(word, MapAtoB)...)
		lastHKL = LayoutBHKL
	} else if lt == LtB {
		out = append(out, trns(word, MapBtoA)...)
		lastHKL = LayoutAHKL
	} else { // unreachable
		out = append(out, word...)
	}
	outs = string(out)
	return
}

func trns(in []rune, mp map[rune]rune) (out []rune) {
	out = make([]rune, len(in))
	for i, r := range in {
		if r1, ok := mp[r]; ok {
			out[i] = r1
		} else {
			out[i] = r
		}
	}
	return
}

func init() {
	MapAtoB = make(map[rune]rune)
	MapBtoA = make(map[rune]rune)

	rA, rB := []rune(LayoutA), []rune(LayoutB)
	cA, cB := []rune{}, []rune{}

	for i, r := range rA {
		MapAtoB[r] = rB[i]
		if !strings.ContainsRune(LayoutB, r) { // руна содержится только в раскладке A
			cA = append(cA, r)
		}
	}
	for i, r := range rB {
		MapBtoA[r] = rA[i]
		if !strings.ContainsRune(LayoutA, r) { // руна содержится только в раскладке B
			cB = append(cB, r)
		}
	}

	CharsA = string(cA)
	CharsB = string(cB)

	list := win32.GetKeyboardLayoutList()
	for _, l := range list {
		win32.ActivateKeyboardLayout(l, 0)
		name := win32.GetKeyboardLayoutName()
		if strings.HasSuffix(name, LayoutASuffix) {
			LayoutAHKL = l
		}
		if strings.HasSuffix(name, LayoutBSuffix) {
			LayoutBHKL = l
		}
	}
}
