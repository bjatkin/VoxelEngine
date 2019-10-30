package main

import (
	"syscall/js"
)

type mouse struct {
	x           float64
	y           float64
	leftClick   bool
	rightClick  bool
	middleClick bool
}

const (
	leftMouseButton   = 0
	middleMouseButton = 1
	rightMouseButton  = 2
)

func (m *mouse) init(doc, canvas js.Value) {
	mouseDownEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		if evt.Get("target") != canvas {
			return nil
		}
		button := evt.Get("button").Float()
		if button == leftMouseButton {
			m.leftClick = true
		}
		if button == middleMouseButton {
			m.middleClick = true
		}
		if button == rightMouseButton {
			m.rightClick = true
		}

		m.x = evt.Get("clientX").Float()
		m.y = evt.Get("clientY").Float()
		return nil
	})

	mouseUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m.leftClick = false
		m.rightClick = false
		return nil
	})

	mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		if evt.Get("target") != canvas {
			return nil
		}
		m.x = evt.Get("clientX").Float()
		m.y = evt.Get("clientX").Float()
		return nil
	})

	doc.Call("addEventListener", "mousedown", mouseDownEvt)
	doc.Call("addEventListener", "mouseup", mouseUpEvt)
	doc.Call("addEventListener", "mousemove", mouseMoveEvt)
}
