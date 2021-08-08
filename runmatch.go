package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	//memory "github.com/go-git/go-git/storage/memory"
	"github.com/vnkeaton/biometric-match-client/matchclient"
)

const (
	facesURL  = "https://github.com/TheMdTF/mdtf-public/tree/master/rally2-matching-system/tests/test-routine-images/face"
	imagesDir = "/images/"
)

type MatchScoreData struct {
	FileName1  string
	FileName2  string
	FatchScore float64
}

func reverseArray(arr []os.FileInfo) []os.FileInfo {
	// reverse file name array
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func RemoveIndex(arr []os.FileInfo, index int) []os.FileInfo {
	return append(arr[:index], arr[index+1:]...)
}

//func mustOpen(f string) *os.File {
func mustOpen(p string, f string) string {
	filename := p + imagesDir + f
	fmt.Println("check to see if filename exists with mustOpen for: " + filename)
	_, err := os.Open(filename)
	if err != nil {
		fmt.Println("This file failed upon mustOpen: " + filename)
		log.Fatal(err)
	}
	return filename
}

func main() {
	matchclient.Hello()

	//get the faces repo
	err := downloadFaces()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("faces downloaded")

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("error in determing working directory")
		log.Fatal(err)
	}

	fmt.Println("working directory is: " + path)

	//read in list of images
	files, err := ioutil.ReadDir(path + imagesDir)
	if err != nil {
		log.Fatal(err)
	}

	revFiles, err := ioutil.ReadDir(path + imagesDir)
	if err != nil {
		log.Fatal(err)
	}
	//reverse the list of names
	reverseArray(revFiles)

	//match each file
	for _, f := range files {
		//triangular comparison for comparing unique files - do not assume the match operation is symmetric
		//revFiles = RemoveIndex(revFiles, len(revFiles)-1)
		for _, r := range revFiles {

			fmt.Println("Comparing " + f.Name() + " with " + r.Name())
			mediaFiles := []string{mustOpen(path, f.Name()), mustOpen(path, r.Name())}

			fmt.Println("calling match client....")
			matchScore, err := matchclient.MatchFiles(mediaFiles)
			if err != nil {
				fmt.Println("error from matchclient")
				log.Fatal(err)
			}
			s, _ := json.MarshalIndent(matchScore, "", "\t")
			fmt.Print(string(s))
			fmt.Print("\n")
		}
	}

	//get match score results
	allMatchScores, err := matchclient.GetAllMatchScores()
	if err != nil {
		fmt.Println("error from matchclient getting all match score data")
	}
	s, _ := json.MarshalIndent(allMatchScores, "", "\t")
	fmt.Print(string(s))
	fmt.Print("\n")

}

func downloadFaces() error {
	resp, err := http.Get(facesURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}

	//Create an empty file
	file, err := os.Create()

}
