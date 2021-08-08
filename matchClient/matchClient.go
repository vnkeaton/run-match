package matchClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"time"
)

const (
	BaseURL = "http://localhost:8080/biometric"
)

type ErrorResponse struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type MatchScoreResponse struct {
	MatchResult float64 `json:"matchResult"`
	FileName1   string  `json:"fileName1"`
	FileName2   string  `json:"fileName2"`
}

type AllMatchScoresResponse []struct {
	ID         string  `json:"id"`
	Dir1       string  `json:"dir1"`
	File1Name  string  `json:"file1Name"`
	Dir2       string  `json:"dir2"`
	File2Name  string  `json:"file2Name"`
	MatchScore float64 `json:"matchScore"`
}

type MatchScoreData struct { //must be a capital letter to be exported and the fields
	FileName1  string
	FileName2  string
	MatchScore float64
}

func Hello() {

	//Basic HTTP Get request
	url := BaseURL + "/hello/IDSL"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}
	defer resp.Body.Close()

	//Read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
	}

	fmt.Printf("%s\n", body)
}

func MatchFiles(values []string) (MatchScoreData, error) {
	dst := BaseURL + "/image/match"
	fmt.Println("call upload files")
	matchScore, err := UploadFiles(dst, values)
	if err != nil {
		return MatchScoreData{}, fmt.Errorf("failed to get match score data: %w", err)
	}
	return matchScore, nil
}

func UploadFiles(dst string, values []string) (MatchScoreData, error) {

	u, err := url.Parse(dst)
	if err != nil {
		return MatchScoreData{}, fmt.Errorf("failed to parse destination url: %w", err)
	}

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	//buffSize := 0
	for _, fname := range values {
		fd, err := os.Open(fname)
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to open file to upload: %w", err)
		}
		defer fd.Close()

		stat, err := fd.Stat()
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to query file info: %w", err)
		}

		hdr := make(textproto.MIMEHeader)
		cd := mime.FormatMediaType("form-data", map[string]string{
			"name":     "files",
			"filename": fname,
		})
		hdr.Set("Content-Disposition", cd)
		hdr.Set("Content-Type", "image/png")
		//hdr.Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

		part, err := writer.CreatePart(hdr)
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to creae new form part: %w", err)
		}
		n, err := io.Copy(part, fd)
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to write form part: %w", err)
		}
		if int64(n) != stat.Size() {
			return MatchScoreData{}, fmt.Errorf("file size changed while writing: %s", fd.Name())
		}
	}
	// Don't forget to close the multipart writer.
	// If you don't close it, your request will be missing the terminating boundary.
	writer.Close()

	//q := u.Query()
	//q.Set("access_token", token)
	//u.RawQuery = q.Encode()

	hdr := make(http.Header)
	hdr.Set("Content-Type", writer.FormDataContentType())

	// Now that you have a form, you can submit it to your handler.
	req := http.Request{
		Method: "POST",
		URL:    u,
		Header: hdr,
		Body:   ioutil.NopCloser(&b),
		//ContentLength: int64(form.contentLen),
	}

	//DumpRequest(&req) //for debugging

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//req.Header.Set("Content-Length", b.Len())

	// Call the api client
	fmt.Println("make the api call")
	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return MatchScoreData{}, fmt.Errorf("failed to perform http request: %w", err)
	}
	if resp.Body != nil {
		fmt.Println("body not nil")
		defer resp.Body.Close()
	}

	// resp body is []byte
	//_, _ = io.Copy(os.Stdout, resp.Body) //print to stdOut for debugging

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code not ok")
		return MatchScoreData{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("create json response")

	//create json response, response body is []bytes to the go struct ptr
	matchScoreBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("matchScore body has issues")
		return MatchScoreData{}, fmt.Errorf("error reading body: %w", err)
	}

	fmt.Println("unmarshall byte code")

	var matchScore MatchScoreResponse
	jsonErr := json.Unmarshal(matchScoreBody, &matchScore)
	if jsonErr != nil {
		fmt.Println("am i here?")
		return MatchScoreData{}, fmt.Errorf("can not unmarshal Json: %w", err)
	}

	fmt.Println("return match score structure")

	return MatchScoreData{
		FileName1:  matchScore.FileName1,
		FileName2:  matchScore.FileName2,
		MatchScore: matchScore.MatchResult,
	}, nil
}

func DumpRequest(req *http.Request) {

	output, err := httputil.DumpRequest(req, false)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println(string(output))
}

func GetAllMatchScores() (AllMatchScoresResponse, error) {

	//Basic HTTP Get request
	url := BaseURL + "/matchscore/downloadFile/all"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return AllMatchScoresResponse{}, fmt.Errorf("failed to perform http request: %w", err)
	}
	if resp.Body != nil {
		fmt.Println("body not nil")
		defer resp.Body.Close()
	} // Check the status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code not ok")
		return AllMatchScoresResponse{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	fmt.Println("create json response")

	//create json response, response body is []bytes to the go struct ptr
	matchScoreBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("matchScore body has issues")
		return AllMatchScoresResponse{}, fmt.Errorf("error reading body: %w", err)
	}

	fmt.Println("unmarshall byte code")

	var matchScores AllMatchScoresResponse
	jsonErr := json.Unmarshal(matchScoreBody, &matchScores)
	if jsonErr != nil {
		fmt.Println("am i here?")
		return AllMatchScoresResponse{}, fmt.Errorf("can not unmarshal Json: %w", err)
	}

	return matchScores, nil
}

//references:
//stackoverflow.com/questions/20205796/post-data-using-content-type-multipart-form-data
//https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
//https://pkg.go.dev/net/http
//https://ayada.dev/posts/multipart-requests-in-go/
//https://stackoverflow.com/questions/63636454/golang-multipart-file-form-request
//https://blog.alexellis.io/golang-json-api-client/
