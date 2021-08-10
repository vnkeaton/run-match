package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	tabby "github.com/cheynewallace/tabby"
	billy "github.com/go-git/go-billy/v5"
	memfs "github.com/go-git/go-billy/v5/memfs"

	//afero "github.com/spf13/afero"
	git "github.com/go-git/go-git/v5"
	memory "github.com/go-git/go-git/v5/storage/memory"
	matchclient "github.com/vnkeaton/biometric-match-client"
)

var storer *memory.Storage
var origin billy.Filesystem

const (
	facesURL  = "https://github.com/TheMdTF/mdtf-public"
	imagesDir = "/tmp/images/"
	repoDir   = "/rally2-matching-system/tests/test-routine-images/face"
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
	//return append(arr[:index], arr[index+1:]...)
	ret := make([]os.FileInfo, 0)
	ret = append(ret, arr[:index]...)
	return append(ret, arr[index+1:]...)
}

func main() {
	storer = memory.NewStorage()
	origin = memfs.New()

	matchclient.Hello("IDSL")

	//get the faces repo
	err := downloadFaces()
	if err != nil {
		fmt.Println("Error from downloadFaces  ", err)
		log.Fatal(err)
	}

	//read in list of images
	imageFaces, err := ioutil.ReadDir(imagesDir)
	if err != nil {
		fmt.Println("Error from read dir  ", err)
		log.Fatal(err)
	}

	revFiles, err := ioutil.ReadDir(imagesDir)
	if err != nil {
		fmt.Println("Error from read dir rev files ", err)
		log.Fatal(err)
	}

	//reverse the list of names
	reverseArray(revFiles)

	//match each file
	for _, f := range imageFaces {
		//triangular comparison for comparing unique files - do not assume the match operation is symmetric
		//revFiles = RemoveIndex(revFiles, len(revFiles)-1)
		for _, r := range revFiles {
			//fmt.Println("Comparing " + imagesDir + f.Name() + " with " + imagesDir + r.Name())
			mediaFiles := []string{imagesDir + f.Name(), imagesDir + r.Name()}
			_, err := matchclient.MatchFiles(mediaFiles)
			if err != nil {
				fmt.Println("error from matchclient")
				log.Fatal(err)
			}
			//print out json
			/*s, _ := json.MarshalIndent(matchScore, "", "\t")
			fmt.Print(string(s))
			fmt.Print("\n")
			*/
		}
	}

	//get match score results
	allMatchScores, err := matchclient.GetAllMatchScores()
	if err != nil {
		fmt.Println("error from matchclient getting all match score data")
	}
	//Print out Json
	/*s, _ := json.MarshalIndent(allMatchScores, "", "\t")
	fmt.Print(string(s))
	fmt.Print("\n")
	*/

	//Print out table
	ShowTable(allMatchScores)

}

func downloadFaces() error {

	storer = memory.NewStorage()
	origin = memfs.New()

	_, err := git.Clone(storer, origin, &git.CloneOptions{
		URL: facesURL,
	})
	if err != nil {
		fmt.Println("Error from clone  ", err)
		log.Fatal(err)
	}

	//read in list of images
	memFaces, err := origin.ReadDir(repoDir)
	if err != nil {
		fmt.Println("Error from read dir  ", err)
		log.Fatal(err)
	}

	if _, err := os.Stat(imagesDir); os.IsNotExist(err) {
		err := os.Mkdir(imagesDir, 0755)
		if err != nil {
			fmt.Println("error mkdir for " + imagesDir)
			log.Fatal(err)
		}
	}

	for _, f := range memFaces {
		ext := filepath.Ext(f.Name())
		if ext != ".png" {
			break
		}

		//open a file
		src, err := origin.Open(repoDir + "/" + f.Name())
		if err != nil {
			fmt.Println("Error from open: ", err)
			log.Fatal(err)
		}

		//create a new file
		dst, err := os.Create(imagesDir + "/" + f.Name())
		if err != nil {
			fmt.Println("Error from create: ", err)
			log.Fatal(err)
		}

		//copy file to disk
		_, err = io.Copy(dst, src)
		if err != nil {
			fmt.Println("Error from copy: ", err)
			log.Fatal(err)
		}

		if err := dst.Close(); err != nil {
			fmt.Println("Error from close (dist): ", err)
			log.Fatal(err)
		}

		if err := src.Close(); err != nil {
			fmt.Println("Error from close (src): ", err)
			log.Fatal(err)
		}
	}
	return nil
}

func ShowTable(allMatchScores matchclient.AllMatchScoresResponse) {

	fmt.Println("")
	fmt.Println("Match Score Comparisons:")

	t := tabby.New()
	t.AddHeader("FILENAME_1", "FILENAME_2", "MATCH_SCORE")
	for _, line := range allMatchScores {
		t.AddLine(line.File1Name, line.File2Name, line.MatchScore)
	}
	t.Print()

}
