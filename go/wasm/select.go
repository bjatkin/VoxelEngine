package main

import (
	"fmt"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type selection struct {
	startVox, endVox *Voxel
	face             int
	cubes            []*Voxel
	cubesLen         int
	init             bool
}

func (s *selection) emptySelection() {
	s.cubesLen = 0
	for _, c := range s.cubes {
		c.deselectFace(s.face)
	}
	fmt.Printf("deselect faces\n")
}

func (s *selection) addCube(v *Voxel) {
	len := len(s.cubes)
	if s.cubesLen >= len {
		s.cubes = append(s.cubes, v)
	} else {
		s.cubes[s.cubesLen] = v
	}
	s.cubesLen++
}

var currentSelection selection

func hilightSelection(s *Scene, startCorner, endCorner mgl32.Vec2) {
	if !currentSelection.init {
		currentSelection.cubes = []*Voxel{}
		currentSelection.init = true
	}

	vox1, face1, intersect1 := intersectVoxel(s, startCorner)
	currentSelection.face = face1
	if intersect1 {
		vox1.selectFace(face1)
	} else {
		return
	}

	vox2, face2, intersect2 := intersectVoxel(s, endCorner)
	if intersect2 {
		vox2.selectFace(face2)
	}

	if !intersect2 {
		return
	}

	if vox1 == currentSelection.startVox && vox2 == currentSelection.endVox {
		return
	}

	currentSelection.emptySelection()
	currentSelection.addCube(vox1)
	currentSelection.addCube(vox2)

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

	// l := len(currentSelection.cubes)
	// i := 0
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

			var add *Voxel
			for _, c := range s.voxels {
				if voxDist(c.x, c.y, c.z, cube[0], cube[1], cube[2]) <= 0.01 {
					add = c
					break
				}
			}

			currentSelection.addCube(add)
			// if i >= l {
			// 	currentSelection.cubes = append(currentSelection.cubes, add)
			// } else {
			// 	currentSelection.cubes[i] = add
			// }
			// i++
		}
	}
	// i++
	// if i >= l {
	// 	currentSelection.cubes = append(currentSelection.cubes, vox1)
	// } else {
	// 	currentSelection.cubes[i] = vox1
	// }
	// i++
	// if i >= l {
	// 	currentSelection.cubes = append(currentSelection.cubes, vox2)
	// } else {
	// 	currentSelection.cubes[i] = vox2
	// }

	// currentSelection.cubesLen = i
	for i, c := range currentSelection.cubes {
		if i >= currentSelection.cubesLen {
			break
		}
		c.selectFace(face1)
	}

	s.update = true
}

func clearSelection(s *Scene, startCorner, endCorner mgl32.Vec2) {
	del1, del2 := false, false
	vox1, face1, intersect1 := intersectVoxel(s, startCorner)
	if intersect1 {
		del1 = vox1.deselectFace(face1)
	}

	vox2, face2, intersect2 := intersectVoxel(s, endCorner)
	if intersect2 {
		del2 = vox2.deselectFace(face2)
	}

	if del1 || del2 {
		s.update = true
	}
}

func voxDist(x1, y1, z1, x2, y2, z2 float32) float64 {
	return math.Abs(float64(x1-x2)) +
		math.Abs(float64(y1-y2)) +
		math.Abs(float64(z1-z2))
}

func intersectVoxel(s *Scene, point mgl32.Vec2) (*Voxel, int, bool) {
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
		return s.voxels[sel], selFace, true
	}
	return &Voxel{}, -1, false
}
