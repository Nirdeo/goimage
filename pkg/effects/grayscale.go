package effects

import (
	"image"
	"image/color"
)

type GrayscaleEffect struct{}

func (g *GrayscaleEffect) Name() string        { return "Niveaux de gris" }
func (g *GrayscaleEffect) Description() string { return "Convertit l'image en niveaux de gris" }

func (g *GrayscaleEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			gray := uint8(0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8))
			result.Set(x, y, color.RGBA{
				R: gray,
				G: gray,
				B: gray,
				A: uint8(a >> 8),
			})
		}
	}
	return result
} 