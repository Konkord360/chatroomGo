package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
    db *sql.DB
}

func (s *Sqlite) OpenDBConnection(path string) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        fmt.Println(err)
    }

    s.db = db
}

func (s *Sqlite) CheckIfUserExists(username string) bool {
    statement := "SELECT * FROM USER WHERE NAME = ?"
    row := s.db.QueryRow(statement, username)
    fmt.Printf("Row %v", row)
    return row.Scan() == nil 
}

func (s *Sqlite) CreateUser(username string) bool {
    statement := "INSERT INTO USER(username, id) VALUES(?, 1)"
    result, err := s.db.Exec(statement, username)
    if err != nil {
        fmt.Printf("Error inserting user %s", username)
    }
    fmt.Printf("Created User %s", username)

    return result != nil
}

func executeStatement(statement string) bool {

    return true;
}
