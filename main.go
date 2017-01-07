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

func Distance(c1 Point, c2 Point) float64 {
	r1, g1, b1 := c1.GetCoords()
	r2, g2, b2 := c2.GetCoords()
	r := r1 - r2
	g := g1 - g2
	b := b1 - b2
	return math.Sqrt(float64(r*r + g*g + b*b))
}

func NearestCentroid(c1 Point, centroids []Point) Point {
	lowestDistance := float64(0.0)
	index := -1
	for i, centroid := range centroids {
		distance := Distance(c1, centroid)
		if distance <= lowestDistance {
			lowestDistance = distance
			index = i
		}
	}
	return centroids[index]
}

func generateCentroids(r int32) []Point {
	centroids := make([]Centroid, r)
	for i, _ := range centroids {
		r := uint32(rand.Intn(65535))
		g := uint32(rand.Intn(65535))
		b := uint32(rand.Intn(65535))
		centroids[i] = Centroid{RGBPoint{Red: r, Green: g, Blue: b}}
	}
	return centroids
}

func main() {
	reader, err := os.Open("testdata.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	centroids := generateCentroids(6)

	// Add colors to slice
	colors := make([]Color, (bounds.Min.Y * bounds.Min.X))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			cluster := NearestCentroid(RGBPoint{
				Red:   r,
				Green: g,
				Blue:  b,
			},
				[]Point(centroids))
			colors[y*x] = Color{
				RGBPoint: RGBPoint{
					Red:   r,
					Green: g,
					Blue:  b,
				},
				Cluster: cluster,
			}
		}
	}
}
