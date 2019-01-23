package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// FunctionReadfull using function as variable
var FunctionReadfull = io.ReadFull

// FunctionEncryptStream using function as variable
var FunctionEncryptStream = encryptStream

// FunctionDecryptStream using function as variable
var FunctionDecryptStream = decryptStream

// FunctionNewCipherblock using function as variable
var FunctionNewCipherblock = newCipherBlock

func encryptStream(key string, iv []byte) (cipher.Stream, error) {

	block, err := FunctionNewCipherblock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBEncrypter(block, iv), nil
}

// Encrypt will take the key & plaintext as input & convert the plaintext in ciphertext & return the same
func Encrypt(key string, plaintext string) (string, error) {

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := FunctionReadfull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream, err := FunctionEncryptStream(key, iv)
	if err != nil {
		return "", err
	}
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", ciphertext), nil
}

// EncryptWriter function takes key & io.writer & returns the streamwriter with error
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {

	iv := make([]byte, aes.BlockSize)
	if _, err := FunctionReadfull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream, err := FunctionEncryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write full iv")
	}

	return &cipher.StreamWriter{S: stream, W: w}, nil
}

func decryptStream(key string, iv []byte) (cipher.Stream, error) {

	block, err := FunctionNewCipherblock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBDecrypter(block, iv), nil
}

// Decrypt will take the key & cipherhex which is hex representation of ciphertext
func Decrypt(key string, cipherHex string) (string, error) {

	ciphertext, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("encrypt: cipher too short")
	}
	// iv - initialisation vector
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream, err := FunctionDecryptStream(key, iv)
	if err != nil {
		return "", nil
	}

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// DecryptReader will return a reader  that will decrypt the data from the reader & gives a way to user to read
// the data as if it was not encrypted
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {

	iv := make([]byte, aes.BlockSize)
	stream, err := FunctionDecryptStream(key, iv)
	if err != nil {
		return nil, err
	}

	n, err := r.Read(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to read full iv")
	}

	return &cipher.StreamReader{S: stream, R: r}, nil

}

func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}
