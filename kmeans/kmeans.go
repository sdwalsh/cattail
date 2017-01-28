package kmeans

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// Image type is a decomposed image
type Image struct {
	ImportedImage image.Image
	Colors        []Color
	Centroids     []*Centroid
	//OldCentroids  []Centroid
	//CentroidPop   map[*Centroid]int
}

// PrintCentroids prints out the current centroids to Standard Output
func (m *Image) PrintCentroids() {
	for i, centroid := range m.Centroids {
		fmt.Printf("number: %v ", i)
		fmt.Printf("centroid: %v \n", centroid)
		centroidColors := m.filter(centroid, compareCentroid)
		total := len(centroidColors)
		fmt.Printf("count: %v \n", total)
	}
}

// CreateAndRun is a combination of the Create and Run functions. It runs k-means
// on the supplied image - provided by name in a relative to the executable when
// given the number of centroids and iterations desired
func CreateAndRun(filename string, centroids int, iterations int) (*Image, error) {
	image, err := Create(filename, centroids)
	if err != nil {
		return nil, err
	}
	image.Run(iterations)
	return image, nil
}

// Create returns a new image when provided with a filename and the number of
// centroids desired
func Create(filename string, n int) (*Image, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	m, err := importImage(filename)
	if err != nil {
		return nil, err
	}
	return convertImage(m, n), nil
}

// Run runs the update function until convergence or max iterations is reached
func (m *Image) Run(maxIterations int) {
	for i := 0; i < maxIterations; i++ {
		m.Update()
	}
}

// CreateColorSwatch creates a color swatch from the Image type provided
// the color swatch is a color list of the centroids (each centroid is 60px x 60px)
func (m *Image) CreateColorSwatch() error {
	length := len(m.Centroids)

	// new image that will serve as the generated picture (x0, y0, x1, y1)
	img := image.NewRGBA(image.Rect(0, 0, 60, length*60))

	for i, centroid := range m.Centroids {
		draw.Draw(img, image.Rect(0, i*60, 60, ((i+1)*60)), &image.Uniform{centroid.Color}, image.ZP, draw.Src)
	}

	toimg, err := os.Create("colorswatch.png")
	if err != nil {
		return err
	}
	defer toimg.Close()
	err = png.Encode(toimg, img)
	if err != nil {
		return err
	}

	return nil
}

// CreateColorImage recolors the original image provided using each pixel's nearest
// centroid
func (m *Image) CreateColorImage() error {
	bounds := m.ImportedImage.Bounds()
	x0 := 0
	y0 := 0
	height := bounds.Max.Y - bounds.Min.Y
	width := bounds.Max.X - bounds.Min.X

	img := image.NewRGBA(image.Rect(x0, y0, width, height))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.ImportedImage.At(x, y).RGBA()
			color := Color{Color: colorful.Color{R: float64(r) / float64(65535), G: float64(g) / float64(65535), B: float64(b) / float64(65535)}, Cluster: nil}
			cluster := color.nearestCentroid(m.Centroids)
			draw.Draw(img, image.Rect(x0, y0, x0+1, y0+1), &image.Uniform{cluster.Color}, image.ZP, draw.Src)
			x0++
		}
		x0 = 0
		y0++
	}

	toimg, err := os.Create("colorblend.png")
	if err != nil {
		return err
	}
	defer toimg.Close()
	err = png.Encode(toimg, img)
	if err != nil {
		return err
	}

	return nil
}

// Update is the main loop of the k-means algorithm
func (m *Image) Update() {
	m.recalculateCentroids()
	m.recalculateColors()
}

// Takes an image and a list of pregenerated centroids to generate
// a list of colors (points with an assigned cluster)
// bounds do not necessarily start at 0, use bounds.Min and bounds.Max instead
// colorful.Color takes float64 [0..1] divide by alpha pre-multiplied value provided by RGBA()
func addColors(m image.Image, centroids []*Centroid) []Color {
	bounds := m.Bounds()
	var colors []Color
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			color := Color{Color: colorful.Color{R: float64(r) / float64(65535), G: float64(g) / float64(65535), B: float64(b) / float64(65535)}, Cluster: nil}
			color.setCluster(color.nearestCentroid(centroids))
			colors = append(colors, color)
		}
	}
	return colors
}

// Generate n number of centroids randomly initialized between 0 - 1
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

func (m *Image) containsEmptyCentroid() bool {
	for _, centroid := range m.Centroids {
		if centroid.count(m) == 0 {
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
	image := Image{ImportedImage: m, Colors: colors, Centroids: centroids}

	for image.containsEmptyCentroid() {
		image.reroll()
		image.recalculateColors()
	}
	return &image
}

// Compare the values of a color's cluster to a given centroid
func compareCentroid(color Color, centroid *Centroid) bool {
	if color.Cluster == centroid {
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
			centroid.setColor(colorful.Lab(lt/float64(total), at/float64(total), bt/float64(total)))
		}
	}
}

// recalculateColors recalculates the nearestCentroid (replaces centroid address)
func (m *Image) recalculateColors() {
	var newColors []Color
	for _, color := range m.Colors {
		cluster := color.nearestCentroid(m.Centroids)
		newColors = append(newColors, Color{Color: color.Color, Cluster: cluster})
	}
	m.Colors = newColors
}

// reroll rerolls empty centroids in an image
func (m *Image) reroll() {
	var c []*Centroid
	for _, centroid := range m.Centroids {
		if centroid.isEmpty(m) {
			c = append(c, generateCentroid())
		} else {
			c = append(c, centroid)
		}
	}
	m.Centroids = c
}

// importImage given an image location, open image and return an image.Image
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
