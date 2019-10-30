package main

import (
	"fmt"
	"math/rand"
	"syscall/js"

	"github.com/bobcat7/wasm-rotating-cube/gltypes"
)

var (
	glTypes    gltypes.GLTypes
	tmark      float64
	mouseInput mouse
)

func main() {
	fmt.Printf("Start\n")
	//Get the canvas element
	doc := js.Global().Get("document")
	goCanvas := doc.Call("getElementById", "gocanvas")

	//Create a new scene
	scene := newScene(goCanvas, newRGB(28, 46, 43))

	//Start listening for input
	mouseInput = mouse{}
	mouseInput.init(doc, goCanvas)

	//Create the base scene
	baseColor := newRGBSet(61, 191, 189)
	scene.addVoxel(
		newVoxel(0, 0, 0, baseColor),
		newVoxel(-1, 0, 0, baseColor),
		newVoxel(-1, 0, 1, baseColor),
		newVoxel(0, 0, 1, baseColor),
	)

	//Start the render loop
	start(newWebGL(scene))

	//stop from exiting, only nessisary if not compling
	//with tiny go
	done := make(chan bool)

	<-done
}

var Rot = float32(0.0)

func update(deltaT float64, scenes []*Scene) {
	S := scenes[0]
	if mouseInput.leftClick {
		v := len(S.voxels)
		r := rand.Intn(v)
		new := S.voxels[r].newVoxelNeighbor(upFace)
		S.addVoxel(new)
		mouseInput.leftClick = false
	}

	Rot += float32(deltaT / 900.0)
	S.setModelMat(0.5*Rot, 0.2*Rot, 0.3*Rot)
}

func start(context WebGl) {
	globalConetxt = context

	//Needs to be added in if you're not compiling with tinygo
	js.Global().Set("renderFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		renderFrame(args[0].Float())
		return nil
	}))

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
