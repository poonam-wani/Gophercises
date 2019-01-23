package db

import (
	"fmt"
	"path/filepath"
	"testing"

	homeDir "github.com/mitchellh/go-homedir"
)

// TestInitNegative This function intialize db by passing the directory only
func TestInitNegative(t *testing.T) {

	err := Init("/")
	if err == nil {
		fmt.Printf("Db initialisation negative scenario, got: %s\n", err)
	}
}

// TestInitPositive This function intialize db
func TestInitPositive(t *testing.T) {

	home, _ := homeDir.Dir()
	dbPath := filepath.Join(home, "tests.db")
	err := Init(dbPath)
	if err != nil {
		fmt.Printf("Expected: nil, got: %s\n", err)
	}
}

// TestCreateTaskpositive This function
func TestCreateTaskpositive(t *testing.T) {

	_, err := CreateTask("temp task")
	if err == nil {
		fmt.Printf("getting error \"%s\" while creating task", err)
	}

}

// TestGetAllTaskLists This function covers the unit test of list of tasks
func TestGetAllTaskLists(t *testing.T) {

	_, err := GetAllLists()
	if err != nil {
		fmt.Printf("getting error \"%s\" while listing all tasks", err)
	}
}

// TestDeleteTasks This function covers unit test of deleted task
func TestDeleteTasks(t *testing.T) {

	err := DeleteTasks(3)
	if err != nil {
		fmt.Println("Error while deleting tasks ")
	}
}

// TestItob This function covers unit test of integer to byte conversion
func TestItob(t *testing.T) {
	var actual []byte
	actual = itob(5)
	if actual != nil {
		fmt.Printf("got: %d\n", actual)
	}
}

// TestBtoi This function covers unit tests of byte to integer conversion
func TestBtoi(t *testing.T) {

	var exp []byte
	exp = make([]byte, 8)
	act := btoi(exp)
	if act == 0 {
		fmt.Printf("\nact : %d", act)
	}

}
