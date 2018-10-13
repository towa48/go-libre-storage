package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/towa48/go-libre-storage/internal/pkg/config"
	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

func crawl() {
	rootFolder := config.Get().Storage
	fmt.Println("Crawl mode is enabled.")
	fmt.Printf("Root folder is '%s'.\n", rootFolder)

	fi, err := os.Stat(rootFolder)
	if err != nil {
		fmt.Println(err)
		return
	}

	mode := fi.Mode()
	if !mode.IsDir() {
		fmt.Println("Error: root folder is not directory.")
		return
	}

	items, err := ioutil.ReadDir(rootFolder)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, fi = range items {
		if fi.IsDir() {
			dirName := fi.Name()
			_, found := users.GetUserIdByLogin(dirName)
			if !found {
				fmt.Printf("Found directory for unknown account: '%s'. This directory would be skipped.\n", dirName)
			} else {
				fmt.Printf("Found directory for account: '%s'.\n", dirName)
			}
		}
	}
}
