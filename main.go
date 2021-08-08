package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	//memory "github.com/go-git/go-git/storage/memory"
	"github.com/vnkeaton/run-match/matchClient"
)

const (
	facesURL  = "https://github.com/TheMdTF/mdtf-public/tree/master/rally2-matching-system/tests/test-routine-images/face"
	imagesDir = "/images/"
)

//var storer *memory.Storage
//var fs billy.Filesystem

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
	matchClient.Hello()

	// TODO cloning the faces repo is not working.  Come back to this.
	/*err := os.Mkdir("/tmp/foo", 0755)
	if err != nil {
		fmt.Println("error in making tmp foo directory")
	}
	_, err2 := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:      facesURL,
		Progress: os.Stdout,
	})

	if err2 != nil {
		fmt.Println("error in pullig repository for faces")
	}
	*/

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("error in determing working directory")
		log.Fatal(err)
	}

	fmt.Println("working directory is: " + path)

	// read in list of images
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

	for _, f := range files {
		////
		//matchClient.PostThisFile(mustOpen(path, f.Name()))
		/////

		revFiles = RemoveIndex(revFiles, len(revFiles)-1)
		for _, r := range revFiles {

			fmt.Println("Comparing " + f.Name() + " with " + r.Name())
			mediaFiles := []string{mustOpen(path, f.Name()), mustOpen(path, r.Name())}
			matchClient.MatchFiles(mediaFiles)
			//err := matchClient.PostThisFile(mustOpen(path, f.Name()), mustOpen(path, r.Name())
			//if err != nil {
			//	panic(err)
			//}

		}

	}
}
