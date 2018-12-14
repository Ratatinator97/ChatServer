package main

import "net"
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
	connecte := false
	for connecte != true {
		for EntreeServeur(message) != true {
			//delay
		}
		fmt.Print("Recu msg serveur")
		fmt.Fprintf(connexion, "TCCHAT_REGISTER"+"\t"+EntreeMsg(1)+"\n")
		fmt.Print(connexion, "TCCHAT_REGISTER"+"notre pseudo")
		connecte = true //Il faudra handle les erreurs ici
	}
	fmt.Print("Etape connexion terminee")
	for {
		//Lecture msg serveur
		MsgType := LectureMsgServeur(message)
		TabMsg := traitermsg(message)
		if MsgType == 0 {
			fmt.Print(TabMsg[1], " :", TabMsg[2])
		} else if MsgType == 1 {

		}

		//Ecriture msg serveur
		EntreeMsg(2)
		//On lit notre buffer d entree jusqu a ENTER

		fmt.Print("Serveur: " + message) //On affiche le msg recu
	}
}
func traitermsg(msg string) []string {
	tab_msg := strings.Split(msg, "\t")
	return tab_msg
}

func EntreeServeur(S string) bool {
	TabS := strings.Split(S, "\t")
	if TabS[0] == "TCCHAT_WELCOME" {
		fmt.Print("On a recu le welcome")
		return true
	} else {
		return false
	}
}
func LectureMsg(Str string) []string {
	TabStr := strings.Split(Str, "\t")
	return TabStr
}
func LectureMsgServeur(Msg string) int {
	TypeMsg := LectureMsg(Msg)
	if TypeMsg[0] == "TCCHAT_BCAST" {
		fmt.Print("msg_broadcast")
		return 0
	} else if TypeMsg[0] == "TCCHAT_USEROUT" {
		fmt.Print("msg_Userout")
		return 1
	} else if TypeMsg[0] == "TCCHAT_USERIN" {
		fmt.Print("msg UserIn")
		return 2
	} else {
		fmt.Print("Bizzare")
		return 666
	}
}
func EntreeMsg(MsgType int) string {
	valide := false
	reader := bufio.NewReader(os.Stdin)
	if MsgType == 1 {
		//Nickname msg
		fmt.Print("Qui etes vous ?: ")
		texte, _ := reader.ReadString('\n')
		for valide == false {
			if texte != "\n" {
				valide = true
			}
		}
		fmt.Print("Votre nom est " + texte)
		NvTexte := strings.TrimSuffix(texte, "\n")
		fmt.Print("C est mieux comme ca " + NvTexte)
		return NvTexte
	} else if MsgType == 2 {
		texte, _ := reader.ReadString('\n')
		if texte != "\n" {
			texte := strings.TrimSuffix(texte, "\n")
			fmt.Print("Envoi de message" + texte)
			return texte
		}
	}
	return ""
}
