[![Build Status](https://travis-ci.org/sdwalsh/cattail.svg?branch=master)](https://travis-ci.org/sdwalsh/cattail)
[![Go Report Card](https://goreportcard.com/badge/github.com/sdwalsh/cattail)](https://goreportcard.com/report/github.com/sdwalsh/cattail)
[![GoDoc](https://godoc.org/github.com/sdwalsh/cattail/kmeans?status.svg)](https://godoc.org/github.com/sdwalsh/cattail/kmeans)

cattail
=======
categorize colors in an image using k-means clustering

### This project was created in order to better learn golang. Lacks idiomatic golang and best practices.

Algorithm / Structure
---------------------
### k-means clustering
- `n` centroids are randomly generated in RGB color space
- Images are converted to a slice of colors (see below) which includes a pointer to the nearest centroid (when centroids do not contain colors they are rerolled)
- While convergence or `nth` iteration has not occurred
    - Using euclidean distance and CIE L\*a\*b\* color space recalculate centroids
    - Update `[]Color` for new centroids (CIE94 distance)
- Once convergence or `nth` loops has occurred render loop over image and update colors with nearest centroid color

### Types
```go
// Color, Centroid, and Image types
type Image struct {
	ImportedImage image.Image
	Colors        []Color
	Centroids     []*Centroid
	OldCentroids  []Centroid                // To be implemented
	CentroidPop   map[*Centroid]int         // To be implemented
}
type Color struct {
	Color   colorful.Color
	Cluster *Centroid
}

type Centroid struct {
	Color colorful.Color
}
```

### Image generation
```go
// Recolor image using centroid colors
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
```

Future development
------------------

- Fix comparison of colors and centroids -- O(n^2) current
- Speed recoloring of image
- Create convergence method
- Allow choice of CIE94, CIEDE2000


Libraries used
--------------
- [**colorful**](https://github.com/lucasb-eyer/go-colorful)
