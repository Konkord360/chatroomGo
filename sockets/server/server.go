package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"runtime"
	"strings"
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
    db Sqlite

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

var listener net.Listener

func main() {
    numberOfChatters = 0
    var chatters []Chatter 
    var messages []Message

    var server Server 
    server.db = Sqlite{}
    server.db.OpenDBConnection("file:./prodChatroom.db")
    server.chatters = chatters
    server.messages = messages
    server.synchro = sync.Mutex{}

    server.RunServer()
    defer server.db.db.Close()
    defer listener.Close()

    for {
        newChatter := server.AddAChatter()
        go server.HandleClientConnection(newChatter)
    }
}

func (s *Server) AddAChatter() Chatter {
    log.Println("Waiting for connection")
    s.synchro.Lock()
    conn, err := listener.Accept()
    if err != nil {
        log.Fatalf("error accepting connection: %s", err)
    }

    //newChatter := Chatter{"chatter " + strconv.Itoa(len(s.chatters) + 1), conn}

    newChatter := Chatter{fmt.Sprintf("chatter %d", numberOfChatters), conn}
    s.WriteToChatter("Provide username:\n", newChatter)
    data, err := bufio.NewReader(newChatter.connection).ReadString('\n')
    if err != nil {
        log.Fatalf("Error reading username")
    }
    newChatter.name = strings.TrimRight(data, "\n")
    s.WriteToChatter(fmt.Sprintf("Connected to chatroom as %s:\n", newChatter.name), newChatter)

    if !s.db.CheckIfUserExists(newChatter.name) {
        fmt.Printf("Creating new user %s\n", newChatter.name)
        s.db.CreateUser(newChatter.name) 
    } else {
        fmt.Printf("User found in database %s\n", newChatter.name)
        s.db.getUser(newChatter.name)
    }


    numberOfChatters++
    log.Printf("Chatter %s connected", newChatter.name)

    s.chatters = append(s.chatters, newChatter)
    log.Printf("Number of chatters %d", len(s.chatters))
    s.synchro.Unlock()
    s.displayChatHistory(newChatter)

    return newChatter
}

func (s *Server) RunServer() {
    runtime.GOMAXPROCS(runtime.NumCPU())

    if !s.db.CheckIfTableExist("User") {
        log.Println("Required tables are missing. Trying to create them")
        s.db.CreateTables()
        log.Println("Tables created successfully")
    }
    log.Println("Required tables are ready")

    log.Println("starting server")
    list, err := net.Listen("tcp", "localhost:1234")
    if err != nil {
        log.Fatalf("error listening: %s", err)
    }
    log.Println("server is running")
    listener = list;
}

func (s *Server) HandleClientConnection(chatter Chatter) {
    log.Println("reading from the client")

    //TODO move reader and write to the chatter. 
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

func (s *Server) WriteToChatter(message string, chatter Chatter) {
    log.Println("Waiting for user to provide username")
    _, err := io.WriteString(chatter.connection, message)
    log.Printf("Writing %s to conn %v\n", message, chatter.connection)
    if err != nil {
        log.Println("error writing message to the chatter: ", err)
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

