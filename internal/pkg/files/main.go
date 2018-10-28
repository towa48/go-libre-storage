package files

import (
	"database/sql"
	"fmt"
	"strings"
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

type DbHierarchyItem struct {
	Id       int64
	IsDir    bool
	Name     string
	OwnerId  int
	ParentId *int64
	Child    *DbHierarchyItem
}

func GetFolderContent(url string, userId int, urlPrefix string, includeContent bool) (items []DbFileInfo, hasAccess bool) {
	var result []DbFileInfo

	db := GetDbConnection()
	defer db.Close()

	folder, found := getFolderInfo(db, url, urlPrefix, userId)
	if !found {
		return nil, false
	}

	if url == "/" {
		folder.Name = config.Get().SystemName
	}

	result = append(result, folder)

	if includeContent {
		content := getFolderContent(db, userId, folder.Id, urlPrefix)
		for _, c := range content {
			result = append(result, c)
		}
	}

	return result, true
}

func GetFolderInfo(url string, userId int) (item *DbFileInfo, hasAccess bool) {
	c, f := GetFolderContent(url, userId, "", false)
	if !f {
		return nil, false
	}
	return &c[0], true
}

func GetFileInfo(url string, userId int, urlPrefix string) (file DbFileInfo, hasAccess bool) {
	db := GetDbConnection()
	defer db.Close()

	file, found := getFileInfo(db, url, urlPrefix, userId)
	if !found {
		return file, false
	}

	return file, true
}

func GetFileHierarchy(fileId int64) (root DbHierarchyItem, found bool) {
	db := GetDbConnection()
	defer db.Close()

	it, f := getFileById(db, fileId)
	if !f {
		return it, f
	}

	rows, err := db.Query("with t as (select id, name, owner_id, parent_id from folders where id=? union all select f.id, f.name, f.owner_id, f.parent_id from t join folders as f on f.id = t.parent_id) select t.id, t.name, t.owner_id, t.parent_id from t;", *it.ParentId)
	checkErr(err)
	defer rows.Close()

	var folders []DbHierarchyItem
	for rows.Next() {
		var folder DbHierarchyItem
		err = rows.Scan(&folder.Id, &folder.Name, &folder.OwnerId, &folder.ParentId)
		checkErr(err)

		folder.IsDir = true
		folders = append(folders, folder)
	}

	err = rows.Close()
	checkErr(err)

	//fmt.Println("Folders count:", len(folders))
	r := buildHierarchy(folders, it)

	return r, true
}

func RemoveFolder(folderId int64) {
	db := GetDbConnection()
	defer db.Close()

	_, err := db.Exec("delete from folders where id=?;", folderId)
	checkErr(err)
}

func RemoveFile(fileId int64) {
	db := GetDbConnection()
	defer db.Close()

	_, err := db.Exec("delete from files where id=?;", fileId)
	checkErr(err)
}

func ClearUserStorage(userId int) {
	db := GetDbConnection()
	defer db.Close()

	_, err := db.Exec("delete from folders where owner_id=?;", userId)
	checkErr(err)

	_, err = db.Exec("delete from files where owner_id=?;", userId)
	checkErr(err)
}

func CheckDatabase() {
	db := GetDbConnection()
	defer db.Close()

	if !isTableExists(db, "schema") {
		fmt.Println("Creating DB schema table")
		createSchemaTable(db, "1.1")
	}

	if !isTableExists(db, "files") {
		fmt.Println("Creating DB filesystem tables")
		createFilesTable(db)
	}
}

func createFilesTable(db *sql.DB) {
	_, err := db.Exec("CREATE TABLE folders (id integer not null primary key autoincrement, name text, parent_id integer, url text, created_date_utc datetime, changed_date_utc datetime, owner_id integer, constraint fk_parent_folder foreign key (parent_id) references folders(id) on delete cascade);")
	checkErr(err)

	_, err = db.Exec("CREATE TABLE files (id integer not null primary key autoincrement, name text, folder_id integer not null, url text, created_date_utc datetime, changed_date_utc datetime, etag string, mime_type string, size integer, owner_id integer, constraint fk_parent_folder foreign key (folder_id) references folders(id) on delete cascade);")
	checkErr(err)
}

func isTableExists(db *sql.DB, name string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE name=? and type='table';", name).Scan(&count)
	checkErr(err)

	return count == 1
}

func createSchemaTable(db *sql.DB, version string) {
	_, err := db.Exec("CREATE TABLE schema (id integer not null primary key autoincrement, version text);")
	checkErr(err)

	_, err = db.Exec("insert into schema(version) values(?)", version)
	checkErr(err)
}

func getFolderContent(db *sql.DB, userId int, folderId int64, urlPrefix string) []DbFileInfo {
	rows, err := db.Query("SELECT id, name, url, created_date_utc, changed_date_utc FROM folders WHERE parent_id=? and owner_id=?;", folderId, userId)
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

		it.Path = urlJoin(urlPrefix, it.Path)
		result = append(result, it)
	}

	rows, err = db.Query("SELECT id, name, url, created_date_utc, changed_date_utc, etag, mime_type, size FROM files WHERE folder_id=? and owner_id=?;", folderId, userId)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		it := DbFileInfo{
			IsDir:   false,
			OwnerId: userId,
		}
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &it.ETag, &it.Mime, &it.Size)
		checkErr(err)

		it.Path = urlJoin(urlPrefix, it.Path)
		result = append(result, it)
	}

	return result
}

