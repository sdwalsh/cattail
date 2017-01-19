package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sdwalsh/cattail/kmeans"
)

func readTrim(reader *bufio.Reader, s string, e error) string {
	fmt.Printf(s, e)
	c, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing")
	}
	t := strings.TrimSpace(c)
	return t
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter image name: ")
	c, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing name")
	}
	t := strings.TrimSpace(c)

	rand.Seed(time.Now().UTC().UnixNano())

	m, err := kmeans.ImportImage(t)
	for err != nil {
		t = readTrim(reader, "Error! (%v) Try entering an image again: ", err)
		m, err = kmeans.ImportImage(t)
	}
	fmt.Printf("image imported\n")

	fmt.Print("enter number of centroids: ")
	c, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing number")
	}
	t = strings.TrimSpace(c)
	nCentroids, err := strconv.Atoi(t)
	for err != nil {
		t = readTrim(reader, "Error! (%v) Try entering a number again: ", err)
		nCentroids, err = strconv.Atoi(t)
	}

	fmt.Print("enter max number of iterations: ")
	c, err = reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing number")
	}
	t = strings.TrimSpace(c)
	nIterations, err := strconv.Atoi(t)
	for err != nil {
		t = readTrim(reader, "Error! (%v) Try entering a number again: ", err)
		nIterations, err = strconv.Atoi(t)
	}

	fmt.Println("\nbegin timer")
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)

	colors, centroids := kmeans.ConvertImage(m, nCentroids)

	fmt.Printf("image converted\n")
	fmt.Printf("colors size :%v \n", len(colors))
	fmt.Println("/////////////////////////////////////////////////")
	for i, centroid := range centroids {
		fmt.Printf("centroid #%v: ", i)
		fmt.Printf("%v \n", len(kmeans.Filter(colors, centroid, kmeans.CompareCentroid)))
	}
	fmt.Println("/////////////////////////////////////////////////")
	var oldCentroids []kmeans.Centroid

	fmt.Printf("begin convergence tests\n")
	fmt.Println("-------------------------------------------------")
	for i := 0; i < nIterations; i++ {
		fmt.Printf("loop #%v\n", i)
		t1 := time.Now()

		// replace with copy
		for _, centroid := range centroids {
			oldCentroids = append(oldCentroids, *centroid)
		}

		centroids = kmeans.RecalculateCentroids(colors, centroids)
		colors = kmeans.RecalculateColors(colors, centroids)

		t2 := time.Now()
		fmt.Printf("time: %v \n", t2.Sub(t1))
		for i, centroid := range centroids {
			fmt.Printf("centroid #%v: ", i)
			fmt.Printf("%v \t", len(kmeans.Filter(colors, centroid, kmeans.CompareCentroid)))
			fmt.Printf("%v \n", centroid.Color.Hex())
		}
		fmt.Println("\n/////////////////////////////////////////////////")
	}
	fmt.Println("begin rendering")
	err = kmeans.CreateColorImage(m, centroids)
	if err != nil {
		fmt.Printf("Error! (%v) cannot create image", err)
	}
	err = kmeans.CreateColorSwatch(centroids)
	if err != nil {
		fmt.Printf("Error! (%v) cannot create swatch", err)
	}
	endTime := time.Now()
	fmt.Printf("total time to complete: %v \n", endTime.Sub(startTime))
}
