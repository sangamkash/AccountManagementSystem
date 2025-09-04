package test

import (
	"AccountManagementSystem/log_color"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"log/slog"
)

func GetDB() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", dbUser, dbPwd, host, port, dbName, schema)
	slog.Info(log_color.Yellow("connecting :" + connStr))
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		slog.Error(log_color.Red("Make use postgres is running with defined table"))
		log.Fatal(err)
		return nil
	}
	slog.Info(log_color.Green("connected :" + connStr))
	return db
}

func ClearAccount(db *sql.DB) {
	if _, delErr := db.Exec("DELETE FROM accounts"); delErr != nil {
		slog.Error(log_color.Red("failed to delete accounts"))
		log.Fatal(delErr)
	}
}
func ClearTransaction(db *sql.DB) {
	if _, delErr := db.Exec("DELETE FROM transactions"); delErr != nil {
		slog.Error(log_color.Red("failed to delete accounts"))
		log.Fatal(delErr)
	}
}
