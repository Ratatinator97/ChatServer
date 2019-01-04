package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	actif bool
	conn net.Conn
}

func main(){

	fmt.Println("Lancement du serveur...")


		ln, err := net.Listen("tcp",":8081")
		if err != nil {
			fmt.Println("Read error:", err)
			os.Exit(1)
		}

		buffer:=make(chan string,20)

		var tabClient [10]Client

		var tabConn [10]bool

		for i:=0;i<len(tabConn);i++{
			tabConn[i]=true
		}

		var pc [10]*Client

		for i:=0 ; i<len(tabClient);i++{
			pc[i]=&tabClient[i]
		}

		go answer(buffer,pc)

		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Read error:", err)
				os.Exit(1)
			} else {
				fmt.Println("Connection detected")
			}

			index := -1

			for i := 0; i < len(tabConn); i++ {
				if tabConn[i]==true {
					index = i
				}
			}

			if (index != -1) {
				fmt.Println("Connection available. Index : "+strconv.Itoa(index))
				tabConn[index] = false
				pC := &tabConn[index]
				tabClient[index] = Client{true,conn}
				go handleConnection(conn, buffer, pC, index)

			} else {
				fmt.Println("No connection available...")
				go overSize(conn)
			}
		}
}

func overSize(conn net.Conn){

	conn.Write([]byte("Le navire est plein a craquer!!! Nous ne pouvons pas vous accepter!"))
	conn.Close()
}

func answer(buffer chan string, pointer [10]*Client ) {

	for {

		select {
			case answer, ok := <-buffer:
			if ok {

				formatAnsewer := join("> " ,answer,"\n")

				for i:=0;i<len(pointer);i++{

					if pointer[i]!=nil {

						pC :=pointer[i]

						client:=*pC


						if client.actif == true {
							if _, err := client.conn.Write([]byte(formatAnsewer)); err != nil {
								fmt.Println("Read error:", err)
							}else{
								fmt.Println("Message envoyé")
							}
						}
					}
				}

			} else {
				fmt.Println("Channel closed!")
			}
		default:

		}

	}
}

func handleConnection (conn net.Conn, buffer chan string, pC *bool , index int){

	conn.Write([]byte("TCCHAT_WELCOME\tLE_CHALUTIER_DE_L'ENFER\n"))

	whiteList := readFile("whiteList.txt")

	blackList := readFile("blackList.txt")

	blackListed := false

	name:=""

	myReader:=bufio.NewReader(conn)

	for{

		message,err := myReader.ReadString('\n')

		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
				os.Exit(1)
			}
		}


		data:=strings.Split(message,"\t")

		name=strings.TrimSuffix(data[1],"\n")

		compteur:=strings.Count(name, "\n")

		reponse:="Le nom du Client est: "+name

		fmt.Println(reponse)

		fmt.Println("Le nombre de putain d'espace est "+strconv.Itoa(compteur))

		if(name!=""){

			if(identification(name, blackList)){
				conn.Write([]byte("TCCHAT_USERIN\tMais oui c'est ca! Tu es Barbe Rousse le Maquereau! Au secours! A l'aide\n"))
				blackListed = true
				break
			}else{
				if(identification(name, whiteList)){
					conn.Write([]byte("TCCHAT_USERIN\tEt bien, tu en as mis du temps pour decuver, vieux loup de mer!!!\n"))
					answer := "Un vieux loup de mer vient de se connecter : "+name
					fmt.Println(answer)
					buffer<- answer
					break
				}else{
					conn.Write([]byte("TCCHAT_USERIN\tBienvenu à bord, marin d'eau douce!!!\n"))
					ioutil.WriteFile("whiteList.txt",[]byte(name+";"),0644)
					answer := "Bienvenu à "+name+", un nouveau moussaillon"
					fmt.Println(answer)
					buffer<- answer
					break
				}
			}

		}
	}

	if(!blackListed) {

		fmt.Println("Je rentre")
		exit := false

		myReader := bufio.NewReader(conn)

		for exit == false {

			message, err := myReader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					fmt.Println("Read error:", err)
					fmt.Println("Pouet!")
					break
				}
			}

			message = strings.TrimSuffix(message, "\n")

			fmt.Println(message)

			if (message != "") {

			if strings.EqualFold(message,"exit") {

					exit = true

				}else {

					/*in:="["
					upper:= strings.ToUpper(name)
					out :="] ->"

					message+=in+upper+out*/

					buffer <-message
				}
			}
		}
	}


	conn.Write([]byte("Bon vent mon gas!!!"))
	*pC=true
	conn.Close()

}

func identification(name string, data [] string) bool{

	for i:=0 ; i< len(data); i++{
		if(data[i] == name){
			return true
		}
	}

	return false
}


func readFile(filePath string) [] string{

	var _, err = os.Stat(filePath)

	if os.IsNotExist(err) {
		var file, err = os.Create(filePath)
		if err!=nil {
			fmt.Println(err)
			data :=[]string{""}
			return data
		}
		defer file.Close()
		fmt.Println("File Created Successfully", filePath)
	}

	reader, err := ioutil.ReadFile(filePath) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}

	str := string(reader)

	data:= strings.Split(str,";")

	return data
}

func join(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}

