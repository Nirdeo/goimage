package effects

import (
	"image"
	"image/color"
)

type SquareEffect struct {
	X, Y, Size int
	Color      color.Color
}

func (s *SquareEffect) Name() string        { return "Carré" }
func (s *SquareEffect) Description() string { return "Dessine un carré rempli" }
func (s *SquareEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}
	for y := s.Y; y < s.Y+s.Size && y < bounds.Max.Y; y++ {
		for x := s.X; x < s.X+s.Size && x < bounds.Max.X; x++ {
			if x >= bounds.Min.X && y >= bounds.Min.Y {
				result.Set(x, y, s.Color)
			}
		}
	}
	return result
}

type CircleEffect struct {
	CenterX, CenterY, Radius int
	Color                    color.Color
}

func (c *CircleEffect) Name() string        { return "Cercle" }
func (c *CircleEffect) Description() string { return "Dessine un cercle rempli" }
func (c *CircleEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}
	drawFilledCircle(result, c.CenterX, c.CenterY, c.Radius, c.Color)
	return result
}

func drawFilledCircle(img *image.RGBA, centerX, centerY, radius int, color color.Color) {
	bounds := img.Bounds()
	for y := centerY - radius; y <= centerY + radius; y++ {
		for x := centerX - radius; x <= centerX + radius; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx + dy*dy <= radius*radius {
				if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
					img.Set(x, y, color)
				}
			}
		}
	}
}

type TriangleEffect struct {
	X1, Y1, X2, Y2, X3, Y3 int
	Color                   color.Color
}

func (t *TriangleEffect) Name() string        { return "Triangle" }
func (t *TriangleEffect) Description() string { return "Dessine un triangle rempli" }
func (t *TriangleEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}
	drawFilledTriangle(result, t.X1, t.Y1, t.X2, t.Y2, t.X3, t.Y3, t.Color)
	return result
}

func drawFilledTriangle(img *image.RGBA, x1, y1, x2, y2, x3, y3 int, color color.Color) {
	bounds := img.Bounds()
	minX := min(min(x1, x2), x3)
	maxX := max(max(x1, x2), x3)
	minY := min(min(y1, y2), y3)
	maxY := max(max(y1, y2), y3)
	minX = max(minX, bounds.Min.X)
	maxX = min(maxX, bounds.Max.X-1)
	minY = max(minY, bounds.Min.Y)
	maxY = min(maxY, bounds.Max.Y-1)
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInTriangle(x, y, x1, y1, x2, y2, x3, y3) {
				img.Set(x, y, color)
			}
		}
	}
}

func pointInTriangle(px, py, x1, y1, x2, y2, x3, y3 int) bool {
	denom := float64((y2-y3)*(x1-x3) + (x3-x2)*(y1-y3))
	if denom == 0 {
		return false
	}
	w1 := float64((y2-y3)*(px-x3) + (x3-x2)*(py-y3)) / denom
	w2 := float64((y3-y1)*(px-x3) + (x1-x3)*(py-y3)) / denom
	w3 := 1.0 - w1 - w2
	return w1 >= 0 && w2 >= 0 && w3 >= 0
}

type LineEffect struct {
	X1, Y1, X2, Y2 int
	Color           color.Color
}

func (l *LineEffect) Name() string        { return "Ligne" }
func (l *LineEffect) Description() string { return "Dessine une ligne droite" }
func (l *LineEffect) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}
	drawLine(result, l.X1, l.Y1, l.X2, l.Y2, l.Color)
	return result
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, color color.Color) {
	bounds := img.Bounds()
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	var sx, sy int
	if x1 < x2 { sx = 1 } else { sx = -1 }
	if y1 < y2 { sy = 1 } else { sy = -1 }
	err := dx - dy
	for {
		if x1 >= bounds.Min.X && x1 < bounds.Max.X && y1 >= bounds.Min.Y && y1 < bounds.Max.Y {
			img.Set(x1, y1, color)
		}
		if x1 == x2 && y1 == y2 { break }
		e2 := 2 * err
		if e2 > -dy { err -= dy; x1 += sx }
		if e2 < dx { err += dx; y1 += sy }
	}
}

func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }
func abs(x int) int { if x < 0 { return -x }; return x } 