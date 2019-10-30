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
	width, height int      //width and height of the canvas
	bufferData    []float32
	buffLen       int

	projMat   js.TypedArray //Uniforms
	uProjMat  js.Value
	viewMat   js.TypedArray
	uViewMat  js.Value
	modelMat  js.TypedArray
	uModelMat js.Value

	aColor    js.Value //attributes
	aPosition js.Value
	aNormal   js.Value

	update bool //does the data need to be updated?
}

const vShaderCode = `
uniform mat4 u_Pmatrix;
uniform mat4 u_Vmatrix;
uniform mat4 u_Mmatrix;

attribute vec3 a_Color;
attribute vec3 a_Position;
attribute vec3 a_Normal;

varying vec3 v_Color;
varying vec3 v_Normal;

void main(void) {
	gl_Position = u_Pmatrix*u_Vmatrix*u_Mmatrix*vec4(a_Position, 1.0);
	v_Normal = vec3( u_Pmatrix*u_Vmatrix*u_Mmatrix*vec4(a_Normal, 0.0));
	v_Color = a_Color;
}
`
const fShaderCode = `
precision mediump float;

varying vec3 v_Color;
varying vec3 v_Normal;

void main(void) {
	
	vec3 to_light;
	float cos_angle;
	float ambient_light;

	//normalized 1, 1, 1
	to_light = vec3(0.57735026919, 0.57735026919, 0.57735026919);

	//amount of ambient light in the scene
	ambient_light = 0.3;

	//for difuse light calculation
	cos_angle = dot(v_Normal, to_light);
	cos_angle = clamp(cos_angle, 0.0, 1.0);

	gl_FragColor = vec4(v_Color*ambient_light + (v_Color*cos_angle) * (1.0-ambient_light), 1.0);
}
`

func newScene(canvas js.Value, color RGB) *Scene {
	ret := Scene{}
	ret.gl, ret.width, ret.height = getContext(canvas)

	vert := complieShader(ret.gl, glTypes.VertexShader, vShaderCode)
	frag := complieShader(ret.gl, glTypes.FragmentShader, fShaderCode)

	ret.program = linkProgram(ret.gl, vert, frag)
	ret.uProjMat = ret.gl.Call("getUniformLocation", ret.program, "u_Pmatrix")
	ret.uViewMat = ret.gl.Call("getUniformLocation", ret.program, "u_Vmatrix")
	ret.uModelMat = ret.gl.Call("getUniformLocation", ret.program, "u_Mmatrix")

	ret.aPosition = ret.gl.Call("getAttribLocation", ret.program, "a_Position")
	ret.gl.Call("enableVertexAttribArray", ret.aPosition)
	ret.aColor = ret.gl.Call("getAttribLocation", ret.program, "a_Color")
	ret.gl.Call("enableVertexAttribArray", ret.aColor)
	ret.aNormal = ret.gl.Call("getAttribLocation", ret.program, "a_Normal")
	ret.gl.Call("enableVertexAttribArray", ret.aNormal)

	ret.setProjMat(45)
	ret.setViewMat(mgl32.Vec3{3.0, 3.0, 3.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 1.0, 0.0})
	ret.setModelMat(0, 0, 0)

	bgColor := color.vec3()
	ret.gl.Call("useProgram", ret.program)
	ret.gl.Call("clearColor", bgColor[0], bgColor[1], bgColor[2], 1.0) // Color the screen is cleared to
	ret.gl.Call("clearDepth", 1.0)                                     // Z value that is set to the Depth buffer every frame
	ret.gl.Call("viewport", 0, 0, ret.width, ret.height)               // Viewport size
	ret.gl.Call("depthFunc", glTypes.LEqual)
	ret.gl.Call("enable", 2884) //gl.CULL_FACE
	ret.gl.Call("enable", glTypes.DepthTest)

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
	float32Bytes := 4
	posSize := 3
	colSize := 3
	normSize := 3
	vertexByteSize := posSize*float32Bytes + colSize*float32Bytes + normSize*float32Bytes
	s.buffLen /= posSize + colSize + normSize

	s.gl.Call("bindBuffer", glTypes.ArrayBuffer, dataBuff)
	s.gl.Call("bufferData", glTypes.ArrayBuffer, tArray, glTypes.StaticDraw)
	s.gl.Call("vertexAttribPointer", s.aPosition, posSize, glTypes.Float, false, vertexByteSize, 0)
	s.gl.Call("vertexAttribPointer", s.aColor, colSize, glTypes.Float, false, vertexByteSize, posSize*float32Bytes)
	s.gl.Call("vertexAttribPointer", s.aNormal, colSize, glTypes.Float, false, vertexByteSize, posSize*float32Bytes+colSize*float32Bytes)
}

func (s *Scene) addVoxel(voxels ...*Voxel) {
	for _, v := range voxels {
		s.voxels = append(s.voxels, v)
	}
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

	//TODO this call is leaking memory, find out why/ how to stop it!
	s.modelMat = js.TypedArrayOf([]float32((*modelMatrixBuffer)[:]))
}

func (s *Scene) render() {
	//build all the data and attach the interleved buffered data
	s.buildBufferData()

	// Bind all the uniforms for this draw call
	s.gl.Call("uniformMatrix4fv", s.uProjMat, false, s.projMat)
	s.gl.Call("uniformMatrix4fv", s.uViewMat, false, s.viewMat)
	s.gl.Call("uniformMatrix4fv", s.uModelMat, false, s.modelMat)

	// Clear the screen
	s.gl.Call("clear", glTypes.ColorBufferBit)
	s.gl.Call("clear", glTypes.DepthBufferBit)

	// Make the draw call
	s.gl.Call("drawArrays", glTypes.Triangles, 0, s.buffLen)
}
