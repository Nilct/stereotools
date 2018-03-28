package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Configuration data
type Configuration struct {
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_subfolder"`
	Perc       int    `json:"percentage"`
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
		if !file.IsDir() && strings.HasSuffix(file.Name(), "JPG") {
			files = append(files, file.Name())
		}
	}
	return files, nil
}

func selectPanoramics(config *Configuration, images *[]string) {
	fmt.Println("Selection in progress")
	// create folder
	err := os.Mkdir(config.OutputPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	// move files (pick %d images every 100 pics)

}

func main() {
	fmt.Printf("-- Select a percentage of panoramic files and move them to VISU subfolder. --\n")
	args := os.Args[1:]
	config := setup(args[0])
	images, err := listFiles(config)
	if err != nil {
		fmt.Println("FAIL to list files")
	}
	fmt.Printf("\tFound %d panoramics in %s\n", len(images), config.InputPath)
	noOfFiles := (int)((len(images) / 100.0) * config.Perc)
	fmt.Printf("Number of files to check : %d (%d perc.)\n\t ok to proceed ? Y/n \n", noOfFiles, config.Perc)
	var input string
	fmt.Scanln(&input)
	if len(input) == 0 || strings.TrimRight(input, "\n") == "Y" || strings.TrimRight(input, "\n") == "y" {
		selectPanoramics(config, &images)
	}

	fmt.Printf("The end.\n")

	// ffmpeg -i small_movie.avi  -c:v libx264 -c:a copy -crf 18 -preset slower small_compressed.avi

	// mencoder small_movie.avi -ovc x264 -x264encopts ratetol=100:preset=medium:crf=20:pass=1 -nosound -o video1.h264
}
