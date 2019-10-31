package main

import (
	"math/rand"
	"syscall/js"

	"github.com/bobcat7/wasm-rotating-cube/gltypes"
)

var (
	glTypes    gltypes.GLTypes
	tmark      float32
	mouseInput mouse
	keyInput   keyboard
)

func main() {
	//Get the canvas element
	doc := js.Global().Get("document")
	goCanvas := doc.Call("getElementById", "gocanvas")

	//Create a new scene
	scene := newScene(goCanvas, newRGB(28, 46, 43))

	//Start listening for input
	mouseInput = mouse{}
	mouseInput.init(doc, goCanvas)
	keyInput = keyboard{}
	keyInput.init(doc, goCanvas)

	//Create the base scene
	baseColor := newRGBSet(114, 237, 235)
	scene.addVoxel(
		newVoxel(0, 0, 0, baseColor),
		newVoxel(-1, 0, 0, baseColor),
		newVoxel(-1, 0, 1, baseColor),
		newVoxel(0, 0, 1, baseColor),
	)
	scene.moveCamera(0, 1, 7)
	scene.rotateCamera(0, 0, 0)

	//Start the render loop
	start(newWebGL(scene))

	//stop from exiting, only nessisary if not compling
	//with tiny go
	done := make(chan bool)

	<-done
}

var (
	zoomStart    float32
	originalZoom float32
	applyZoom    float32
	zooming      bool
)

func update(deltaT float32, scenes []*Scene) {
	S := scenes[0]

	mdx, mdy := mouseInput.dx/deltaT, mouseInput.dy/deltaT

	if (keyInput.keys[leftShift] && mouseInput.leftClick) || mouseInput.middleClick {
		//Pan
		S.moveCamera(-mdx, mdy, 0)
	}

	if keyInput.keys[leftAlt] && mouseInput.leftClick {
		//Rotate
		S.rotateCamera(0, mdx*0.5, mdy*0.5)
	}

	if mouseInput.rightClick {
		//Zoom
		if !zooming {
			originalZoom = S.cameraLoc[2]
			zoomStart = mouseInput.y
			zooming = true
		}
		cy := zoomStart - mouseInput.y
		cy = cy * 0.5 / deltaT

		dx := cy
		applyZoom = dx
		S.setCameraLoc(S.cameraLoc[0], S.cameraLoc[1], originalZoom+applyZoom)
	}

	if !mouseInput.rightClick && zooming {
		zooming = false
	}

	if !keyInput.keys[leftAlt] && !keyInput.keys[leftShift] && mouseInput.leftClick {
		//Select Voxel face
		if mouseInput.leftClick {
			v := len(S.voxels)
			r := rand.Intn(v)
			new := S.voxels[r].newVoxelNeighbor(upFace)
			S.addVoxel(new)
			mouseInput.leftClick = false
		}
	}
}

func start(context WebGl) {
	globalConetxt = context

	//Needs to be added in if you're not compiling with tinygo
	//This will export the renderFrame function
	js.Global().Set("renderFrame", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		renderFrame(float32(args[0].Float()))
		return nil
	}))

	//Kick the process off
	js.Global().Call("requestAnimationFrame", js.Global().Get("renderFrame"))
}

//go:export renderFrame
func renderFrame(now float32) {
	deltaT := now - tmark
	tmark = now

	//Run the update loop
	update(deltaT, globalConetxt.scenes)
	for _, s := range globalConetxt.scenes {
		s.render()
	}
	js.Global().Call("requestAnimationFrame", js.Global().Get("renderFrame"))
}
