package primitive

import (
	"errors"
	"io"
	"os"
	"testing"
)

func TestWithMode(t *testing.T) {

	WithMode(ModeBeziers)

}
func TestTransform(t *testing.T) {
	var img io.Reader
	img, err := os.Open("../monalisa.png") // Positive scenario
	img, err = Transform(img, ".png", 10, []string{"-m", "1"})
	if err != nil {
		t.Errorf("Expected nil but got %s error ", err)
	}

	img, err = os.Open("monalisa.png") // Negative scenario which fails while copying image to file
	img, err = Transform(img, ".png", 10, []string{"-m", "1"})
	if err == nil {
		t.Errorf("%s error while copying image to file ", err)
	}

	img, err = os.Open("../monalisa.png") // Negative scenario where not able to run primitive cmd
	img, err = Transform(img, "", 10, []string{"-m", "1"})
	if err == nil {
		t.Errorf("%s error while running primitive command ", err)
	}

	tmpcopy := FunctionCopy
	defer func() {
		FunctionCopy = tmpcopy

	}()

	FunctionCopy = func(dst io.Writer, src io.Reader) (int64, error) {
		return -1, errors.New("Error while copy output file in buffer")
	}

	img, err = os.Open("../monalisa.png") // Negative scenario which fails while copying output file to buffer
	img, err = Transform(img, "png", 10, []string{"-m", "1"})
	if err == nil {
		t.Errorf("%s error while copy output file in buffer ", err)
	}

	img, err = os.Open("../monalisa.png") // Negative scenario which fails while creating temp input files
	img, err = Transform(img, "/", 10, []string{"-m", "1"})
	if err == nil {
		t.Errorf("%s error while creating temp input files ", err)
	}

	tmpinoutfile := FunctioninoutFile
	defer func() {
		FunctioninoutFile = tmpinoutfile

	}()
	var f *os.File
	FunctioninoutFile = func(prefix, ext string) (*os.File, error) {
		return f, errors.New("Error while creating input file")
	}

	img, err = os.Open("../monalisa.png") // Negative scenario which fails while creating temp output files
	img, err = Transform(img, " ", 10, []string{"-m", "1"})
	if err == nil {
		t.Errorf("%s error while creating temp output files ", err)
	}

}

func TestPrimitive(t *testing.T) {

	_, err := primitive("monalisa.png", "out.png", 10, "temp")
	if err == nil {
		t.Errorf("%s got error for positive sceanrio ", err)
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
