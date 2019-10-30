package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Voxel struct {
	x, y, z float32
	color   [6]RGB
	data    [324]float32
}

func newVoxel(x, y, z float32, col [6]RGB) *Voxel {
	ret := Voxel{
		x:     x,
		y:     y,
		z:     z,
		color: col,
	}

	ret.buildData()
	return &ret
}

func (v *Voxel) buildData() {
	x, y, z := v.x, v.y, v.z
	a, b, c := x+1.0, y+1.0, z-1.0

	var col [6]mgl32.Vec3
	for i, c := range v.color {
		col[i] = c.vec3()
	}
	r1, g1, b1 := col[0][0], col[0][1], col[0][2]
	r2, g2, b2 := col[1][0], col[1][1], col[1][2]
	r3, g3, b3 := col[2][0], col[2][1], col[2][2]
	r4, g4, b4 := col[3][0], col[3][1], col[3][2]
	r5, g5, b5 := col[4][0], col[4][1], col[4][2]
	r6, g6, b6 := col[5][0], col[5][1], col[5][2]

	v.data = [324]float32{
		x, y, z,
		r1, g1, b1,
		0, 0, 1,
		a, y, z,
		r1, g1, b1,
		0, 0, 1,
		a, b, z,
		r1, g1, b1,
		0, 0, 1, //1
		x, y, z,
		r1, g1, b1,
		0, 0, 1,
		a, b, z,
		r1, g1, b1,
		0, 0, 1,
		x, b, z,
		r1, g1, b1,
		0, 0, 1, //2
		a, y, z,
		r2, g2, b2,
		1, 0, 0,
		a, b, c,
		r2, g2, b2,
		1, 0, 0,
		a, b, z,
		r2, g2, b2,
		1, 0, 0, //3
		a, y, z,
		r2, g2, b2,
		1, 0, 0,
		a, y, c,
		r2, g2, b2,
		1, 0, 0,
		a, b, c,
		r2, g2, b2,
		1, 0, 0, //4
		a, y, c,
		r3, g3, b3,
		0, 0, -1,
		x, b, c,
		r3, g3, b3,
		0, 0, -1,
		a, b, c,
		r3, g3, b3,
		0, 0, -1, //5
		a, y, c,
		r3, g3, b3,
		0, 0, -1,
		x, y, c,
		r3, g3, b3,
		0, 0, -1,
		x, b, c,
		r3, g3, b3,
		0, 0, -1, //6
		x, y, c,
		r4, g4, b4,
		-1, 0, 0,
		x, b, z,
		r4, g4, b4,
		-1, 0, 0,
		x, b, c,
		r4, g4, b4,
		-1, 0, 0, //7
		x, y, c,
		r4, g4, b4,
		-1, 0, 0,
		x, y, z,
		r4, g4, b4,
		-1, 0, 0,
		x, b, z,
		r4, g4, b4,
		-1, 0, 0, //8
		a, b, c,
		r5, g5, b5,
		0, 1, 0,
		x, b, c,
		r5, g5, b5,
		0, 1, 0,
		x, b, z,
		r5, g5, b5,
		0, 1, 0, //9
		a, b, c,
		r5, g5, b5,
		0, 1, 0,
		x, b, z,
		r5, g5, b5,
		0, 1, 0,
		a, b, z,
		r5, g5, b5,
		0, 1, 0, //10
		a, y, z,
		r6, g6, b6,
		0, -1, 0,
		x, y, c,
		r6, g6, b6,
		0, -1, 0,
		a, y, c,
		r6, g6, b6,
		0, -1, 0, //11
		a, y, z,
		r6, g6, b6,
		0, -1, 0,
		x, y, z,
		r6, g6, b6,
		0, -1, 0,
		x, y, c,
		r6, g6, b6,
		0, -1, 0, //12
	}
}
