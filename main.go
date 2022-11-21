package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
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

type Trait struct {
	Rarity string
	Path   string
	Name   string
	Weight int
	Group  string
}

func main() {

	doReset := true
	totalNum := 500

	if doReset {
		reset()
	}

	//img 폴더에 저장되어 있는 rarity 이름, trait 이름 설정
	var rarities = [6]string{"1. legend", "2. prime", "3. master", "4. expert", "5. junior", "6. rookie"}
	var traits = [14]string{"0. legend", "1. background", "2. backpack", "3. sleepbag", "4. pet", "5. body", "6. outfit", "7. ring", "8. expression", "9. lhair", "10. eye acc", "11. rhair", "12. mouth", "13. rarity"}
	var dnaArr []string

	// 가능한 모든 rarity, trait를 이중배열로 구성
	var traitsArr [6][14][]string
	rarityCounter := make([]int, 6)
	//조합의 파일명을 넣는 슬라이스
	images := make([]string, len(traits))
	dnaExist := false
	rarityFull := false
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for i, rarity := range rarities {

		for j, trait := range traits {
			//피부색 + 눈커풀 컨트롤
			var files []fs.FileInfo
			if trait == "8. expression" {
				files, _ = ioutil.ReadDir("./imgs/" + rarity + "/" + trait + "/Yellow/")
			} else {
				files, _ = ioutil.ReadDir("./imgs/" + rarity + "/" + trait)
			}

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
			if i == 11 {
				images[i] = images[9]
			}
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
		fmt.Println(images)

		// fmt.Println("1. DNA Created")

		//레어리티별 최대 갯수 설정
		switch raritySelecter {
		case 0:
			if rarityCounter[raritySelecter] == 0 {
				rarityFull = true
			}
		case 1:
			if rarityCounter[raritySelecter] == 100 {
				rarityFull = true
			}
		case 2:
			if rarityCounter[raritySelecter] == 100 {
				rarityFull = true
			}
		case 3:
			if rarityCounter[raritySelecter] == 100 {
				rarityFull = true
			}
		case 4:
			if rarityCounter[raritySelecter] == 100 {
				rarityFull = true
			}
		case 5:
			if rarityCounter[raritySelecter] == 100 {
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
		// fmt.Println("2. DNA, rarity checked")
		// fmt.Print(rarityCounter)

		//csv 배열에 push
		csvItem := []string{
			strconv.Itoa(i),
			strings.Split(images[0], ".")[0],
			strings.Split(images[2], ".")[0],
			strings.Split(images[1], ".")[0],
			strings.Split(images[3], ".")[0],
			strings.Split(images[4], ".")[0],
			strings.Split(images[5], ".")[0],
			strings.Split(images[6], ".")[0],
			strings.Split(images[7], ".")[0],
			strings.Split(images[8], ".")[0],
			strings.Split(images[9], ".")[0],
			strings.Split(images[10], ".")[0],
			strings.Split(images[11], ".")[0],
			strings.Split(images[12], ".")[0],
			strings.Split(images[13], ".")[0],
		}
		csvCell = append(csvCell, csvItem)
		// fmt.Println("3. csv appended")

		dnaArr = append(dnaArr, dna)
		fmt.Println("ID: ", i, "DNA : ", dna)

		fmt.Println(images)

		decodedImages := make([]image.Image, len(images))

		for i, v := range images {
			//피부색과 일치하는 눈커풀 넣기
			if traits[i] == "8. expression" {
				decodedImages[i] = openAndDecode("./imgs/" + rarities[raritySelecter] + "/" + traits[i] + "/" + strings.Split(images[5], ".")[0] + "/" + v)
			} else {
				decodedImages[i] = openAndDecode("./imgs/" + rarities[raritySelecter] + "/" + traits[i] + "/" + v)
			}
		}
		// fmt.Println("4. img decoded")

		bounds := decodedImages[0].Bounds()
		newImage := image.NewRGBA(bounds)

		for _, img := range decodedImages {
			draw.Draw(newImage, img.Bounds(), img, image.ZP, draw.Over)
		}
		// fmt.Println("5. img drawed")

		//디렉토리 생성 후 이미지 제작
		os.Mkdir("./result/"+strconv.Itoa(i), 0777)
		// result, err := os.Create("./result/" + strconv.Itoa(i) + "/image.png")
		result2, _ := os.Create("./result2/" + strconv.Itoa(i) + ".png")

		// fmt.Println("6. img created")
		//json 제작
		module.Json_generator(images, i)
		// fmt.Println("7. json created")

		if err != nil {
			log.Fatalf("Failed to create: %s", err)
		}

		// png.Encode(result, newImage)
		png.Encode(result2, newImage)
		// fmt.Println("8. img incoded")

		// defer result.Close()
		defer result2.Close()

		sort.Strings(dnaArr)

	}
	defer csvw.WriteAll(csvCell)
}
