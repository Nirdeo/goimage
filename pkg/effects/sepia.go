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
			gray := 0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)
			sepiaR := uint8(min(255, int(0.393*gray+0.769*gray+0.189*gray)))
			sepiaG := uint8(min(255, int(0.349*gray+0.686*gray+0.168*gray)))
			sepiaB := uint8(min(255, int(0.272*gray+0.534*gray+0.131*gray)))
			result.Set(x, y, color.RGBA{R: sepiaR, G: sepiaG, B: sepiaB, A: uint8(a >> 8)})
		}
	}
	return result
}

 