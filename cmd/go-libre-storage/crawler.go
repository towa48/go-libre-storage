package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"mime"
	"os"
	"path"
	"time"

	"github.com/towa48/go-libre-storage/internal/pkg/config"
	"github.com/towa48/go-libre-storage/internal/pkg/files"
	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

const DefaultMimeType string = "application/octet-stream"

func crawl() {
	rootFolder := config.Get().Storage
	fmt.Println("Crawl mode is enabled.")
	fmt.Printf("Root folder is '%s'.\n", rootFolder)

	crawlRootFolder(rootFolder)
}

func crawlRootFolder(rootPath string) {
	fi, err := os.Stat(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	mode := fi.Mode()
	if !mode.IsDir() {
		fmt.Println("Error: root folder is not directory.")
		return
	}

	items, err := ioutil.ReadDir(rootPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, fi = range items {
		if fi.IsDir() {
			dirName := fi.Name()
			id, found := users.GetUserIdByLogin(dirName)
			if !found {
				fmt.Printf("Found directory for unknown account: '%s'. You should create that account first.\n", dirName)
			} else {
				fmt.Printf("Found directory for account: '%s'.\n", dirName)
				files.ClearUserStorage(id)

				db := files.GetDbConnection()
				crawlUserDirectory(db, id, path.Join(rootPath, dirName), 0)
			}
		}
	}
}

func crawlUserDirectory(db *sql.DB, userId int, dirPath string, parentId int64) {
	items, err := ioutil.ReadDir(dirPath)
	if err != nil {
		db.Close()
		fmt.Println(err)
		return
	}

	var fileName string
	var modTime time.Time
	var dbFile files.DbFileInfo
	for _, fi := range items {
		fileName = fi.Name()
		modTime = fi.ModTime()

		if fi.IsDir() {
			dbFile = files.DbFileInfo{
				IsDir:           true,
				Name:            fileName,
				CreatedDateUtc:  modTime,
				ModifiedDateUtc: modTime,
				OwnerId:         userId,
			}
			id := files.AppendFolder(db, dbFile, parentId)
			crawlUserDirectory(db, userId, path.Join(dirPath, fileName), id) // recurse
		} else {
			ext := path.Ext(fileName)
			mime := mime.TypeByExtension(ext)
			if mime == "" {
				mime = DefaultMimeType
			}
			dbFile = files.DbFileInfo{
				IsDir:           false,
				Name:            fileName,
				CreatedDateUtc:  modTime,
				ModifiedDateUtc: modTime,
				ETag:            "",
				Mime:            mime,
				Size:            fi.Size(),
				OwnerId:         userId,
			}
			files.AppendFile(db, dbFile, parentId)
		}
	}
}
