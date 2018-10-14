package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/towa48/go-libre-storage/internal/pkg/config"
	"github.com/towa48/go-libre-storage/internal/pkg/files"
	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

const DefaultMimeType string = "application/octet-stream"
const UrlSeparator string = "/"

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
				crawlUserDirectory(db, id, rootPath, dirName, 0)
			}
		}
	}
}

func crawlUserDirectory(db *sql.DB, userId int, rootPath string, dirPath string, parentId int64) {
	absPath := path.Join(rootPath, dirPath)
	items, err := ioutil.ReadDir(absPath)
	if err != nil {
		db.Close()
		fmt.Println(err)
		return
	}

	var fileName string
	var fileRelPath string
	var urlPath string
	var modTime time.Time
	var dbFile files.DbFileInfo

	for _, fi := range items {
		fileName = fi.Name()
		modTime = fi.ModTime()
		fileRelPath = path.Join(dirPath, fileName)

		urlPath = urlJoin(dirPath, fileName, fi.IsDir())

		if fi.IsDir() {
			dbFile = files.DbFileInfo{
				IsDir:           true,
				Name:            fileName,
				Path:            urlPath,
				CreatedDateUtc:  modTime,
				ModifiedDateUtc: modTime,
				OwnerId:         userId,
			}

			id := files.AppendFolder(db, dbFile, parentId)
			crawlUserDirectory(db, userId, rootPath, fileRelPath, id) // recurse
		} else {
			mime := getFileMime(fileName)
			etag, err := getFileChecksum(path.Join(rootPath, fileRelPath))
			if err != nil {
				fmt.Println("Cannot culculate file checksum "+fileRelPath+".", err)
				return
			}

			dbFile = files.DbFileInfo{
				IsDir:           false,
				Name:            fileName,
				Path:            urlPath,
				CreatedDateUtc:  modTime,
				ModifiedDateUtc: modTime,
				ETag:            etag,
				Mime:            mime,
				Size:            fi.Size(),
				OwnerId:         userId,
			}

			files.AppendFile(db, dbFile, parentId)
		}
	}
}

func getFileMime(fileName string) string {
	ext := path.Ext(fileName)
	mime := mime.TypeByExtension(ext)
	if mime == "" {
		mime = DefaultMimeType
	}
	return mime
}

func getFileChecksum(filePath string) (checksum string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func urlJoin(base string, item string, isDir bool) string {
	items := filepath.SplitList(base)
	items = append(items, item)

	result := UrlSeparator + strings.Join(items, UrlSeparator)
	if isDir {
		return result + UrlSeparator
	}
	return url.PathEscape(result)
}