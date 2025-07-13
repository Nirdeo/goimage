package effects

import (
	"image"
	"image/color"
)

type ContrastEffect struct {
	Factor float64
}

func (c *ContrastEffect) Name() string        { return "Contraste" }
func (c *ContrastEffect) Description() string { return "Ajuste le contraste de l'image" }

func (c *ContrastEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			
			newR := uint8(clampContrast((float64(r>>8)-128)*c.Factor + 128))
			newG := uint8(clampContrast((float64(g>>8)-128)*c.Factor + 128))
			newB := uint8(clampContrast((float64(b>>8)-128)*c.Factor + 128))
			
			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}
	return result
}

func clampContrast(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
} 