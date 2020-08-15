package main

import (
	"syscall/js"

	"github.com/go-gl/mathgl/mgl32"
)

//TODO fix this code so it does not need to use js.TypedArrays as these were removed
//https://github.com/golang/go/issues/31980
//Instead create a new Float32Array
//Then use the js.CopyBytesToJS to copy the data into the array
//Then use the float32 array in the bufferData call

type scene struct {
	voxels        []*voxel
	gl            js.Value //Rendering context
	program       js.Value //shader program
	width, height int      //width and height of the canvas
	bufferData    []float32
	buffLen       int

	dataBuff    js.Value //data buffer
	dataBuffSet bool

	// projMat        js.TypedArray //Uniforms
	projMat    js.Value
	uProjMat   js.Value
	rawProjMat mgl32.Mat4
	// viewMat        js.TypedArray
	viewMat    js.Value
	uViewMat   js.Value
	rawViewMat mgl32.Mat4
	// modelMat       js.TypedArray
	modelMat       js.Value
	uModelMat      js.Value
	rawModelMat    mgl32.Mat4
	rawModelRotMat mgl32.Mat4

	aColor    js.Value //attributes
	aPosition js.Value
	aNormal   js.Value

	cameraLoc mgl32.Vec3
	cameraRot mgl32.Vec3

	update bool //does the data buffer need to be updated?
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
	v_Normal = vec3(u_Pmatrix*u_Vmatrix*u_Mmatrix*vec4(a_Normal, 0.0));
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
	to_light = vec3(0.57735026919, 0.57735026919, -0.57735026919);

	//amount of ambient light in the scene
	ambient_light = 0.3;

	//for difuse light calculation
	cos_angle = dot(v_Normal, to_light);
	cos_angle = clamp(cos_angle, 0.0, 1.0);

	gl_FragColor = vec4(v_Color*ambient_light + (v_Color*cos_angle) * (1.0-ambient_light), 1.0);
}
`

func newScene(canvas js.Value, color rgb) *scene {
	ret := scene{}
	ret.gl, ret.width, ret.height = getContext(canvas)

	vert := complieShader(ret.gl, glTypes.VertexShader, vShaderCode)
	frag := complieShader(ret.gl, glTypes.FragmentShader, fShaderCode)

	// ret.projMat = js.Global().Get("Float32Array").New(16)
	// ret.viewMat = js.Global().Get("Float32Array").New(16)
	// ret.modelMat = js.Global().Get("Float32Array").New(16)

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
	ret.setViewMat(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, -1.0}, mgl32.Vec3{0.0, 1.0, 0.0})
	ret.setModelMat(0, 0, 0, 0, 0, 0)
	ret.cameraLoc = mgl32.Vec3{0.0, 0.0, 0.0}

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

func (s *scene) buildBufferData() {
	if !s.update {
		return
	}

	s.update = false
	s.buffLen = 0
	for i, v := range s.voxels {
		for j, d := range v.data {
			spot := i*324 + j
			if spot >= len(s.bufferData) {
				s.bufferData = append(s.bufferData, d)
			} else {
				s.bufferData[spot] = d
			}
			s.buffLen++
		}
	}

	// Create a data buffer
	tArray := js.Global().Get("Float32Array").New(toIArray(s.bufferData)) //js.TypedArrayOf(s.bufferData)
	dataBuff := s.dataBuff
	if !s.dataBuffSet {
		dataBuff = s.gl.Call("createBuffer")
		s.dataBuff = dataBuff
		s.dataBuffSet = true
	}

	// Bind the data into the buffer
	float32Bytes := 4
	posSize := 3
	colSize := 3
	normSize := 3
	vertexByteSize := posSize*float32Bytes + colSize*float32Bytes + normSize*float32Bytes
	s.buffLen /= posSize + colSize + normSize

	s.gl.Call("bindBuffer", glTypes.ArrayBuffer, dataBuff)
	s.gl.Call("bufferData", glTypes.ArrayBuffer, tArray, glTypes.DynamicDraw)
	s.gl.Call("vertexAttribPointer", s.aPosition, posSize, glTypes.Float, false, vertexByteSize, 0)
	s.gl.Call("vertexAttribPointer", s.aColor, colSize, glTypes.Float, false, vertexByteSize, posSize*float32Bytes)
	s.gl.Call("vertexAttribPointer", s.aNormal, colSize, glTypes.Float, false, vertexByteSize, posSize*float32Bytes+colSize*float32Bytes)
}

func (s *scene) removeVoxel(index ...int) {
	l := len(index)
	for d, i := range index {
		s.voxels[i] = s.voxels[len(s.voxels)-(1+d)]
	}
	s.voxels = s.voxels[:len(s.voxels)-l]
	s.update = true
}

func (s *scene) addVoxel(voxels ...*voxel) []int {
	ret := []int{}
	for _, v := range voxels {
		s.voxels = append(s.voxels, v)
		ret = append(ret, len(s.voxels)-1)
	}
	s.update = true
	return ret
}

func (s *scene) moveCamera(x, y, z float32) {
	s.cameraLoc[0] += x
	s.cameraLoc[1] += y
	s.cameraLoc[2] += z
	s.setModelMat(-s.cameraLoc[0], -s.cameraLoc[1], -s.cameraLoc[2],
		-s.cameraRot[0], s.cameraRot[1], s.cameraRot[2])
}

func (s *scene) setCameraLoc(x, y, z float32) {
	s.cameraLoc[0] = x
	s.cameraLoc[1] = y
	s.cameraLoc[2] = z
	s.setModelMat(-s.cameraLoc[0], -s.cameraLoc[1], -s.cameraLoc[2],
		-s.cameraRot[0], s.cameraRot[1], s.cameraRot[2])
}

func (s *scene) rotateCamera(x, y, z float32) {
	s.cameraRot[0] += x
	s.cameraRot[1] += y
	s.cameraRot[2] += z
	s.setModelMat(-s.cameraLoc[0], -s.cameraLoc[1], -s.cameraLoc[2],
		-s.cameraRot[0], s.cameraRot[1], s.cameraRot[2])
}

func (s *scene) setCameraRot(x, y, z float32) {
	s.cameraRot[0] = x
	s.cameraRot[1] = y
	s.cameraRot[2] = z
	s.setModelMat(-s.cameraLoc[0], -s.cameraLoc[1], -s.cameraLoc[2],
		-s.cameraRot[0], s.cameraRot[1], s.cameraRot[2])
}

func (s *scene) setProjMat(deg float32) {
	// Generate a projection matrix
	ratio := float32(s.width) / float32(s.height)

	projMatrix := mgl32.Perspective(mgl32.DegToRad(deg), ratio, 1, 1000.0)

	// var projMatrixBuffer *[16]float32
	// projMatrixBuffer = (*[16]float32)(unsafe.Pointer(&projMatrix))
	// s.projMat = js.TypedArrayOf([]float32((*projMatrixBuffer)[:]))

	mat := [16]float32(projMatrix)
	s.projMat = js.Global().Get("Float32Array").New(toIArray(mat[:]))

	// ret.projMat = js.Global().Get("Float32Array").New(16)
	// ret.viewMat = js.Global().Get("Float32Array").New(16)
	// ret.modelMat = js.Global().Get("Float32Array").New(16)
	s.rawProjMat = projMatrix
}

func (s *scene) setViewMat(eye, center, up mgl32.Vec3) {
	// Generate a view matrix
	viewMatrix := mgl32.LookAtV(eye, center, up)

	mat := [16]float32(viewMatrix)
	s.viewMat = js.Global().Get("Float32Array").New(toIArray(mat[:]))
	// var viewMatrixBuffer *[16]float32
	// viewMatrixBuffer = (*[16]float32)(unsafe.Pointer(&viewMatrix))
	// s.viewMat = js.TypedArrayOf([]float32((*viewMatrixBuffer)[:]))
	s.rawViewMat = viewMatrix
}

func (s *scene) setModelMat(tranX, tranY, tranZ, rotX, rotY, rotZ float32) {
	// Generate a model matrix
	movMatrix := mgl32.HomogRotate3DX(rotX)
	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DY(rotY))
	movMatrix = movMatrix.Mul4(mgl32.HomogRotate3DZ(rotZ))
	s.rawModelRotMat = movMatrix
	movMatrix[12] = tranX
	movMatrix[13] = tranY
	movMatrix[14] = tranZ

	mat := [16]float32(movMatrix)
	s.modelMat = js.Global().Get("Float32Array").New(toIArray(mat[:]))
	// var modelMatrixBuffer *[16]float32
	// modelMatrixBuffer = (*[16]float32)(unsafe.Pointer(&movMatrix))
	//js.CopyBytesToJS(s.modelMat, toBytes(movMatrix))
	//s.modelMat = js.TypedArrayOf([]float32((*modelMatrixBuffer)[:]))
	s.rawModelMat = movMatrix
}

func (s *scene) render() {
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

func toIArray(data []float32) []interface{} {
	var ret []interface{}
	for _, f := range data {
		ret = append(ret, interface{}(f))
	}
	return ret
}
