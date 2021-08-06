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
	facesURL = "https://github.com/TheMdTF/mdtf-public/tree/master/rally2-matching-system/tests/test-routine-images/face"
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

func main() {
	matchClient.HelloViki()

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

	// read in list of images
	files, err := ioutil.ReadDir("./images")
	if err != nil {
		log.Fatal(err)
	}

	revFiles, err := ioutil.ReadDir("./images")
	if err != nil {
		log.Fatal(err)
	}

	//reverse the list of names
	reverseArray(revFiles)

	for _, f := range files {
		revFiles = RemoveIndex(revFiles, len(revFiles)-1)
		for _, r := range revFiles {
			fmt.Println("Comparing " + f.Name() + " with " + r.Name())
		}
	}

}
