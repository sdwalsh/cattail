package kmeans

import (
	"image"
	_ "image/jpeg"

	"github.com/lucasb-eyer/go-colorful"
)

type Image struct {
	ImportedImage image.Image
	Colors        []Color
	Centroids     []*Centroid
	//OldCentroids  []Centroid
	//CentroidPop   map[*Centroid]int
}

type Color struct {
	Color   colorful.Color
	Cluster *Centroid
}

type Centroid struct {
	Color colorful.Color
}
