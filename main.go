package main

import (
	"fmt"
	"log"
	"main/Core"
	"net"
	"os"
)

func main() {
	bot := Core.Bot{}
	err := bot.Start()

	if err != nil {
		log.Fatal(err)
	}
}

func listenPort() error {
	host := fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	buf := make([]byte, 1024)
	data, err := conn.Read(buf)
	if err != nil {
		log.Println("Data reading error: ", err)
	}

	log.Println("Data received: ", string(data))

	conn.Write([]byte("OK"))
	conn.Close()
}


