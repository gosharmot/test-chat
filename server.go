package main

import (
	"fmt"
	"net"
	"strconv"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")

	go broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	cl := client{make(chan string), ID, ""}
	ID++

	go sendMessage(conn, cl.ch)

	cl.ch <- "User_" + strconv.Itoa(cl.id) + ", Hi!"

	enter <- cl

	for {
		input := make([]byte, (1024 * 4))
		n, _ := conn.Read(input)
		msg := string(input[0:n])
		if msg == `\q` {
			break
		} else {
			cl.msg = msg
			message <- cl
		}
	}

	leave <- cl
}

func broadcast() {
	clients := make(map[int]client)

	for {
		select {
		case cl := <-enter:
			clients[cl.id] = cl
			for id, client := range clients {
				if id != cl.id {
					client.ch <- "User_" + strconv.Itoa(cl.id) + " connect"
				}
			}
		case cl := <-leave:
			delete(clients, cl.id)
			close(cl.ch)
			ID--
			for id, client := range clients {
				if id != cl.id {
					client.ch <- "User_" + strconv.Itoa(cl.id) + " leave"
				}
			}
		case cl := <-message:
			for id, client := range clients {
				if id != cl.id {
					client.ch <- "User_" + strconv.Itoa(cl.id) + ": " + cl.msg
				}
			}
		}
	}
}

func sendMessage(conn net.Conn, ch <-chan string) {
	for str := range ch {
		fmt.Fprintln(conn, str)
	}
}

var (
	ID      = 1
	enter   = make(chan client)
	leave   = make(chan client)
	message = make(chan client)
)

type client struct {
	ch  chan string
	id  int
	msg string
}
