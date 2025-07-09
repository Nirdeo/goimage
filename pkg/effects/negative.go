package effects

import (
	"image"
	"image/color"
)

type NegativeEffect struct{}

func (n *NegativeEffect) Name() string        { return "Négatif" }
func (n *NegativeEffect) Description() string { return "Inverse toutes les couleurs de l'image" }

func (n *NegativeEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			result.Set(x, y, color.RGBA{
				R: uint8(255 - uint8(r>>8)),
				G: uint8(255 - uint8(g>>8)),
				B: uint8(255 - uint8(b>>8)),
				A: uint8(a >> 8),
			})
		}
	}
	return result
} 