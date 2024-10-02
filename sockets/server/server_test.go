package main

import (
    "log"
    "net"
    "sync"
    "testing"
)


func TestAddAChatterExtendTheChatterSliceInTheServer(t *testing.T) {
    var chatters []Chatter 
    var wg sync.WaitGroup
    var chatHistory []Message

    server := Server{chatters, chatHistory, sync.Mutex{}}
    listener := server.RunServer();
    defer listener.Close()

    wg.Add(1)
    go func(net.Listener) {
        defer wg.Done()
        server.AddAChatter(listener)
    }(listener)

    conn, _ := net.Dial("tcp", "localhost:1234")
    defer conn.Close()

    wg.Wait()
    numberOfChatters := len(server.chatters)
    if numberOfChatters != 1 {
        t.Fatalf("Expected number of chatters: %d, got %d", 1, numberOfChatters)
    }
}

func TestAddingTwoChattersAtTheSameTimeAdsBothOfThemCorrectly(t *testing.T) {
    var chatters []Chatter 
    var wg sync.WaitGroup
    var chatHistory []Message

    server := Server{chatters, chatHistory, sync.Mutex{}}
    listener := server.RunServer();
    defer listener.Close()

    log.Println("Adding chatters")
    wg.Add(2)

    conn, _ := net.Dial("tcp", "localhost:1234")
    defer conn.Close()

    conn1, _ := net.Dial("tcp", "localhost:1234")
    defer conn1.Close()

    go func() {
        defer wg.Done()
        server.AddAChatter(listener)
    }()

    go func() {
        defer wg.Done()
        server.AddAChatter(listener)
    }()

    wg.Wait()
    numberOfChatters := len(server.chatters)
    if numberOfChatters != 2 {
        t.Fatalf("Expected number of chatters: %d, got %d", 2, numberOfChatters)
    }
}


func TestAddingChattersWithTwoTousandConnections(t *testing.T) {
    var chatters []Chatter 
    var chatHistory []Message
    var wg sync.WaitGroup
    var wgRun sync.WaitGroup
    expectedNumberOfChatters := 3000

    server := Server{chatters, chatHistory, sync.Mutex{}}
    listener := server.RunServer();
    defer listener.Close()

    log.Println("Adding chatters")
    wg.Add(expectedNumberOfChatters)
    wgRun.Add(expectedNumberOfChatters)

    for i := 0; i < expectedNumberOfChatters; i++ {
        conn, _ := net.Dial("tcp", "localhost:1234")
        defer conn.Close()
    }

    for i := 0; i < expectedNumberOfChatters; i++ {
        go func() {
            defer wg.Done()
            wgRun.Done()
            wgRun.Wait()
            server.AddAChatter(listener)
        }()
    }

    wg.Wait()
    numberOfChatters := len(server.chatters)
    if numberOfChatters != expectedNumberOfChatters {
        t.Fatalf("Expected number of chatters: %d, got %d", expectedNumberOfChatters, numberOfChatters)
    }
}
