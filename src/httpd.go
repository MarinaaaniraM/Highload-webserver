package main

import "net"
import "fmt"
import "strings"
import "io/ioutil"
import "runtime"
import "os"
import "strconv"
import "parser"
import "net/url"

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
    var path string
    var responseStatus string = "200 OK"
    var fileString string

    fmt.Println("Header: ", header)

    if header[0] == "GET" || header[0] == "HEAD" {

        url, err := url.Parse(header[1])
        if err != nil {
            fmt.Println("Error parse url")
        }

        fmt.Println("Request: " + url.Path)


        //TODO date!!!
        

        if strings.Count(url.Path, "../") > strings.Count(rootDir, "/") {
            responseStatus = "403 Forbidden"                                    //TODO

        } else if strings.HasSuffix(url.Path, "/") {
            path = rootDir + strings.Split(url.Path, "?")[0] + "index.html"

            file, err := ioutil.ReadFile(path)
            if err != nil {
                fmt.Println("Error readFile: ", err)
                responseStatus = "403 Forbidden"                                //TODO
            }
            fileString = string(file)

	    } else {
            path= rootDir + strings.Split(url.Path, "?")[0]

            file, err := ioutil.ReadFile(path)
            if err != nil {
                fmt.Println("Error readFile: ", err)
                responseStatus = "404 Not Found"                                //TODO
            }
            fileString = string(file)
        }
        fmt.Println("Path: " + path)

	} else {
        responseStatus = "405 Method Not Allowed"                               //TODO
    }

    var httpState string = "HTTP/1.1 " + responseStatus + "\r\n"
    var httpContentType string = "Content-Type: " + parser.GetContentType(path) + "\r\n"
    var httpContentLength string = "Content-Length: " + strconv.Itoa(len(fileString)) + "\r\n"
    var serverName string = "Server: Cheburashka_v0.1\r\n"

    var responseHeader string = httpState + httpContentType + httpContentLength + serverName + "\r\n"
    fmt.Println("Response header: \r\n", responseHeader)

    if header[0] == "GET" {
        responseHeader = responseHeader + fileString
    }
    _, errr := conn.Write([]byte(responseHeader))
    if errr != nil {
        fmt.Println("Error write: ", errr)
    }

    fmt.Println("\n")
}













