package main

import (
"io"
"net"
)
import "fmt"

func main() {

	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")

	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
		}
		//fmt.Println("got", n, "bytes.")
		buf = append(buf, tmp[:n]...)

		message := string(tmp)
		fmt.Println(message)

	}
	fmt.Println("total size:", len(buf))

}


