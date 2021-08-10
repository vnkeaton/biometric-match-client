package matchclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vnkeaton/biometric-match-client/mocks"
)

func init() {
	Client = &mocks.MockClient{}
}

type AllMatchScores []struct {
	ID         string
	Dir1       string
	File1Name  string
	Dir2       string
	File2Name  string
	MatchScore float64
}

// TestHelloName calls matchclient.Hello with a name,
// checking for an error.
func TestHelloName(t *testing.T) {
	expected := "Hello Viki"
	r := ioutil.NopCloser(bytes.NewReader([]byte(expected)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	resp, err := Hello("Viki")
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, expected, resp)
}

// TestHelloEmpty calls matchclient.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
	expected := "Hello World"
	r := ioutil.NopCloser(bytes.NewReader([]byte(expected)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	resp, err := Hello("") // empty
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, expected, resp)
}

// TestUploadFiles calls matchclient.UploadFiles
func TestMatchFiles(t *testing.T) {
	expectedJson := `{"matchResult": 6,"fileName1": "1.png","fileName2": "2.png"}`
	var matchScore MatchScoreData
	matchScore.FileName1 = "1.png"
	matchScore.FileName2 = "2.png"
	matchScore.MatchScore = 6
	r := ioutil.NopCloser(bytes.NewReader([]byte(expectedJson)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	images := []string{"./testData/1.png", "./testData/2.png"}
	resp, err := MatchFiles(images)
	fmt.Println(resp)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, matchScore.FileName1, resp.FileName1)
	assert.EqualValues(t, matchScore.FileName2, resp.FileName2)
	assert.EqualValues(t, matchScore.MatchScore, resp.MatchScore)
}

// TestGetAllMatchScores calls matchclient.GetAllMatchScores
func TestGetAllMatchScores(t *testing.T) {
	expectedJson := `[{"id": "1","dir1": "images","file1Name": "1.png","dir2": "images","file2Name": "6.png","matchScore": 3.00}]`
	fmt.Println(expectedJson)
	r := ioutil.NopCloser(bytes.NewReader([]byte(expectedJson)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	resp, _ := GetAllMatchScores()
	fmt.Println(resp)
	assert.NotNil(t, resp)
	//assert.Nil(t, err)
	/*for _, r := range resp {
		assert.EqualValues(t, "images", r.Dir1)
		assert.EqualValues(t, "images", r.Dir2)
		assert.EqualValues(t, "1,png", r.File1Name)
		assert.EqualValues(t, "2,png", r.File2Name)
		assert.EqualValues(t, "3.00", r.MatchScore)
	}*/
}

//TODO Should have more tests but I am running out of time

//references:
//https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/
