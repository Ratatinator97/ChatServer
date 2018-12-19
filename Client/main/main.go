package main

import (
	"io"
	"net"
)
import "fmt"
import "bufio"
import "os"
import "strings"

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

	ecritureMsgServeur(1, connexion)

	fmt.Println("Etape connexion terminee")

	go read(connexion)
	go write(connexion)

	exit:=false
	for exit==false{
		exit=false
	}

}


func read(conn net.Conn){

	for {

		tmp := make([]byte, 256)
		_, err := conn.Read(tmp)

		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
		}

		message := string(tmp)

		fmt.Println(message)

	}
}



func write(conn net.Conn){
	for{
		ecritureMsgServeur(2, conn)
	}

}

func ecritureMsgServeur(msgType int,conn net.Conn)  {

	reader := bufio.NewReader(os.Stdin)

	switch msgType {

	case 1:
		fmt.Print("Qui etes vous ?: ")
		texte, _ := reader.ReadString('\n')
		for  {
			if texte != "\n" {
				break
			}
		}
		fmt.Println("Votre nom est " + texte)
		NvTexte := strings.TrimSuffix(texte, "\n")

		fmt.Fprintf(conn, "TCCHAT_REGISTER"+"\t"+NvTexte+"\n")


	case 2:
		texte, _ := reader.ReadString('\n')
		if texte != "\n" && texte!="" {
			texte := strings.TrimSuffix(texte, "\n")
			//fmt.Print("Envoi de message" + texte)
			fmt.Fprintf(conn, "Prout\t"+texte+"\n")
			fmt.Println(texte)

		}
	}


}