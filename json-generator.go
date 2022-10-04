package json_generator

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

func json_generator() {
	var result Metadata
	result.Attributes = make([]Attributes, 3)
	result.Properties.Files = make([]Files, 1)

	result.Symbol = "999"
	result.Description = "99999"
	result.Attributes[0].Trait_type = "head"
	result.Attributes[0].Value = "good head"
	result.Attributes[1].Trait_type = "body"
	result.Attributes[1].Value = "good body"
	result.Attributes[2].Trait_type = "background"
	result.Attributes[2].Value = "good background"
	result.Properties.Files[0].Type = "image/png"
	for i := 0; i < 1000; i++ {
		result.Name = "999NFTs#" + strconv.Itoa(i)
		result.Image = strconv.Itoa(i) + ".png"
		result.Properties.Files[0].Uri = strconv.Itoa(i) + ".png"
		doc, _ := json.Marshal(result)

		err := ioutil.WriteFile("./nft/assets/"+strconv.Itoa(i)+".json", doc, os.FileMode(0644))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
