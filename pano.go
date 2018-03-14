package main

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "image/jpeg"

	"github.com/icza/mjpeg"
)

// Configuration data
type Configuration struct {
	InputPath      string `json:"input_path"`
	InputPattern   string `json:"input_pattern"`
	OutputFullPath string `json:"output_fullpath"`
	Scale          int    `json:"scale"`
	Fps            int    `json:"fps"`
}

func setup(filename string) *Configuration {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("FAIL to read config file")
		fmt.Println("error:", err)
	}
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("FAIL to decode config file")
		fmt.Println("error:", err)
	}

	return &configuration
}

func listFiles(config *Configuration) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(config.InputPath)
	if err != nil {
		return files, err
	}
	for _, file := range fileInfo {
		if !file.IsDir() && strings.HasPrefix(file.Name(), config.InputPattern) && strings.HasSuffix(file.Name(), "JPG") {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func getDim(path string, filename string) (int, int) {
	file, err := os.Open(path + filename)
	defer file.Close()
	if err != nil {
		log.Println(err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Println(path+filename, err)
	}
	return image.Width, image.Height
}

func makeMovie(config *Configuration, files *[]string, originalWidth int, originalHeight int) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	width := int32(originalWidth / config.Scale)
	height := int32(originalHeight / config.Scale)

	// Video size
	aw, err := mjpeg.New(config.OutputFullPath, width, height, int32(config.Fps))
	checkErr(err)

	// // Create a movie from images: 1.jpg, 2.jpg, ..., 10.jpg
	// for i := 1; i <= 10; i++ {
	// 	data, err := ioutil.ReadFile(fmt.Sprintf("%d.jpg", i))
	// 	checkErr(err)
	// 	checkErr(aw.AddFrame(data))
	// }

	for _, f := range *files {
		data, err := ioutil.ReadFile(config.InputPath + f)
		checkErr(err)
		checkErr(aw.AddFrame(data))
	}

	checkErr(aw.Close())
}

func main() {
	fmt.Printf("Build a movie from pano files.\n")
	args := os.Args[1:]
	config := setup(args[0])
	images, err := listFiles(config)
	if err != nil {
		fmt.Println("FAIL to list files")
	}
	fmt.Printf("\tFound %d files in %s starting with `%s`\n", len(images), config.InputPath, config.InputPattern)
	width, height := getDim(config.InputPath, images[0])
	fmt.Printf("\tDimension of panoramics : %d x %d\n", width, height)
	makeMovie(config, &images, width, height)
	fmt.Printf("The end.\n")

	// ffmpeg -i small_movie.avi  -c:v libx264 -c:a copy -crf 18 -preset slower small_compressed.avi
}
