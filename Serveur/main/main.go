package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	active bool
	conn net.Conn
	name string
	mute bool
}

func main(){




	/*~~~~~~~~~~~~INITIALISATION DU SERVEUR~~~~~~~~~~~~~*
	*Dans cette première partie de code, on va:        	*
	*		>ouvrir la connection TCP sur le port 8081 	*
	*		>mettre en place le channel des reponses cli*
	*		ent										   	*
	*		>mettre en place la table des clients		*
	*		>lancer l'automate de réponse automatique	*
	*		>lancer la goroutine d'administration du ser*
	*		veur										*
	*On fait cela avant déccepter la moindre connection	*
	*client.											*
	*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

	fmt.Println("\n[---------------------Serveur Start---------------------]\n")

	ln, err := net.Listen("tcp",":8081") //initilisation du serveur en TCP sur le port 8081
	if err != nil {
		fmt.Println("Read error:", err)
		os.Exit(1)
	}

	buffer:=make(chan string,20) //creation du buffer qui va contenir les  messages des clients (20 max.)

	var clientArray [10]Client //creation du tableau contennant les clients connectés

	var clientPointers [10]*Client  //creation d'un tableau de pointeur vers les clients

	for i:=0 ; i<len(clientArray);i++{  //initialisation du tableau de pointeur
		clientPointers[i]=&clientArray[i]
	}

	go answer(buffer,clientPointers)  //lancement de la goroutine permettant de relancer les messages vers tous les clients connectés

	go adminServer(clientPointers, buffer)  //lancement de la goroutine permettant à l'admin de moderer le serveur




	/*~~~~~~~~~~FIN INITIALISATION DU SERVEUR~~~~~~~~~~~*
	*Le serveur est maintenant près à recevoir des conn	*
	*extions clients									*
	*		>mise en encoute du serveur pour voir si un *
	*client veut se connecter							*
	*Si un client se connecte:							*
	*	>regarde si un slot est disponible!				*
	*		>si un slot est disponible, alors on traite	*
	*		le client dans un goroutine à part			*
	*		>sinon, le client est rejeté				*
	*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

	for {
		conn, err := ln.Accept() //serveur en attente d'une connection client
		if err != nil {  //erreur lors de la connection du client
			fmt.Println("Read error:", err)
			//nop
			os.Exit(1)
		} else { //erreur le client se connecte sans defauts
			fmt.Println("[------------------Connection detected------------------]\n")
		}

		index := -1

		for i := 0; i < len(clientArray); i++ { //boucle verifiant si un slot est disponible
			if clientArray[i].active==false {
				index = i
			}
		}

		if index != -1 { //un slot est disponible au numero index

			fmt.Println("Connection available. Connection index : "+strconv.Itoa(index)) //on affiche du coté serveur les infos
			clientArray[index] = Client{true,conn,"Barbe Rousse le Maquereau", false} //on enregistre notre client dans le tableau de Client et on indique le slot comme étant occupé(active =true)
			personnalPointer := &clientArray[index] //on recupère un pointeur vers le client
			go handleConnection(conn, buffer, personnalPointer, clientPointers) //on s'occupe du client dans une goroutine

		} else {  //tous les slots du serveur sont occupés

			fmt.Println("No connection available...")
			if _,err := conn.Write([]byte("Le navire est plein a craquer!!! Nous ne pouvons pas vous accepter!")); err!=nil{
				fmt.Println("Read error:", err)
				os.Exit(1)
			}

			if err := conn.Close(); err !=nil {//fermeture de la connection avec le client: le serveur est full
				fmt.Println("Read error:", err)
				os.Exit(1)
			}
		}
	}
}




	/*~~~~~~~~~~~~~~~AUTOMATE DE REPONSE~~~~~~~~~~~~~~~~*
	*Le serveur recoit des messages de la part des clien*
	*ts connectés										*
	*Ces messages sont encoyer par les goroutines de tra*
	*itement des clients dans un channel				*
	*Les messages du channel sont traités dans cette met*
	*hode qui va redistribuer le message à tous les clie*
	*nts connectés										*
	*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func answer(buffer chan string, clientPointers [10]*Client ) {

	for {

		select {
			case message, ok := <-buffer: //on regarde s'il y a quelque chose en attente dans la buffer
			if ok { //si oui...

				formatedAnsewer := message+"\n" //on met en forme la reponse pour les clients

				for i:=0;i<len(clientPointers);i++{ //on parcour la liste des clients connecté à l'aide de pointeurs

					if clientPointers[i]!=nil {

						personnalPointer :=clientPointers[i] //on récupère le pointeur vers un client contenu dans le tableau

						client:=*personnalPointer //on recupère le client lié au pointeur

						if client.active == true && client.mute==false{  //on verifie s'il est toujours connecté

							fmt.Println(client)

							if _, err := client.conn.Write([]byte(formatedAnsewer)); err != nil { //on envoie le message au client
								fmt.Println("Read error:", err)
								fmt.Println("Hip")
							}
						}else{
							if(client.mute==true){
								fmt.Println("Message non envoyé à "+client.name+". Cause : mute serveur!")
							}
						}
					}
				}
				fmt.Println("\n")
			} else {
				fmt.Println("Channel closed!")
				os.Exit(1)
			}
		}
	}
}




	/*~~~~~~~~~~~~~ADMINISATRATION SERVEUR~~~~~~~~~~~~~~*
	*Cette partie permet à l'adminisatrateur d'effectuer diverse action de modération sur sont serveur. Ainsi, il peut :
		>mutter une personne durant une durée définie
		>déconnecter une personne du serveur
		>bannir une personne du serveur
	Pour cela il lui suiffit de taper quelques commandes dan son terminal:
		>TCCHAT_MUTE <nom utilisateur> <durée>
		>TCCHAT_KICK <nom utilsateur>
		>TCCHAT_BLACKLIST <nom utilisateur>*
	*~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~*/

func adminServer(clientPointers [10]*Client,buffer chan string)  {

	reader := bufio.NewReader(os.Stdin)  //on place un buffer sur l'invite de commande

	for {
		texte, _ := reader.ReadString('\n') //tout ce qui est écrit va dans la variable texte

		if texte != "\n" && texte != "" {

			texte := strings.TrimSuffix(texte, "\r\n")
			data := strings.Split(texte, " ") //on divise le texte pour pouvoir reconnaitre la commande
			name := ""

			for i := 1; i < len(data)-1; i++ { //on determine le nom qui peut être composé, on fait data.lenght-1 à cause de la commande mute qui a le temps en dernière position
				name += data[i] + " "
			}

			switch data[0] {

			case "TCCHAT_BLACKLIST":
				name += data[len(data)-1]                //on récupère la dernière partie du pseudo: il n'y a pas de temps dans la commande
				writeIntoFile(name+";", "blacklist.txt") //on écrit le nom dans la blackListe
				fmt.Println("La personne " + name + " à bien été blacklisté du serveur")
				goKick(name,clientPointers)
				buffer <- "TCCHAT_USEROUT\t"+name+" a été pendu au mat du navire!" //message pour les autres clients

			case "TCCHAT_KICK":
				name += data[len(data)-1] //on récupère la dernière partie du pseudo: il n'y a pas de temps dans la commande
				fmt.Println(name)
				if (goKick(name, clientPointers)) {
					fmt.Println(name + " a bien été kické du serveur!")
					buffer <- "TCCHAT_USEROUT\t"+name+" a été jeter par dessus bord!" //message pour les autres clients
				} else {
					fmt.Println("Impossible de kicker " + name + " : soit la personne n'existe pas soit elle est deconnecté du serveur")
				}
			case "TCCHAT_MUTE":
				time, err := strconv.Atoi(data[len(data)-1]) //la dernière chaine du string est un temps normalement
				if err != nil { // l'admin n'a pas rentrer de temps ou la valeur n'est pas correcte !!!
					fmt.Println("Read error:", err)
				} else {
					name=strings.TrimSuffix(name," ")
					fmt.Println("nom : "+name+"|")

					goMute(name, time, clientPointers,buffer)

				}

			}
		}
	}
}

func goMute(name string, temps int, clientPointers [10]*Client, buffer chan string ){
	for i:=0;i<len(clientPointers);i++ {
		if clientPointers[i]!=nil {
			personnalPointer :=clientPointers[i]
			client:=*personnalPointer
			if client.active == true {
				if client.name==name {
					buffer <- "TCCHAT_USEROUT\t"+name+" a trop bu pour pouvoir parler!" //message pour les autres clients
					fmt.Println(name+" a bien été mute")
					*personnalPointer = Client{true, client.conn,client.name, true}
					time.Sleep(time.Duration(temps)*time.Second)
					*personnalPointer = Client{true, client.conn,client.name, false}
					fmt.Println(name+" a bien été démute")
					return
				}
			}
		}
	}
	fmt.Println("Impossible de mute "+name+" : soit la personne n'existe pas soit elle est deconnecté du serveur")
	return
}

func goKick(name string, clientPointers [10]*Client) bool{
	for i:=0;i<len(clientPointers);i++ {
		if clientPointers[i]!=nil {
			personnalPointer :=clientPointers[i]
			client:=*personnalPointer
			if client.active == true {
				if client.name==name {
					fmt.Println(client)
					if _, err := client.conn.Write([]byte("TCCHAT_PERSO\tVous avez subit le suplice de la plance et avez quitté le navire!")); err != nil { //on envoie le message au client
						fmt.Println("Read error:", err)

					}
					if err := client.conn.Close(); err !=nil {
						fmt.Println("Read error:", err)
						os.Exit(1)
					}
					*personnalPointer =  Client{false, nil,"", false}//on reinitiallise le slot du client patit pour qu'il accèpte une nouvelle connection
					return true
				}
			}
		}
	}
	return false
}

func handleConnection (conn net.Conn, buffer chan string, personnalPointer *Client , clientPointers [10]*Client){

	if _, err := conn.Write([]byte("TCCHAT_WELCOME\tLE_CHALUTIER_DE_L'ENFER\n")); err != nil {  //on salue le client qui vient de se connecter
		fmt.Println("Read error:", err)
	}

	whiteList := readFile("whiteList.txt") //on lit le fichier whiteList en recuperant la liste des personnes de la white liste
	blackList := readFile("blackList.txt") //on lit le fichier blackList en recuperant la liste des personnes de la black liste

	blackListed := false  // on suppose que le client n'est pas blackListé

	name:=""

	myReader:=bufio.NewReader(conn) //on crée un nouveau buffer à l'écoute de la connection Client

	for{

		message,err := myReader.ReadString('\n') //on lit les messages provenants du client
		if err != nil {
			fmt.Println("Read error:", err)
		}

		//le premier message client contient forcement le nom du client
		name = messageCleaning(message)  //on traite le message selon si on est sur UNIX ou Windows sinon on à des problèmes d'affichage

		if name!="" {

			if identification(name, blackList) { //on verifie si le client est blacklisté

				if _, err := conn.Write([]byte("TCCHAT_PERSO\tMais oui c'est ca! Tu es Barbe Rousse le Maquereau! Au secours! A l'aide\n")); err != nil {  //on "salurt" le mechant client
					fmt.Println("Read error:", err)
				}
				blackListed = true  //le client ne peut pas entrée sur le chat
				break

			}else{ //le client n'est pas blacklisté
				if !alreadyConnected(name,clientPointers) { //on regarde si un client avec le même pseudo n'est pas déja connecté

					if identification(name, whiteList) { //on regarde si le client s'est déja connecté dans le passé

						answer := "TCCHAT_USERIN\tUn vieux loup de mer vient de se connecter : " + name //le client s'est deja connécté dans la passé
						buffer <- answer
						if _, err := conn.Write([]byte("TCCHAT_PERSO\tEt bien, tu en as mis du temps pour decuver, vieux loup de mer!!!\n")); err != nil {  //on salue le client gentil ancien client qui vient de se connecter
							fmt.Println("Read error:", err)
						}
						break

					} else {

						answer := "TCCHAT_USERIN\tUn nouveau moussaillon vient de se connecter : " + name
						buffer <- answer
						if _, err := conn.Write([]byte("TCCHAT_PERSO\tBienvenu à bord, marin d'eau douce!!!\n")); err != nil {  //on salue le nouveau client qui vient de se connecter
							fmt.Println("Read error:", err)
						}
						writeIntoFile(name+";","whiteList.txt" ) //on ecrit le nom du nouveau client dans la WhiteListe pour se souvenir de lui lors de sa prochaine connection
						break
					}

				}else{  //un client avec le même pseudo est déja connecté!!

					fmt.Println("Le client est deja connecté ou le nom existe deja")
					if _, err := conn.Write([]byte("TCCHAT_PERSO\tJe t'ai reconnu Barbe Rousse le Marquereau!!! "+name+" est deja sur le navire!!!\n")); err != nil {  //on salue le nouveau client qui vient de se connecter
						fmt.Println("Read error:", err)
					}
					blackListed=true //on expulse le doublon du serveur!
					break
					
				}
			}

		}
	}

	if !blackListed { //si la personne n'est pas blacklisté ou n'est pas déja connecté...

		*personnalPointer =  Client{true, conn,name, false}  //on met à jour le profit du client

		client := *personnalPointer

		exit := false

		myReader := bufio.NewReader(conn)

		for exit == false || client.active { //tant que le client ne veut pas quitter le serveur, on l'écoute

			message, err := myReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Println("Read error:", err)
					break
				}
			}

			message = messageCleaning(message)  //on néttoye le message pour eviter les problèmes d'affichage

			if message != "" {

				switch message {
					case "TCCHAT_EXIT": //le client veut se déconnecter

						fmt.Println("Demande de DECONNECTION recu de ["+strings.ToUpper(name)+"]")
						exit= true //on passe exit à true pour quitter la boucle for exit{}
						break

					case "TCCHAT_INFO"://le client veut savoir qui est connecter en ce moment sur le serveur

						fmt.Println("Demande d'INFO recu de ["+strings.ToUpper(name)+"]")

						info:=getConnectedUser(clientPointers) //on recupère les personnes connectées sur le serveur
						fmt.Println(info)
						if _,err := conn.Write([]byte(info)); err!=nil{ //on envoie les infos au client
							fmt.Println("Read error:", err)

						}
						break

					case "TCCHAT_BLACKLIST": //le client veut savoir qui est blacklisté

						fmt.Println("Demande de BLACKLIST recu de ["+strings.ToUpper(name)+"]")

						info:="TCCHAT_PERSO\tVoici les membres de l'équipage de Barbe Rousse le Maquereau : "
						for i := 0; i<len(blackList);i++ {//boucle for pour récupérer les valeurs du tableau déterminé en entrée
							info+=blackList[i]+" | "
						}
						info+="\n"
						if _,err := conn.Write([]byte(info)); err!=nil{//on envoie les infos au client
							fmt.Println("Read error:", err)
							os.Exit(1)
						}

						break

					case "TCCHAT_WHITELIST"://le client veut savoir qui est dans la whitelist

						fmt.Println("Demande de WHITELIST recu de ["+strings.ToUpper(name)+"]")

						info:="TCCHAT_PERSO\tVoici les membres d'équipage du CHALUTIER DE L'ENFER : "
						for i := 0; i<len(whiteList);i++ {//boucle for pour récupérer les valeurs du tableau déterminé en entrée
							info+=whiteList[i]+" |"
						}
						info+="\n"
						if _,err := conn.Write([]byte(info)); err!=nil{//on envoie les infos au client
							fmt.Println("Read error:", err)
							os.Exit(1)
						}
						break

					default: //par defaut c'est un message basique, déstiné aux autres clients

					client:=*personnalPointer

						if client.mute==false{
							formatedMessage := "TCCHAT_BCAST\t["+name+"]: "+message
							fmt.Println("Message recu de ["+strings.ToUpper(name)+"]: \""+message+"\"")
							buffer <-formatedMessage //on envoie le message dans le channel pour qu'il soit envoyé à tous les clients
						}
				}
			}
		}
	}


	if err := conn.Close(); err !=nil {//fermeture de la connection avec le client
		fmt.Println("Read error:", err)
	}

	*personnalPointer =  Client{false, conn,"", false}//on reinitiallise le slot du client patit pour qu'il accèpte une nouvelle connection
	buffer <- "TCCHAT_USEROUT\t"+name+" a changé de pavillon" //message pour les autres clients

}

