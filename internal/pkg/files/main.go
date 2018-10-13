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
	ContentType     string
	ContentLength   string
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

func stub() {
	db, err := sql.Open("sqlite3", config.Get().FilesDb)
	checkErr(err)
	defer db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
