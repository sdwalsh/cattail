package kmeans

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"

	colorful "github.com/lucasb-eyer/go-colorful"
)

type Image struct {
	ImportedImage image.Image
	Colors        []Color
	Centroids     []*Centroid
	//OldCentroids  []Centroid
	//CentroidPop   map[*Centroid]int
}

func (m *Image) PrintCentroids() {
	for i, centroid := range m.Centroids {
		fmt.Printf("number: %v ", i)
		fmt.Printf("centroid: %v \n", centroid)
		centroidColors := m.filter(centroid, compareCentroid)
		total := len(centroidColors)
		fmt.Printf("count: %v \n", total)
	}
}

// Easy public combined function that will return an image upon completion of k-means
func CreateAndRun(filename string, centroids int, iterations int) (*Image, error) {
	image, err := Create(filename, centroids)
	if err != nil {
		return nil, err
	}
	image.Run(iterations)
	return image, nil
}

// Create a new Image given a file name and number of centroids
func Create(filename string, n int) (*Image, error) {
	rand.Seed(time.Now().UTC().UnixNano())
	m, err := importImage(filename)
	if err != nil {
		return nil, err
	}
	return convertImage(m, n), nil
}

// FUTURE :: replace with convergence test
func (m *Image) Run(maxIterations int) {
	for i := 0; i < maxIterations; i++ {
		m.Update()
	}
}

// FUTURE :: order colors in swatch by most common
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
