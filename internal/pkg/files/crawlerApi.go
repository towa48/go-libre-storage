package files

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/towa48/go-libre-storage/internal/pkg/config"
)

func GetDbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", config.Get().FilesDb+"?_foreign_keys=true")
	checkErr(err)

	return db
}

func AppendFolder(db *sql.DB, fi DbFileInfo, parentId int64) int64 {
	// select last_insert_rowid();
	cmd := "insert into folders (name, parent_id, url, created_date_utc, changed_date_utc, owner_id) values(?, ?, ?, ?, ?, ?);"
	var result sql.Result
	var err error
	if parentId == 0 {
		result, err = db.Exec(cmd, fi.Name, nil, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.OwnerId)
	} else {
		result, err = db.Exec(cmd, fi.Name, parentId, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.OwnerId)
	}
	checkDbErr(db, err)

	id, err := result.LastInsertId()
	checkDbErr(db, err)

	return id
}

func AppendFile(db *sql.DB, fi DbFileInfo, folderId int64) {
	cmd := "insert into files (name, folder_id, url, created_date_utc, changed_date_utc, etag, mime_type, size, owner_id) values(?, ?, ?, ?, ?, ?, ?, ?, ?);"
	var err error
	if folderId == 0 {
		_, err = db.Exec(cmd, fi.Name, nil, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.ETag, fi.Mime, fi.Size, fi.OwnerId)
	} else {
		_, err = db.Exec(cmd, fi.Name, folderId, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.ETag, fi.Mime, fi.Size, fi.OwnerId)
	}
	checkDbErr(db, err)
}

func checkDbErr(db *sql.DB, err error) {
	if err != nil {
		db.Close()
		panic(err)
	}
}
