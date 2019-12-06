package main

import (
	"fmt"
	"syscall/js"
)

type colorPicker struct {
	colorChanged bool
	color        RGB
}

var ColorPicker = colorPicker{}

func (cp *colorPicker) init() {
	js.Global().Set("updateColorPicker", js.FuncOf(updateColorPicker))
}

func (cp *colorPicker) newColor() (bool, RGB) {
	if cp.colorChanged {
		cp.colorChanged = false
		return true, cp.color
	}
	return false, cp.color
}

func updateColorPicker(this js.Value, args []js.Value) interface{} {
	fmt.Printf("The inputs were: %v\n", args)
	r := args[0].Float()
	g := args[1].Float()
	b := args[2].Float()

	ColorPicker.color = newRGB(float32(r), float32(g), float32(b))
	ColorPicker.colorChanged = true
	return nil
}
