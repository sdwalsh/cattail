package main

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
)

type Point struct {
	Red   uint32
	Green uint32
	Blue  uint32
}

type Color struct {
	Point
}

type Centroid struct {
	Point
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

	centroids := make([]Centroid, 6)
	for i, _ := range centroids {
		r := uint32(rand.Intn(65535))
		g := uint32(rand.Intn(65535))
		b := uint32(rand.Intn(65535))

		centroids[i] = Centroid{Point{Red: r, Green: g, Blue: b}}
	}

	colors := make([]Color, (bounds.Min.Y * bounds.Min.X))
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			colors[y*x] = Color{Point{Red: r, Green: g, Blue: b}}
		}
	}
}
