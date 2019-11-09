package main

import (
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
	scene := newScene(goCanvas, newRGB(230, 230, 230))

	//Start listening for input
	mouseInput = mouse{}
	mouseInput.init(doc, goCanvas)
	keyInput = keyboard{}
	keyInput.init(doc, goCanvas)

	//Create the base scene
	baseColor := newRGBSet(235, 254, 255)
	vox := []*Voxel{}
	voxCount := 50
	for x := 0; x < voxCount; x++ {
		for y := 0; y < voxCount; y++ {
			vox = append(vox, newVoxel(float32(x)-float32(voxCount)/2, 0, float32(y)-float32(voxCount)/2, baseColor))
		}
	}
	scene.addVoxel(vox...)

	scene.moveCamera(0, 0, 85)
	scene.rotateCamera(-1, 0, 0)

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
		cy = cy * 1.5 / deltaT

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
			closest := float32(99999999.0)
			r := newRay(S, mouseInput.x, mouseInput.y)
			sel := -1
			selFace := 0
			for i, v := range S.voxels {
				face, dist, hit := v.intersect(&r, closest)
				if hit && dist < closest {
					closest = dist
					sel = i
					selFace = face
				}
			}

			if sel >= 0 {
				v := S.voxels[sel]
				S.addVoxel(v.newVoxelNeighbor(selFace))
				// S.voxels[sel].setColor(newRGB(255, 0, 0), selFace)
				// S.update = true //force the buffers to rebuild
			}
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
