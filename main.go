package main

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No argument provided")
		os.Exit(1)
	}

	fileName := os.Args[1]

	stat, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		os.Exit(1)
	}

	if stat.IsDir() {
		if err := traversDir(fileName); err != nil {
			fmt.Println("Error traversing directory:", err)
			os.Exit(1)
		}
	} else {
		if err := ParseFile(fileName); err != nil {
			fmt.Println("Error parsing file:", err)
			os.Exit(1)
		}
	}
}

func traversDir(baseDir string) error {
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".jack" {
			if err := ParseFile(path); err != nil {
				return errors.Wrap(err, "Failed to parse")
			}
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "Error walking the path")
	}

	return nil
}
