package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	str := ""
	go A(conn)
	for {
		fmt.Scan(&str)
		if str == `\q` {
			conn.Write([]byte(str))
			break
		} else {
			conn.Write([]byte(str))
		}
	}
}
func A(conn net.Conn) {
	io.Copy(os.Stdout, conn)
}
