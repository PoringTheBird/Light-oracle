package main

import (
	"fmt"
	"log"
	"main/Core"
	"net"
	"os"
)

func main() {
	go startBot()
	startSocket()
}

func startSocket() {
	err := listenPort()
	if err != nil {
		log.Println("Failed to listen to port: ", err)
		return
	}
}

func startBot() {
	bot := Core.Bot{}
	err := bot.Start()

	if err != nil {
		log.Fatal(err)
		return
	}
}

func listenPort() error {
	host := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	log.Println("Start listening of ", host)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go handleConnection(conn)
	}

	ln.Close()
	return nil
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


