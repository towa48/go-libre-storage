package files

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/towa48/go-libre-storage/internal/pkg/config"
)

type DbFileInfo struct {
	Id              int64
	IsDir           bool
	Name            string
	Path            string
	ETag            string
	Mime            string
	Size            int64
	CreatedDateUtc  time.Time
	ModifiedDateUtc time.Time
	OwnerId         int
}

func GetPathInfo(path string, userId int, includeContent bool) []DbFileInfo {
	var result []DbFileInfo

	if path == "/" {
		// TODO: get root folder info from DB
		time := time.Now()
		result = append(result, DbFileInfo{
			Id:              0,
			IsDir:           true,
			Name:            config.Get().SystemName,
			Path:            path,
			CreatedDateUtc:  time,
			ModifiedDateUtc: time,
			OwnerId:         userId,
		})
		return result
	} else {
		// TODO: add root folder info?
	}

	if !includeContent {
		return result
	}

	// TODO: include folder content

	return result
}

func ClearUserStorage(userId int) {
	db, err := sql.Open("sqlite3", config.Get().FilesDb)
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare("delete from folders where owner_id=?;")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	checkErr(err)

	stmt, err = db.Prepare("delete from files where owner_id=?;")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec(userId)
	checkErr(err)
}

func CheckDatabase() {
	db, err := sql.Open("sqlite3", config.Get().FilesDb)
	checkErr(err)
	defer db.Close()

	if !isTableExists(db, "schema") {
		createSchemaTable(db)
	}

	if !isTableExists(db, "files") {
		createFilesTable(db)
	}
}

func createFilesTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE folders (id integer not null primary key autoincrement, name text, parent_id integer, created_date_utc datetime, changed_date_utc datetime, owner_id integer);")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("CREATE TABLE files (id integer not null primary key autoincrement, name text, folder_id integer, created_date_utc datetime, changed_date_utc datetime, etag string, mime_type string, size integer, owner_id integer);")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	checkErr(err)
}

func isTableExists(db *sql.DB, name string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE name=? and type='table';", name).Scan(&count)
	checkErr(err)

	return count == 1
}

func createSchemaTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE schema (id integer not null primary key autoincrement, version text);")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("insert into schema(version) values(?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec("1.0")
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
