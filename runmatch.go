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

const (
	//facesURL  = "https://github.com/TheMdTF/mdtf-public/tree/master/rally2-matching-system/tests/test-routine-images/face"
	facesURL  = "https://github.com/TheMdTF/mdtf-public"
	imagesDir = "/images/"
	repoDir   = "/tree/master/rally2-matching-system/tests/test-routine-images/face"
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

//func mustOpen(f string) *os.File {
func mustOpen(p string, f string) string {
	filename := p + imagesDir + f
	//fmt.Println("check to see if filename exists with mustOpen for: " + filename)
	_, err := os.Open(filename)
	if err != nil {
		fmt.Println("This file failed upon mustOpen: " + filename)
		log.Fatal(err)
	}
	return filename
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	matchclient.Hello("IDSL")

	path, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting pwd ", err)
		log.Fatal(err)
	}

	//get the faces repo
	/*downloadFaces(path)
	if err != nil {
		fmt.Println("Error from downloadFaces  ", err)
		log.Fatal(err)
	}
	fmt.Println("faces downloaded")
	*/

	//fmt.Println("working directory is: " + path)

	//read in list of images
	//fmt.Println("read in list of images from: " + path + imagesDir)
	imageFaces, err := ioutil.ReadDir(path + imagesDir)
	if err != nil {
		fmt.Println("Error from read dir  ", err)
		log.Fatal(err)
	}

	revFiles, err := ioutil.ReadDir(path + imagesDir)
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
			//fmt.Println("Comparing " + f.Name() + " with " + r.Name())
			mediaFiles := []string{mustOpen(path, f.Name()), mustOpen(path, r.Name())}
			//matchScore, err := matchclient.MatchFiles(mediaFiles)
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

func downloadFaces(path string) error {
	var storer *memory.Storage
	var origin billy.Filesystem

	storer = memory.NewStorage()
	origin = memfs.New()

	_, err := git.Clone(storer, origin, &git.CloneOptions{
		URL: facesURL,
	})
	checkError(err)
	fmt.Println("Repository cloned")

	//read in list of images
	memFaces, err := origin.ReadDir("/rally2-matching-system/tests/test-routine-images/face")
	if err != nil {
		fmt.Println("Error from read dir  ", err)
		log.Fatal(err)
	}
	//file, err := fs.Open("/rally2-matching-system/tests/test-routine-images/face/1.png")
	//fmt.Println("we have an image file: " + file.Name())
	//checkError(err)

	fmt.Println("we have a list of image files from the repository")

	if _, err := os.Stat(path + imagesDir); os.IsNotExist(err) {
		err := os.Mkdir(path+imagesDir, 0755)
		checkError(err)
	}
	fmt.Println("new images directory created")

	for _, f := range memFaces {
		fmt.Println("file is: " + f.Name())
		fmt.Println(f.Size())

		ext := filepath.Ext(f.Name())
		if ext == ".png" {
			fmt.Println("ext is .png")

			//open a file
			src, err := origin.Open(f.Name())
			if err != nil {
				fmt.Println("Error from open: ", err)
				log.Fatal(err)
			}

			//create a new file
			dst, err := os.Create(f.Name())
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

			/*fmt.Println("do the sizes match?  f.size and n")
			fmt.Println(f.Size())
			fmt.Println(nsize)
			*/

			fmt.Println("new image file created in : " + path + imagesDir + f.Name())
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
