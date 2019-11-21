package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

const (
	frontFace = 0
	rightFace = 1
	backFace  = 2
	leftFace  = 3
	upFace    = 4
	downFace  = 5
)

type Voxel struct {
	x, y, z   float32
	color     [6]RGB
	prevColor [6]RGB
	selected  [6]bool
	data      [324]float32
	centers   [6]mgl32.Vec3
}

func (v *Voxel) intersect(r *Ray, closest float32) (int, float32, bool) {
	retFace := -1
	retLen := float32(0.0)
	retHit := false
	for i := 0; i < 6; i++ {
		var newR mgl32.Vec3
		center := v.centers[i]

		if i == frontFace || i == backFace {
			if r.dir[2] == 0 {
				continue
			}
			newZ := center[2] - r.orig[2]
			scale := newZ / r.dir[2]
			newR = mgl32.Vec3{r.dir[0] * scale, r.dir[1] * scale, newZ}
		}
		if i == leftFace || i == rightFace {
			if r.dir[0] == 0 {
				continue
			}
			newX := center[0] - r.orig[0]
			scale := newX / r.dir[0]
			newR = mgl32.Vec3{newX, r.dir[1] * scale, r.dir[2] * scale}
		}
		if i == upFace || i == downFace {
			if r.dir[1] == 0 {
				continue
			}
			newY := center[1] - r.orig[1]
			scale := newY / r.dir[1]
			newR = mgl32.Vec3{r.dir[0] * scale, newY, r.dir[2] * scale}
		}

		len := newR.Len()
		if len > closest || len <= 0 {
			continue
		}

		newDest := newR.Add(r.orig)

		if i == frontFace || i == backFace {
			if newDest[0] > center[0]-0.5 && newDest[0] < center[0]+0.5 &&
				newDest[1] > center[1]-0.5 && newDest[1] < center[1]+0.5 {
				retFace = i
				retLen = len
				closest = len
				retHit = true
				continue
			}
		}
		if i == leftFace || i == rightFace {
			if newDest[2] > center[2]-0.5 && newDest[2] < center[2]+0.5 &&
				newDest[1] > center[1]-0.5 && newDest[1] < center[1]+0.5 {
				retFace = i
				retLen = len
				closest = len
				retHit = true
				continue
			}
		}
		if i == upFace || i == downFace {
			if newDest[0] > center[0]-0.5 && newDest[0] < center[0]+0.5 &&
				newDest[2] > center[2]-0.5 && newDest[2] < center[2]+0.5 {
				retFace = i
				retLen = len
				closest = len
				retHit = true
				continue
				// return i, len, true
			}
		}
	}

	return retFace, retLen, retHit
}

func newVoxel(x, y, z float32, col [6]RGB) *Voxel {
	centers := [6]mgl32.Vec3{
		mgl32.Vec3{x + 0.5, y + 0.5, z},
		mgl32.Vec3{x + 1, y + 0.5, z - 0.5},
		mgl32.Vec3{x + 0.5, y + 0.5, z - 1},
		mgl32.Vec3{x, y + 0.5, z - 0.5},
		mgl32.Vec3{x + 0.5, y + 1, z - 0.5},
		mgl32.Vec3{x + 0.5, y, z - 0.5},
	}

	ret := Voxel{
		x:       x,
		y:       y,
		z:       z,
		color:   col,
		centers: centers,
	}

	ret.buildVertexData()
	return &ret
}

func (v *Voxel) setColor(rgb RGB, face ...int) {
	for i := 0; i < 6; i++ {
		for _, f := range face {
			if i == f {
				v.color[i] = rgb
			}
		}
	}

	v.buildVertexData()
}

func (v *Voxel) buildVertexData() {
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
		// loc  | color     | norm  |
		x, y, z, r1, g1, b1, 0, 0, 1,
		a, y, z, r1, g1, b1, 0, 0, 1,
		a, b, z, r1, g1, b1, 0, 0, 1, //1
		x, y, z, r1, g1, b1, 0, 0, 1,
		a, b, z, r1, g1, b1, 0, 0, 1,
		x, b, z, r1, g1, b1, 0, 0, 1, //2
		a, y, z, r2, g2, b2, 1, 0, 0,
		a, b, c, r2, g2, b2, 1, 0, 0,
		a, b, z, r2, g2, b2, 1, 0, 0, //3
		a, y, z, r2, g2, b2, 1, 0, 0,
		a, y, c, r2, g2, b2, 1, 0, 0,
		a, b, c, r2, g2, b2, 1, 0, 0, //4
		a, y, c, r3, g3, b3, 0, 0, -1,
		x, b, c, r3, g3, b3, 0, 0, -1,
		a, b, c, r3, g3, b3, 0, 0, -1, //5
		a, y, c, r3, g3, b3, 0, 0, -1,
		x, y, c, r3, g3, b3, 0, 0, -1,
		x, b, c, r3, g3, b3, 0, 0, -1, //6
		x, y, c, r4, g4, b4, -1, 0, 0,
		x, b, z, r4, g4, b4, -1, 0, 0,
		x, b, c, r4, g4, b4, -1, 0, 0, //7
		x, y, c, r4, g4, b4, -1, 0, 0,
		x, y, z, r4, g4, b4, -1, 0, 0,
		x, b, z, r4, g4, b4, -1, 0, 0, //8
		a, b, c, r5, g5, b5, 0, 1, 0,
		x, b, c, r5, g5, b5, 0, 1, 0,
		x, b, z, r5, g5, b5, 0, 1, 0, //9
		a, b, c, r5, g5, b5, 0, 1, 0,
		x, b, z, r5, g5, b5, 0, 1, 0,
		a, b, z, r5, g5, b5, 0, 1, 0, //10
		a, y, z, r6, g6, b6, 0, -1, 0,
		x, y, c, r6, g6, b6, 0, -1, 0,
		a, y, c, r6, g6, b6, 0, -1, 0, //11
		a, y, z, r6, g6, b6, 0, -1, 0,
		x, y, z, r6, g6, b6, 0, -1, 0,
		x, y, c, r6, g6, b6, 0, -1, 0, //12
	}
}

func (v *Voxel) newVoxelNeighbor(face int) *Voxel {
	x := v.x
	y := v.y
	z := v.z
	if face == frontFace {
		z++
	}
	if face == backFace {
		z--
	}
	if face == rightFace {
		x++
	}
	if face == leftFace {
		x--
	}
	if face == upFace {
		y++
	}
	if face == downFace {
		y--
	}

	fCol := v.color[face]
	if v.selected[face] {
		fCol = v.prevColor[face]
	}
	vox := newVoxel(x, y, z, newRGBSet(fCol.r, fCol.g, fCol.b))
	vox.selectFace(face)
	return vox
}

func (v *Voxel) selectFace(face int) {
	if v.selected[face] {
		return
	}
	v.selected[face] = true
	v.prevColor[face] = v.color[face]
	v.setColor(RGB{255, 0, 0}, face)
}

func (v *Voxel) deselectFace(face int) bool {
	if !v.selected[face] {
		return false
	}
	v.selected[face] = false
	v.setColor(v.prevColor[face], face)
	return true
}
