package main

import (
	"bufio"
	"fmt"
	_ "image/jpeg"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/sdwalsh/cattail/kmeans"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("enter image name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing name")
	}
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("\nbegin timer")
	startTime := time.Now()

	m := kmeans.ImportImage(name)
	fmt.Printf("image imported\n")

	fmt.Print("enter number of centroids: ")
	c, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing number")
	}
	i, err := strconv.Atoi(c)
	if err != nil {
		log.Fatal(err)
	}

	colors, centroids := kmeans.ConvertImage(m, i)

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
	for i := 0; !kmeans.Convergence(centroids, oldCentroids) && i < 12; i++ {
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
	kmeans.CreateColorImage(m, centroids)
	endTime := time.Now()
	//fmt.Printf("%+v \n", centroids)
	fmt.Printf("total time to complete: %v \n", endTime.Sub(startTime))
}
