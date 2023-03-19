package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/fsnotify/fsnotify"
)

func main() {
	//创建一个监控对象
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()

	fmt.Println("请输入监控目录：")
	var watch_path string
	fmt.Scanf("%s", &watch_path)
	//添加要监控的对象，文件或文件夹
	err = watch.Add(watch_path)
	fmt.Println("监控目录：", watch_path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("请输入server ip:port")
	var server_addr string
	fmt.Scanf("%s", &server_addr)

	//我们另启一个goroutine来处理监控对象的事件
	for {
		select {
		case ev := <-watch.Events:
			{
				//判断事件发生的类型，如下5种
				// Create 创建
				// Write 写入
				// Remove 删除
				// Rename 重命名
				// Chmod 修改权限
				if ev.Op&fsnotify.Create == fsnotify.Create {
					log.Println("创建文件 : ", ev.Name)
					go transfer(ev.Name, server_addr)
				}
				if ev.Op&fsnotify.Write == fsnotify.Write {
					log.Println("写入文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Remove == fsnotify.Remove {
					log.Println("删除文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("重命名文件 : ", ev.Name)
				}
				if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("修改权限 : ", ev.Name)
				}
			}
		case err := <-watch.Errors:
			{
				log.Println("error : ", err)
				return
			}
		}
	}
}

func transfer(filepath string, server_addr string) {
	conn, err := net.Dial("tcp", server_addr)
	if err != nil {
		fmt.Println("连接服务器失败:", err)
		return
	}
	defer conn.Close()

	fmt.Println("连接服务器成功")
	// 打开文件
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("打开文件失败:", err)
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("获取文件信息失败:", err)
	}

	// 发送文件名和大小
	filename := fileInfo.Name()
	fileSize := fileInfo.Size()
	_, err = conn.Write([]byte(fmt.Sprintf("%s,%d", filename, fileSize)))
	if err != nil {
		fmt.Println("发送文件信息失败:", err)
	}

	// 发送文件内容
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("读取文件内容失败:", err)
			continue
		}
		_, err = conn.Write(buffer[:n])
		if err != nil {
			fmt.Println("发送文件内容失败:", err)
			continue
		}
	}
	fmt.Println("文件传输完成")
}
