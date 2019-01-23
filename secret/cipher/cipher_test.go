package cipher

import (
	"crypto/cipher"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	homedir "github.com/mitchellh/go-homedir"
)

var tempcipher string

// Vault It is the structure of vault
type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyValues   map[string]string
}

// File taking the input as encodingKey & filepath
func File(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
	}
}

func TestEncryptCipher(t *testing.T) {

	tmpnewcipherblock := FunctionNewCipherblock
	defer func() {
		FunctionNewCipherblock = tmpnewcipherblock
	}()
	var a cipher.Block
	FunctionNewCipherblock = func(key string) (cipher.Block, error) {
		return a, errors.New("Error while encryptstream")
	}
	var b []byte
	_, err := encryptStream("test", b)
	if err == nil {
		t.Errorf("key blank then also it plaintext gets encrypted")
	}
}

func TestDecryptCipher(t *testing.T) {

	tmpnewcipherblock := FunctionNewCipherblock
	defer func() {
		FunctionNewCipherblock = tmpnewcipherblock
	}()
	var a cipher.Block
	FunctionNewCipherblock = func(key string) (cipher.Block, error) {
		return a, errors.New("Error while encryptstream")
	}
	var b []byte
	_, err := decryptStream("test", b)
	if err == nil {
		t.Errorf("key blank then also it plaintext gets encrypted")
	}
}
func TestEncrypt(t *testing.T) {

	var err error
	tempcipher, err = Encrypt("test", "sample")
	if err != nil {
		t.Errorf("Error while encrypting plaintext using key")
	}

	tmpEncryptStream := FunctionEncryptStream
	defer func() {
		FunctionEncryptStream = tmpEncryptStream
	}()
	var a cipher.Stream
	FunctionEncryptStream = func(key string, b []byte) (cipher.Stream, error) {
		return a, errors.New("Error while encryptstream")
	}
	_, err = Encrypt(" ", "dsds")
	if err == nil {
		t.Errorf("key blank then also it plaintext gets encrypted")
	}

	tmpReadfull := FunctionReadfull
	defer func() {
		FunctionReadfull = tmpReadfull
	}()

	FunctionReadfull = func(a io.Reader, b []byte) (int, error) {
		return -1, errors.New("Error while encrypting plaintext using key")
	}
	_, err = Encrypt(" ", "dsds")
	if err == nil {
		t.Errorf("key blank then also it plaintext gets encrypted")
	}

}

func TestEncryptWriter(t *testing.T) {

	var w io.Writer

	file := secretsPath()
	v := File("test", file)
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("got error")
	}
	defer f.Close()
	_, err = EncryptWriter(v.encodingKey, f) // positive
	if err != nil {
		t.Errorf("Error in encrypt writer")
	}

	tmpEncryptStream := FunctionEncryptStream
	defer func() {
		FunctionEncryptStream = tmpEncryptStream
	}()
	var a cipher.Stream
	FunctionEncryptStream = func(key string, b []byte) (cipher.Stream, error) {
		return a, errors.New("Error while encryptstream")
	}

	_, err = EncryptWriter("", w)
	if err == nil {
		t.Errorf("Error in encryptstream")
	}

	tmpReadfull := FunctionReadfull
	defer func() {
		FunctionReadfull = tmpReadfull
	}()

	FunctionReadfull = func(a io.Reader, b []byte) (int, error) {
		return -1, errors.New("Error while Listing passing empty list")
	}

	_, err = EncryptWriter("test", w) // positive
	if err == nil {
		t.Errorf("Error in encrypt writer")
	}

}

func TestEncryptWriterNeg(t *testing.T) {

	file := secretsPath()
	File("test", file)
	f, err := os.Open(file)
	if err != nil {
		t.Errorf("got error")
	}
	defer f.Close()
	_, err = EncryptWriter("", f) // positive
	if err == nil {
		t.Errorf("Error in encrypt writer")
	}
}

func TestDecrypt(t *testing.T) {

	_, err := Decrypt("test", tempcipher) // Positive - Provided valid key & cipherhex
	if err != nil {
		t.Errorf("error while decrypting the ciphertext")
	}
	tempdecryptStream := FunctionDecryptStream
	defer func() {
		FunctionDecryptStream = tempdecryptStream
	}()
	var a cipher.Stream
	FunctionDecryptStream = func(key string, b []byte) (cipher.Stream, error) {
		return a, errors.New("Error while Listing passing empty list")
	}
	_, err = Decrypt("test", tempcipher) // Positive - Provided valid key & cipherhex
	if err != nil {
		t.Errorf("error while decrypting the ciphertext")
	}

	_, err = Decrypt("test", "3242") // Negative - Provided valid key but ciphertext length is short
	if err == nil {
		t.Errorf("length of ciphertext is too short")
	}

	_, err = Decrypt("test", "Go is an open source lang")
	if err == nil {
		t.Errorf("length of ciphertext is too short")
	}

}

func TestDecryptReader(t *testing.T) {

	file := secretsPath()
	v := File("test", file)
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Errorf("got error")
	}
	defer f.Close()
	_, err = DecryptReader(v.encodingKey, f) // positive
	if err != nil {
		t.Errorf("Error in decrypt reader")
	}

	tempdecryptStream := FunctionDecryptStream
	defer func() {
		FunctionDecryptStream = tempdecryptStream
	}()
	var a cipher.Stream
	FunctionDecryptStream = func(key string, b []byte) (cipher.Stream, error) {
		return a, errors.New("Error while Listing passing empty list")
	}
	var r io.Reader
	_, err = DecryptReader("test", r) // positive
	if err == nil {
		t.Errorf("Error in decrypt reader")
	}

}

func TestDecryptReaderNeg(t *testing.T) {

	file := secretsPath()
	File("test", file)
	f, err := os.Open(file)
	if err != nil {
		t.Errorf("got error")
	}
	defer f.Close()
	_, err = DecryptReader("4", f) // neagtive
	if err != nil {
		t.Errorf("Error in decrypt reader")
	}
	os.Remove(file)

	file = secretsPath()
	File("test", file)
	f, err = os.Open(file)
	defer f.Close()
	_, err = DecryptReader("pqr", f) // neagtive
	if err == nil {
		t.Errorf("Expected Error but got nil in decrypt reader")
	}
	os.Remove(file)
}

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, "remp324das.txt")
}
