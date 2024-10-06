package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "modernc.org/sqlite"
)

type Sqlite struct {
    db *sql.DB
}

type User struct {
    id int
    username string
}

func (s *Sqlite) CreateTables() {
    userTableStatement := "CREATE TABLE User (username string)"
    _, err := s.db.Exec(userTableStatement)
    if err != nil {
        log.Fatalf("Error creating tables")
    }
}

func (s *Sqlite) ClearTable(tableName string) {
    if strings.ToLower(tableName) == "user" {
        statement := "DELETE FROM USER"
        _, err := s.db.Exec(statement, tableName)
        if err != nil {
            log.Fatalf("Error deleting data from table %v", err)
        }
    }
}

func (s *Sqlite) CheckIfTableExist(tableName string) bool {
    statement := "SELECT name from SQLITE_MASTER WHERE type='table' AND name=?"
    var result string
    err := s.db.QueryRow(statement, tableName).Scan(&result)
    if err != nil {
        log.Println("Error querying mastser table")
    }
    return result == tableName
}

func (s *Sqlite) OpenDBConnection(path string) {
    db, err := sql.Open("sqlite", path)
    if err != nil {
        log.Fatalln(err)
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
    statement := "INSERT INTO USER(username) VALUES(?)"
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
