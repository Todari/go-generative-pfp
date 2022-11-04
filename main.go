package main

import (
	"bufio"
	"encoding/csv"
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

func reset() {
	if !os.IsExist(os.RemoveAll("./result/")) {
		os.Mkdir("./result/", 0777)
	}
	if !os.IsExist(os.RemoveAll("./json/")) {
		os.Mkdir("./json/", 0777)
	}
}

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

	doReset := true
	totalNum := 1000

	if doReset {
		reset()
	}
	var rarities = [6]string{"1. legend", "2. prime", "3. master", "4. expert", "5. junior", "6. rookie"}
	var traits = [11]string{"0. legend", "1. background", "2. body", "3. outfit", "4. backpack", "5. pet", "6. hat", "7. eye", "8. mouth", "9. ring", "10. rarity"}
	var dnaArr []string

	//trait이름 별 weight를 map으로 만들어 줌
	// weightForTrait := make(map[string]int)

	var traitsArr [6][11][]string
	rarityCounter := make([]int, 6)
	images := make([]string, len(traits))
	dnaExist := false
	rarityFull := false
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for i, rarity := range rarities {
		for j, trait := range traits {
			files, _ := ioutil.ReadDir("./imgs/" + rarity + "/" + trait)
			for _, file := range files {
				traitsArr[i][j] = append(traitsArr[i][j], file.Name())
				// traitName := strings.Split(file.Name(), "#")[0]
				// traitWeight, _ := strconv.Atoi(strings.Split(strings.Split(file.Name(), "#")[1], ".")[0])
				// weightForTrait[traitName] = traitWeight
			}
		}
	}
	fmt.Print(traitsArr)

	file, err := os.Create("./metadata.csv")
	if err != nil {
		panic(err)
	}
	csvw := csv.NewWriter(bufio.NewWriter((file)))
	var csvCell [][]string
	defer csvw.Flush()

	for i := 0; i < totalNum; i++ {
		var dna string
		var raritySelecter int
		dnaExist = false
		rarityFull = false
		rand.Seed(time.Now().UnixNano())
		raritySelecter = rand.Intn(6)

		for i, _ := range traitsArr[raritySelecter] {
			rand.Seed(time.Now().UnixNano())
			selecter := rand.Intn(len(traitsArr[raritySelecter][i]))
			images[i] = traitsArr[raritySelecter][i][selecter]
			var dna2 string
			if selecter < 16 {
				dna2 = "0" + strconv.FormatInt(int64(selecter), 16)
			} else if selecter >= 16 {
				dna2 = strconv.FormatInt(int64(selecter), 16)
			}
			dna += dna2
		}
		dna = dna[0:len(dna)-2] + strconv.Itoa(raritySelecter)
		for _, v := range dnaArr {
			if dna == v {
				dnaExist = true
				break
			}
		}
		switch raritySelecter {
		case 0:
			if rarityCounter[raritySelecter] == 1 {
				rarityFull = true
			}
		case 1:
			if rarityCounter[raritySelecter] == 10 {
				rarityFull = true
			}
		case 2:
			if rarityCounter[raritySelecter] == 70 {
				rarityFull = true
			}
		case 3:
			if rarityCounter[raritySelecter] == 120 {
				rarityFull = true
			}
		case 4:
			if rarityCounter[raritySelecter] == 300 {
				rarityFull = true
			}
		case 5:
			if rarityCounter[raritySelecter] == 499 {
				rarityFull = true
			}
		}

		if rarityFull {
			fmt.Println("RARITY FULL")
			i--
			continue
		} else {
			rarityCounter[raritySelecter]++
		}

		if dnaExist {
			fmt.Println("DNA EXIST")
			i--
			continue
		}
		fmt.Print(rarityCounter)

		csvItem := []string{
			strconv.Itoa(i),
			strings.Split(images[0], "#")[0],
			strings.Split(images[1], "#")[0],
			strings.Split(images[2], "#")[0],
			strings.Split(images[3], "#")[0],
			strings.Split(images[4], "#")[0],
			strings.Split(images[5], "#")[0],
			strings.Split(images[6], "#")[0],
			strings.Split(images[7], "#")[0],
			strings.Split(images[8], "#")[0],
			strings.Split(images[9], "#")[0],
			strings.Split(images[10], "#")[0],
		}
		csvCell = append(csvCell, csvItem)
		// fmt.Print(csvCell)

		module.Json_generator(images, i)

		dnaArr = append(dnaArr, dna)
		fmt.Println(i, dna)
		// fmt.Fprintln(w, dnaArr)

		decodedImages := make([]image.Image, len(images))

		for i, v := range images {
			decodedImages[i] = openAndDecode("./imgs/" + rarities[raritySelecter] + "/" + traits[i] + "/" + v)
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
	defer csvw.WriteAll(csvCell)
}
