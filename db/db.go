package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

var MySQL *sql.DB

func InitMySQL() {
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = ""
	cfg.Net = "tcp"
	cfg.Addr = "localhost:3306"
	cfg.DBName = "fanchiikawa"
	cfg.ParseTime = true

	var err error
	MySQL, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := MySQL.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	log.Println("Connected to Database")
}
