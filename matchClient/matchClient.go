package matchClient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	FileName2   string `json:"fileName1"`
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

func MatchFiles(values []string) {
	dst := BaseURL + "/image/match"
	UploadFiles(dst, values)
}

func UploadFiles(dst string, values []string) (err error) {

	u, err := url.Parse(dst)
	if err != nil {
		return fmt.Errorf("failed to parse destination url: %w", err)
	}

	// Prepare a form that you will submit to that URL.
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	//buffSize := 0
	for _, fname := range values {
		fd, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("failed to open file to upload: %w", err)
		}
		defer fd.Close()

		stat, err := fd.Stat()
		if err != nil {
			return fmt.Errorf("failed to query file info: %w", err)
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
			return fmt.Errorf("failed to creae new form part: %w", err)
		}
		n, err := io.Copy(part, fd)
		if err != nil {
			return fmt.Errorf("failed to write form part: %w", err)
		}
		if int64(n) != stat.Size() {
			return fmt.Errorf("file size changed while writing: %s", fd.Name())
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

	DumpRequest(&req)

	// Don't forget to set the content type, this will contain the boundary.
	req.Header.Set("Content-Type", writer.FormDataContentType())
	//req.Header.Set("Content-Length", b.Len())

	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return fmt.Errorf("failed to perform http request: %w", err)
	}
	defer resp.Body.Close()

	_, _ = io.Copy(os.Stdout, resp.Body)

	// Check the response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", resp.Status)
		return err
	}
	return nil
}

func DumpRequest(req *http.Request) {

	output, err := httputil.DumpRequest(req, false)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println(string(output))
}

//references:
//stackoverflow.com/questions/20205796/post-data-using-content-type-multipart-form-data
//https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
//https://pkg.go.dev/net/http
//https://ayada.dev/posts/multipart-requests-in-go/
//https://stackoverflow.com/questions/63636454/golang-multipart-file-form-request
