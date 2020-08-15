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

func getContext(canvas js.Value) (js.Value, int, int) {
	// Init Canvas
	width := canvas.Get("clientWidth").Int()
	height := canvas.Get("clientHeight").Int()
	canvas.Call("setAttribute", "width", width)
	canvas.Call("setAttribute", "height", height)
	canvas.Set("tabIndex", 0) // Not sure if this is needed

	gl := canvas.Call("getContext", "webgl")
	if gl == js.Undefined() {
		gl = canvas.Call("getContext", "experimental-webgl")
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

type webGl struct {
	scenes []*scene
}

var globalConetxt webGl

func newWebGL(scenes ...*scene) webGl {
	return webGl{
		scenes: scenes,
	}
}

func (w *webGl) addScene(scene *scene) {
	w.scenes = append(w.scenes, scene)
}
