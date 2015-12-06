package main

import (
	"strings"

	"github.com/davidmz/k-switcher.v2/errors"
	"github.com/davidmz/k-switcher.v2/win32"
	"github.com/davidmz/mustbe"
	"golang.org/x/sys/windows/registry"
)

type KLayout struct {
	Key    win32.HKL
	Name   string
	Title  string
	Layout string
}

func GetSystemLayouts() (outList []*KLayout, outErr error) {
	defer mustbe.Catched(func(err error) { outList, outErr = nil, err })

	list := win32.GetKeyboardLayoutList()
	for _, k := range list {
		win32.ActivateKeyboardLayout(k, 0)
		name := win32.GetKeyboardLayoutName()

		func() {
			key, err := registry.OpenKey(
				registry.LOCAL_MACHINE,
				`SYSTEM\CurrentControlSet\Control\Keyboard Layouts\`+strings.ToLower(name),
				registry.QUERY_VALUE,
			)
			mustbe.OK(errors.Wrap("error reading registry", err))
			defer key.Close()

			title, _, err := key.GetStringValue(`Layout Text`)
			mustbe.OK(errors.Wrap("error reading registry", err))

			outList = append(outList, &KLayout{Key: k, Name: name, Title: title})
		}()
	}
	return
}
