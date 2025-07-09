package effects

import "image"

type Effect interface {
	Apply(img image.Image) image.Image
	Name() string
	Description() string
} 