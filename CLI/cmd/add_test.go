package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	homeDir "github.com/mitchellh/go-homedir"
	"github.com/poonam-wani/gophercises/CLI/db"
)

// TestAddCmd This function is covering the negative unit tests for adding the task in DB
func TestAddCmd(t *testing.T) {

	home, _ := homeDir.Dir()
	dbPath := filepath.Join(home, "tests.db")
	err := db.Init(dbPath)
	if err != nil {
		fmt.Println("Failed to open db")
	}
	tmpcreate := FunctionCreate
	defer func() {
		FunctionCreate = tmpcreate
	}()

	FunctionCreate = func(a1 string) (int, error) {
		return 0, errors.New("Error while creation")
	}
	a := []string{"new task"}
	AddCmd.Run(AddCmd, a)

}

func TestAddPositiveCmd(t *testing.T) {

	a := []string{"new task 2"}
	AddCmd.Run(AddCmd, a)

}
