package main

import (
	"syscall/js"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/bobcat7/wasm-rotating-cube/gltypes"
)

var (
	glTypes      gltypes.GLTypes
	tmark        float32
	mouseInput   mouse
	keyInput     keyboard
	curSelection selection
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

	//Set up the save/load button
	initSaveBtn(scene)
	initLoadBtn(scene)

	//Set up the selection tracker
	curSelection = newSelection()

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
	selectColor       RGB
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
		rotateSpeed := float32(500.0)
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

	if !keyInput.keys[leftAlt] && !keyInput.keys[leftShift] && !addMode && !subMode && mouseInput.leftClick {
		if !selectMode {
			selectStartCorner = mgl32.Vec2{mouseInput.x, mouseInput.y}
			selectEndCorner = selectStartCorner
			curSelection.deselectAll(S)
		}

		selectMode = true
		//Hilight the new selection
		selectEndCorner = mgl32.Vec2{mouseInput.x, mouseInput.y}
		curSelection.newSelection(S, selectStartCorner, selectEndCorner, RGB{171, 129, 126})
	}

	if !mouseInput.leftClick && selectMode {
		// lock in a selection
		selectMode = false
	}

	if mouseInput.leftClick && addMode {
		// check if we're clicking on a selected voxel face
		i, f, succ := intersectVoxel(S, mgl32.Vec2{mouseInput.x, mouseInput.y})
		if succ && S.voxels[i].selected[f] {
			for i, j := range curSelection.voxels() {
				v := S.voxels[j]
				v.deselectFace(f)
				vox := v.newVoxelNeighbor(curSelection.face, selectColor)
				index := S.addVoxel(vox)
				curSelection.allVox[i] = index[0]
			}
		}
	}

	if mouseInput.leftClick && subMode {
		i, f, succ := intersectVoxel(S, mgl32.Vec2{mouseInput.x, mouseInput.y})
		remove := []int{}
		newSel := []int{}
		if succ && S.voxels[i].selected[f] {
			for _, v := range curSelection.voxels() {
				vox := S.voxels[v]
				x, y, z := faceToShift(f)
				add := -1
				for i, c := range S.voxels {
					if voxDist(c.x, c.y, c.z, vox.x+x, vox.y+y, vox.z+z) <= 0.01 {
						add = i
						break
					}
				}
				if add != -1 {
					newSel = append(newSel, add)
				}
				remove = append(remove, v)
			}
			curSelection.deselectAll(S)
			for _, i := range newSel {
				curSelection.addVox(i)
			}
			curSelection.selectAll(S, RGB{0, 0, 200})
			S.removeVoxel(remove...)
			S.update = true
		}

		//check if we have no more selection
		if curSelection.isEmpty() {
			keyInput.keys[esc] = true
		}
	}

	if keyInput.keys[aKey] && !curSelection.isEmpty() {
		keyInput.keys[aKey] = false
		addMode = !addMode
		subMode = false
		selectColor = RGB{200, 0, 0}
		curSelection.colorSelection(S, selectColor)
		S.update = true
	}

	if keyInput.keys[sKey] && !curSelection.isEmpty() {
		keyInput.keys[sKey] = false
		subMode = !subMode
		addMode = false
		selectColor = RGB{0, 0, 200}
		curSelection.colorSelection(S, selectColor)
		S.update = true
	}

	if keyInput.keys[esc] {
		addMode = false
		subMode = false
		curSelection.deselectAll(S)
		selectColor = RGB{171, 129, 126}
		curSelection.colorSelection(S, selectColor)
		S.update = true
	}

	//check if we need up update the colors of the current selection
	if new, color := ColorPicker.newColor(); new {
		curSelection.color(S, color)
	}
}

func start(context WebGl) {
	selectColor = RGB{171, 129, 126}
	globalConetxt = context
	ColorPicker.init()

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