//Methode qui permet de verifier si un nom appartient à un tableau de string
func identification(name string, data [] string) bool {
	if data != nil {
		for i := 0; i < len(data); i++ {//on parcoure le tableau
			if data[i] == name { //une valeur du tableau correspond au nom
				return true //on retourne que le nom est contenue dans la tableau
			}
		}
	}
	return false //on retourne que le nom n'est pas contenue dans la tableau
}

//methode permettant de récuperer les valeurs d'un fichier .txt séparés par des virgule
//retourne un tableau de string contenant les valeurs demandées
func readFile(filePath string) [] string{

	var _, err = os.Stat(filePath)
	if os.IsNotExist(err) {//on regard si le fichier existe ou pas
		var file, err = os.Create(filePath) //s'il n'existe pas on le créée
		if err!=nil {
			fmt.Println(err)
			data :=[]string{""}
			return data
		}
		defer file.Close() //fermture automatique du fichier au return de la méthode
		fmt.Println("File Created Successfully", filePath)
	}

	reader, err := ioutil.ReadFile(filePath) // on met un reader pour lire le fichier
	if err != nil {
		fmt.Print(err)
	}

	str := string(reader)

	if strings.Contains(str,";"){ //on separe les noms contenues dans le fichier grace aux ";"
		data := strings.Split(str,";")
		return data
	}

	return nil
}

