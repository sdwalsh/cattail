package main

import (
	"bufio"
	"fmt"
	"log"
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

func importImage(reader *bufio.Reader, centroids int) (*kmeans.Image, error) {
	fmt.Print("enter image name: ")
	c, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing name")
	}
	t := strings.TrimSpace(c)

	m, err := kmeans.Create(t, centroids)
	for err != nil {
		t = readTrim(reader, "Error! (%v) Try entering an image again: ", err)
		m, err = kmeans.Create(t, centroids)
	}
	fmt.Printf("image imported\n")
	return m, nil
}

func nIterations(reader *bufio.Reader) (int, error) {
	fmt.Print("enter max number of iterations: ")
	c, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing number")
	}
	t := strings.TrimSpace(c)
	nIterations, err := strconv.Atoi(t)
	for err != nil {
		t = readTrim(reader, "Error! (%v) Try entering a number again: ", err)
		nIterations, err = strconv.Atoi(t)
	}
	return nIterations, err
}

func nCentroids(reader *bufio.Reader) (int, error) {
	fmt.Print("enter number of centroids: ")
	c, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal("Error parsing number")
	}
	t := strings.TrimSpace(c)
	nCentroids, err := strconv.Atoi(t)
	for err != nil {
		t = readTrim(reader, "Error! (%v) Try entering a number again: ", err)
		nCentroids, err = strconv.Atoi(t)
	}
	return nCentroids, err
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	nIterations, err := nIterations(reader)
	nCentroids, err := nCentroids(reader)
	m, err := importImage(reader, nCentroids)

	fmt.Println("\nbegin timer")
	startTime := time.Now()
	fmt.Printf("start time: %v \n", startTime)

	for i := 0; i < nIterations; i++ {
		fmt.Printf("Loop: #%v \n", i)
		m.PrintCentroids()
		m.Update()
	}

	fmt.Println("begin rendering")
	err = m.CreateColorImage()
	if err != nil {
		fmt.Printf("Error! (%v) cannot create image", err)
	}
	err = m.CreateColorSwatch()
	if err != nil {
		fmt.Printf("Error! (%v) cannot create swatch", err)
	}
	endTime := time.Now()
	fmt.Printf("total time to complete: %v \n", endTime.Sub(startTime))
}
