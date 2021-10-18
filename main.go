package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type message struct {
	user         net.Conn
	body         string
	notification string
	time         string
}

var (
	msg_repo         []string
	numOfConnections int
	conn_repo        map[string]net.Conn
)

const maxConnections = 10
const timeFormat = "2006-01-02 15:04:05"

func main() {
	port := ""
	args := os.Args[1:]
	conn_repo = make(map[string]net.Conn)

	if len(args) == 0 {
		port = ":8989"
	} else if len(args) == 1 {
		port = ":" + args[0]
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port_number")
		return
	}

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Listening on the port %s\n", port)

	ch := make(chan message)

	go sendToChannels(ch)

	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			log.Fatal(err)
			return
		}
		numOfConnections++
		if numOfConnections <= maxConnections {
			go handleConnection(conn, ch)
		} else {
			conn.Write([]byte("Too many users, try later.\n"))
			conn.Close()
		}
	}
}

func handleConnection(conn net.Conn, ch chan<- message) {
	defer conn.Close()

	printTux(conn)

	name, err := getName(conn)

	if err != nil {
		log.Fatal(err, " : can't get name")
		return
	}
	conn_repo[name] = conn

	for _, message := range msg_repo {
		conn.Write([]byte(message))
		conn.Write([]byte("\n"))
	}

	notif_time := time.Now().Format(timeFormat)
	connMessage := message{notification: "\n" + name + " has joined the chat...\n",
		body: "",
		time: "[" + notif_time + "]",
		user: conn}

	ch <- connMessage

	for {
		cur_time := time.Now().Format(timeFormat)
		conn.Write([]byte("[" + cur_time + "]" + "[" + name + "]" + ": "))
		msg, _, err := bufio.NewReader(conn).ReadLine()

		if err != nil {
			numOfConnections--
			log.Println(name + " disconnected")
			connMessage := message{notification: "\n" + name + " has left the chat...\n",
				body: "",
				time: "[" + cur_time + "]",
				user: conn}

			ch <- connMessage
			break
		}

		if string(msg) != "" {
			connMessage := message{body: "[" + name + "]" + ": " + string(msg) + "\n",
				time: "\n[" + cur_time + "]",
				user: conn}

			ch <- connMessage
			msg_repo = append(msg_repo, connMessage.time[1:]+connMessage.body[:len(connMessage.body)-1])
		}
	}
}
