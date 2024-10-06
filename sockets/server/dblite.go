package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
    db *sql.DB
}

type User struct {
    id int
    username string
}

func (s *Sqlite) OpenDBConnection(path string) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        fmt.Println(err)
    }

    s.db = db
}

func (s *Sqlite) CheckIfUserExists(username string) bool {
    statement := "SELECT rowid, USERNAME FROM USER WHERE USERNAME = ?"
    var user User
    error := s.db.QueryRow(statement, username).Scan(&user.id, &user.username)
    if error != nil {
        fmt.Printf("error querrying user %v\n", error)
    }
    return error == nil 
}

func (s *Sqlite) CreateUser(username string) bool {
    statement := "INSERT INTO USER(username, id) VALUES(?, 1)"
    result, err := s.db.Exec(statement, username)
    if err != nil {
        fmt.Printf("Error inserting user %s\n", username)
    }
    fmt.Printf("Created User %s\n", username)

    return result != nil
}

func (s *Sqlite) getUser(username string) User{
    statement := "SELECT rowid, username FROM USER WHERE USERNAME = ?"
    var user User
    error := s.db.QueryRow(statement, username).Scan(&user.id, &user.username)
    if error != nil {
        fmt.Println(error)
    }
    return user
}

func executeStatement(statement string) bool {

    return true;
}
