package main

import (
    //"fmt"
    //"io"
    //"log"
    "fmt"
    "io"
    "log"
    "net"

    //"os"
    //"os/exec"
    "sync"
    "testing"
)

const testDatabase = "file:./test.db"

func TestAddAChatterExtendTheChatterSliceInTheServer(t *testing.T) {
    var chatters []Chatter 
    var wg sync.WaitGroup
    var chatHistory []Message

    server := Server{chatters, chatHistory, sync.Mutex{}, Sqlite{}}
    server.db.OpenDBConnection(testDatabase)
    server.RunServer();
    defer listener.Close()

    wg.Add(1)
    go func(net.Listener) {
        defer wg.Done()
        server.AddAChatter()
    }(listener)

    conn, _ := net.Dial("tcp", "localhost:1234")
    _, err := io.WriteString(conn, "test\n")
    if err != nil {
        t.Fatalf("Could notwrite username")
    }

    defer conn.Close()

    wg.Wait()
    numberOfChatters := len(server.chatters)
    if numberOfChatters != 1 {
        t.Fatalf("Expected number of chatters: %d, got %d", 1, numberOfChatters)
    }
    server.db.ClearTable("User")
}

func TestAddingTwoChattersAtTheSameTimeAdsBothOfThemCorrectly(t *testing.T) {
    var chatters []Chatter 
    var wg sync.WaitGroup
    var chatHistory []Message

    server := Server{chatters, chatHistory, sync.Mutex{}, Sqlite{}}
    server.db.OpenDBConnection(testDatabase)
    server.RunServer();
    defer listener.Close()

    log.Println("Adding chatters")
    wg.Add(2)

    conn, _ := net.Dial("tcp", "localhost:1234")
    _, err := io.WriteString(conn, "test\n")
    if err != nil {
        t.Fatalf("Could notwrite username")
    }
    defer conn.Close()

    conn1, _ := net.Dial("tcp", "localhost:1234")
    _, err = io.WriteString(conn1, "test\n")
    if err != nil {
        t.Fatalf("Could notwrite username")
    }
    defer conn1.Close()

    go func() {
        defer wg.Done()
        server.AddAChatter()
    }()

    go func() {
        defer wg.Done()
        server.AddAChatter()
    }()

    wg.Wait()
    numberOfChatters := len(server.chatters)
    if numberOfChatters != 2 {
        t.Fatalf("Expected number of chatters: %d, got %d", 2, numberOfChatters)
    }
    server.db.ClearTable("User")
}

func TestAddingChattersWithTwoTousandConnections(t *testing.T) {
    var chatters []Chatter 
    var chatHistory []Message
    expectedNumberOfChatters := 500
    server := Server{chatters, chatHistory, sync.Mutex{}, Sqlite{}}
    server.db.OpenDBConnection(testDatabase)
    server.RunServer();
    defer listener.Close()

    log.Println("Adding chatters")
    for i := 0; i < expectedNumberOfChatters; i++ {
        conn, err := net.Dial("tcp", "localhost:1234")
        if err != nil {
            log.Panicf("Failed to establish connection %d, %v", i, err)
        }
        _, err = io.WriteString(conn, fmt.Sprintf("test%d\n", i))
        if err != nil {
            t.Fatalf("Could notwrite username")
        }
        defer conn.Close()
        fmt.Print("\033[H\033[2J")
        server.AddAChatter()
        log.Printf("making connection %d", i)
        defer conn.Close()
    }

    numberOfChatters := len(server.chatters)
    if numberOfChatters != expectedNumberOfChatters {
        t.Fatalf("Expected number of chatters: %d, got %d", expectedNumberOfChatters, numberOfChatters)
    }
    server.db.ClearTable("User")
}
