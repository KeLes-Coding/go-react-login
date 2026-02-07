package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	const (
		host     = "localhost"
		port     = 5433
		user     = "keles"
		password = "c05022007"
		dbname   = "login_demo"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("无法连接到数据库: %v", err)
	}

	log.Println("成功连接到数据库")
	return db
}
