package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type RGB struct {
	r, g, b float32
}

func newRGB(r, g, b float32) RGB {
	return RGB{
		r: r,
		g: g,
		b: b,
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
		rgb.r / 255.0,
		rgb.g / 255.0,
		rgb.b / 255.0,
	}
}
