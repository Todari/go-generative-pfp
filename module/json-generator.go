package module

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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
	var result Metadata
	result.Attributes = make([]Attributes, 7)
	result.Properties.Files = make([]Files, 1)

	result.Symbol = "999"
	result.Description = "99999"
	result.Attributes[0].Trait_type = "background"
	result.Attributes[0].Value = traitArr[0]
	result.Attributes[1].Trait_type = "item"
	result.Attributes[1].Value = traitArr[1]
	result.Attributes[2].Trait_type = "body"
	result.Attributes[2].Value = traitArr[2]
	result.Attributes[3].Trait_type = "clothes"
	result.Attributes[3].Value = traitArr[3]
	result.Attributes[4].Trait_type = "hair"
	result.Attributes[4].Value = traitArr[4]
	result.Attributes[5].Trait_type = "eye"
	result.Attributes[5].Value = traitArr[5]
	result.Attributes[6].Trait_type = "hat"
	result.Attributes[6].Value = traitArr[6]

	result.Properties.Files[0].Type = "image/png"
	result.Name = "HUHU#" + strconv.Itoa(num)
	result.Image = strconv.Itoa(num) + ".png"
	result.Properties.Files[0].Uri = strconv.Itoa(num) + ".png"
	doc, _ := json.Marshal(result)

	err := ioutil.WriteFile("./json/"+strconv.Itoa(num)+".json", doc, os.FileMode(0644))
	if err != nil {
		fmt.Println(err)
		return
	}
}
