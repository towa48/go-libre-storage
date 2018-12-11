package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/towa48/go-libre-storage/internal/pkg/files"

	"github.com/towa48/go-libre-storage/internal/pkg/users"
)

func checkCliCommands() bool {
	doCrawl := flag.Bool("crawl", false, "Clear DB file metadata and restore it from filesystem. All shared items will be dropped.")
	newUser := flag.String("add-user", "", "Add new user. Password will be promted.")

	shareFolder := flag.Int64("share-folder", -1, "Share folder by id to another user (used with --to argument)")
	toUser := flag.String("to", "", "User login should be specified for some commands (used with --share-folder).")
	write := flag.Bool("write", false, "Allow modifications of the item to user (used with --share-folder).")
	flag.Parse()

	if doCrawl != nil && *doCrawl {
		crawl()
		return true
	}

	if newUser != nil && *newUser != "" {
		addNewUser(*newUser)
		return true
	}

	if shareFolder != nil && *shareFolder != -1 {
		if write == nil {
			w := false
			write = &w
		}
		shareFolderToUser(*shareFolder, toUser, *write)
		return true
	}

	return false
}

func trim(s string) string {
	return strings.Trim(s, " \n")
}

func isEmptyString(s string) bool {
	return trim(s) == ""
}

func addNewUser(login string) {
	reader := bufio.NewReader(os.Stdin)
	pass := ""
	for isEmptyString(pass) {
		fmt.Print("Enter password: ")
		pass, _ = reader.ReadString('\n')
		if isEmptyString(pass) {
			fmt.Println("Password is invalid.")
		}
	}

	pass = trim(pass)
	userId := users.AddUser(login, pass)

	if userId != 0 {
		createUserRoot(userId, login)
	}
}

func shareFolderToUser(folderId int64, login *string, write bool) {
	if folderId < 0 {
		fmt.Println("Folder id should be positive.")
		return
	}

	if login == nil || *login == "" {
		fmt.Println("User login should be specified (see --to argument).")
		return
	}

	fi, found := files.GetFolderInfoById(folderId)
	if !found {
		fmt.Printf("Folder with id '%d' not found.\n", folderId)
		return
	}

	u, found := users.GetUserByLogin(*login)
	if !found {
		fmt.Printf("User with login '%s' not found.\n", *login)
		return
	}

	files.ShareFolderToUser(fi.Id, u.Id, !write)
	fmt.Printf("Folder '%s' successfully shared.\n", fi.Name)
}
