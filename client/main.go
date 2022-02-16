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

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// read host menu
	reply := make([]byte, 1024)
	_, err = conn.Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", reply)

	// read opc
	n, err := bufio.NewReader(os.Stdin).Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	// write opc on connection
	_, err = conn.Write(reply[:n-1])
	if err != nil {
		log.Fatal(err)
	}

	// read name file
	n, err = conn.Read(reply)
	if err != nil {
		log.Fatal(err)
	}

	// create file
	f, err := os.Create(string(reply[:n]))
	if err != nil {
		log.Print(err)
	}
	defer f.Close()

	for {
		// read content file
		n, err = conn.Read(reply)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		// write file
		n, err = f.Write(reply[:n])
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("File downloaded")

}
