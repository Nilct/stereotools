package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"
)

// CHUNKSIZE is the size of one chunck
const CHUNKSIZE = 1000.0

var csvFields []string
var csvData []string

// Configuration data
type Configuration struct {
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_subfolder"`
	Perc       int    `json:"percentage"`
	Fields     string `json:"fields"`
}

func initRandom() {
	rand.Seed(time.Now().UnixNano())
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

	// get all fields for csv file
	s := strings.Trim(configuration.Fields, " ")
	csvFields = strings.Split(s, ",")
	csvData = make([]string, len(csvFields))

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

func selectPanoramics(config *Configuration, images []string) {
	debug := true
	fmt.Println("Selection in progress")
	totalPick := 0
	estimatedPick := (int)((len(images) / 100.0) * config.Perc)
	// create folder
	os.Mkdir(config.OutputPath, os.ModePerm)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// maximum chunck size is 1000 --> split accordingly and select randomly on full last chunks
	noOfSplits := (int)(len(images)/CHUNKSIZE) + 1
	noOfPickedImages := CHUNKSIZE * config.Perc / 100
	if debug {
		fmt.Printf("noOfSplits: %d, noOfPickedImages in one chunk: %d\n", noOfSplits, noOfPickedImages)
	}

	fmt.Printf("selected panoramics are saved in %s\n", path.Join(config.OutputPath, "panoramics.csv"))
	file, err := os.OpenFile(path.Join(config.OutputPath, "panoramics.csv"), os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()

	if err != nil {
		fmt.Printf("Problem creating csv file (check if it already exists)\n")
		os.Exit(1)
	}
	writer := csv.NewWriter(file)
	writer.Write(csvFields)
	writer.Flush()

	var picName string
	for i := 0; i < noOfSplits; i++ {
		startIndex := CHUNKSIZE * i
		for j := 0; j < noOfPickedImages; j++ {
			pick := rand.Intn(CHUNKSIZE) // warning : allow for duplicate
			if startIndex+pick < len(images) {
				picName = images[startIndex+pick]
				// check file has not been moved already
				if _, err := os.Stat(path.Join(config.InputPath, picName)); err == nil {
					err = os.Rename(path.Join(config.InputPath, picName), path.Join(config.OutputPath, picName))
					if err == nil {
						totalPick++
						csvData[0] = picName
						writer.Write(csvData)
						writer.Flush()
					}
				}
			}
		}
	}
	fmt.Printf("Picked files : %d (estimated %d)\n", totalPick, estimatedPick)
}

func main() {
	fmt.Printf("-- Select a percentage of panoramic files and move them to VISU subfolder (WARNING: it must be on the SAME partition !). --\n")
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("Please add a json ")
		return
	}
	initRandom()
	config := setup(args[0])
	if !strings.HasPrefix(config.OutputPath, config.InputPath[:2]) {
		fmt.Println("ERROR: input and output path not on same partition !")
		return
	}

	images, err := listFiles(config)
	if err != nil {
		fmt.Println("FAIL to list files")
		return
	}
	fmt.Printf("\tFound %d panoramics in %s\n", len(images), config.InputPath)
	noOfFiles := (int)((len(images) / 100.0) * config.Perc)
	fmt.Printf("Number of files to check : %d (%d perc.)\n\t ok to proceed ? Y/n \n", noOfFiles, config.Perc)
	var input string
	fmt.Scanln(&input)
	if len(input) == 0 || strings.TrimRight(input, "\n") == "Y" || strings.TrimRight(input, "\n") == "y" {
		selectPanoramics(config, images)
	}

	fmt.Printf("The end.\n")
}
