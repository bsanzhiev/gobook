package main

import (
	"fmt"
	"os"
	"path/filepath"
	"taskman/cmd"
	"taskman/db"

	homedir "github.com/mitchellh/go-homedir"
)

func main() {
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")
	must(db.Init(dbPath))
	// if err != nil {
	// 	panic(err)
	// }
	must(cmd.RootCmd.Execute())

}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
