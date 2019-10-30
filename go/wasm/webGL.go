package main

import (
	"syscall/js"
)

func complieShader(gl js.Value, shaderType js.Value, code string) js.Value {
	shader := gl.Call("createShader", shaderType)
	gl.Call("shaderSource", shader, code)
	gl.Call("compileShader", shader)
	return shader
}

func linkProgram(gl, vertShader, fragShader js.Value) js.Value {
	program := gl.Call("createProgram")
	gl.Call("attachShader", program, vertShader)
	gl.Call("attachShader", program, fragShader)
	gl.Call("linkProgram", program)
	return program
}

func getCanvas(id string) (js.Value, int, int) {
	// Init Canvas
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", id)
	width := canvasEl.Get("clientWidth").Int()
	height := canvasEl.Get("clientHeight").Int()
	canvasEl.Call("setAttribute", "width", width)
	canvasEl.Call("setAttribute", "height", height)
	canvasEl.Set("tabIndex", 0) // Not sure if this is needed

	gl := canvasEl.Call("getContext", "webgl")
	if gl == js.Undefined() {
		gl = canvasEl.Call("getContext", "experimental-webgl")
	}
	// once again
	if gl == js.Undefined() {
		js.Global().Call("alert", "browser might not support webgl")
		return gl, 0, 0
	}

	// Get some WebGL bindings
	glTypes.New(gl)
	return gl, width, height
}

type WebGl struct {
	scenes []*Scene
}

var globalConetxt WebGl

func newWebGL(scenes ...*Scene) WebGl {
	return WebGl{
		scenes: scenes,
	}
}

func (w *WebGl) addScene(scene *Scene) {
	w.scenes = append(w.scenes, scene)
}
