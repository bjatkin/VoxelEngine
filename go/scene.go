package main

import "syscall/js"

type Scene struct {
	voxels     []Voxel
	gl         js.Value //Rendering context
	program    js.Value //shader program
	update     bool
	bufferData []float32
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

func newScene(canvasId string) *Scene {
	ret := Scene{}
	ret.gl = getCanvas(canvasId)

	vert := complieShader(glTypes.VertexShader, vShaderCode)
	frag := complieShader(glTypes.FragmentShader, fShaderCode)

	ret.program = linkProgram(vert, frag)
	//TODO finish this. Should bind fragment uniforms and attributes
	return &ret
}

func (s *Scene) buildBufferData() {
	if !s.update {
		return
	}

	s.update = false
	s.bufferData = []float32{}
	for _, v := range s.voxels {
		for _, d := range v.data {
			s.bufferData = append(s.bufferData, d)
		}
	}
}

func (s *Scene) addVoxel(v Voxel) {
	s.voxels = append(s.voxels, v)
	s.update = true
}

func (s *Scene) render(now float64) {

}
