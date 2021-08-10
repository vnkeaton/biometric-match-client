package matchclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vnkeaton/biometric-match-client/mocks"
)

func init() {
	Client = &mocks.MockClient{}
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
func TestUploadFiles(t *testing.T) {

	expectedJson := `{"matchResult": 6.00,"fileName1": "3.png","fileName2": "5.png"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(expectedJson)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	resp, err := UploadFiles(gomock.Any(), gomock.Any)
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedJson, resp)

}

// TestGetAllMatchScores calls matchclient.GetAllMatchScores
func TestGetAllMatchScores(t *testing.T) {

	expectedJson := `[{"id": "1","dir1": "images","file1Name": "6.png",dir2": "images","file2Name": "2.png","matchScore": 3}]`
	r := ioutil.NopCloser(bytes.NewReader([]byte(expectedJson)))
	mocks.GetDoFunc = func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	resp, err := GetAllMatchScores()
	assert.NotNil(t, resp)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedJson, resp)
}

//TODO Should have more tests but I am running out of time

//references:
//https://www.thegreatcodeadventure.com/mocking-http-requests-in-golang/
