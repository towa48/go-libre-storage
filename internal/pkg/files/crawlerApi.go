package files

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/towa48/go-libre-storage/internal/pkg/config"
)

func GetDbConnection() *sql.DB {
	db, err := sql.Open("sqlite3", config.Get().FilesDb)
	checkErr(err)

	return db
}

func AppendFolder(db *sql.DB, fi DbFileInfo, parentId int64) int64 {
	// select last_insert_rowid();
	stmt, err := db.Prepare("insert into folders (name, parent_id, path, created_date_utc, changed_date_utc, owner_id) values(?, ?, ?, ?, ?, ?);")
	checkDbErr(db, err)
	defer stmt.Close()

	var result sql.Result
	if parentId == 0 {
		result, err = stmt.Exec(fi.Name, nil, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.OwnerId)
	} else {
		result, err = stmt.Exec(fi.Name, parentId, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.OwnerId)
	}
	checkDbErr(db, err)

	id, err := result.LastInsertId()
	checkDbErr(db, err)

	return id
}

func AppendFile(db *sql.DB, fi DbFileInfo, folderId int64) {
	stmt, err := db.Prepare("insert into files (name, folder_id, path, created_date_utc, changed_date_utc, etag, mime_type, size, owner_id) values(?, ?, ?, ?, ?, ?, ?, ?, ?);")
	checkDbErr(db, err)
	defer stmt.Close()

	if folderId == 0 {
		_, err = stmt.Exec(fi.Name, nil, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.ETag, fi.Mime, fi.Size, fi.OwnerId)
	} else {
		_, err = stmt.Exec(fi.Name, folderId, fi.Path, fi.CreatedDateUtc, fi.ModifiedDateUtc, fi.ETag, fi.Mime, fi.Size, fi.OwnerId)
	}
	checkDbErr(db, err)
}

func checkDbErr(db *sql.DB, err error) {
	if err != nil {
		db.Close()
		panic(err)
	}
}
