package main

import (
    "database/sql"
    _ "github.com/lib/pq"
    "fmt"
)

func GetDatabase() (*sql.DB, error) {
    dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
        DB_USER, DB_PASSWORD, DB_NAME)
    db, err := sql.Open("postgres", dbinfo)
    checkErr(err)
    return db, err
}
