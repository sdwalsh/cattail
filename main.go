package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
)

type Point interface {
	GetCoords() (uint32, uint32, uint32)
}

type RGBPoint struct {
	Red   uint32
	Green uint32
	Blue  uint32
}

func (rgb RGBPoint) GetCoords() (uint32, uint32, uint32) {
	return rgb.Red, rgb.Green, rgb.Blue
}

type Color struct {
	RGBPoint
	Cluster Centroid
}

type Centroid struct {
	RGBPoint
}

func setCentroid(color *Color, centroid Centroid) {
	color.Cluster = centroid
}

func multiply(c1 Point, c2 Point) (uint32, uint32, uint32) {
	r1, g1, b1 := c1.GetCoords()
	r2, g2, b2 := c2.GetCoords()
	r := r1 * r2
	g := g1 * g2
	b := b1 * b2

	return r, g, b
}

// Find distance between two points
func distance(c1 Point, c2 Point) float64 {
	r1, g1, b1 := c1.GetCoords()
	r2, g2, b2 := c2.GetCoords()
	r := r1 - r2
	g := g1 - g2
	b := b1 - b2
	return math.Sqrt(float64(r*r + g*g + b*b))
}

// Iterate through centroids to deterimine which is closest
// O(n) running time
func nearestCentroid(c1 Point, centroids []Centroid) Centroid {
	lowestDistance := float64(0.0)
	index := -1
	for i, centroid := range centroids {
		distance := distance(c1, centroid)
		if distance <= lowestDistance {
			lowestDistance = distance
			index = i
		}
	}
	return centroids[index]
}

// Generate n number of centroids
func generateCentroids(n int32) []Centroid {
	centroids := make([]Centroid, n)
	for i, _ := range centroids {
		r := uint32(rand.Intn(65535))
		g := uint32(rand.Intn(65535))
		b := uint32(rand.Intn(65535))
		centroids[i] = Centroid{RGBPoint{Red: r, Green: g, Blue: b}}
	}
	return centroids
}

// Takes an image and a list of pregenerated centroids to generate
// a list of colors - points with an assigned cluster
func addColors(m image.Image, centroids []Centroid) []Color {
	bounds := m.Bounds()
	var colors []Color
	i := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			cluster := nearestCentroid(RGBPoint{
				Red:   r,
				Green: g,
				Blue:  b,
			},
				centroids)
			colors = append(colors, Color{
				RGBPoint: RGBPoint{
					Red:   r,
					Green: g,
					Blue:  b,
				},
				Cluster: cluster,
			})
		}
	}
	return colors
}

func filter(vs *[]Color, cs Centroid, f func(Color, Centroid) bool) []Color {
	vsf := make([]Color, 0)
	for _, v := range *vs {
		if f(v, cs) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func compareCentroid(color Color, centroid Centroid) bool {
	if color.Cluster == centroid {
		return true
	}
	return false
}

func calcPosition(colors []Color, centroid *Centroid) *Centroid {
	var rt, gt, bt uint32
	total := uint32(len(colors))
	for i, color := range colors {
		r, g, b := color.GetCoords()
		rt += r
		gt += g
		bt += b
	}
	centroid.Red = rt / total
	centroid.Green = gt / total
	centroid.Blue = bt / total

	return centroid
}

func recalculateCentroids(colors *[]Color, centroids *[]Centroid) ([]Color, []Centroid) {
	for _, centroid := range *centroids {
		centroidColors := filter(colors, centroid, compareCentroid)
		calcPosition(*colors, &centroid)
	}
	for _, color := range *colors {
		colorCentroid := nearestCentroid(color, *centroids)
		setCentroid(&color, colorCentroid)
	}
}

// Given an image location, open image and return an image.Image
func importImage(i string) image.Image {
	reader, err := os.Open(i)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	return m
}

// Given an image and a number of desired clusters generate a slice
// of colors
func convertImage(m image.Image, n int32) []Color {
	centroids := generateCentroids(n)
	return addColors(m, centroids)
}

func main() {
	m := importImage("testimage.png")
	colors := convertImage(m, 6)

}
