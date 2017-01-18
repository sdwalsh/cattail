package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/lucasb-eyer/go-colorful"
)

type Color struct {
	Color   colorful.Color
	Cluster *Centroid
}

type Centroid struct {
	Color colorful.Color
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

// Iterate through centroids to deterimine which is closest using CIE94
// CIE94 provides distance method that is closer to human perception
// Delta E	Perception
// <= 1.0	Not perceptible by human eyes.
// 1 - 2	Perceptible through close observation.
// 2 - 10	Perceptible at a glance.
// 11 - 49	Colors are more similar than opposite
// 100	Colors are exact opposite
// http://zschuessler.github.io/DeltaE/learn/
func nearestCentroid(c1 colorful.Color, centroids []*Centroid) *Centroid {
	// set lowest distance to max value (opposites)
	lowestDistance := 100.0
	var address *Centroid
	for _, centroid := range centroids {
		distance := c1.DistanceCIE94(centroid.Color)
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
			color := colorful.Color{R: float64(r) / float64(65535), G: float64(g) / float64(65535), B: float64(b) / float64(65535)}
			cluster := nearestCentroid(color, centroids)
			colors = append(colors, Color{Color: color, Cluster: cluster})
		}
	}
	return colors
}

// Generate n number of centroids between 0 - 1
func generateCentroids(n int32) []*Centroid {
	centroids := make([]*Centroid, 0)
	for i := int32(0); i < n; i++ {
		r := rand.Float64()
		g := rand.Float64()
		b := rand.Float64()
		centroids = append(centroids, &Centroid{Color: colorful.Color{R: r, G: g, B: b}})
	}
	if len(centroids) != 0 {
		return centroids
	} else {
		log.Fatal("centroids incorrectly generated")
		return centroids
	}
}

// Given an image and a number of desired clusters generate a slice of colors and
// a slice of centroid pointers
func convertImage(m image.Image, n int32) ([]Color, []*Centroid) {
	centroids := generateCentroids(n)
	return addColors(m, centroids), centroids
}

// Given a slice of colors, a centroid, and a function to compareCentroid
// a color and a centroid return a slice of colors that satisfy the function
func filter(vs []Color, cs *Centroid, f func(Color, *Centroid) bool) []Color {
	vsf := make([]Color, 0)
	for _, v := range vs {
		if f(v, cs) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

// Compare the values of a color's cluster to a given centroid
func compareCentroid(color Color, centroid *Centroid) bool {
	if *color.Cluster == *centroid {
		return true
	}
	return false
}

// Recalculate the centroids using average of l, a, b components of the colors
func recalculateCentroids(colors []Color, centroids []*Centroid) {
	for _, centroid := range centroids {
		var lt, at, bt float64
		centroidColors := filter(colors, centroid, compareCentroid)
		//fmt.Printf("filtered colors:: %v \n", len(centroidColors))
		total := float64(len(centroidColors))
		if total == 0.0 {
			fmt.Println("------------EMPTY-----------")
		} else {
			for _, color := range centroidColors {
				l, a, b := color.Color.Lab()
				lt += l
				at += a
				bt += b
			}
			if centroid.Color != colorful.Lab(lt/total, at/total, bt/total) {
				fmt.Println("------------UPDATED-----------")
				fmt.Printf("old:: %v \t", centroid.Color.Hex())
				fmt.Printf("new:: %v \t", colorful.Lab(lt/total, at/total, bt/total).Hex())
			}
			centroid.Color = colorful.Lab(lt/total, at/total, bt/total)
			fmt.Printf("current:: %v \n", centroid.Color.Hex())
			fmt.Println("------------END-----------")
			//fmt.Printf("%v \n", colorful.Lab(lt/total, at/total, bt/total).Hex())
		}
	}
}

// Alter colors in place (colors is a very, very large slice - copying to a new
// slice would be costly)
func recalculateColors(colors []Color, centroids []*Centroid) []Color {
	var newColors []Color
	for _, color := range colors {
		cluster := nearestCentroid(color.Color, centroids)
		newColors = append(newColors, Color{Color: color.Color, Cluster: cluster})
	}
	return newColors
}

func createColorImage(m image.Image, centroids []*Centroid) {
	bounds := m.Bounds()
	x_0 := 0
	y_0 := 0
	height := bounds.Max.Y - bounds.Min.Y
	width := bounds.Max.X - bounds.Min.X

	// new image that will serve as the generated picture (x0, y0, x1, y1)
	img := image.NewRGBA(image.Rect(x_0, y_0, width, height))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			color := colorful.Color{R: float64(r) / float64(65535), G: float64(g) / float64(65535), B: float64(b) / float64(65535)}
			cluster := nearestCentroid(color, centroids)
			draw.Draw(img, image.Rect(x_0, y_0, x_0+1, y_0+1), &image.Uniform{cluster.Color}, image.ZP, draw.Src)
			x_0++
		}
		x_0 = 0
		y_0++
	}

	toimg, err := os.Create("colorblend.png")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	defer toimg.Close()

	png.Encode(toimg, img)
}

func convergence(centroids []*Centroid, oldCentroids []Centroid) bool {
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

// Reroll the random position of a centroid (useful if empty)
func (c *Centroid) reroll(colors []Color) {

}

func main() {
	//reader := bufio.NewReader(os.Stdin)
	//fmt.Print("enter image name: ")
	//name, err := reader.ReadString('\n')
	//if err != nil {
	//	log.Fatal("Error parsing name")
	//}
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("begin timer")
	startTime := time.Now()

	m := importImage("./lexington.jpg")
	fmt.Printf("image imported\n")

	colors, centroids := convertImage(m, 16)
	// detect and remove empty centroids
	/*for _, centroid := range centroids {
		// reroll centroid until it at least one color
		for len(filter(colors, centroid, compareCentroid)) != 0 {
		}
	} */
	fmt.Printf("image converted\n")
	fmt.Printf("colors size :%v \n", len(colors))
	fmt.Println("/////////////////////////////////////////////////")
	for i, centroid := range centroids {
		fmt.Printf("centroid #%v: ", i)
		fmt.Printf("%v \n", len(filter(colors, centroid, compareCentroid)))
	}
	fmt.Println("/////////////////////////////////////////////////")
	var oldCentroids []Centroid

	i := 0
	fmt.Printf("begin convergence tests\n")
	fmt.Println("-------------------------------------------------")
	for !convergence(centroids, oldCentroids) && i < 12 {
		fmt.Printf("loop #%v\n", i)
		t1 := time.Now()

		for _, centroid := range centroids {
			oldCentroids = append(oldCentroids, *centroid)
		}

		recalculateCentroids(colors, centroids)
		colors = recalculateColors(colors, centroids)

		t2 := time.Now()
		fmt.Printf("time: %v \n", t2.Sub(t1))
		for i, centroid := range centroids {
			fmt.Printf("centroid #%v: ", i)
			fmt.Printf("%v \t", len(filter(colors, centroid, compareCentroid)))
			fmt.Printf("%v \n", centroid.Color.Hex())
		}
		fmt.Println("\n/////////////////////////////////////////////////")
		i++
	}
	fmt.Println("begin rendering")
	createColorImage(m, centroids)
	endTime := time.Now()
	//fmt.Printf("%+v \n", centroids)
	fmt.Printf("total time to complete: %v \n", endTime.Sub(startTime))
}
