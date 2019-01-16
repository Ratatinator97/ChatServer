package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
)

func main() {
	//On se connecte
	connexion, _ := net.Dial("tcp", "127.0.0.1:8081") //LocalHost
	//On se connecte au serveur
	//Penser a handle les erreurs
	message, _ := bufio.NewReader(connexion).ReadString('\n')

	for {
		TabS := strings.Split(message, "\t")
		if TabS[0] == "TCCHAT_WELCOME" {
			fmt.Println(message)
			break
		}
	}

	nameOfUser := messageCleaning(ecritureMsgServeur(1, connexion))


	go read(connexion, nameOfUser)
	go ecritureMsgServeur(2,connexion)

	exit:=false
	for exit==false{
		exit=false
	}

}

func read(conn net.Conn, nameOfUser string) {

	for {
		message1, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
				os.Exit(1)
			}
		}
		if message1 != "" {
			tabS := strings.Split(message1, "\t")
			switch tabS[0] {
			case "TCCHAT_BCAST":
				identifiant:= strings.Split(tabS[1], ":")
				inputName := "[" + nameOfUser + "]"
				if inputName == identifiant[0] {
				} else {
					fmt.Println(tabS[1])
				}
			case "TCCHAT_USERIN":
				fmt.Println(tabS[1])
			case "TCCHAT_USEROUT":
				fmt.Println(tabS[1])
			case "TCCHAT_PERSO":
				fmt.Println(tabS[1])
			default:
				fmt.Println("Unexpected type of msg")
			}
		}
	}
}


func ecritureMsgServeur(msgType int, conn net.Conn) (string) {

	reader := bufio.NewReader(os.Stdin)

	switch msgType {

	case 1:
		fmt.Print("Identifiant : ")
		texte, _ := reader.ReadString('\n')
		for {
			if texte != "\n" {
				break
			}
		}

		name := strings.TrimSuffix(texte, "\r\n")

		if _,err := conn.Write([]byte("TCCHAT_REGISTER"+"\t"+name+"\n")); err!=nil{
			fmt.Println("Read error : ")
		}

		return name

	case 2:
		exit:=false
		for exit==false{
			texte, _ := reader.ReadString('\n')
			if texte != "\n" && texte != "" {
				texte := messageCleaning(texte)

				if(texte=="exit"){
					exit=true
				}
				//fmt.Print("Envoi de message" + texte)
				if _,err := conn.Write([]byte("TCCHAT_MESSAGE\t"+texte+"\n")); err!=nil{//A le reception du serveur corriger ca
					os.Exit(-1)
				}
			}
		}
		fmt.Println("Vous vous êtes deconnecté du serveur!")
		os.Exit(-1)
	}
	return ""
}

func messageCleaning(message string) string{
	newMessage:=""
	if runtime.GOOS == "windows" {
		newMessage = strings.TrimRight(message, "\r\n")
	} else {
		newMessage = strings.TrimRight(message, "\n")
	}
	return newMessage
}