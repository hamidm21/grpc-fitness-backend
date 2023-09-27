package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"gopkg.in/yaml.v2"

	"gitlab.com/mefit/mefit-server/entity"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func copy(src string, dst string) {
	// Read all content of src to data
	data, err := ioutil.ReadFile(src)
	checkErr(err)
	// Write data to dst
	err = ioutil.WriteFile(dst, data, 0644)
	checkErr(err)
}

func main() {
	var (
		subPath  string
		videoDir string
		imageDir string
	)
	movementsDict := make(map[string]entity.Movement)
	flag.StringVar(&subPath, "http-subpath", "/media", "related http path for fetching media(video/image) from e.g. media for example.com/media/{media-file}")
	flag.StringVar(&videoDir, "video-dir", "tmp/video", "the directory containing our video files")
	flag.StringVar(&imageDir, "image-dir", "tmp/image", "the directory containing cover image for our video files")

	// First check videos
	videos, err := ioutil.ReadDir(videoDir)
	if err != nil {
		log.Fatal(err)
	}

	for idx, f := range videos {
		m := entity.Movement{}
		name := f.Name()
		arr := strings.Split(name, ".")
		key, fileExt := arr[0], arr[1]
		if strings.HasSuffix(key, "-1") {
			key = strings.Split(name, "-1")[0]
		}
		//We save lower case string only
		key = strings.ToLower(key)
		fileName := strings.Join(strings.Split(key, " "), "-")
		//Copy this file with new name --> build/assets
		copy(fmt.Sprintf("tmp/video/%s", name), fmt.Sprintf("build/assets/%s.%s", fileName, fileExt))
		fmt.Printf("%s --key--> %s --file--> %s.%s", f.Name(), key, fileName, fileExt)
		fmt.Println("")
		//Save to seed data
		m.Name = key
		m.ID = uint(idx)
		m.CreatedAt = time.Now()
		m.UpdatedAt = time.Now()
		m.VideoUrl = fmt.Sprintf("%s/%s.%s", subPath, fileName, fileExt)
		movementsDict[key] = m

	}

	//Add picture url
	images, err := ioutil.ReadDir(imageDir)
	if err != nil {
		log.Fatal(err)
	}

	//Add images to asset and seed data
	for _, f := range images {
		m := entity.Movement{}
		name := f.Name()
		arr := strings.Split(name, ".")
		key, fileExt := arr[0], arr[len(arr)-1]
		if strings.HasSuffix(key, "-1") {
			key = strings.Split(name, "-1")[0]
		}
		//We save lower case string only
		key = strings.ToLower(key)
		fileName := strings.Join(strings.Split(key, " "), "-")
		//Save to seed file
		m, ok := movementsDict[key]
		if !ok {
			panic("Something goes wrong with key: " + key)
		}
		//Copy this file with new name --> build/assets
		copy(fmt.Sprintf("tmp/image/%s", name), fmt.Sprintf("build/assets/%s.%s", fileName, fileExt))
		fmt.Printf("%s --key--> %s --file--> %s.%s", f.Name(), key, fileName, fileExt)
		fmt.Println("")
		m.ThumbnailUrl = fmt.Sprintf("%s/%s.%s", subPath, fileName, fileExt)
		movementsDict[key] = m
	}

	//Convert data to yaml format
	movements := make([]entity.Movement, 0, len(movementsDict))
	for _, v := range movementsDict {
		movements = append(movements, v)
	}
	out, err := yaml.Marshal(movements)
	if err != nil {
		checkErr(err)
	}
	//Save to yaml format
	checkErr(ioutil.WriteFile("seed/seed.yaml", out, 0644))
}
