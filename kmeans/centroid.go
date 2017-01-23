package kmeans

import colorful "github.com/lucasb-eyer/go-colorful"

type Centroid struct {
	Color colorful.Color
}

func (c *Centroid) setColor(color colorful.Color) {
	c.Color = color
}

// Given an image, a centroid, and a function to compareCentroid
// a color and a centroid return a slice of colors that satisfy the function
func (c *Centroid) count(image *Image) int {
	var count int
	for _, color := range image.Colors {
		if compareCentroid(color, c) {
			count++
		}
	}
	return count
}

func (centroid *Centroid) isEmpty(image *Image) bool {
	return centroid.count(image) == 0
}
