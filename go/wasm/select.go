package main

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type selection struct {
	face             int
	startVox, endVox int
	allVox           []int
	len              int
}

func (sel *selection) voxels() []int {
	return sel.allVox[:sel.len]
}

func (sel *selection) deselectAll(s *scene) {
	for _, v := range sel.voxels() {
		s.voxels[v].deselectFace(sel.face)
	}
	sel.len = 0
}

func (sel *selection) selectAll(s *scene, color rgb) {
	for _, v := range sel.voxels() {
		s.voxels[v].selectFace(sel.face, color)
	}
}

func (sel *selection) colorSelection(s *scene, color rgb) {
	for _, v := range sel.voxels() {
		s.voxels[v].selectFaceColor(color)
	}
}

func (sel *selection) color(s *scene, color rgb) {
	for _, v := range sel.voxels() {
		s.voxels[v].colorSelectedFace(color)
	}
}

func (sel *selection) addVox(v int) {
	len := len(sel.allVox)
	if sel.len >= len {
		sel.allVox = append(sel.allVox, v)
	} else {
		sel.allVox[sel.len] = v
	}
	sel.len++
}

func (sel *selection) newSelection(s *scene, start, end mgl32.Vec2, color rgb) {
	v1, f1, s1 := intersectVoxel(s, start)
	if !s1 {
		sel.deselectAll(s)
		s.update = true
		return
	}

	v2, f2, s2 := intersectVoxel(s, end)
	if !s2 || f1 != f2 {
		return
	}

	if v1 == sel.startVox && v2 == sel.endVox {
		return
	}

	sel.deselectAll(s)
	sel.startVox = v1
	sel.endVox = v2
	sel.face = f1

	vox1 := s.voxels[v1]
	vox2 := s.voxels[v2]
	diffX := vox1.x - vox2.x
	diffY := vox1.y - vox2.y
	diffZ := vox1.z - vox2.z
	dir1 := diffX
	dir2 := diffY

	if diffX == 0 {
		dir1 = diffY
		dir2 = diffZ
	}
	if diffY == 0 {
		dir1 = diffX
		dir2 = diffZ
	}
	if diffZ == 0 {
		dir1 = diffX
		dir2 = diffY
	}

	//Build new selection
	for x := 0.0; x < math.Abs(float64(dir1))+0.1; x++ {
		for y := 0.0; y < math.Abs(float64(dir2))+0.1; y++ {
			nx, ny := x, y
			if dir1 > 0 {
				nx = -x
			}
			if dir2 > 0 {
				ny = -y
			}

			cube := mgl32.Vec3{vox1.x, vox1.y + float32(nx), vox1.z + float32(ny)}
			if diffY == 0 {
				cube = mgl32.Vec3{vox1.x + float32(nx), vox1.y, vox1.z + float32(ny)}
			}
			if diffZ == 0 {
				cube = mgl32.Vec3{vox1.x + float32(nx), vox1.y + float32(ny), vox1.z}
			}

			var add int
			for i, c := range s.voxels {
				if voxDist(c.x, c.y, c.z, cube[0], cube[1], cube[2]) <= 0.01 {
					add = i
					break
				}
			}

			sel.addVox(add)
		}
	}

	sel.selectAll(s, color)
	s.update = true
}

func (sel *selection) hilightSelection(s *scene, color rgb) {
	for i, v := range sel.allVox {
		if i >= sel.len {
			break
		}
		s.voxels[v].selectFace(sel.face, color)
	}
}

func (sel *selection) isEmpty() bool {
	return len(sel.allVox) == 0
}

func newSelection() selection {
	return selection{
		allVox: []int{},
	}
}

func voxDist(x1, y1, z1, x2, y2, z2 float32) float64 {
	return math.Abs(float64(x1-x2)) +
		math.Abs(float64(y1-y2)) +
		math.Abs(float64(z1-z2))
}

func intersectVoxel(s *scene, point mgl32.Vec2) (int, int, bool) {
	closest := float32(99999999.0)
	r := newRay(s, point[0], point[1])
	sel := -1
	selFace := 0
	for i, v := range s.voxels {
		face, dist, hit := v.intersect(&r, closest)
		if hit && dist < closest {
			closest = dist
			sel = i
			selFace = face
		}
	}

	if sel >= 0 {
		return sel, selFace, true
	}
	return -1, -1, false
}

func faceToShift(face int) (float32, float32, float32) {
	switch face {
	case 0:
		return 0, 0, -1
	case 1:
		return -1, 0, 0
	case 2:
		return 0, 0, 1
	case 3:
		return 1, 0, 0
	case 4:
		return 0, -1, 0
	case 5:
		return 0, 1, 0
	default:
		return 0, 0, 0
	}
}
