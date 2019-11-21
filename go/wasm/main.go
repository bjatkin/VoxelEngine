package main

import (
	"fmt"
	"syscall/js"

	"github.com/go-gl/mathgl/mgl32"

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
	zoomStart         float32
	originalZoom      float32
	applyZoom         float32
	zooming           bool
	selectMode        bool
	selectStartCorner mgl32.Vec2
	selectEndCorner   mgl32.Vec2
	addMode           bool
	subMode           bool
)

func update(deltaT float32, scenes []*Scene) {
	S := scenes[0]

	mdx, mdy := (mouseInput.dx/deltaT)/float32(S.width), (mouseInput.dy/deltaT)/float32(S.height)

	if (keyInput.keys[leftShift] && mouseInput.leftClick) || mouseInput.middleClick {
		//Pan
		panSpeed := float32(5000.0)
		S.moveCamera(-mdx*panSpeed, mdy*panSpeed, 0)
	}

	if keyInput.keys[leftAlt] && mouseInput.leftClick {
		//Rotate
		rotateSpeed := float32(700.0)
		S.rotateCamera(0, mdx*rotateSpeed, mdy*rotateSpeed)
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

	if !keyInput.keys[leftAlt] && !keyInput.keys[leftShift] && !addMode && mouseInput.leftClick {
		if !selectMode {
			selectStartCorner = mgl32.Vec2{mouseInput.x, mouseInput.y}
			selectEndCorner = selectStartCorner
			currentSelection.emptySelection()
			fmt.Printf("empty select\n")
		}

		selectMode = true
		//Release the old selection
		clearSelection(S, selectStartCorner, selectEndCorner)

		//Hilight the new selection
		selectEndCorner = mgl32.Vec2{mouseInput.x, mouseInput.y}
		hilightSelection(S, selectStartCorner, selectEndCorner)
	}

	if !mouseInput.leftClick && selectMode {
		// lock in a selection
		selectMode = false
	}

	if mouseInput.leftClick && addMode {
		for i, v := range currentSelection.cubes {
			if i >= currentSelection.cubesLen {
				break
			}
			vox := v.newVoxelNeighbor(currentSelection.face)
			currentSelection.cubes[i] = vox
			S.addVoxel(vox)
		}
	}

	if mouseInput.leftClick && subMode {
		//TODO implement this
	}

	if keyInput.keys[aKey] {
		keyInput.keys[aKey] = false
		addMode = !addMode
		subMode = !addMode
	}

	if keyInput.keys[sKey] {
		keyInput.keys[sKey] = false
		subMode = !subMode
		addMode = !subMode
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
