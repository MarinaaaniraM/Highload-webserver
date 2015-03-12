package main

import ( 
    "net"
    "fmt"
    "strings"
    "io/ioutil"
    "runtime"
    "os"
    "strconv"
    "parser"
    "net/url"
    "time"
)

var rootDir string = "static"
var ncpu int = 1
var isLogging bool =  false

func main() {
    i := 2
    for i <= len(os.Args) {
        if i == len(os.Args) {
            switch os.Args[i - 1] {
            case "-l":
                isLogging = true
                break
            default:
                fmt.Println("________________________________________")
                fmt.Println("\n'-r /path/to/dir' path to root directory\n'-c 2'            num of CPU\n")
                fmt.Println("Or run without arguments in defalt\n")
                fmt.Println("________________________________________")
                os.Exit(0)
            }
        }
        switch os.Args[i - 1] {
        case "-r": 
            rootDir = strings.TrimRight(os.Args[i], "/")
        case "-c":
            i, er := strconv.Atoi(os.Args[i])
            if er != nil {
                logging("Error strconv: " + er.Error())
            }
            ncpu = i
        }
        i = i + 2
    }

    runtime.GOMAXPROCS(ncpu)

    logging("\n")
    logging("Launching server...")

    ln, err := net.Listen("tcp", ":80")                                       //TODO
    if err != nil {
        logging("Error listen: " + err.Error())
    }
    for {
        conn, err := ln.Accept()
        if err != nil {
            logging("Error accept: " + err.Error())
        }
        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    // runtime.Gosched()
    
    var buf = make ([]byte, 1024)
    conn.Read(buf)

    header := strings.Split(string(buf), "\n")
    header = strings.Split(header[0], " ")
    var path string
    var responseStatus string = "200 OK"                                        //TODO
    var fileString string

    // fmt.Println("Header: ", header)

    if header[0] == "GET" || header[0] == "HEAD" {

        url, err := url.Parse(header[1])
        if err != nil {
            logging("Error parse url")
        }
        logging("Request: " + url.Path)

        if strings.Count(url.Path, "../") > strings.Count(rootDir, "/") {
            responseStatus = "403 Forbidden"                                    //TODO

        } else if strings.HasSuffix(url.Path, "/") {
            path = rootDir + strings.Split(url.Path, "?")[0] + "index.html"

            file, err := ioutil.ReadFile(path)
            if err != nil {
                logging("Error readFile: " + err.Error())
                responseStatus = "403 Forbidden"                                //TODO
            }
            fileString = string(file)

	    } else {
            path= rootDir + strings.Split(url.Path, "?")[0]

            file, err := ioutil.ReadFile(path)
            if err != nil {
                logging("Error readFile: " + err.Error())
                responseStatus = "404 Not Found"                                //TODO
            }
            fileString = string(file)
        }
        logging("Path: " + path)

	} else {
        responseStatus = "405 Method Not Allowed"                               //TODO
    }

    var httpState string = "HTTP/1.1 " + responseStatus + "\r\n"                //TODO
    var httpContentType string = "Content-Type: " + parser.GetContentType(path) + "\r\n"
    var httpContentLength string = "Content-Length: " + strconv.Itoa(len(fileString)) + "\r\n"
    var currentDate string = "Date: " + time.Now().Format(time.RFC850) + "\r\n"
    var serverName string = "Server: Cheburashka_v0.1\r\n"

    var responseHeader string = httpState + httpContentType + httpContentLength + currentDate + serverName + "\r\n"
    // fmt.Println("Response header: \r\n", responseHeader)

    if header[0] == "GET" {
        responseHeader = responseHeader + fileString
    }
    _, errr := conn.Write([]byte(responseHeader))
    if errr != nil {
        logging("Error write: " + errr.Error())
    }

    logging("\n")
}


func logging(str string) {
    if isLogging {
        fmt.Println(str)
    }
}













