package web

import (
	"fmt"
	"os"
)

type Web struct {
}

func (web *Web) Read() {
	root := "./" // change to any path

	entries, err := os.ReadDir(root)
	if err != nil {

		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Println(entry.Name())
		}
	}
}
