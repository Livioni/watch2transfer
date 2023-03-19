package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080") // 将地址绑定到任意地址
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Server is listening on 0.0.0.0:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("Client %s connected\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Read file from disk
	file, err := os.Open("file.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}

	fileSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(fileSize, uint32(fileInfo.Size()))

	// Send file size to client
	_, err = conn.Write(fileSize)
	if err != nil {
		log.Println(err)
		return
	}

	// Send file contents to client
	_, err = io.Copy(conn, file)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("File sent to %s\n", conn.RemoteAddr().String())
}
