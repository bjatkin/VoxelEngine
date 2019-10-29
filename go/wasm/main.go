package main

import (
	"github.com/bobcat7/wasm-rotating-cube/gltypes"
)

var glTypes gltypes.GLTypes

func main() {
	//Create a new scene
	scene := newScene("gocanvas", newRGB(28, 46, 43))

	test := newVoxel(0.0, 0.0, 0.0,
		[6]RGB{
			newRGB(1.0, 0.0, 0.0),
			newRGB(0.0, 1.0, 0.0),
			newRGB(0.0, 0.0, 1.0),
			newRGB(1.0, 0.0, 1.0),
			newRGB(1.0, 1.0, 0.0),
			newRGB(0.0, 1.0, 1.0),
		},
	)
	scene.addVoxel(test)

	wgl := newWebGL(scene)
	start(wgl)
}
