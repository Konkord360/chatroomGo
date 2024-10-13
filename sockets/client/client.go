package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	//"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var reader *bufio.Reader
var scanner *bufio.Scanner
var conn net.Conn

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Block struct {
    Id int
}

type Blocks struct {
    Start int
    Next int
    More bool
    Blocks []Block
}

type Count struct{
 Count int
}

type Contacts = []Contact

type Data struct {
    Contacts Contacts
}

func newData() Data{
    return Data{
        Contacts: []Contact{
            newContact("John", "asdf"),
            newContact("Emily", "Emy@gmail.com"),
            newContact("Billy", "BJ@gmail.com"),
        },
    }
}

type Contact struct {
    Name string
    Email string
}

func newContact(name, email string) Contact {
    return Contact{
        Name: name,
        Email: email,
    }
}

func (d Data) hasEmail(email string) bool {
    for _, contact := range d.Contacts {
        if contact.Email == email {
            return true
        }
    }
    return false
}

type FormData struct {
    Values map[string]string
    Errors map[string]string
}

func newFormData() FormData {
    return FormData {
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

type Page struct {
    Data Data
    Form FormData
}

func newPage() Page {
    return Page {
        Data: newData(),
        Form: newFormData(),
    }
}

func main() {
	e := echo.New()
    e.Renderer = NewTemplates()
    e.Use(middleware.Logger())

    page := newPage()
    e.POST("/contacts", func(c echo.Context) error {
        
        name := c.FormValue("name")
        email := c.FormValue("email")
        
        if page.Data.hasEmail(email) {
            formData := newFormData()
            formData.Values["name"] = name
            formData.Values["email"] = email
            formData.Errors["email"] = "Email already exists"

            return c.Render(422, "form", formData)
        }

        page.Data.Contacts = append(page.Data.Contacts, newContact(name, email))

        return c.Render(200, "display", page)
    })

    e.GET("/", func(c echo.Context) error {
        return c.Render(200, "index", page)
    })

    e.GET("/blocks", func(c echo.Context) error {
        startStr := c.QueryParam("start")
        start, err := strconv.Atoi(startStr)
        if err != nil {
            start = 0
        }

        blocks := []Block{}
        for i := start; i < start + 10; i++ {
            blocks = append(blocks, Block{Id: i})
        }

        template := "blocks"
        if start == 0 {
            template = "blocks-index"
        }
        return c.Render(http.StatusOK, template, Blocks{
            Start: start,
            Next: start + 10,
            More: start + 10 < 100,
            Blocks: blocks,
        });
    });

    e.Logger.Fatal(e.Start(":42069"))
}

func mainConsoleInput() {
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
