package main

type Voxel struct {
	x, y, z float32
	color   [6][3]float32
	data    [216]float32
}

func newVoxel(x, y, z float32, col [6][3]float32) *Voxel {
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
	col := v.color
	r1, g1, b1 := col[0][0], col[0][1], col[0][2]
	r2, g2, b2 := col[1][0], col[1][1], col[1][2]
	r3, g3, b3 := col[2][0], col[2][1], col[2][2]
	r4, g4, b4 := col[3][0], col[3][1], col[3][2]
	r5, g5, b5 := col[4][0], col[4][1], col[4][2]
	r6, g6, b6 := col[5][0], col[5][1], col[5][2]

	v.data = [216]float32{
		x, y, z,
		r1, g1, b1,
		a, y, z,
		r1, g1, b1,
		a, b, z,
		r1, g1, b1, //1
		x, y, z,
		r1, g1, b1,
		a, b, z,
		r1, g1, b1,
		x, b, z,
		r1, g1, b1, //2
		a, y, z,
		r2, g2, b2,
		a, b, c,
		r2, g2, b2,
		a, b, z,
		r2, g2, b2, //3
		a, y, z,
		r2, g2, b2,
		a, y, c,
		r2, g2, b2,
		a, b, c,
		r2, g2, b2, //4
		a, y, c,
		r3, g3, b3,
		x, b, c,
		r3, g3, b3,
		a, b, c,
		r3, g3, b3, //5
		a, y, c,
		r3, g3, b3,
		x, y, c,
		r3, g3, b3,
		x, b, c,
		r3, g3, b3, //6
		x, y, c,
		r4, g4, b4,
		x, b, z,
		r4, g4, b4,
		x, b, c,
		r4, g4, b4, //7
		x, y, c,
		r4, g4, b4,
		x, y, z,
		r4, g4, b4,
		x, b, z,
		r4, g4, b4, //8
		a, b, c,
		r5, g5, b5,
		x, b, c,
		r5, g5, b5,
		x, b, z,
		r5, g5, b5, //9
		a, b, c,
		r5, g5, b5,
		x, b, z,
		r5, g5, b5,
		a, b, z,
		r5, g5, b5, //10
		a, y, z,
		r6, g6, b6,
		x, y, c,
		r6, g6, b6,
		a, y, c,
		r6, g6, b6, //11
		a, y, z,
		r6, g6, b6,
		x, y, z,
		r6, g6, b6,
		x, y, c,
		r6, g6, b6, //12
	}
}
