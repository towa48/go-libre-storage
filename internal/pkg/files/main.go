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

func GetFolderInfo(url string, userId int, includeContent bool) (items []DbFileInfo, hasAccess bool) {
	var result []DbFileInfo

	db := GetDbConnection()
	defer db.Close()

	if url == "/" {
		time := time.Now()
		result = append(result, DbFileInfo{
			Id:              0,
			IsDir:           true,
			Name:            config.Get().SystemName,
			Path:            url,
			CreatedDateUtc:  time,
			ModifiedDateUtc: time,
			OwnerId:         userId,
		})

		if includeContent {
			content := getFolderContent(db, userId, 0)
			for _, c := range content {
				result = append(result, c)
			}
		}

		return result, true
	}

	folder, found := getFolderInfo(db, url, userId)
	if !found {
		return nil, true
	}

	result = append(result, folder)

	if includeContent {
		content := getFolderContent(db, userId, folder.Id)
		for _, c := range content {
			result = append(result, c)
		}
	}

	return result, true
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
	stmt, err := db.Prepare("CREATE TABLE folders (id integer not null primary key autoincrement, name text, parent_id integer, url text, created_date_utc datetime, changed_date_utc datetime, owner_id integer);")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("CREATE TABLE files (id integer not null primary key autoincrement, name text, folder_id integer, url text, created_date_utc datetime, changed_date_utc datetime, etag string, mime_type string, size integer, owner_id integer);")
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

func getFolderContent(db *sql.DB, userId int, folderId int64) []DbFileInfo {
	var rows *sql.Rows
	var err error

	fq := "SELECT id, name, url, created_date_utc, changed_date_utc FROM folders "
	if folderId == 0 {
		rows, err = db.Query(fq+"WHERE parent_id=NULL and owner_id=?;", userId)
	} else {
		rows, err = db.Query(fq+"WHERE parent_id=? and owner_id=?;", folderId, userId)
	}
	checkErr(err)
	defer rows.Close()

	var result []DbFileInfo
	for rows.Next() {
		it := DbFileInfo{
			IsDir:   true,
			OwnerId: userId,
		}
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc)
		checkErr(err)

		result = append(result, it)
	}

	ffq := "SELECT id, name, url, created_date_utc, changed_date_utc, etag, mime_type, size FROM files "
	if folderId == 0 {
		rows, err = db.Query(ffq+"WHERE folder_id=NULL and owner_id=?;", userId)
	} else {
		rows, err = db.Query(ffq+"WHERE folder_id=? and owner_id=?;", folderId, userId)
	}
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		it := DbFileInfo{
			IsDir:   false,
			OwnerId: userId,
		}
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &it.ETag, &it.Mime, &it.Size)
		checkErr(err)

		result = append(result, it)
	}

	return result
}

func getFolderInfo(db *sql.DB, url string, userId int) (item DbFileInfo, found bool) {
	rows, err := db.Query("SELECT id, name, url, created_date_utc, changed_date_utc FROM folders WHERE url=? and owner_id=?;", url, userId)
	checkErr(err)
	defer rows.Close()

	var it DbFileInfo
	var f bool
	for rows.Next() {
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc)
		checkErr(err)
		it.OwnerId = userId
		f = true
		break
	}

	return it, f
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
