package kmeans

import colorful "github.com/lucasb-eyer/go-colorful"

// Color type is a RGBA point (of a pixel) and the centroid nearest to the color
type Color struct {
	Color   colorful.Color
	Cluster *Centroid
}

// Getters and Setters
func (c Color) getColor() colorful.Color {
	return c.Color
}

func (c Color) getCluster() *Centroid {
	return c.Cluster
}

func (c *Color) setColor(color colorful.Color) {
	c.Color = color
}

func (c *Color) setCluster(cluster *Centroid) {
	c.Cluster = cluster
}

// CIE94 provides distance method that is closer to human perception
// Delta E	Perception
// <= 1.0	Not perceptible by human eyes.
// 1 - 2	Perceptible through close observation.
// 2 - 10	Perceptible at a glance.
// 11 - 49	Colors are more similar than opposite
// 100	Colors are exact opposite
// http://zschuessler.github.io/DeltaE/learn/
func (c Color) nearestCentroid(centroids []*Centroid) *Centroid {
	// set lowest distance to max value (opposites)
	lowestDistance := 100.0
	var address *Centroid
	for _, centroid := range centroids {
		distance := c.Color.DistanceCIE94(centroid.Color)
		if distance < lowestDistance {
			lowestDistance = distance
			address = centroid
		}
	}
	return address
}
