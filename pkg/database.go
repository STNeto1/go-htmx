package pkg

import (
	"database/sql"
	"log"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func PrepareDB() {
	var dbUrl = "file:file.db"
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatalf("failed to open db %s: %s", dbUrl, err)
	}

	DB = db

	_, err = DB.Exec(base_table_sql)
	if err != nil {
		log.Fatalf("failed to create table: %s", err)
	}

}

var base_table_sql string = `
    CREATE TABLE IF NOT EXISTS links (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    title TEXT NOT NULL,
	    description TEXT,
	    link TEXT NOT NULL,
	    thumbnail text
    );
    `
