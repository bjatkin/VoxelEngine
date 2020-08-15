package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type ray struct {
	orig mgl32.Vec3
	dir  mgl32.Vec3
}

func newRay(s *scene, x, y float32) ray {
	// orig := s.cameraLoc
	orig := s.rawModelMat.Inv().Mul4x1(mgl32.Vec4{0, 0, 0, 1})
	newX := 2*x/float32(s.width) - 1
	newY := -2*y/float32(s.height) + 1

	dir := mgl32.Vec4{newX, newY, -1, 1}

	inv := s.rawProjMat.Mul4(s.rawViewMat).Inv()
	dir = s.rawModelRotMat.Inv().Mul4x1(inv.Mul4x1(dir))
	return Ray{
		orig: orig.Vec3(),
		dir:  dir.Vec3(),
	}
}
