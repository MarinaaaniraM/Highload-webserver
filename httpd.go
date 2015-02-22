package main

import "net"
import "fmt"
import "strings"
import "io/ioutil"


func main() {

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
    defer conn.Close()

    var buf = make ([]byte, 1024)
    conn.Read(buf)

    header := strings.Split(string(buf), "\n")
    fmt.Println(header[0])
    request := strings.Split(header[0], " ")


    if request[0] == "GET" || request[0] == "HEAD" {
	    if request[1] == "/" {

	    	file, err := ioutil.ReadFile("index.html")
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

            var response_header string = "HTTP/1.1 200 OK\nContent-Type: image/jpg\nDate: Sun, 22 Feb 2015 20:40:57 GMT\nContent-Length:" + string(len(string(file))) + "\n\n"
            _, errr := conn.Write([]byte(response_header + string(file)))
            if errr != nil {
                fmt.Println("Error write: ", errr)
            }
        }


        
	}
}













