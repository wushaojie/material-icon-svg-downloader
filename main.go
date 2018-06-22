package main

import (
	"log"
	"fmt"
	"strconv"
	"net/http"
	"os"
	"io"
	"strings"
	"io/ioutil"
	"encoding/json"
)

type ImageUrl struct {
	Twotone  string `json:"twotone,omitempty"`
	Sharp    string `json:"sharp,omitempty"`
	Outline  string `json:"outline,omitempty"`
	Round    string `json:"round,omitempty"`
	Baseline string `json:"baseline,omitempty"`
}

type ICON struct {
	ID        string    `json:"id"`
	ImageUrls *ImageUrl `json:"imageUrls,omitempty"`
}

type ICONS struct {
	Icons []ICON `json:"icons"`
	Name  string `json:"name"`
}

type MaterialIcon struct {
	BaseUrl    string  `json:"baseUrl"`
	Categories []ICONS `json:"categories"`
}

const (
	BASE_URL           = "https://material.io/tools/icons/static/icons/"
	DATA_URL           = "https://material.io/tools/icons/static/data.json"
	MATERIAL_DIRECTORY = "material"
	ICON_TYPE          = "baseline"
	ICON_SIZE          = 24
	PERM               = 0755
)


func save(materialIcon MaterialIcon, iconType string, iconSize int) {
	var categories = materialIcon.Categories
	var count int
	fmt.Println("downloading...")
	for i := 0; i < len(categories); i++ {
		category := categories[i]
		directoryName := category.Name
		icons := category.Icons
		err := os.Mkdir(directoryName, PERM)

		if err != nil {
			fmt.Println(err)
		}

		for j := 0; j < len(icons); j++ {
			var fileName string
			item := icons[j]
			fileName = iconType + "-" + item.ID + "-" + strconv.Itoa(iconSize) + "px.svg"
			imageUrls := item.ImageUrls

			if imageUrls != nil {
				fileName = imageUrls.Baseline
			}
			imageUrl := BASE_URL + fileName

			resp, err := http.Get(imageUrl)
			if err != nil {
				log.Fatal(err)
			}

			file, err := os.Create(directoryName + "/" + strings.Replace(fileName, iconType+"-", "", 1))
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(file, resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			count++

			file.Close()
			resp.Body.Close()
		}
	}
	fmt.Printf("Complete %d\n icons", count)
}

func getMaterialIconData(url string) MaterialIcon {
	r, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer r.Body.Close()

	var materialIcon MaterialIcon
	json.Unmarshal(b, &materialIcon)
	return materialIcon
}

func main() {
	os.Mkdir(MATERIAL_DIRECTORY, PERM)
	os.Chdir(MATERIAL_DIRECTORY)

	var data = getMaterialIconData(DATA_URL)
	save(data, ICON_TYPE, ICON_SIZE)
}
