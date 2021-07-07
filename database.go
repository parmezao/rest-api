package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func GetDatabase() sql.DB {
	db, err := sql.Open("mysql", "root:123456@/contacts")
	if err != nil {
		log.Fatal(err)
	}
	return *db
}
