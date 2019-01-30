package files

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/towa48/go-libre-storage/internal/pkg/config"
)

const UrlSeparator string = "/"

type DbFileInfo struct {
	Id              int64
	IsDir           bool
	IsShared        bool
	IsReadOnly      bool
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

func GetFolderContent(url string, userId int, urlPrefix string, includeContent bool) (items []DbFileInfo, found bool) {
	var result []DbFileInfo

	db := GetDbConnection()
	defer db.Close()

	folder, found := getFolderInfo(db, url, urlPrefix, userId)

	if !found {
		folder, found = getSharedFolderInfo(db, url, urlPrefix, userId)
		if !found {
			return nil, false
		}
	}

	var isRoot bool
	if url == "/" {
		folder.Name = config.Get().SystemName
		isRoot = true
	}

	result = append(result, folder)

	if includeContent {
		var content []DbFileInfo
		if folder.IsShared {
			content = getSharedFolderContent(db, folder, urlPrefix)
		} else {
			content = getFolderContent(db, userId, folder.Id, urlPrefix, isRoot)
		}
		for _, c := range content {
			result = append(result, c)
		}
	}

	return result, true
}

func GetFolderInfo(url string, userId int) (item *DbFileInfo, found bool) {
	c, f := GetFolderContent(url, userId, "", false)
	if !f {
		return nil, false
	}
	return &c[0], true
}

func GetFolderInfoById(folderId int64) (item *DbFileInfo, found bool) {
	db := GetDbConnection()
	defer db.Close()

	return getFolderInfoById(db, folderId, "")
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

	rows, err := db.Query(`with t
	as
	(
		select id, name, owner_id, parent_id
		from folders
		where id=?
		union all
		select f.id, f.name, f.owner_id, f.parent_id
		from t
		join folders as f
		on f.id = t.parent_id
	)
	select t.id, t.name, t.owner_id, t.parent_id from t;`, *it.ParentId)

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

func ShareFolderToUser(folderId int64, userId int, readOnly bool) {
	db := GetDbConnection()
	defer db.Close()

	shareFolderToUser(db, folderId, userId, readOnly)
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
		createSchemaTable(db, "1.2")
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

	_, err = db.Exec("CREATE TABLE shared_folders (folder_id integer not null, read_only integer, target_user_id integer not null, constraint fk_shared_folder foreign key (folder_id) references folders(id) on delete cascade);")
	checkErr(err)

	_, err = db.Exec("CREATE TABLE shared_files (file_id integer not null, read_only integer, target_user_id integer not null, constraint fk_shared_file foreign key (file_id) references files(id) on delete cascade);")
	checkErr(err)

	_, err = db.Exec("CREATE INDEX idx_shared_folders on shared_folders(folder_id, target_user_id);")
	checkErr(err)

	_, err = db.Exec("CREATE INDEX idx_shared_folders_user on shared_folders(target_user_id);")
	checkErr(err)

	_, err = db.Exec("CREATE INDEX idx_shared_files on shared_files(file_id, target_user_id);")
	checkErr(err)

	_, err = db.Exec("CREATE INDEX idx_shared_files_user on shared_files(target_user_id);")
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

func getFolderContent(db *sql.DB, userId int, folderId int64, urlPrefix string, isRoot bool) []DbFileInfo {
	var result []DbFileInfo

	// include user folders
	rows, err := db.Query(`SELECT f.id, f.name, f.url, f.created_date_utc, f.changed_date_utc
		FROM folders as f
		WHERE f.parent_id=? and f.owner_id=?;`, folderId, userId)

	checkErr(err)
	defer rows.Close()

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

	// include root shared folders
	if isRoot {
		rows, err = db.Query(`SELECT f.id, f.name, f.url, f.created_date_utc, f.changed_date_utc, f.owner_id, s.read_only, f2.url as parentUrl
			FROM shared_folders s
			INNER JOIN folders as f
				ON f.id = s.folder_id
			LEFT OUTER JOIN folders as f2
				ON f2.id = f.parent_id AND f.parent_id IS NOT NULL
			WHERE s.target_user_id=?;`, userId)

		checkErr(err)
		defer rows.Close()

		var readOnly int
		var parentUrl *string
		for rows.Next() {
			it := DbFileInfo{
				IsDir:    true,
				IsShared: true,
			}
			err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &it.OwnerId, &readOnly, &parentUrl)
			checkErr(err)

			// exclude parent url from shared folder
			if parentUrl != nil {
				it.Path = UrlSeparator + strings.TrimPrefix(it.Path, *parentUrl)
			}

			it.IsReadOnly = readOnly > 0
			it.Path = urlJoin(urlPrefix, it.Path)

			result = append(result, it)
		}
	}

	// include user files
	rows, err = db.Query(`SELECT id, name, url, created_date_utc, changed_date_utc, etag, mime_type, size
		FROM files
		WHERE folder_id=? and owner_id=?;`, folderId, userId)
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

func getSharedFolderContent(db *sql.DB, parentFolder DbFileInfo, urlPrefix string) []DbFileInfo {
	var result []DbFileInfo

	rows, err := db.Query(`SELECT id, name, url, created_date_utc, changed_date_utc
		FROM folders
		WHERE owner_id=? AND parent_id=?;`, parentFolder.OwnerId, parentFolder.Id)

	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		it := DbFileInfo{
			IsDir:    true,
			IsShared: true,
		}
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc)
		checkErr(err)

		// cut path started from parent path
		startPos := strings.LastIndex(it.Path, strings.TrimPrefix(parentFolder.Path, urlPrefix))
		folderRelativePath := it.Path[startPos:]

		it.OwnerId = parentFolder.OwnerId
		it.IsReadOnly = parentFolder.IsReadOnly
		it.Path = urlJoin(urlPrefix, folderRelativePath)
		result = append(result, it)
	}

	rows, err = db.Query(`SELECT id, name, url, created_date_utc, changed_date_utc, etag, mime_type, size
		FROM files
		WHERE folder_id=? and owner_id=?;`, parentFolder.Id, parentFolder.OwnerId)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		it := DbFileInfo{
			IsDir:    false,
			IsShared: true,
			OwnerId:  parentFolder.OwnerId,
		}
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &it.ETag, &it.Mime, &it.Size)
		checkErr(err)

		// cut path started from parent path
		startPos := strings.LastIndex(it.Path, strings.TrimPrefix(parentFolder.Path, urlPrefix))
		fileRelativePath := it.Path[startPos:]

		it.Path = urlJoin(urlPrefix, fileRelativePath)
		result = append(result, it)
	}

	return result
}

