package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"sort"
)

func max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(x int, y int) int {
	if x < y {
		return x
	}
	return y
}

func bucket(vals []int, numBuckets int) [][]int {
	minv := vals[0]
	maxv := vals[0]
	for _, val := range vals {
		minv = min(minv, val)
		maxv = max(maxv, val)
	}

	var buckets [][]int
	buckets = make([][]int, numBuckets)
	for _, val := range vals {
		bucket := (val - minv) / (maxv - minv + 1)
		buckets[bucket] = append(buckets[bucket], val)
	}

	return buckets
}

type manipulator func([]int) int

// Treatment 1
func meanManip(vals []int) int {
	sum := 0
	for _, val := range vals {
		sum += (val / len(vals))
	}

	//log.Printf("%v -> %v", vals, sum)

	return sum
}

// Treatment 2
func medianManip(vals []int) int {
	sort.Ints(vals)
	return vals[len(vals)/2]
}

// Treatment 3
func modalBucketHighManip(vals []int) int {
	buckets := bucket(vals, 4)
	maxBucket := 0
	maxLen := len(buckets[0])
	for i, buck := range buckets {
		if len(buck) > maxLen {
			maxBucket = i
			maxLen = len(buck)
		}
	}

	return buckets[maxBucket][rand.Intn(len(buckets[maxBucket]))]
}

// Treatment 4
func modalBucketLowManip(vals []int) int {
	buckets := bucket(vals, 4)
	minBucket := 0
	minLen := len(buckets[0])
	for i, buck := range buckets {
		if (len(buck) > 0 && len(buck) < minLen) || minLen == 0 {
			minBucket = i
			minLen = len(buck)
		}
	}

	return buckets[minBucket][rand.Intn(len(buckets[minBucket]))]
}

// Treatment 5
func triangleHighManip(vals []int) int {
	triangle := [][]int{
		[]int{1},
		[]int{1, 1},
		[]int{1, 2, 1},
		[]int{1, 2, 2, 1},
		[]int{1, 2, 3, 2, 1},
		[]int{1, 2, 3, 3, 2, 1},
		[]int{1, 2, 3, 4, 3, 2, 1},
		[]int{1, 2, 3, 4, 4, 3, 2, 1},
		[]int{1, 2, 3, 4, 5, 4, 3, 2, 1},
		[]int{1, 2, 3, 4, 5, 5, 4, 3, 2, 1},
	}
	weights := triangle[len(vals)-1]
	weightedSum := 0
	sumWeights := 0

	for i := range vals {
		weightedSum += weights[i] * vals[i]
		sumWeights += weights[i]
	}

	return weightedSum / sumWeights
}

// Treatment 6
func triangleLowManip(vals []int) int {
	triangle := [][]int{
		[]int{1},
		[]int{1, 1},
		[]int{2, 1, 2},
		[]int{2, 1, 1, 2},
		[]int{3, 2, 1, 2, 3},
		[]int{3, 2, 1, 1, 2, 3},
		[]int{4, 3, 2, 1, 2, 3, 4},
		[]int{4, 3, 2, 1, 1, 2, 3, 4},
		[]int{5, 4, 3, 2, 1, 2, 3, 4, 5},
		[]int{5, 4, 3, 2, 1, 1, 2, 3, 4, 5},
	}
	weights := triangle[len(vals)-1]
	weightedSum := 0
	sumWeights := 0

	for i := range vals {
		weightedSum += weights[i] * vals[i]
		sumWeights += weights[i]
	}

	return weightedSum / sumWeights
}

// Treatment 7
func minManip(vals []int) int {
	retVal := vals[0]
	for _, val := range vals {
		if val < retVal {
			retVal = val
		}
	}
	return retVal
}

// Treatment 8
func maxManip(vals []int) int {
	retVal := vals[0]
	for _, val := range vals {
		if val > retVal {
			retVal = val
		}
	}
	return retVal
}

func makeFinalImage(images []string, fn manipulator, mult float64) image.Image {

	log.Printf("WEIGHT %v", mult)
	loadedImages := make([]image.Image, len(images))
	for i, name := range images {
		log.Printf("LOADING %v", name)
		reader, err := os.Open(name)
		if err != nil {
			panic(err)
		}
		img, err := png.Decode(reader)
		if err != nil {
			panic(err)
		}
		loadedImages[i] = img
	}

	//Produce the compressed image
	medImage := image.NewGray(loadedImages[0].Bounds())
	for x := 0; x <= loadedImages[0].Bounds().Dx(); x++ {
		for y := 0; y <= loadedImages[0].Bounds().Dy(); y++ {
			var pixVals []int
			for _, img := range loadedImages {
				gVal := color.GrayModel.Convert(img.At(x, y))
				r, _, _, _ := gVal.RGBA()
				pixVals = append(pixVals, int(float64(r)*mult))
			}
			colorVal := color.Gray16{Y: uint16(fn(pixVals))}
			medImage.Set(x, y, colorVal)
		}
	}

	return medImage
}

func getFileName(num int) string {
	zeros := ""
	if num < 10 {
		zeros += "0"
	}
	if num < 100 {
		zeros += "0"
	}
	if num < 1000 {
		zeros += "0"
	}
	return fmt.Sprintf("output_%v%v.png", zeros, num)
}

func makeLastFrame(directory string, starts []int, fn manipulator, name string) {
	i := len(starts) - 1
	for frame := 0; frame <= 299; frame++ {
		var imagePath []string
		for j := 0; j <= i; j++ {
			imagePath = append(imagePath, fmt.Sprintf("%s/%s", directory, getFileName(starts[j]+frame)))
		}
		newImage := makeFinalImage(imagePath, fn, float64(1.0))
		writer, err := os.Create(fmt.Sprintf("testout-%s-%v.png", name, frame+(i*300)))
		if err != nil {
			panic(err)
		}
		png.Encode(writer, newImage)
		writer.Close()

		os.Exit(1)
	}
}

func fadeInFrames(numFrames int, directory string) {
	for i := 0; i < numFrames; i++ {
		imagePath := []string{fmt.Sprintf("%s/%s", directory, getFileName(844))}
		newImage := makeFinalImage(imagePath, meanManip, float64((float64(i))/float64(numFrames)))
		writer, err := os.Create(fmt.Sprintf("final-%v.png", i))
		if err != nil {
			panic(err)
		}
		png.Encode(writer, newImage)
		writer.Close()
	}
}

func produceVideo(directory string, starts []int, fn manipulator, name string) {
	frameNum := 0
	for i := range starts {
		for frame := 30 * 5; frame > 0; frame-- {
			var imagePath []string
			for j := 9 - i; j >= 0; j-- {
				imagePath = append(imagePath, fmt.Sprintf("%s/%s", directory, getFileName(starts[j]+frame)))
			}
			newImage := makeFinalImage(imagePath, fn, float64(1.0))
			log.Printf("Writing %v (%v)", fmt.Sprintf("testout-%s-%v.png", name, frameNum), i)
			writer, err := os.Create(fmt.Sprintf("testout-%s-%v.png", name, frameNum))
			if err != nil {
				panic(err)
			}
			png.Encode(writer, newImage)
			writer.Close()
			frameNum++
		}
	}
}

func main() {
	framestarts := []int{844, 1738, 2641, 3350, 4071, 4798, 5510, 6232, 6953, 7671}
	produceVideo("/Users/simon/movie/set1/", framestarts, medianManip, "set8")
	//fadeInFrames(30*2, "/Users/simon/movie/set1/")
}
