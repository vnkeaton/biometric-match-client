package matchclient

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

type MatchScoreData struct {
	FileName1  string
	FileName2  string
	MatchScore float64
}

type AllMatchScoresResponse []struct {
	ID         string  `json:"id"`
	Dir1       string  `json:"dir1"`
	File1Name  string  `json:"file1Name"`
	Dir2       string  `json:"dir2"`
	File2Name  string  `json:"file2Name"`
	MatchScore float64 `json:"matchScore"`
}

type AllMatchScoreData struct {
	ID         string
	Dir1       string
	File1Name  string
	Dir2       string
	File2Name  string
	MatchScore float64
}

// HTTPClient interface - makes testing easier
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

//only at start up
func init() {
	Client = &http.Client{}
}

// Hello calls the api endpoint to say hello
func Hello(name string) (string, error) {

	//basic HTTP Get request
	url := BaseURL + "/hello/" + name
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}

	//set header and call client api
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := Client.Do(req)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}
	defer resp.Body.Close()

	//read body from response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body. ", err)
	}

	//print repsonse
	fmt.Printf("%s\n", body)
	return string(body), nil
}

// MatchFiles sets up call for matching images
func MatchFiles(values []string) (MatchScoreData, error) {
	dst := BaseURL + "/image/match"
	//fmt.Println("call upload files")
	matchScore, err := UploadFiles(dst, values)
	if err != nil {
		return MatchScoreData{}, fmt.Errorf("failed to get match score data: %w", err)
	}
	return matchScore, nil
}

// UploadFiles will take 2 image files and compare them
//TODO desparately needs refactoring
func UploadFiles(dst string, values []string) (MatchScoreData, error) {
	//parse the url
	u, err := url.Parse(dst)
	if err != nil {
		return MatchScoreData{}, fmt.Errorf("failed to parse destination url: %w", err)
	}

	//prepare a form for submitting to the URL
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	//get each image file
	for _, fname := range values {
		//open image file
		fd, err := os.Open(fname)
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to open file to upload: %w", err)
		}
		defer fd.Close()

		//get image file info
		stat, err := fd.Stat()
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to query file info: %w", err)
		}

		//init multipart header for image file
		hdr := make(textproto.MIMEHeader)
		cd := mime.FormatMediaType("form-data", map[string]string{
			"name":     "files",
			"filename": fname,
		})
		hdr.Set("Content-Disposition", cd)
		hdr.Set("Content-Type", "image/png")

		//create mulitpart section for image file
		part, err := writer.CreatePart(hdr)
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to creae new form part: %w", err)
		}

		//copy image file into multipart section
		n, err := io.Copy(part, fd)
		if err != nil {
			return MatchScoreData{}, fmt.Errorf("failed to write form part: %w", err)
		}
		if int64(n) != stat.Size() {
			return MatchScoreData{}, fmt.Errorf("file size changed while writing: %s", fd.Name())
		}
	}
	//close the multipart writer to ensure terminating boundary.
	writer.Close()

	//init http post header
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

	//set conent type and boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	//call the api client
	resp, err := Client.Do(&req)
	if err != nil {
		return MatchScoreData{}, fmt.Errorf("failed to perform http request: %w", err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	// resp body is []byte
	//_, _ = io.Copy(os.Stdout, resp.Body) //print to stdOut for debugging

	//check the status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code not ok")
		return MatchScoreData{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	//create json response, response body is []bytes to the go struct ptr
	matchScoreBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("matchScore body has issues")
		return MatchScoreData{}, fmt.Errorf("error reading body: %w", err)
	}

	//unmarshall json
	var matchScore MatchScoreData
	jsonErr := json.Unmarshal(matchScoreBody, &matchScore)
	if jsonErr != nil {
		return MatchScoreData{}, fmt.Errorf("can not unmarshal Json: %w", err)
	}

	//return match score structure
	return MatchScoreData{
		FileName1:  matchScore.FileName1,
		FileName2:  matchScore.FileName2,
		MatchScore: matchScore.MatchScore,
	}, nil
}

// DumpRequest is for debugging purposes
func DumpRequest(req *http.Request) {
	output, err := httputil.DumpRequest(req, false)
	if err != nil {
		fmt.Println("Error dumping request:", err)
		return
	}
	fmt.Println(string(output))
}

// GetAllMatchScores calls the api endpoint to retrieve all match scores
func GetAllMatchScores() ([]AllMatchScoreData, error) {
	//basic HTTP Get request
	url := BaseURL + "/matchscore/downloadFile/all"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error reading response. ", err)
	}

	//set header and call api endpoint
	req.Header.Set("Cache-Control", "no-cache")
	resp, err := Client.Do(req)
	if err != nil {
		return []AllMatchScoreData{}, fmt.Errorf("failed to perform http request: %w", err)
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	//check the status code
	if resp.StatusCode != http.StatusOK {
		fmt.Println("status code not ok")
		return []AllMatchScoreData{}, fmt.Errorf("bad status: %s", resp.Status)
	}

	//create json response, response body is []bytes to the go struct ptr
	matchScoreBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("matchScore body has issues")
		return []AllMatchScoreData{}, fmt.Errorf("error reading body: %w", err)
	}

	//unmarshal json response array into the data structure array
	var matchScores []AllMatchScoreData
	jsonErr := json.Unmarshal(matchScoreBody, &matchScores)
	if jsonErr != nil {
		return []AllMatchScoreData{}, fmt.Errorf("can not unmarshal Json: %w", err)
	}
	return matchScores, nil
}

//TODO Need more client functions here, running out of time....

//references:
//stackoverflow.com/questions/20205796/post-data-using-content-type-multipart-form-data
//https://gist.github.com/mattetti/5914158/f4d1393d83ebedc682a3c8e7bdc6b49670083b84
//https://pkg.go.dev/net/http
//https://ayada.dev/posts/multipart-requests-in-go/
//https://stackoverflow.com/questions/63636454/golang-multipart-file-form-request
//https://blog.alexellis.io/golang-json-api-client/
