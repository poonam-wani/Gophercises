package main

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/poonam-wani/gophercises/image/primitive"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {

	temp := ListenAndServeFunction
	defer func() {
		ListenAndServeFunction = temp
	}()
	ListenAndServeFunction = func(port string, handle http.Handler) error {
		panic("testing")
	}
	assert.PanicsWithValuef(t, "testing", main, "they should be equal")
}

func TestStartProgram(t *testing.T) {
	request, _ := http.NewRequest("GET", "localhost:4000/", nil)
	response := httptest.NewRecorder()
	startProgram(response, request)
	_, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

}

func TestUploadImage(t *testing.T) {

	// This function is written to cover all positive sceanrio of upload image
	file, err := os.Open("monalisa.png")
	if err != nil {
		t.Errorf("%s error got while uploading image file", err)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("image", "monalisa.png") // creates form header with new field name & file name
	if err != nil {
		t.Errorf("%s error got while creating form file", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Errorf("%s error got while copying image file ", err)
	}
	writer.Close()

	request, _ := http.NewRequest("POST", "localhost:4000/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	uploadImage(response, request)
	res := response.Result()
	if res.StatusCode != 302 {
		t.Errorf("got response code %d", res.StatusCode)
	}

}

func TestUploadNegative(t *testing.T) {

	// This function is to cover the negative of temporary file
	file, _ := os.Open("")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "mona.png")

	tmpinoutfile := FunctioninoutFile
	defer func() {
		FunctioninoutFile = tmpinoutfile

	}()
	var f *os.File
	FunctioninoutFile = func(string, string) (*os.File, error) {
		return f, errors.New("Error while creating temp file")
	}
	io.Copy(part, file)
	writer.Close()

	request, _ := http.NewRequest("POST", "localhost:4000/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	uploadImage(response, request)
	res := response.Result()
	if res.StatusCode != 500 {
		t.Errorf("got response code %d", res.StatusCode)
	}

	// passing the invalid field name in createFormFile
	file, _ = os.Open("monalisa.png")
	body = &bytes.Buffer{}
	writer = multipart.NewWriter(body)
	part, _ = writer.CreateFormFile("profile", "mona.png")

	io.Copy(part, file)
	writer.Close()

	request, _ = http.NewRequest("POST", "localhost:4000/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response = httptest.NewRecorder()
	uploadImage(response, request)
	res = response.Result()
	if res.StatusCode != 400 {
		t.Errorf("got response code %d", res.StatusCode)
	}

}

func TestUploadNeg(t *testing.T) {

	// This function is to cover the negative of copy function
	file, _ := os.Open("mona.png")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "mona.png")
	tmpcopy := FunctionCopy
	defer func() {
		FunctionCopy = tmpcopy

	}()

	FunctionCopy = func(io.Writer, io.Reader) (int64, error) {
		return -1, errors.New("Error while copying image file")
	}
	_, err := io.Copy(part, file)
	if err == nil {
		t.Errorf("%s error got while copying image file ", err)
	}
	writer.Close()

	request, _ := http.NewRequest("POST", "localhost:4000/upload", body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	response := httptest.NewRecorder()
	uploadImage(response, request)
	res := response.Result()
	if res.StatusCode != 500 {
		t.Errorf("got response code %d", res.StatusCode)
	}

}

func TestModifyNeg(t *testing.T) {

	request, _ := http.NewRequest("GET", "localhost:4000/modify/523912166.png?mode=7", nil)
	response := httptest.NewRecorder()
	modifyImage(response, request)
	res := response.Result()
	if res.StatusCode != 500 {
		t.Errorf("%d response got while modifying image", res.StatusCode)
	}

	// Covers string to integer conversion part for mode
	request, _ = http.NewRequest("GET", "localhost:4000/modify/523912166.png?mode=aaaa", nil)
	response = httptest.NewRecorder()
	modifyImage(response, request)
	res = response.Result()
	if res.StatusCode != 400 {
		t.Errorf("%d response got while modifying image", res.StatusCode)
	}

	// Covers string to integer conversion part for n
	request, _ = http.NewRequest("GET", "localhost:4000/modify/523912166.png?mode=7&n=wer", nil)
	response = httptest.NewRecorder()
	modifyImage(response, request)
	res = response.Result()
	if res.StatusCode != 400 {
		t.Errorf("%d response got while modifying image", res.StatusCode)
	}

	// Covers modestr part
	request, _ = http.NewRequest("GET", "localhost:4000/img/modify/", nil)
	response = httptest.NewRecorder()
	modifyImage(response, request)
	res = response.Result()
	if res.StatusCode != 500 {
		t.Errorf("%d response got while modifying image", res.StatusCode)
	}

	// Covers the open image file part
	request, _ = http.NewRequest("GET", "localhost:4000/img/modify/", nil) // pass Invalid file path to open
	tmpopen := FunctionOpen
	defer func() {
		FunctionOpen = tmpopen

	}()
	var f *os.File
	FunctionOpen = func(string) (*os.File, error) {
		return f, errors.New("Error while opening the file")
	}

	response = httptest.NewRecorder()
	modifyImage(response, request)
	res = response.Result()
	if res.StatusCode != 400 {
		t.Errorf("%d response got while modifying image", res.StatusCode)
	}

}

func TestModifyImage(t *testing.T) {

	request, _ := http.NewRequest("GET", "localhost:4000/img/modify/523912166.png?mode=0&n=10", nil)
	response := httptest.NewRecorder()
	modifyImage(response, request)
	res := response.Result()
	if res.StatusCode != 302 {
		t.Errorf("%d response got while modifying image", res.StatusCode)
	}

}

func TestRenderNumShapeChoices(t *testing.T) {

	var rs1 io.ReadSeeker
	rs1, _ = os.Open("monalisa.png") // Positive scenario
	request, _ := http.NewRequest("GET", "localhost:4000/img/modify/523912166.png?mode=0&n=10", nil)
	response := httptest.NewRecorder()
	renderNumShapeChoices(response, request, rs1, ".png", primitive.ModeBeziers)
	res := response.Result()
	if res.StatusCode != 200 {
		t.Errorf("%d response got while rendering shapes of image", res.StatusCode)
	}

	// covers the negative scenrio of template
	tmptemplate := FunctionTemplate
	defer func() {
		FunctionTemplate = tmptemplate

	}()

	FunctionTemplate = func(t *template.Template, err error) *template.Template {
		tp := template.Must(template.New("").Parse(""))
		tp.Tree = nil
		return tp
	}
	rs1, _ = os.Open("monalisa.png") // Negative scenario
	request, _ = http.NewRequest("GET", "localhost:4000/img/modify/523912166.png", nil)
	response = httptest.NewRecorder()
	renderNumShapeChoices(response, request, rs1, ".png", primitive.ModeEllipse)

	res = response.Result()
	if res.StatusCode != 500 {
		t.Errorf("%d response in template of rendernum shapes", res.StatusCode)
	}

}

func TestRenderModeChoices(t *testing.T) {
	var rs1 io.ReadSeeker
	rs1, _ = os.Open("monalisa.png") // Positive scenario
	request, _ := http.NewRequest("GET", "localhost:4000/img/modify/523912166.png?mode=0&n=10", nil)
	response := httptest.NewRecorder()
	renderModeChoices(response, request, rs1, ".png")
	res := response.Result()
	if res.StatusCode != 200 {
		t.Errorf("%d response got while rendering modes of image", res.StatusCode)
	}

	// covers the negative scenario of template
	tmptemplate := FunctionTemplate
	defer func() {
		FunctionTemplate = tmptemplate

	}()

	FunctionTemplate = func(t *template.Template, err error) *template.Template {
		tp := template.Must(template.New("").Parse(""))
		tp.Tree = nil
		return tp
	}
	rs1, _ = os.Open("monalisa.png") // Negative scenario
	request, _ = http.NewRequest("GET", "localhost:4000/img/modify/523912166.png", nil)
	response = httptest.NewRecorder()
	renderModeChoices(response, request, rs1, ".png")

	res = response.Result()
	if res.StatusCode != 500 {
		t.Errorf("%d response in template in render mode choices", res.StatusCode)
	}

}

func TestGenImage(t *testing.T) {
	var r1 io.Reader
	r1, err := os.Open("monalisa.png") // Positive scenario
	_, err = genImage(r1, ".png", 10, 1)
	if err != nil {
		t.Errorf("%s error occured while generating image", err)
	}

	r1, err = os.Open("mona.png") // Negative scenario when error occured while transforming image
	_, err = genImage(r1, ".png", 10, 1)
	if err == nil {
		t.Errorf("%s error occured while transforming image ", err)
	}

	tmpinoutfile := FunctioninoutFile
	defer func() {
		FunctioninoutFile = tmpinoutfile

	}()
	var f *os.File
	FunctioninoutFile = func(prefix, ext string) (*os.File, error) {
		return f, errors.New("Error while creating output file")
	}
	r1, err = os.Open("monalisa.png") // Negative scenario when error occured while creating output image
	_, err = genImage(r1, ".png", 10, 3)
	if err == nil {
		t.Errorf("%s error occured while creating output file ", err)
	}
}
func TestTempfile(t *testing.T) {
	_, err := tempfile("", "") // Positive scenario
	if err != nil {
		t.Errorf("%s error while creating temporary file", err)
	}

	tmpfile := FunctionTempfile
	defer func() {
		FunctionTempfile = tmpfile

	}()
	var f *os.File
	FunctionTempfile = func(dir string, prefix string) (*os.File, error) {
		return f, errors.New("Error while creating temp file")
	}
	_, err = tempfile("", "in_") //  Negative scenario
	if err == nil {
		t.Errorf("%s error while creating temporary file", err)
	}
}
