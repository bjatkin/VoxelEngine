package main

import (
	"syscall/js"
	"unsafe"

	"github.com/go-gl/mathgl/mgl32"
)

type Scene struct {
	voxels        []*Voxel
	gl            js.Value //Rendering context
	program       js.Value //shader program
	update        bool
	bufferData    []float32
	buffLen       int
	projMat       js.TypedArray
	uProjMat      js.Value
	viewMat       js.TypedArray
	uViewMat      js.Value
	modelMat      js.TypedArray
	uModelMat     js.Value
	aColor        js.Value
	aPosition     js.Value
	width, height int
}

const vShaderCode = `
attribute vec3 position;
uniform mat4 Pmatrix;
uniform mat4 Vmatrix;
uniform mat4 Mmatrix;
attribute vec3 color;
varying vec3 vColor;
void main(void) {
	gl_Position = Pmatrix*Vmatrix*Mmatrix*vec4(position, 1.);
	vColor = color;
}
`
const fShaderCode = `
precision mediump float;
varying vec3 vColor;
void main(void) {
	gl_FragColor = vec4(vColor, 1.);
}
`

func newScene(canvasId string, color RGB) *Scene {
	ret := Scene{}
	ret.gl, ret.width, ret.height = getCanvas(canvasId)

	vert := complieShader(ret.gl, glTypes.VertexShader, vShaderCode)
	frag := complieShader(ret.gl, glTypes.FragmentShader, fShaderCode)

	ret.program = linkProgram(ret.gl, vert, frag)
	ret.uProjMat = ret.gl.Call("getUniformLocation", ret.program, "Pmatrix")
	ret.uViewMat = ret.gl.Call("getUniformLocation", ret.program, "Vmatrix")
	ret.uModelMat = ret.gl.Call("getUniformLocation", ret.program, "Mmatrix")

	ret.aColor = ret.gl.Call("getAttribLocation", ret.program, "color")
	ret.aPosition = ret.gl.Call("getAttribLocation", ret.program, "position")

	ret.setProjMat(45)
	ret.setViewMat(mgl32.Vec3{3.0, 3.0, 3.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 1.0, 0.0})
	ret.setModelMat(0, 0, 0)

	bgColor := color.vec3()
	ret.gl.Call("useProgram", ret.program)
	ret.gl.Call("clearColor", bgColor[0], bgColor[1], bgColor[2], 1.0) // Color the screen is cleared to
	ret.gl.Call("clearDepth", 1.0)                                     // Z value that is set to the Depth buffer every frame
	ret.gl.Call("viewport", 0, 0, ret.width, ret.height)               // Viewport size
	ret.gl.Call("depthFunc", glTypes.LEqual)

	return &ret
}

func (s *Scene) buildBufferData() {
	if !s.update {
		return
	}

	s.update = false
	s.bufferData = []float32{}
	s.buffLen = 0
	for _, v := range s.voxels {
		for _, d := range v.data {
			s.bufferData = append(s.bufferData, d)
			s.buffLen++
		}
	}
	// Create a data buffer
	tArray := js.TypedArrayOf(s.bufferData)
	dataBuff := s.gl.Call("createBuffer")

	// Bind the data into the buffer
	floatBytes := 4
	s.gl.Call("bindBuffer", glTypes.ArrayBuffer, dataBuff)
	s.gl.Call("bufferData", glTypes.ArrayBuffer, tArray, glTypes.StaticDraw)

	s.gl.Call("vertexAttribPointer", s.aPosition, 3, glTypes.Float, false, floatBytes*6, 0)
	s.gl.Call("enableVertexAttribArray", s.aPosition)

	s.gl.Call("vertexAttribPointer", s.aColor, 3, glTypes.Float, false, floatBytes*6, floatBytes*3)
	s.gl.Call("enableVertexAttribArray", s.aColor)
}

func (s *Scene) addVoxel(v *Voxel) {
	s.voxels = append(s.voxels, v)
	s.update = true
}

func (s *Scene) setProjMat(deg float32) {
	// Generate a projection matrix
	ratio := float32(s.width) / float32(s.height)

	projMatrix := mgl32.Perspective(mgl32.DegToRad(deg), ratio, 1, 100.0)
	var projMatrixBuffer *[16]float32
	projMatrixBuffer = (*[16]float32)(unsafe.Pointer(&projMatrix))
	s.projMat = js.TypedArrayOf([]float32((*projMatrixBuffer)[:]))
}

func (s *Scene) setViewMat(eye, center, up mgl32.Vec3) {
	// Generate a view matrix
	viewMatrix := mgl32.LookAtV(eye, center, up)
	var viewMatrixBuffer *[16]float32
	viewMatrixBuffer = (*[16]float32)(unsafe.Pointer(&viewMatrix))
	s.viewMat = js.TypedArrayOf([]float32((*viewMatrixBuffer)[:]))
}

func (s *Scene) setModelMat(rotX, rotY, rotZ float32) {
	// Generate a model matrix
	movMatrix := mgl32.HomogRotate3DX(rotX)
	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DY(rotY))
	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DZ(rotZ))

	var modelMatrixBuffer *[16]float32
	modelMatrixBuffer = (*[16]float32)(unsafe.Pointer(&movMatrix))
	s.modelMat = js.TypedArrayOf([]float32((*modelMatrixBuffer)[:]))
}

func (s *Scene) render(tdiff float64) {
	s.buildBufferData() //build all the data and attach the interleved buffered data
	// Bind all the uniforms for this draw call
	s.gl.Call("uniformMatrix4fv", s.uProjMat, false, s.projMat)
	s.gl.Call("uniformMatrix4fv", s.uViewMat, false, s.viewMat)
	s.gl.Call("uniformMatrix4fv", s.uModelMat, false, s.modelMat)

	// Clear the screen
	s.gl.Call("enable", glTypes.DepthTest)
	s.gl.Call("clear", glTypes.ColorBufferBit)
	s.gl.Call("clear", glTypes.DepthBufferBit)

	// Make the draw call
	s.gl.Call("drawArrays", glTypes.Triangles, 0, s.buffLen)
}
