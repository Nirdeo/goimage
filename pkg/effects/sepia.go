package effects

import (
	"image"
	"image/color"
)

type SepiaEffect struct{}

func (s *SepiaEffect) Name() string        { return "Sépia" }
func (s *SepiaEffect) Description() string { return "Applique un filtre sépia vintage" }

func (s *SepiaEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)
			
			sepiaR := uint8(min(255, int(0.393*r8+0.769*g8+0.189*b8)))
			sepiaG := uint8(min(255, int(0.349*r8+0.686*g8+0.168*b8)))
			sepiaB := uint8(min(255, int(0.272*r8+0.534*g8+0.131*b8)))
			
			result.Set(x, y, color.RGBA{R: sepiaR, G: sepiaG, B: sepiaB, A: uint8(a >> 8)})
		}
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

 