func getFileInfo(db *sql.DB, url string, urlPrefix string, userId int) (item DbFileInfo, found bool) {
	rows, err := db.Query("SELECT id, name, url, created_date_utc, changed_date_utc, etag, mime_type, size FROM files WHERE url=? and owner_id=?;", url, userId)
	checkErr(err)
	defer rows.Close()

	var it DbFileInfo
	var f bool
	for rows.Next() {
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &it.ETag, &it.Mime, &it.Size)
		checkErr(err)
		it.Path = urlJoin(urlPrefix, it.Path)
		it.OwnerId = userId
		f = true
		break
	}

	return it, f
}

func getFileById(db *sql.DB, fileId int64) (file DbHierarchyItem, found bool) {
	rows, err := db.Query("SELECT id, name, owner_id, folder_id FROM files WHERE id=?;", fileId)
	checkErr(err)
	defer rows.Close()

	var it DbHierarchyItem
	var f bool
	for rows.Next() {
		err = rows.Scan(&it.Id, &it.Name, &it.OwnerId, &it.ParentId)
		checkErr(err)
		f = true
		break
	}

	return it, f
}

func getFolderInfo(db *sql.DB, url string, urlPrefix string, userId int) (item DbFileInfo, found bool) {
	rows, err := db.Query("SELECT id, name, url, created_date_utc, changed_date_utc FROM folders WHERE url=? and owner_id=?;", url, userId)
	checkErr(err)
	defer rows.Close()

	var it DbFileInfo
	var f bool
	for rows.Next() {
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc)
		checkErr(err)
		it.Path = urlJoin(urlPrefix, it.Path)
		it.IsDir = true
		it.OwnerId = userId
		f = true
		break
	}

	return it, f
}

func buildHierarchy(folders []DbHierarchyItem, file DbHierarchyItem) DbHierarchyItem {
	var m = map[int64]*DbHierarchyItem{}

	for _, f := range folders {
		node := f
		m[f.Id] = &node
	}

	var root *DbHierarchyItem
	for _, f := range m {
		node := f

		if node.ParentId == nil {
			root = node
			continue
		}

		parent, found := m[*node.ParentId]
		if !found {
			continue
		}

		parent.Child = node
	}

	// assign file
	parent, found := m[*file.ParentId]
	if found {
		parent.Child = &file
	}

	return *root
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func urlJoin(base string, item string) string {
	if !strings.HasSuffix(base, "/") && !strings.HasPrefix(item, "/") {
		base = base + "/"
	}
	return base + item
}
