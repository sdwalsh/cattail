package kmeans

// Iterate through centroids to deterimine which is closest using CIE94
import (
	"image"
	"math/rand"
	"os"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// Main loop of k-means
func (image *Image) Update() {
	image.recalculateColors()
	image.recalculateCentroids()
}

// Takes an image and a list of pregenerated centroids to generate
// a list of colors (points with an assigned cluster)
// bounds do not necessarily start at 0, use bounds.Min and bounds.Max instead
// colorful.Color takes float64 [0..1] divide by alpha pre-multiplied value provided by RGBA()
func addColors(m image.Image, centroids []*Centroid) *[]Color {
	bounds := m.Bounds()
	var colors []Color
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			color := Color{Color: colorful.Color{R: float64(r) / float64(65535), G: float64(g) / float64(65535), B: float64(b) / float64(65535)}, Cluster: nil}
			color.Cluster = color.nearestCentroid(centroids)
			colors = append(colors, color)
		}
	}
	return &colors
}

// Generate n number of centroids randomly initalized between 0 - 1
func generateCentroids(n int) []*Centroid {
	centroids := make([]*Centroid, 0)
	for i := 0; i < n; i++ {
		r := rand.Float64()
		g := rand.Float64()
		b := rand.Float64()
		centroids = append(centroids, &Centroid{Color: colorful.Color{R: r, G: g, B: b}})
	}
	return centroids
}

// Generate a single centroid
func generateCentroid() *Centroid {
	var centroid *Centroid
	r := rand.Float64()
	g := rand.Float64()
	b := rand.Float64()
	centroid = &Centroid{Color: colorful.Color{R: r, G: g, B: b}}
	return centroid
}

func (image *Image) containsEmptyCentroid() bool {
	for _, centroid := range image.Centroids {
		if centroid.count(image) == 0 {
			return true
		}
	}
	return false
}

// Given an image and a number of desired clusters generate a new Image
func convertImage(m image.Image, n int) *Image {
	// Generate components of an Image
	centroids := generateCentroids(n)
	colors := addColors(m, centroids)
	image := Image{ImportedImage: m, Colors: *colors, Centroids: centroids}

	for image.containsEmptyCentroid() {
		image.reroll()
		image.recalculateColors()
	}
	return &image
}

// Compare the values of a color's cluster to a given centroid
func compareCentroid(color Color, centroid *Centroid) bool {
	if *color.Cluster == *centroid {
		return true
	}
	return false
}

// Given a slice of colors, a centroid, and a function to compareCentroid
// a color and a centroid return a slice of colors that satisfy the function
func (m Image) filter(cs *Centroid, f func(Color, *Centroid) bool) []Color {
	vsf := make([]Color, 0)
	for _, v := range m.Colors {
		if f(v, cs) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Recalculate the centroids using average of l, a, b components of the colors
func (m *Image) recalculateCentroids() {
	for _, centroid := range m.Centroids {
		var lt, at, bt float64
		centroidColors := m.filter(centroid, compareCentroid)
		total := len(centroidColors)
		if total != 0.0 {
			for _, color := range centroidColors {
				l, a, b := color.Color.Lab()
				lt += l
				at += a
				bt += b
			}
			centroid.Color = colorful.Lab(lt/float64(total), at/float64(total), bt/float64(total))
		}
	}
}

func (image *Image) recalculateColors() {
	for _, color := range image.Colors {
		cluster := color.nearestCentroid(image.Centroids)
		*color.Cluster = *cluster
	}
}

/*
func (m *Image) convergence() bool {
	if len(oldCentroids) == 0 {
		return false
	}
	for i, centroid := range centroids {
		distance := oldCentroids[i].Color.DistanceCIE94(centroid.Color)
		if distance > 2.0 {
			return false
		}
	}
	return true
}
*/

// Reroll empty centroids in an image
func (image *Image) reroll() {
	var newCentroids []*Centroid
	for _, centroid := range image.Centroids {
		if centroid.isEmpty(image) {
			c := generateCentroid()
			newCentroids = append(newCentroids, c)
		} else {
			newCentroids = append(newCentroids, centroid)
		}
	}
	image.Centroids = newCentroids
}

// Given an image location, open image and return an image.Image
func importImage(filename string) (image.Image, error) {
	reader, err := os.Open(filename)
	defer reader.Close()
	if err != nil {
		return nil, err
	}

	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}
	return m, err
}
