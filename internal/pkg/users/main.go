package users

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/towa48/go-libre-storage/internal/pkg/config"
)

type User struct {
	Id             int
	Login          string
	PasswordHash   string
	Salt           string
	CreatedDateUtc time.Time
}

func GetAll() []User {
	db, err := sql.Open("sqlite3", config.Get().UsersDb)
	checkErr(err)
	defer db.Close()

	return nil
}

func GetUserIdByLogin(login string) (userId int, found bool) {
	db, err := sql.Open("sqlite3", config.Get().UsersDb)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT id FROM users WHERE login=?;", login)
	checkErr(err)
	defer rows.Close()

	var id int
	var f bool
	for rows.Next() {
		err = rows.Scan(&id)
		checkErr(err)
		f = true
		break
	}

	return id, f
}

func GetUserByLogin(login string) (user User, found bool) {
	db, err := sql.Open("sqlite3", config.Get().UsersDb)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT id, login, password_hash, salt, created_date_utc FROM users WHERE login=?;", login)
	checkErr(err)
	defer rows.Close()

	var u User
	var f bool
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Login, &u.PasswordHash, &u.Salt, &u.CreatedDateUtc)
		checkErr(err)
		f = true
		break
	}

	return u, f
}

func GetUserById(userId int) (user User, found bool) {
	db, err := sql.Open("sqlite3", config.Get().UsersDb)
	checkErr(err)
	defer db.Close()

	rows, err := db.Query("SELECT id, login, password_hash, salt, created_date_utc FROM users WHERE id=?;", userId)
	checkErr(err)
	defer rows.Close()

	var u User
	var f bool
	for rows.Next() {
		err = rows.Scan(&u.Id, &u.Login, &u.PasswordHash, &u.Salt, &u.CreatedDateUtc)
		checkErr(err)
		f = true
		break
	}

	return u, f
}

func CheckDatabase() {
	db, err := sql.Open("sqlite3", config.Get().UsersDb)
	checkErr(err)
	defer db.Close()

	if !isTableExists(db, "schema") {
		createSchemaTable(db)
	}

	if !isTableExists(db, "users") {
		createUsersTable(db)
	}
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

func createUsersTable(db *sql.DB) {
	stmt, err := db.Prepare("CREATE TABLE users (id integer not null primary key autoincrement, login text, password_hash text, salt text, created_date_utc datetime);")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("insert into users (login, password_hash, salt, created_date_utc) values(?, ?, ?, ?)")
	checkErr(err)
	defer stmt.Close()

	_, err = stmt.Exec("admin", "PUhWAYaboem3IuUl40kOa1GzDM2pSSUR9OvNW217UnI=", "gYVa2GgdDYbR6R4AFnk5y2aU0sQirNIIoAcpOUh/aZk=", time.Now())
	checkErr(err)
}

func addUser(user User) int {
	db, err := sql.Open("sqlite3", config.Get().UsersDb)
	checkErr(err)
	defer db.Close()

	result, err := db.Exec("insert into users (login, password_hash, salt, created_date_utc) values(?, ?, ?, ?)", user.Login, user.PasswordHash, user.Salt, user.CreatedDateUtc)
	checkErr(err)

	id, err := result.LastInsertId()
	checkErr(err)

	return int(id)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
