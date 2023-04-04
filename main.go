package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/Todari/go-generative-pfp/module"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type Metadata struct {
	Name        string       `json:"name"`
	Symbol      string       `json:"symbol"`
	Description string       `json:"description"`
	Image       string       `json:"image"`
	Attributes  []Attributes `json:"attributes"`
	// Properties  Properties   `json:"properties"`
}

type Attributes struct {
	Trait_type string `json:"trait_type"`
	Value      string `json:"value"`
}

type Properties struct {
	Files []Files `json:"files"`
}

type Files struct {
	Uri  string `json:"uri"`
	Type string `json:"type"`
}

func reset() {
	if !os.IsExist(os.RemoveAll("./result/")) {
		os.Mkdir("./result/", 0777)
	}
	if !os.IsExist(os.RemoveAll("./result2/")) {
		os.Mkdir("./result2/", 0777)
	}
}

func openAndDecode(imgPath string) image.Image {
	img, err := os.Open(imgPath)
	if err != nil {
		log.Fatalf("Failed to open %s", err)
	}

	decoded, err := png.Decode(img)
	if err != nil {
		log.Fatalf("Failed to decode %s : %s", err, imgPath)
	}
	defer img.Close()

	return decoded
}

func main() {

	doReset := true
	totalNum := 10000
	traitNum := [6]int{10, 100, 700, 1200, 3000, 4990}
	traitNumSum := 0
	for _, v := range traitNum {
		traitNumSum += v
	}

	if totalNum != traitNumSum {
		log.Fatalf("totalNum doesn't match with sum of traitNum")
	}

	if doReset {
		reset()
	}

	//img 폴더에 저장되어 있는 rarity 이름, trait 이름 설정
	var rarities = [6]string{"1. legend", "2. prime", "3. master", "4. expert", "5. junior", "6. rookie"}
	var traits = [13]string{"0. legend", "1. background", "2. backpack", "3. sleepbag", "4. pet", "5. body", "6. outfit", "7. ring", "8. lhair", "9. eye", "10. rhair", "11. mouth", "12. rarity"}
	var dnaArr []string

	// 가능한 모든 rarity, trait를 이중배열로 구성
	var traitsArr [6][13][]string
	rarityCounter := make([]int, 6)
	//조합의 파일명을 넣는 슬라이스
	images := make([]string, len(traits))
	var imagesArr [][]string
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	var raritySelecterArr []int
	for j, v := range traitNum {
		for i := 0; i < v; i++ {
			raritySelecterArr = append(raritySelecterArr, j)
		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(totalNum, func(i int, j int) {
		raritySelecterArr[i], raritySelecterArr[j] = raritySelecterArr[j], raritySelecterArr[i]
	})

	fmt.Print(raritySelecterArr)
	for i, rarity := range rarities {

		for j, trait := range traits {
			//피부색 + 눈커풀 컨트롤
			var files []fs.FileInfo
			// if trait == "8. expression" {
			// 	files, _ = ioutil.ReadDir("./imgs/" + rarity + "/" + trait + "/Yellow/")
			// } else {
			files, _ = ioutil.ReadDir("./imgs/" + rarity + "/" + trait)
			// }

			for _, file := range files {
				//.DS_Store 파일 존재하면 삭제
				if file.Name() == ".DS_Store" {
					os.Remove(".DS_Store")
					continue
				} else {
					traitsArr[i][j] = append(traitsArr[i][j], file.Name())
				}
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
	csvCategory := []string{
		"ID", "legend", "background", "backpack", "sleepbag", "pet", "outfit", "ring", "head", "eye", "mouth", "rarity",
	}
	csvCell = append(csvCell, csvCategory)
	defer csvw.Flush()

	for tokenId, rarityIndex := range raritySelecterArr {

	resetTraits:
		dnaExist := false
		var dna string
		for i, _ := range traitsArr[rarityIndex] {
			rand.Seed(time.Now().UnixNano())
			selecter := rand.Intn(len(traitsArr[rarityIndex][i]))
			images[i] = traitsArr[rarityIndex][i][selecter]
			if i == 10 {
				images[i] = images[8]
			}
			var dna2 string
			if selecter < 16 {
				dna2 = "0" + strconv.FormatInt(int64(selecter), 16)
			} else if selecter >= 16 {
				dna2 = strconv.FormatInt(int64(selecter), 16)
			}
			dna += dna2
		}
		dna = dna[0:len(dna)-2] + strconv.Itoa(rarityIndex)
		for _, v := range dnaArr {
			if dna == v {
				dnaExist = true
				break
			}
		}

		if dnaExist {
			fmt.Println("DNA EXIST")
			goto resetTraits
		}

		fmt.Println(images)
		imagesArr = append(imagesArr, images)

		rarityCounter[rarityIndex]++
		fmt.Print(rarityCounter)

		//csv 배열에 push
		csvItem := []string{
			strconv.Itoa(tokenId),
			strings.Split(images[0], ".")[0],
			strings.Split(images[1], ".")[0],
			strings.Split(images[2], ".")[0],
			strings.Split(images[3], ".")[0],
			strings.Split(images[4], ".")[0],
			// strings.Split(images[5], ".")[0],
			strings.Split(images[6], ".")[0],
			strings.Split(images[7], ".")[0],
			strings.Split(images[8], ".")[0],
			strings.Split(images[9], ".")[0],
			// strings.Split(images[10], ".")[0],
			strings.Split(images[11], ".")[0],
			strings.Split(images[12], ".")[0],
		}
		csvCell = append(csvCell, csvItem)

		dnaArr = append(dnaArr, dna)
		fmt.Println("\tID: ", tokenId, "\tDNA : ", dna)

		decodedImages := make([]image.Image, len(images))

		for i, v := range images {
			//피부색과 일치하는 눈커풀 넣기
			// if traits[i] == "8. expression" {
			// 	decodedImages[i] = openAndDecode("./imgs/" + rarities[raritySelecter] + "/" + traits[i] + "/" + strings.Split(images[5], ".")[0] + "/" + v)
			// } else {
			decodedImages[i] = openAndDecode("./imgs/" + rarities[rarityIndex] + "/" + traits[i] + "/" + v)
			// }
		}

		bounds := decodedImages[0].Bounds()
		newImage := image.NewRGBA(bounds)
		for _, img := range decodedImages {
			draw.Draw(newImage, img.Bounds(), img, image.ZP, draw.Over)
		}
		//디렉토리 생성 후 이미지 제작
		go os.Mkdir("./result/"+strconv.Itoa(tokenId), 0777)
		// result, err := os.Create("./result/" + strconv.Itoa(i) + "/image.png")

		result2, err := os.Create("./result2/" + strconv.Itoa(tokenId) + ".png")

		//json 제작
		module.Json_generator(images, tokenId)

		if err != nil {
			log.Fatalf("Failed to create: %s", err)
		}

		if tokenId == totalNum-1 {
			png.Encode(result2, newImage)
		} else {
			go png.Encode(result2, newImage)
		}

		// defer result.Close()
		defer result2.Close()
	}
	defer csvw.WriteAll(csvCell)
}
