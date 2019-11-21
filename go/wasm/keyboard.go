package main

import (
	"fmt"
	"syscall/js"
)

type keyboard struct {
	keys map[int]bool
}

const (
	leftShift = 16
	leftAlt   = 18
	space     = 32
	leftKey   = 37
	upKey     = 38
	rightKey  = 39
	downKey   = 40
	aKey      = 65
	sKey      = 83
)

func (k *keyboard) init(doc, canvas js.Value) {
	if k.keys == nil {
		k.keys = make(map[int]bool)
	}

	keyDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		if evt.Get("target") != canvas {
			return nil
		}
		k.keys[evt.Get("keyCode").Int()] = true
		// get the key codes
		fmt.Printf("%s: %d\n", evt.Get("key").String(), evt.Get("keyCode").Int())
		return nil
	})

	keyUp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		evt := args[0]
		k.keys[evt.Get("keyCode").Int()] = false
		return nil
	})

	doc.Call("addEventListener", "keydown", keyDown)
	doc.Call("addEventListener", "keyup", keyUp)
}
