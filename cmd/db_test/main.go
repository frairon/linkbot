package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:data.db")

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	db.Exec("insert into users (id, name) values ( 1234, \"Peter\")")

	rows, _ := db.Query("select * from users;")
	for rows.Next() {
		cols, _ := rows.Columns()
		log.Printf("cols: %#v", cols)

		var id int64
		var name string
		rows.Scan(&id, &name)

		log.Printf("id: %d, name: %s", id, name)
	}
	defer db.Close()
}
