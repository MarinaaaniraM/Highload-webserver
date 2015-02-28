package main

import "net"
import "fmt"
import "strings"
import "io/ioutil"
import "runtime"
import "os"
import "strconv"
import "parser"

var rootDir string = "."
var ncpu int = 1

func main() {
    i := 2
    for i <= len(os.Args) {
        if i == len(os.Args) {
            fmt.Println("________________________________________")
            fmt.Println("\n'-r /path/to/dir' path to root directory\n'-c 2'            num of CPU\n")
            fmt.Println("Or run without arguments in defalt\n")
            fmt.Println("________________________________________")
            os.Exit(0)
        }
        switch os.Args[i - 1] {
        case "-r": 
            rootDir = strings.TrimRight(os.Args[i], "/")
        case "-c":
            i, er := strconv.Atoi(os.Args[i])
            if er != nil {
                fmt.Println("Error strconv: ", er)
                // os.Exit(2)
            }
            ncpu = i
        }
        i = i + 2
    }

    runtime.GOMAXPROCS(ncpu)

    fmt.Println("\n")
    fmt.Println("Launching server...")

    ln, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error listen: ", err)
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accept: ", err)
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    runtime.Gosched()
    defer conn.Close()

    var buf = make ([]byte, 1024)
    conn.Read(buf)

    header := strings.Split(string(buf), "\n")
    header = strings.Split(header[0], " ")
  
    fmt.Println("Header: ", header)

    if header[0] == "GET" || header[0] == "HEAD" {
        fmt.Println("Request: " + header[1])
        var path string

	    if strings.HasSuffix(header[1], "/") {
            path = rootDir + header[1] + "index.html"
	    } else {
            path= rootDir + header[1]
        }
        fmt.Println("Path: " + path)

        file, err := ioutil.ReadFile(path)
        if err != nil {
            fmt.Println("Error readFile: ", err)
        }

        var httpState string = "HTTP/1.1 200 OK\n"
        var httpContentType string = "Content-Type: " + parser.GetContentType(path) + "\n"
        var httpContentLength string = "Content-Length: " + strconv.Itoa(len(file)) + "\n"

        var responseHeader string = httpState + httpContentType + httpContentLength + "\n"
        _, errr := conn.Write([]byte(responseHeader + string(file)))
        if errr != nil {
            fmt.Println("Error write: ", errr)
        }
        fmt.Println("Response header: \n", responseHeader)
	}
    fmt.Println("\n")
}













