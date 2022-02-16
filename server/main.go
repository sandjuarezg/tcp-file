package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatal("Insufficient arguments: [host] [port]")
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]))
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Printf("Listening on %s\n", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(conn)
	}
}

// handleRequest Handle client request
//  @param1 (conn): connection between server and client
//
//  @return1 (err): error variable
func handleRequest(conn net.Conn) (err error) {
	defer conn.Close()

	mess := fmt.Sprintln("Select file to download")
	mess += fmt.Sprintln("1. Flower.txt")
	mess += fmt.Sprintln("2. Butterfly.txt")

	// write menu on conn
	_, err = conn.Write([]byte(mess))
	if err != nil {
		log.Fatal(err)
	}

	reply := make([]byte, 1024)

	// read client response
	res := bufio.NewReader(conn)
	n, err := res.Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	mess = "butterfly.txt"
	if string(reply[:n]) == "1" {
		mess = "flower.txt"
	}

	// write filename on conn
	_, err = conn.Write([]byte(mess))
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(fmt.Sprintf("./file/%s", mess))
	if err != nil {
		log.Fatal(err)
	}

	// write content on conn
	for {
		n, err = file.Read(reply)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		_, err = conn.Write(reply[:n])
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("%s file was sent\n", mess)

	return
}
