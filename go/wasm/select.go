package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

func hilightSelection(s *Scene, startCorner, endCorner mgl32.Vec2) {
	sel1, sel2 := false, false
	vox1, face1, intersect1 := intersectVoxel(s, startCorner)
	if intersect1 {
		sel1 = vox1.selectFace(face1)
	}

	vox2, face2, intersect2 := intersectVoxel(s, endCorner)
	if intersect2 {
		sel2 = vox2.selectFace(face2)
	}

	if sel1 || sel2 {
		s.update = true
	}
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
