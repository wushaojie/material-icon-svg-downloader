package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type ICON struct {
	Name                string   `json:"name"`
	Version             int      `json:"version"`
	UnsupportedFamilies []string `json:"unsupported_families"`
	Categories          []string `json:"categories"`
	Tags                []string `json:"tags"`
	SizePx              []string `json:"sizes_px"`
}

type MaterialIcon struct {
	ICONS []ICON `json:"icons"`
}

// https://fonts.gstatic.com/s/i/materialiconsoutlined/accessible/v4/24px.svg?download=true
const (
	BASE_URL  = "https://fonts.gstatic.com/s/i"
	DATA_URL  = "icons.json"
	BASE_PATH = "material"
	ICON_TYPE = "materialiconsoutlined"
	ICON_SIZE = 24
	PERM      = 0755
)

func save(materialIcon MaterialIcon) {
	var icons = materialIcon.ICONS
	var count int
	fmt.Println("downloading...")
	os.RemoveAll(BASE_PATH)

	for i := 0; i < len(icons); i++ {
		icon := icons[i]
		name := icon.Name
		version := "v" + strconv.Itoa(icon.Version)
		path := BASE_PATH + "/" + icon.Categories[0]

		if _, err := os.Stat(path); os.IsNotExist(err) {
			os.MkdirAll(path, PERM)
		}

		var fileName string
		fileName = strconv.Itoa(ICON_SIZE) + "px.svg"
		imageUrl := fmt.Sprintf("%s/%s/%s/%s/%s", BASE_URL, ICON_TYPE, name, version, fileName)
		fmt.Println(imageUrl)

		resp, err := http.Get(imageUrl)
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.Create(path + "/" + name + ".svg")
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
	fmt.Printf("Downloaded %d icons", count)
}

func getMaterialIconData(file string) MaterialIcon {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	var materialIcon MaterialIcon
	json.Unmarshal(b, &materialIcon)
	return materialIcon
}

func main() {
	var data = getMaterialIconData(DATA_URL)
	save(data)
}
