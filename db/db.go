package db

import (
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	xormlog "xorm.io/xorm/log"
)

var Engine *xorm.Engine

func InitMySQL() {
	// Database configuration
	cfg := mysql.NewConfig()
	cfg.User = "root"
	cfg.Passwd = ""
	cfg.Net = "tcp"
	cfg.Addr = "localhost:3306"
	cfg.DBName = "fanchiikawa"
	cfg.ParseTime = true

	// Initialize XORM engine
	dsn := cfg.FormatDSN()
	var err error
	Engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to create XORM engine:", err)
	}

	// Configure XORM engine
	Engine.SetMaxIdleConns(10)
	Engine.SetMaxOpenConns(100)

	// Set log level based on environment
	if os.Getenv("DEBUG") == "true" {
		Engine.SetLogLevel(xormlog.LOG_DEBUG)
		Engine.ShowSQL(true)
	} else {
		Engine.SetLogLevel(xormlog.LOG_WARNING)
	}

	// Test connection
	if err := Engine.Ping(); err != nil {
		log.Fatal("Failed to ping database with XORM:", err)
	}

	// Sync database schema
	if err := SyncSchema(); err != nil {
		log.Fatal("Failed to sync database schema:", err)
	}

	log.Println("Connected to Database with XORM")
}

// SyncSchema synchronizes the database schema with the model structs
func SyncSchema() error {
	return Engine.Sync2(new(User), new(UserDevice), new(Image))
}
