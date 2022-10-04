package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
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
	var images = [7]string{"background", "item", "body", "clothes", "hair", "eye", "hat"}
	decodedImages := make([]image.Image, len(images))

	for i, _ := range images {
		decodedImages[i] = openAndDecode("./imgs/1. background/" + "Daycity#20.png")
	}

	for _, trait := range traits {
		files, _ := ioutil.ReadDir("./imgs/" + trait)
		for _, file := range files {
			fmt.Println(file.Name())
		}
	}

	bounds := decodedImages[0].Bounds()
	newImage := image.NewRGBA(bounds)

	for _, img := range decodedImages {
		draw.Draw(newImage, img.Bounds(), img, image.ZP, draw.Over)
	}

	result, err := os.Create("result.png")
	if err != nil {
		log.Fatalf("Failed to create: %s", err)
	}

	// json_generator()

	png.Encode(result, newImage)
	defer result.Close()
}
