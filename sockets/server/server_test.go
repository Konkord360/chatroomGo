package main

import (
	"fmt"
	"log"
	"net"
    //"os"
	//"os/exec"
	"sync"
	"testing"
)


func TestAddAChatterExtendTheChatterSliceInTheServer(t *testing.T) {
    var chatters []Chatter 
    var wg sync.WaitGroup
    var chatHistory []Message

    server := Server{chatters, chatHistory, sync.Mutex{}}
    server.RunServer();
    defer listener.Close()

    wg.Add(1)
    go func(net.Listener) {
        defer wg.Done()
        server.AddAChatter()
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
    server.RunServer();
    defer listener.Close()

    log.Println("Adding chatters")
    wg.Add(2)

    conn, _ := net.Dial("tcp", "localhost:1234")
    defer conn.Close()

    conn1, _ := net.Dial("tcp", "localhost:1234")
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
}


func TestAddingChattersWithTwoTousandConnections(t *testing.T) {
    var chatters []Chatter 
    var chatHistory []Message
    expectedNumberOfChatters := 100000
    server := Server{chatters, chatHistory, sync.Mutex{}}
    server.RunServer();
    defer listener.Close()
    //cmd := exec.Command("clear")

    log.Println("Adding chatters")
    for i := 0; i < expectedNumberOfChatters; i++ {
        //cmd.Stdout = os.Stdout
        //err := cmd.Run() 
        //if err != nil {
        //    log.Print(err)
        //}

        conn, err := net.Dial("tcp", "localhost:1234")
        if err != nil {
            log.Panicf("Failed to establish connection %d, %v", i, err)
        }
        fmt.Print("\033[H\033[2J")
        server.AddAChatter()
        log.Printf("making connection %d", i)
        defer conn.Close()
    }

    numberOfChatters := len(server.chatters)
    if numberOfChatters != expectedNumberOfChatters {
        t.Fatalf("Expected number of chatters: %d, got %d", expectedNumberOfChatters, numberOfChatters)
    }
}
