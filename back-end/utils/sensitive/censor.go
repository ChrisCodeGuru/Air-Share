package sensitive

import (
	"image"
	"image/color"
)

type UserImg struct {
	image.Image
	custom map[image.Point]color.Color
}

func NewUserImg(img image.Image) *UserImg {
	return &UserImg{img, map[image.Point]color.Color{}}
}

func (m *UserImg) Set(x, y int, c color.Color) {
	m.custom[image.Point{x, y}] = c
}

func (m *UserImg) At(x, y int) color.Color {
	// Explicitly changed part: custom colors of the changed pixels:
	if c := m.custom[image.Point{x, y}]; c != nil {
		return c
	}
	// Unchanged part: colors of the original image:
	return m.Image.At(x, y)
}
