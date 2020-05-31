package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

//type DB struct {
//	DB *sql.DB
//}

var Conn *sql.DB

//const (
//	dbname = "db"
//	dbuser = "root"
//	dbpass = ""
//	dbport = 3306
//)

//func (d *DB) Open() {
//	var err error
//	d.DB, err = sql.Open("mysql", "root@tcp(localhost:3307)/golang")
//	if err != nil {
//		panic(err)
//	}
//
//	err = d.DB.Ping()
//	if err != nil {
//		panic(err.Error()) // proper error handling instead of panic in your app
//	}
//}

func InitDB() {
	var err error
	Conn, err = sql.Open("mysql", "root@tcp(localhost:3307)/golang")
	if err != nil {
		panic(err)
	}

	err = Conn.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}