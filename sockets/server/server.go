package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"time"

//	"strconv"
	"sync"
)

type Chatter struct {
    name string
    connection net.Conn 
}

type Server struct {
    chatters []Chatter
    messages []Message
    synchro sync.Mutex

}

var numberOfChatters int

type chatroom struct {
    chatters []Chatter
    messages []Message
}

type Message struct {
    chatter Chatter
    content string
    timeSent time.Time
}

func main() {
    numberOfChatters = 0
    var chatters []Chatter 
    var messages []Message

    server := Server{chatters, messages, sync.Mutex{}} 
    listener := server.RunServer()
    defer listener.Close()

    for {
        newChatter := server.AddAChatter(listener)
        go server.HandleClientConnection(newChatter)
    }
}

func (s *Server) AddAChatter(listener net.Listener) Chatter {
    log.Println("Waiting for connection")
    conn, err := listener.Accept()
    if err != nil {
        log.Fatalf("error accepting connection: %s", err)
    }

    s.synchro.Lock()
    //newChatter := Chatter{"chatter " + strconv.Itoa(len(s.chatters) + 1), conn}
    newChatter := Chatter{fmt.Sprintf("chatter %d", numberOfChatters), conn}
    numberOfChatters++
    log.Printf("Chatter %s connected", newChatter.name)

    s.chatters = append(s.chatters, newChatter)
    log.Printf("Number of chatters %d", len(s.chatters))
    s.synchro.Unlock()
    s.displayChatHistory(newChatter)

    return newChatter
}

func (s *Server) RunServer() net.Listener {
    runtime.GOMAXPROCS(runtime.NumCPU())

    log.Println("starting server")
    listener, err := net.Listen("tcp", "localhost:1234")
    if err != nil {
        log.Fatalf("error listening: %s", err)
    }
    log.Println("server is running")
    return listener;
}

func (s *Server) HandleClientConnection(chatter Chatter) {
    log.Println("reading from the client")

    for {
        data, err := bufio.NewReader(chatter.connection).ReadString('\n')
        if err != nil {
            log.Printf("Chatter %s disconnected\n", chatter.name)
            chatter.connection.Close()
            log.Printf("Chatters left: %s", s.chatters)
            s.chatters = s.RemoveChatter(chatter)
            break
        }
        log.Printf("message from: %s: %s", chatter.name, data)
        s.messages = append(s.messages, Message{chatter, data, time.Now()})
        s.Broadcast(data, chatter)
    }
}

func (s *Server) RemoveChatter(chatterToDelete Chatter) []Chatter{
    var index int
    for i := range s.chatters {
        if s.chatters[i].name == chatterToDelete.name {
            index = i 
        }
    }

    tempSlice := make([]Chatter, 0)
    tempSlice = append(tempSlice, s.chatters[:index]...)

    if index != len(s.chatters) - 1 {
        return append(tempSlice, s.chatters[index+1])
    }
    return tempSlice
}

func (s *Server) Broadcast(message string, chatter Chatter) {
    s.displayChatters()
    log.Printf("Broadcasting message to chatters:")
    for i := range s.chatters {
        if s.chatters[i].name != chatter.name {
            _, err := io.WriteString(s.chatters[i].connection, message)
            //n, err := bufio.NewWriter(conn).Write([]byte(testString))
            if err != nil {
                log.Println("error broadcasting message: ", err)
            }
        }
    }
}

func (s *Server) displayChatters() {
    for i := range s.chatters {
        log.Println(s.chatters[i].name)
    }
}

func (s *Server) displayChatHistory(chatter Chatter) {
    log.Println(s.messages)
    for _, message := range s.messages {
        log.Printf("Writing %s \n", message.content)
        _, err := io.WriteString(chatter.connection, formatMessage(message))
        if err != nil {
            log.Println("error displaying history: ", err)
        }
        log.Println("Wrote message")

    }
}

func formatMessage(message Message) string {
    return fmt.Sprintf("%s| %s | %s \n", message.timeSent.Format(time.DateTime), message.chatter.name, message.content)
}

