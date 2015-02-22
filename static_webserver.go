package main

import "net"
import "fmt"
import "strings"
// import "io/ioutil"

func main() {

    fmt.Println("Launching server...")

    ln, err := net.Listen("tcp", ":8888")
    if err != nil {
        fmt.Println("Error listen")
    }

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Println("Error accept")
        }

        go handleConnection(conn)
    }
}

func handleConnection(conn net.Conn) {
    defer conn.Close()

    var buf = make ([]byte, 1024)

    // var buffer []byte
    // buffer, err := ioutil.ReadAll(conn)
    // if err != nil {
    //   fmt.Println("Error readAll")
    // }

    conn.Read(buf)
    fmt.Println(string(buf))

    header := strings.Split(string(buf), "\n")
    fmt.Println(header[0])
    request := strings.Split(header[0], " ")
    fmt.Println(request[1])

    if request == "/" {
        // Отдаем index.html
    } else {
        // Отдаем файл по указанному пути
    }
}