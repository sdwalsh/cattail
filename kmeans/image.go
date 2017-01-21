package kmeans

// Iterate through centroids to deterimine which is closest using CIE94
import (
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"os"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// CIE94 provides distance method that is closer to human perception
// Delta E	Perception
// <= 1.0	Not perceptible by human eyes.
// 1 - 2	Perceptible through close observation.
// 2 - 10	Perceptible at a glance.
// 11 - 49	Colors are more similar than opposite
// 100	Colors are exact opposite
// http://zschuessler.github.io/DeltaE/learn/
func (color Color) nearestCentroid(centroids []*Centroid) *Centroid {
	// set lowest distance to max value (opposites)
	lowestDistance := 100.0
	var address *Centroid
	for _, centroid := range centroids {
		distance := color.Color.DistanceCIE94(centroid.Color)
		if distance < lowestDistance {
			lowestDistance = distance
			address = centroid
		}
	}
	return address
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
			color.Cluster = color.nearestCentroid(centroids)
			colors = append(colors, color)
		}
	}
	return colors
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
func convertImage(m image.Image, n int) Image {
	// Generate components of an Image
	centroids := generateCentroids(n)
	colors := addColors(m, centroids)
	image := Image{ImportedImage: m, Colors: colors, Centroids: centroids}

	for image.containsEmptyCentroid() {
		image.reroll()
		image.recalculateColors()
	}
	return image
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

// Compare the values of a color's cluster to a given centroid
func CompareCentroid(color Color, centroid *Centroid) bool {
	if *color.Cluster == *centroid {
		return true
	}
	return false
}

// Recalculate the centroids using average of l, a, b components of the colors
func (m Image) recalculateCentroids() {
	for _, centroid := range m.Centroids {
		var lt, at, bt float64
		centroidColors := m.filter(centroid, CompareCentroid)
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

func CreateColorSwatch(centroids []*Centroid) error {
	length := len(centroids)

	// new image that will serve as the generated picture (x0, y0, x1, y1)
	img := image.NewRGBA(image.Rect(0, 0, 60, length*60))

	for i, centroid := range centroids {
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

// FUTURE :: colors should be properly ordered already in the slice
// remove need to reprocess image
func (m *Image) CreateColorImage() error {
	bounds := m.ImportedImage.Bounds()
	x_0 := 0
	y_0 := 0
	height := bounds.Max.Y - bounds.Min.Y
	width := bounds.Max.X - bounds.Min.X

	// new image that will serve as the generated picture (x0, y0, x1, y1)
	img := image.NewRGBA(image.Rect(x_0, y_0, width, height))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.ImportedImage.At(x, y).RGBA()
			color := Color{Color: colorful.Color{R: float64(r) / float64(65535), G: float64(g) / float64(65535), B: float64(b) / float64(65535)}, Cluster: nil}
			cluster := color.nearestCentroid(m.Centroids)
			draw.Draw(img, image.Rect(x_0, y_0, x_0+1, y_0+1), &image.Uniform{cluster.Color}, image.ZP, draw.Src)
			x_0++
		}
		x_0 = 0
		y_0++
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

func Convergence(centroids []*Centroid, oldCentroids []Centroid) bool {
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

// Given an image, a centroid, and a function to compareCentroid
// a color and a centroid return a slice of colors that satisfy the function
func (cs *Centroid) count(image *Image) int {
	var count int
	for _, v := range image.Colors {
		if CompareCentroid(v, cs) {
			count++
		}
	}
	return count
}

func (centroid *Centroid) isEmpty(image *Image) bool {
	return centroid.count(image) == 0
}

// Main loop of k-means
//filter(vs []Color, cs *Centroid, f func(Color, *Centroid) bool)

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
