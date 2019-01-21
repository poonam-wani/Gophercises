package main

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"runtime/debug"
	"testing"

	"github.com/alecthomas/chroma"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	temp := listenAndServeFunction
	defer func() {
		listenAndServeFunction = temp
	}()
	listenAndServeFunction = func(port string, handle http.Handler) error {
		panic("testing")
	}
	assert.PanicsWithValuef(t, "testing", main, "they should be equal")
}

func TestDevMw(t *testing.T) {
	handler := http.HandlerFunc(panicDemo)
	req, err := http.NewRequest("Get", "localhost:4000/panic/", nil)
	if err != nil {
		t.Errorf("%s error occured", err)
	}
	response := httptest.NewRecorder()
	hand := devMw(handler)
	hand.ServeHTTP(response, req)
}

func TestSource(t *testing.T) {

	tmpconversion := FunctionConversion
	defer func() {
		FunctionConversion = tmpconversion
	}()

	FunctionConversion = func(s string) (int, error) {
		return 0, errors.New("Error while Converting")
	}

	tmpcopy := FunctionCopy
	defer func() {
		FunctionCopy = tmpcopy
	}()

	FunctionCopy = func(a io.Writer, b io.Reader) (int64, error) {
		return 0, errors.New("Error while copying")
	}

	request, _ := http.NewRequest("GET", "localhost:4000/debug/?line=76&path=%2Fhome%2Fpoonam%2Fgithub.com%2Fpoonam-wani%2Fgophercises%2Frecover%2Fmain.go", nil)

	response := httptest.NewRecorder()
	sourceCodeHandler(response, request)

}

func TestSourceCodeHandler(t *testing.T) {

	tmpstyle := FunctionStyle
	defer func() {
		FunctionStyle = tmpstyle
	}()

	FunctionStyle = func(s string) *chroma.Style {
		return nil
	}

	testinput := []struct {
		url          string
		responceCode int
	}{

		{
			url:          "line=76&path=%2Fhome%2Fpoonam%2Fgithub.com%2Fpoonam-wani%2Fgophercises%2Frecover%2Fmain.go",
			responceCode: 200,
		}, {
			url:          "line=19&path=%2Fusr%2Flocal%2Fgo%2Fsrc%2Fhttp%2Fserver.go",
			responceCode: 500,
		},
	}
	for i := 0; i < len(testinput); i++ {

		request, err := http.NewRequest("GET", "localhost:4000/debug/?"+testinput[i].url, nil)
		if testinput[i].url == "" {
			t.Errorf("URL is not null")
		}
		if err != nil {
			t.Errorf("URL is not present, corresponding response code is %d", testinput[i].responceCode)
		}
		response := httptest.NewRecorder()
		sourceCodeHandler(response, request)
		res := response.Result()
		if res.StatusCode != testinput[i].responceCode {
			t.Errorf("Expected response code %d for URL %q, got response code %d", testinput[i].responceCode, testinput[i].url, res.StatusCode)
		}
	}

}
func TestHello(t *testing.T) {
	request, _ := http.NewRequest("GET", "localhost:4000", nil)
	response := httptest.NewRecorder()
	hello(response, request)
	_, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

}
func TestMakeLinks(t *testing.T) {
	stack := debug.Stack()
	hyperlink := makeLinks(string(stack))
	if hyperlink == "" {
		t.Errorf("TestFail beacuse of blank hyperlink")
	}
}
