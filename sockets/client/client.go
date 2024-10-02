package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

var reader *bufio.Reader
var scanner *bufio.Scanner
var conn net.Conn

func main() {
    newConn, err := net.Dial("tcp", "localhost:1234")
    if err != nil {
        log.Fatalf("error dialing server: %s", err)
    }
    conn = newConn
    reader = bufio.NewReader(conn)
    scanner = bufio.NewScanner(os.Stdin)

    go readMessagesFromTheServer()
        //buffer := make([]byte, 1024)
        //testString := "testMessage"
        //copy(buffer[:], testString)
    writeToTheServer()
}

func writeToTheServer() {
    for {
        scanned := scanner.Scan()
        if !scanned {
            return
        }
        line := scanner.Text()
        line = line + "\n"


        _, err := io.WriteString(conn, line)
        //n, err := bufio.NewWriter(conn).Write([]byte(testString))
        if err != nil {
            log.Fatalf("error writing to the server: %s", err)
        }
    }
}

func readMessagesFromTheServer() {
    for {
        data, err := reader.ReadString('\n')
        if err != nil {
            log.Fatalf("Error reading from server: %s \n", err)
        }
        fmt.Print(data)
    }
}
