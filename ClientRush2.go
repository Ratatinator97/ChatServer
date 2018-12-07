package main

import "net"
import "fmt"
import "bufio"
import "os"
import "strings"

func main() {
	//On se connecte
	connexion, _ := net.Dial("tcp","127.0.0.1:8081") //LocalHost
	//On se connecte au serveur
	//Penser a handle les erreurs
	for ConnectionProcess() == true {
		//Mettre un delai
	}
	for {
		//Lecture msg serveur
		LectureMsgServeur()
		//Ecriture msg serveur
		EntreeMsg(2)
	//On lit notre buffer d entree jusqu a ENTER
	message, _ := bufio.NewReader(connexion).ReadString('\n')
	fmt.Print("Serveur: "+message)//On affiche le msg recu
	
}
func traitermsg(msg string ) {
	tab_msg = strings.Split(msg,"\t")

func ConnecionProcess() string {
	fmt.Fprintf(connexion,"TCCHAT_REGISTER"+'\t'+ EntreeMsg(1) + "\n")
	return true //Il faudra handle les erreurs ici

func EntreeMsg(MsgType int) string {
	bool valide = false
	reader := bufio.NewReader(os.Stdin)
	if MsgType == 1 {
		//Nickname msg
		fmt.Print("Qui etes vous ?: ")
		for valide == false {
		texte, _ := reader.ReadString('\n')
		if texte := '\n' {
			valide = true
		}
		texte := strings.TrimSuffix(texte,'\n')
	}
}


