package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
)

func printTux(conn net.Conn) error {
	file, err := ioutil.ReadFile("welcome.txt")
	if err != nil {
		return err
	}

	conn.Write(file)
	return nil
}

func getName(conn net.Conn) (string, error) {
	name := ""

	for name == "" {
		conn.Write([]byte("[ENTER YOUR NAME]: "))
		bufName, _, err := bufio.NewReader(conn).ReadLine()
		if err != nil {
			return "", err
		}
		name = string(bufName)
	}
	log.Println(name + " connected")

	return name, nil
}

func sendToChannels(ch <-chan message) {
	for {
		msg := <-ch
		for name, conn := range conn_repo {
			if conn == msg.user {
				continue
			}

			conn.Write([]byte(msg.notification))
			if len(msg.body) > 0 {
				conn.Write([]byte(msg.time + msg.body[:len(msg.body)-1]))
			}
			conn.Write([]byte(msg.time + "[" + name + "]" + ": "))
		}
		msg.notification = ""
	}
}
