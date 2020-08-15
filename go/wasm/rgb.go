package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

type rgb struct {
	R, G, B float32
}

func newRGB(r, g, b float32) rgb {
	return rgb{
		R: r,
		G: g,
		B: b,
	}
}

func newRGBSet(r, g, b float32) [6]rgb {
	return [6]rgb{
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
		newRGB(r, g, b),
	}
}

func (rgb *rgb) vec3() mgl32.Vec3 {
	return mgl32.Vec3{
		rgb.R / 255.0,
		rgb.G / 255.0,
		rgb.B / 255.0,
	}
}
