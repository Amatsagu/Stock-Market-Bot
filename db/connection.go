package db

import (
	"database/sql"
	"index-bot/logger"
	"index-bot/other"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Conn *sql.DB
var MarketRest other.MarketStackRest

func InitConnectionWithDatabase() {
	// Schema: username:password@tcp(127.0.0.1:3306)/database
	db, err := sql.Open("mysql", os.Getenv("DB_DSL"))
	if err != nil {
		logger.Error.Panic(err) // Unrecoverable!
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 4)

	Conn = db
}

func InitMarketStackRestAPIController() {
	MarketRest = other.CreateMarketStackRest(os.Getenv("MARKET_API_KEY"), 1000)
}
