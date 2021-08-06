package matchClient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	BaseURL = "http://localhost:8080/biometric"
)

// Client
/*type Client struct {
	apiKey     string
	baseURL    string
	userAgent  string
	HTTPClient *http.Client
}

// NewClient .
func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
		baseURL: BaseURL,
	}
}*/

type errorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

type successResponse struct {
	MatchResult int    `json:"matchResult"`
	FileName1   string `json:"fileName1"`
	FileId1     string `json:"fileId1"`
	FileName2   string `json:"fileName1"`
	FileId2     string `json:"fileId1"`
}

func HelloViki() {

	//Basic HTTP Get request
	url := BaseURL + "/hello/Viki"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	//Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
	}

	fmt.Printf("%s\n", body)
}

func post_files() (err error) {

	url := BaseURL + "/match?=files"
	fileDir, _ := os.Getwd()
	fileName := "1.png"
	filePath := path.Join(fileDir, fileName)

	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("files", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		return err
	}
	return

}