//methode permttant d'écrire un string dans un fichier passer en paramètre
func writeIntoFile(str string, fileName string) {

	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()


	if _, err = f.WriteString(str); err != nil {
		panic(err)
	}
}

//methode permttant de mettre en forme un message provenant d'un client windows (apparition d'un /r abscent en UNIX)
func messageCleaning(message string) string{
	tab := strings.Split(message, "\t")
	newMessage:=""
	if runtime.GOOS == "windows" {
		newMessage = strings.TrimRight(tab[1], "\r\n")
	} else {
		newMessage = strings.TrimRight(tab[1], "\n")
	}
	return newMessage
}


func alreadyConnected(name string, clientPointers [10]*Client) bool{

	for i:=0;i<len(clientPointers);i++ {
		if clientPointers[i]!=nil {
			personnalPointer :=clientPointers[i]
			client:=*personnalPointer
			if client.active == true {
				if client.name==name {
					return true
				}
			}
		}
	}
	return false
}

func getConnectedUser(clientPointers [10]*Client)string{
	message :="TCCHAT_PERSO\tLes personnes actuellement connectées sur le serveur sont : "
	for i:=0;i<len(clientPointers);i++ {
		if clientPointers[i]!=nil {
			personnalPointer :=clientPointers[i]
			client:=*personnalPointer
			if client.active == true {
				message+= client.name+" | "
			}
		}
	}
	message+="\n"
	return message
}

