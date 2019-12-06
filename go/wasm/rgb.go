package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type RGB struct {
	R, G, B float32
}

func newRGB(r, g, b float32) RGB {
	return RGB{
		R: r,
		G: g,
		B: b,
	}
}

func newRGBSet(r, g, b float32) [6]RGB {
	return [6]RGB{
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
	}
}

func (rgb *RGB) vec3() mgl32.Vec3 {
	return mgl32.Vec3{
		rgb.R / 255.0,
		rgb.G / 255.0,
		rgb.B / 255.0,
	}
}
