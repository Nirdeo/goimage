package effects

import (
	"image"
	"image/color"
)

type BrightnessEffect struct {
	Factor float64
}

func (br *BrightnessEffect) Name() string        { return "Luminosité" }
func (br *BrightnessEffect) Description() string { return "Ajuste la luminosité de l'image" }

func (br *BrightnessEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			
			newR := uint8(clamp(float64(r>>8) * br.Factor))
			newG := uint8(clamp(float64(g>>8) * br.Factor))
			newB := uint8(clamp(float64(b>>8) * br.Factor))
			
			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}
	return result
}

func clamp(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
} 