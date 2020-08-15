package main

import (
	"fmt"
	"syscall/js"
)

type mouse struct {
	x           float32
	y           float32
	dx          float32
	dy          float32
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
		if !evt.Get("target").Equal(canvas) {
			return nil
		}
		button := evt.Get("button").Int()
		fmt.Printf("CLICK %d\n", button)
		if button == leftMouseButton {
			m.leftClick = true
		}
		if button == middleMouseButton {
			m.middleClick = true
		}
		if button == rightMouseButton {
			m.rightClick = true
		}

		newX := float32(evt.Get("clientX").Float()) - 135
		newY := float32(evt.Get("clientY").Float())
		m.dx = newX - m.x
		m.dy = newY - m.y
		m.x = newX
		m.y = newY
		return nil
	})

	mouseUpEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		m.leftClick = false
		m.rightClick = false
		m.middleClick = false
		return nil
	})

	mouseMoveEvt := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		if !evt.Get("target").Equal(canvas) {
			m.leftClick = false
			m.rightClick = false
			m.middleClick = false
			return nil
		}

		newX := float32(evt.Get("clientX").Float()) - 135
		newY := float32(evt.Get("clientY").Float())
		m.dx = newX - m.x
		m.dy = newY - m.y
		m.x = newX
		m.y = newY
		return nil
	})

	doc.Call("addEventListener", "mousedown", mouseDownEvt)
	doc.Call("addEventListener", "mouseup", mouseUpEvt)
	doc.Call("addEventListener", "mousemove", mouseMoveEvt)
}
