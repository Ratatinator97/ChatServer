package main

import "net"
import "fmt"
import "bufio"
import "strings"


func main() {

	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, _ := net.Listen("tcp", ":8081")

	// accept connection on port
	conn, _ := ln.Accept()
	myReader := bufio.NewReader(conn)
	// run loop forever (or until ctrl-c)
	user := ""
	for {
		// will listen for message to process ending in newline (\n)
		message, _ := myReader.ReadString('\n')
		message = strings.TrimSuffix(message, "\n")
		// output message received
		fmt.Print("Message Received:", string(message), "\n")
		// sample process for string received
		traiteString(message, user)
		newmessage := strings.ToUpper(message)
		// send new string back to client
		conn.Write([]byte(newmessage + "\n"))
	}
}

func traiteString(s string,user string) {
	TableauString := strings.Split(s, "\t")
	fmt.Println(TableauString)
	if TableauString[0] == "TCCHAT_REGISTER" {
		if len(TableauString) == 2{
			user = strings.TrimSuffix(TableauString[1], "\n")
			fmt.Println(user)
		} else {
			fmt.Println("nom invalide")
		}
	} else if TableauString[0] == "TCCHAT_MESSAGE" {
		if len(TableauString) == 2{
			fmt.Println(user)
			if(len(user)>0){
				fmt.Println(user + TableauString[1])
			} else {
				fmt.Println("veuillez vous inscrire")
			}
		} else {
			fmt.Println("pas de tabulations dans le message")
		}
	} else if TableauString[0] == "TCCHAT_DISCONNECT" {
		fmt.Print("disconnect \n")
	} else {
		fmt.Print("yolo")
	}
	return
}
