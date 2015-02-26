package main

import "net"
import "fmt"
import "strings"
import "io/ioutil"
import "runtime"
import "os"
import "strconv"

var rootDir string = ""
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
            rootDir = os.Args[i]
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
    fmt.Println("Launching server...")

    ln, err := net.Listen("tcp", ":8888")
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
    fmt.Println(header[0])
    request := strings.Split(header[0], " ")


    if request[0] == "GET" || request[0] == "HEAD" {
	    if request[1] == "/" {

	    	file, err := ioutil.ReadFile(rootDir + "index.html")
			if err != nil {
		        fmt.Println("Error readFile: ", err)
		    }

			var response_header string = "HTTP/1.1 200 OK\nContent-Type: text/html\nDate: Sun, 22 Feb 2015 20:40:57 GMT\nContent-Length:" + string(len(string(file))) + "\n\n"
	    	_, errr := conn.Write([]byte(response_header + string(file)))
			if errr != nil {
		        fmt.Println("Error write: ", errr)
		    }
	    } else {
            var path string = strings.TrimPrefix(request[1], "/")
            fmt.Println(path)
            
            file, err := ioutil.ReadFile(string(path))
            if err != nil {
                fmt.Println("Error readFile: ", err)
            }

            var response_header string = "HTTP/1.1 200 OK\nContent-Type: image\nDate: Sun, 22 Feb 2015 20:40:57 GMT\nContent-Length:" + string(len(string(file))) + "\n\n"
            _, errr := conn.Write([]byte(response_header + string(file)))
            if errr != nil {
                fmt.Println("Error write: ", errr)
            }
        }



	}
}













