package main

//Importation des bibliotheques requises
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
	//Connexion au serveur
	connexion, _ := net.Dial("tcp", "192.168.43.31:8081") //LocalHost
	//Buffer contenant la réponse du serveur
	message, _ := bufio.NewReader(connexion).ReadString('\n')

	// Attente jusqu'au TCCHAT_WELCOME
	for {
		TabS := strings.Split(message, "\t")
		if TabS[0] == "TCCHAT_WELCOME" {
			fmt.Println(message)
			break
		}
	}
	// L'utilisateur rentre son nom
	nameOfUser := messageCleaning(ecritureMsgServeur(1, connexion))

	// Lancement de go routines pour lire les messages du serveur et du clavier utilisateur
	go read(connexion, nameOfUser)
	go ecritureMsgServeur(2, connexion)

	// Le main ne se termine pas tant que l'utilisateur n'a pas mis fin au programme
	exit := false
	for exit == false {
		exit = false
	}

}

func read(conn net.Conn, nameOfUser string) {

	for {
		// Buffer msg serveur
		message1, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil { // Gestion des erreurs
			if err != io.EOF {
				fmt.Println("Read error:", err)
				os.Exit(1)
			}
		}
		if message1 != "" {
			tabS := strings.Split(message1, "\t")
			switch tabS[0] {

			case "TCCHAT_BCAST":
				identifiant := strings.Split(tabS[1], ":")
				inputName := "[" + nameOfUser + "]"

				// Si le msg entrant correspond au nom de l'utilisateur
				// Alors on n'affiche pas le msg (eviter les doublons)
				if inputName != identifiant[0] {
					fmt.Println(messageCleaning(tabS[1]))
				}
			case "TCCHAT_USERIN":
				fmt.Println("\n" + tabS[1])
			case "TCCHAT_USEROUT":
				fmt.Println("\n" + tabS[1])
			case "TCCHAT_PERSO":
				fmt.Println("\n" + tabS[1])
			default:
				fmt.Println("Unexpected type of msg")
			}
		}
	}
}

func ecritureMsgServeur(msgType int, conn net.Conn) string {

	reader := bufio.NewReader(os.Stdin)

	switch msgType {

	case 1:
		// Demander le nom a l'utilisateur
		fmt.Print("Identifiant : ")
		texte, _ := reader.ReadString('\n')

		// Lire tant que il n'y a pas d'appui sur ENTER
		for {
			if texte != "\n" {
				break
			}
		}

		// Epurer le nom
		name := strings.TrimSuffix(texte, "\r\n")

		// Envoi du msg de connexion au serveur
		if _, err := conn.Write([]byte("TCCHAT_REGISTER" + "\t" + name + "\n")); err != nil {
			fmt.Println("Read error : ")
		}

		return name

	case 2:
		// tant que l'utilisateur n'a pas tape exit
		exit := false
		for exit == false {
			texte, _ := reader.ReadString('\n')
			if texte != "\n" && texte != "" { // Verification msg pas vide
				texte := messageCleaning(texte)

				if texte == "exit" {
					exit = true // l'utilisateur tappe exit
				}

				// Si il n'y a pas d'erreur on demande au serveur de nous deconnecter
				if _, err := conn.Write([]byte("TCCHAT_MESSAGE\t" + texte + "\n")); err != nil {

					// Identique a un CTRL+C
					os.Exit(-1)
				}
			}
		}
		fmt.Println("Vous vous êtes deconnecté du serveur!")
	}
	return ""
}

// Fonction d'epuration
func messageCleaning(message string) string {
	newMessage := ""
	if runtime.GOOS == "windows" {
		newMessage = strings.TrimRight(message, "\r\n")
	} else {
		newMessage = strings.TrimRight(message, "\n")
	}
	return newMessage
}