func getFileInfo(db *sql.DB, url string, urlPrefix string, userId int) (item DbFileInfo, found bool) {
	rows, err := db.Query(`SELECT id, name, url, created_date_utc, changed_date_utc, etag, mime_type, size
		FROM files
		WHERE url=? and owner_id=?;`, url, userId)
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

	if !f {
		root, rootFound := findSharedRoot(db, url, urlPrefix, userId)
		if !rootFound {
			return it, false
		}

		truncatedRootPath := strings.TrimPrefix(root.Path, urlPrefix)

		rows, err = db.Query(`SELECT id, name, url, created_date_utc, changed_date_utc, etag, mime_type, size
			FROM files
			WHERE url like '%' || ? and url like ? || '%' and owner_id=?;`, url, truncatedRootPath, root.OwnerId)
		checkErr(err)
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &it.ETag, &it.Mime, &it.Size)
			checkErr(err)

			it.Path = urlJoin(urlPrefix, it.Path) // TBD fix path
			it.OwnerId = root.OwnerId
			it.IsShared = true
			it.IsReadOnly = root.IsReadOnly

			f = true
			break
		}
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

func getSharedFolderInfo(db *sql.DB, url string, urlPrefix string, userId int) (item DbFileInfo, found bool) {
	root, rootFound := findSharedRoot(db, url, urlPrefix, userId)
	if !rootFound {
		var it DbFileInfo
		return it, false // TODO: return nil
	}

	truncatedRootPath := strings.TrimPrefix(root.Path, urlPrefix)

	rows, err := db.Query(`select id, name, url, created_date_utc, changed_date_utc
		from folders
		WHERE url like '%' || ? and url like ? || '%' and owner_id=?;`, url, truncatedRootPath, root.OwnerId)

	checkErr(err)
	defer rows.Close()

	var it DbFileInfo
	var f bool
	for rows.Next() {
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc)
		checkErr(err)

		it.Path = urlJoin(urlPrefix, url) // get only url suffix from original one
		it.OwnerId = root.OwnerId
		it.IsDir = true
		it.IsShared = true
		it.IsReadOnly = root.IsReadOnly

		f = true
		break
	}

	return it, f
}

func getFolderInfoById(db *sql.DB, folderId int64, urlPrefix string) (item *DbFileInfo, found bool) {
	rows, err := db.Query("SELECT id, name, url, created_date_utc, changed_date_utc, owner_id FROM folders WHERE id=?;", folderId)
	checkErr(err)
	defer rows.Close()

	var it DbFileInfo
	var f bool
	for rows.Next() {
		err = rows.Scan(&it.Id, &it.Name, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &it.OwnerId)
		checkErr(err)
		it.Path = urlJoin(urlPrefix, it.Path)
		it.IsDir = true
		f = true
		break
	}

	if !f {
		return nil, false
	}

	return &it, f
}

func shareFolderToUser(db *sql.DB, folderId int64, userId int, readOnly bool) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM shared_folders WHERE folder_id=? and target_user_id=?;", folderId, userId).Scan(&count)
	checkErr(err)

	if count > 0 {
		return
	}

	cmd := "insert into shared_folders (folder_id, target_user_id, read_only) values(?, ?, ?);"
	_, err = db.Exec(cmd, folderId, userId, readOnly)
	checkErr(err)
}

func findSharedRoot(db *sql.DB, url string, urlPrefix string, userId int) (item DbFileInfo, found bool) {
	var it DbFileInfo

	urlParts := urlSplit(url)
	rootFolder := urlParts[0]
	rootFolderUrl := UrlSeparator + rootFolder + UrlSeparator

	var rootFound bool

	// check root folder is shared
	rows, err := db.Query(`SELECT f.id, f.name, f.owner_id, f.url, f.created_date_utc, f.changed_date_utc, s.read_only
		FROM shared_folders AS s
		INNER JOIN folders AS f
			ON s.folder_id = f.id
		WHERE s.target_user_id=? AND f.url like '%' || ?;`, userId, rootFolderUrl)

	checkErr(err)
	defer rows.Close()

	var readOnlyInt int
	for rows.Next() {
		err = rows.Scan(&it.Id, &it.Name, &it.OwnerId, &it.Path, &it.CreatedDateUtc, &it.ModifiedDateUtc, &readOnlyInt)
		checkErr(err)

		it.Path = urlJoin(urlPrefix, it.Path)
		it.IsDir = true
		it.IsShared = true
		it.IsReadOnly = readOnlyInt > 0

		rootFound = true
	}

	return it, rootFound
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
	if !strings.HasSuffix(base, UrlSeparator) && !strings.HasPrefix(item, UrlSeparator) {
		base = base + UrlSeparator
	}
	return base + item
}

func urlSplit(url string) []string {
	url = strings.TrimLeft(url, UrlSeparator)
	url = strings.TrimRight(url, UrlSeparator)

	return strings.Split(url, UrlSeparator)
}
