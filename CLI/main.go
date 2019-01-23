package main

import (
	"path/filepath"

	homeDir "github.com/mitchellh/go-homedir"
	"github.com/poonam-wani/gophercises/CLI/cmd"
	"github.com/poonam-wani/gophercises/CLI/db"
)

func main() {

	home, _ := homeDir.Dir()
	dbPath := filepath.Join(home, "tests.db") // Created the database with name tests.db
	db.Init(dbPath)
	cmd.RootCmd.Execute()

}
