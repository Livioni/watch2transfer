package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	// 监听端口
	listener, err := net.Listen("tcp", "0.0.0.0:8082")
	fmt.Println("正在监听端口:0.0.0.0:8082")
	if err != nil {
		fmt.Println("监听端口失败:", err)
		return
	}
	defer listener.Close()

	// 循环接收文件
	for {
		// 接受连接
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("接受连接失败:", err)
		}
		// 启动协程处理文件传输
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// 接收文件名和大小
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("接收文件信息失败:", err)
		return
	}
	fileInfo := strings.Split(string(buffer[:n]), ",")
	filename := fileInfo[0]
	fileSize := fileInfo[1]

	// 创建文件
	file, err := os.Create("received_files/" + filename)
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	fmt.Println("创建文件", filename, "成功")

	// 接收文件内容
	var receivedBytes int64
	fileSizeInt64, _ := strconv.ParseInt(fileSize, 10, 64)
	for {
		if (fileSizeInt64 - receivedBytes) < 1024 {
			io.CopyN(file, conn, fileSizeInt64-receivedBytes)
			conn.Read(make([]byte, (receivedBytes+1024)-fileSizeInt64))
			break
		}
		_, err = io.CopyN(file, conn, 1024)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("接收文件内容完成:", err)
			continue
		}
		receivedBytes += 1024
		if receivedBytes >= fileSizeInt64 {
			break
		}
	}

	// 关闭文件
	file.Close()
	fmt.Println("文件接收完成")
}
