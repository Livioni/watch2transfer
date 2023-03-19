package main

import (
	"encoding/binary"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	// Parse command line arguments
	serverAddr := flag.String("addr", "10.1.81.24:8080", "server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Printf("Connected to server at %s\n", *serverAddr)

	// Receive file size from server
	fileSizeBytes := make([]byte, 4)
	_, err = conn.Read(fileSizeBytes)
	if err != nil {
		log.Fatal(err)
	}

	fileSize := binary.LittleEndian.Uint32(fileSizeBytes)

	// Create a file to save received data
	file, err := os.Create("received_file.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Receive file contents from server
	_, err = io.CopyN(file, conn, int64(fileSize))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("File received from server")
}
