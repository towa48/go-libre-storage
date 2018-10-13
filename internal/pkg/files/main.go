package files

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/towa48/go-libre-storage/internal/pkg/config"
)

type FileInfo struct {
	Id              uint64
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

func GetPathInfo(path string, userId int, includeContent bool) []FileInfo {
	var result []FileInfo

	if path == "/" && !includeContent {
		// TODO: get root folder info from DB
		time := time.Now()
		result = append(result, FileInfo{
			Id:              0,
			IsDir:           true,
			Name:            "box",
			Path:            path,
			CreatedDateUtc:  time,
			ModifiedDateUtc: time,
			OwnerId:         userId,
		})
		return result
	}

	return result
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
