package main

import (
	"syscall/js"

	"github.com/bobcat7/wasm-rotating-cube/gltypes"
)

var (
	glTypes gltypes.GLTypes
	tmark   float64
)

func main() {
	//Create a new scene
	scene := newScene("gocanvas", newRGB(28, 46, 43))

	test := newVoxel(-0.5, -0.5, 0.5, newRGBSet(61, 191, 189))
	scene.addVoxel(test)

	start(newWebGL(scene))
}

var Rot = float32(0.0)

func update(deltaT float64, scenes []*Scene) {
	S := scenes[0]
	Rot += float32(deltaT / 500.0)
	S.setModelMat(0.5*Rot, 0.2*Rot, 0.3*Rot)
}

func start(context WebGl) {
	globalConetxt = context
	//Kick the process off
	js.Global().Call("requestAnimationFrame", js.Global().Get("renderFrame"))
}

//go:export renderFrame
func renderFrame(now float64) {
	deltaT := now - tmark
	tmark = now

	//Run the update loop
	update(deltaT, globalConetxt.scenes)
	for _, s := range globalConetxt.scenes {
		s.render()
	}
	js.Global().Call("requestAnimationFrame", js.Global().Get("renderFrame"))
}
