package secret

import (
	"crypto/cipher"
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
)

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, "temp.txt")

}

func TestSet(t *testing.T) {
	file := secretsPath()

	vault1 := File("sample", file)
	err := vault1.Set("apiKey", "123455") // Positive- Passing valid key & value to set function
	if err != nil {
		t.Errorf("%s Error while setting value to key", err)
	}

	vault1 = File("sample", file)
	err = vault1.Set("", "") // Negative- Passing blank value to set function
	if err == nil {
		t.Errorf("Error when u try to set key & value both blank")
	}
	vault1 = File("sample", "") // Negative- Passing invalid file name to set function
	err = vault1.Set("testkey", "45")

	vault1 = File("test", "./") //Negative - Passing invalid file to load
	err = vault1.load()
	if err == nil {
		t.Errorf("Error while getting value from provided key")
	}
	err = vault1.Set("apiKey", "45")

}

func TestLoad(t *testing.T) {

	file := secretsPath()
	vault1 := File("sample", file)
	err := vault1.load() // Positive- Passing valid to load
	if err == nil {
		t.Errorf("Error while loading the provided file")
	}

}

func TestSave(t *testing.T) {

	tmpsave := FunctionSave
	defer func() {
		FunctionSave = tmpsave

	}()
	var w *cipher.StreamWriter
	FunctionSave = func(key string, a io.Writer) (*cipher.StreamWriter, error) {
		return w, errors.New("Error while saving")
	}

	file := secretsPath()
	vault1 := File("sample", file)
	err := vault1.save() // Positive
	if err == nil {
		t.Errorf("Error while loading the provided file")
	}

}

func TestGet(t *testing.T) {

	file := secretsPath()
	f, _ := os.Open(file)
	f.Truncate(0)
	f.Close()

	vault2 := File("sample", file)
	err := vault2.Set("api", "123455")
	_, err = vault2.Get("api") //Positive - Passing valid key to get function
	if err == nil {
		t.Errorf("%s Error while getting value from provided key", err)
	}

	vault1 := File("sample", "") //Negative- Passing invalid file name
	_, err = vault1.Get("testKey")
	if err == nil {
		t.Errorf("Error while getting value from provided key")
	}

	vault1 = File("sample", "./") //Negative - Passing invalid file to load
	err = vault1.load()
	if err == nil {
		t.Errorf("Error while getting value from provided key")
	}
	os.Remove(file)
}
