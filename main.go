package main

import (
	"bufio"
	"fmt"
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

// Given a color pointer update cluster with given centroid
func setCentroid(color *Color, centroid Centroid) {
	color.Cluster = centroid
}

// multiply two points and return component values
func multiply(c1 Point, c2 Point) (uint32, uint32, uint32) {
	r1, g1, b1 := c1.GetCoords()
	r2, g2, b2 := c2.GetCoords()
	return r1 * r2, g1 * g2, b1 * b2
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

// Given a slice of colors, a centroid, and a function to compareCentroid
// a color and a centroid return a slice of colors that satisfy the function
func filter(vs []Color, cs Centroid, f func(Color, Centroid) bool) []Color {
	vsf := make([]Color, 0)
	for _, v := range vs {
		if f(v, cs) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Compare the values of a color's cluster to a given centroid
func compareCentroid(color Color, centroid Centroid) bool {
	if color.Cluster == centroid {
		return true
	}
	return false
}

// Given a slice of (filtered) colors and a pointer to a centroid recalc and update
// the location of the centroid
func calcPosition(colors []Color, centroid *Centroid) {
	var rt, gt, bt uint32
	total := uint32(len(colors))
	for _, color := range colors {
		r, g, b := color.GetCoords()
		rt += r
		gt += g
		bt += b
	}
	centroid.Red = rt / total
	centroid.Green = gt / total
	centroid.Blue = bt / total
}

func recalculateCentroids(colors *[]Color, centroids *[]Centroid) []Centroid {
	var oldCentroids []Centroid
	copy(oldCentroids, *centroids)
	for _, centroid := range *centroids {
		centroidColors := filter(*colors, centroid, compareCentroid)
		calcPosition(centroidColors, &centroid)
	}
	for _, color := range *colors {
		colorCentroid := nearestCentroid(color, *centroids)
		setCentroid(&color, colorCentroid)
	}
	return oldCentroids
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

func convergence(centroids []Centroid, oldCentroids []Centroid) bool {
	if len(oldCentroids) == 0 {
		return false
	}
	//var distances []float64
	for _, centroid := range centroids {
		for _, oldCentroid := range oldCentroids {
			//append(distances, distance(centroid, oldCentroid))
			if distance(centroid, oldCentroid) > 1.0 {
				return false
			}
		}
	}
	return true

}

// Given an image and a number of desired clusters generate a slice
// of colors
func convertImage(m image.Image, n int32) ([]Color, []Centroid) {
	centroids := generateCentroids(n)
	return addColors(m, centroids), centroids
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter image name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing name")
	}
	m := importImage(name)
	colors, centroids := convertImage(m, 6)
	var oldCentroids []Centroid
	i := 0

	for convergence(centroids, oldCentroids) || i < 20 {
		oldCentroids = recalculateCentroids(&colors, &centroids)
		i++
	}
	fmt.Printf("%#v", centroids)

}
