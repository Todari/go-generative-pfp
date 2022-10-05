package main

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/Todari/go-generative-pfp/module"
)

func openAndDecode(imgPath string) image.Image {
	img, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("Failed to open %s", err)
	}

	decoded, err := png.Decode(img)
	if err != nil {
		log.Fatalf("Failed to decode %s", err)
	}
	defer img.Close()

	return decoded
}

func main() {
	// tierNum := [7]int{10, 100, 700, 1200, 3000, 4990}
	// traitGroupNum := 7

	var traits = [7]string{"1. background", "2. item", "3. body", "4. clothes", "5. hair", "6. eye", "7. hat"}
	totalNum := 1000
	var dnaArr []string
	traitsArr := make([][]string, len(traits))
	images := make([]string, len(traits))
	dnaExist := false
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for i, trait := range traits {
		files, _ := ioutil.ReadDir("./imgs/" + trait)
		for _, file := range files {
			traitsArr[i] = append(traitsArr[i], file.Name())
		}
	}

	for i := 0; i < totalNum; i++ {
		var dna string

		for i, _ := range traitsArr {
			selecter := rand.Intn(len(traitsArr[i]))
			images[i] = traitsArr[i][selecter]
			dna2 := strconv.FormatInt(int64(selecter), 16)
			dna += dna2
		}
		for _, v := range dnaArr {
			if dna == v {
				fmt.Fprintln(w, "DNA EXIST")
				dnaExist = true
			}
		}

		if dnaExist == true {
			i--
			continue
		}

		module.Json_generator(images, i)

		dnaArr = append(dnaArr, dna)
		fmt.Fprintln(w, dna)
		fmt.Fprintln(w, dnaArr)

		decodedImages := make([]image.Image, len(images))

		for i, _ := range images {
			decodedImages[i] = openAndDecode("./imgs/" + traits[i] + "/" + images[i])
		}

		bounds := decodedImages[0].Bounds()
		newImage := image.NewRGBA(bounds)

		for _, img := range decodedImages {
			draw.Draw(newImage, img.Bounds(), img, image.ZP, draw.Over)
		}

		result, err := os.Create(strconv.Itoa(i) + ".png")
		if err != nil {
			log.Fatalf("Failed to create: %s", err)
		}

		png.Encode(result, newImage)

		defer result.Close()
	}
}
