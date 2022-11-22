package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type Metadata struct {
	Name        string       `json:"name"`
	Symbol      string       `json:"symbol"`
	Description string       `json:"description"`
	Image       string       `json:"image"`
	Attributes  []Attributes `json:"attributes"`
	Properties  Properties   `json:"properties"`
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

func Json_generator(traitArr []string, num int) {
	// csv

	// json
	var result Metadata
	result.Attributes = make([]Attributes, 14)
	result.Properties.Files = make([]Files, 1)

	result.Symbol = "HUHU"
	result.Description = "HUHU and friends"
	result.Attributes[0].Trait_type = "legend"
	result.Attributes[0].Value = strings.Split(traitArr[0], ".")[0]
	result.Attributes[1].Trait_type = "background"
	result.Attributes[1].Value = strings.Split(traitArr[1], ".")[0]
	result.Attributes[2].Trait_type = "backpack"
	result.Attributes[2].Value = strings.Split(traitArr[2], ".")[0]
	result.Attributes[3].Trait_type = "sleepbag"
	result.Attributes[3].Value = strings.Split(traitArr[3], ".")[0]
	result.Attributes[4].Trait_type = "pet"
	result.Attributes[4].Value = strings.Split(traitArr[4], ".")[0]
	result.Attributes[5].Trait_type = "body"
	result.Attributes[5].Value = strings.Split(traitArr[5], ".")[0]
	result.Attributes[6].Trait_type = "outfit"
	result.Attributes[6].Value = strings.Split(traitArr[6], ".")[0]
	result.Attributes[7].Trait_type = "ring"
	result.Attributes[7].Value = strings.Split(traitArr[7], ".")[0]
	result.Attributes[8].Trait_type = "hair"
	result.Attributes[8].Value = strings.Split(traitArr[8], ".")[0]
	result.Attributes[9].Trait_type = "eye"
	result.Attributes[9].Value = strings.Split(traitArr[9], ".")[0]
	result.Attributes[10].Trait_type = "mouth"
	result.Attributes[10].Value = strings.Split(traitArr[11], ".")[0]
	result.Attributes[11].Trait_type = "rarity"
	result.Attributes[11].Value = strings.Split(traitArr[12], ".")[0]

	result.Properties.Files[0].Type = "image/png"
	result.Name = "HUHU#" + strconv.Itoa(num)
	result.Image = strconv.Itoa(num) + ".png"
	result.Properties.Files[0].Uri = strconv.Itoa(num) + ".png"
	doc, _ := json.Marshal(result)

	err2 := ioutil.WriteFile("./result/"+strconv.Itoa(num)+"/meta.json", doc, os.FileMode(0644))

	if err2 != nil {
		fmt.Println(err2)
		return
	}
}
