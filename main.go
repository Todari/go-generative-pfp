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
	"sort"
	"strconv"
	"strings"
	"time"

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

type Trait struct {
	Rarity string
	Path   string
	Name   string
	Weight int
	Group  string
}

func main() {

	var traits = [7]string{"1. background", "2. item", "3. body", "4. clothes", "5. hair", "6. eye", "7. hat"}
	totalNum := 1000
	var dnaArr []string

	//trait이름 별 weight를 map으로 만들어 줌
	weightForTrait := make(map[string]int)

	traitsArr := make([][]string, len(traits))
	images := make([]string, len(traits))
	dnaExist := false
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for i, trait := range traits {
		files, _ := ioutil.ReadDir("./imgs/" + trait)
		for _, file := range files {
			traitsArr[i] = append(traitsArr[i], file.Name())
			traitName := strings.Split(file.Name(), "#")[0]
			traitWeight, _ := strconv.Atoi(strings.Split(strings.Split(file.Name(), "#")[1], ".")[0])
			weightForTrait[traitName] = traitWeight

		}
	}

	for i := 0; i < totalNum; i++ {
		var dna string

		for i, _ := range traitsArr {
			rand.Seed(time.Now().UnixNano())
			selecter := rand.Intn(len(traitsArr[i]))
			images[i] = traitsArr[i][selecter]
			var dna2 string
			if selecter < 16 {
				dna2 = "0" + strconv.FormatInt(int64(selecter), 16)
			} else if selecter >= 16 {
				dna2 = strconv.FormatInt(int64(selecter), 16)
			}
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
		fmt.Println(i, dna)
		// fmt.Fprintln(w, dnaArr)

		decodedImages := make([]image.Image, len(images))

		for i, v := range images {
			decodedImages[i] = openAndDecode("./imgs/" + traits[i] + "/" + v)
		}

		bounds := decodedImages[0].Bounds()
		newImage := image.NewRGBA(bounds)

		for _, img := range decodedImages {
			draw.Draw(newImage, img.Bounds(), img, image.ZP, draw.Over)
		}

		result, err := os.Create("./result/" + strconv.Itoa(i) + ".png")
		if err != nil {
			log.Fatalf("Failed to create: %s", err)
		}

		png.Encode(result, newImage)

		defer result.Close()

		sort.Strings(dnaArr)

	}

	defer fmt.Println(dnaArr)
}
