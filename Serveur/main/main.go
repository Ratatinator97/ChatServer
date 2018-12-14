package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
}

func main(){

	fmt.Println("Lancement du serveur...")

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Read error:", err)
		os.Exit(1)
	}

	buffer:=make(chan string,20)

	var tabClient [10]Client

	var tabConn [10]bool

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Read error:", err)
		}

		index := -1

		for i:=0;i<len(tabConn);i++{
			if tabConn[i] {
				index=i
			}
		}

		if(index!=-1){
			go handleConnection(conn, buffer, tabConn, index)
			tabClient[index]= Client{conn}
			tabConn[index]=false
		}else{
			go overSize(conn)
		}



	}
}

func overSize(conn net.Conn){

	conn.Write([]byte("Le navire est plein a craquer!!! Nous ne pouvons pas vous accepter!"))

	conn.Close()
}

func answer(buffer chan string,) {

	for {

		select {
			case answer, ok := <-buffer:
			if ok {
				fmt.Printf("Value %d was read.\n", answer)

				for i:=0;i<len(tabClient);i++{

					tabClient[i].conn.Write([]byte(answer))
				}

			} else {
				fmt.Println("Channel closed!")
			}
		default:
			fmt.Println("No value ready")
		}

	}
}

func handleConnection (conn net.Conn, buffer chan string, tabConn [10]bool, index int){

	conn.Write([]byte("Bienvenu sur le serveur Chalutier Moussaillon! Quel est ton nom matelot?")) // send to the server/client

	blackListed := false

	whiteList := readFile("whiteList.txt")

	blackList := readFile("blackList.txt")

	name:=""


	for{

		buf := make([]byte, 0, 4096) // big buffer

		_, err := conn.Read(buf)

		if err != nil {
			if err != io.EOF {
				fmt.Println("Read error:", err)
			}
			break
		}

		name =string(buf)

		fmt.Println("Ah "+name)

		if(name!=""){

			if(identification(name, blackList)){
				conn.Write([]byte("Un pirate nous attaque!!!"))
				blackListed = true
				break
			}else{
				if(identification(name, whiteList)){
					conn.Write([]byte("Et bien, tu en as mis du temps pour decuver, vieux loup de mer!!!"))
					buffer<-name+"viens de se connecter"
					break
				}else{
					conn.Write([]byte("Bienvenu à bord, marin d'eau douce!!!"))
					ioutil.WriteFile("whiteList.txt",[]byte(name+";"),0644)
					buffer<-"Bienvenu à "+name+", un nouveau moussaillon"
					break
				}
			}

		}
	}

	if(!blackListed){

		exit :=false

		for exit==false {

			buf := make([]byte, 0, 4096) // big buffer

			_, err := conn.Read(buf)

			if err != nil {
				if err != io.EOF {
					fmt.Println("Read error:", err)
				}
				break
			}

			blabla:=string(buf)

			fmt.Println(blabla)

			if(blabla!=""){
				switch blabla {
				case "exit":
					exit=true
					break
				default:
					buffer<-"["+strings.ToUpper(name)+"]"+blabla
				}
			}



		}
	}

	conn.Write([]byte("Bon vent mon gas!!!"))
	tabConn[index]=true
	conn.Close();


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

	fmt.Println(reader) // print the content as 'bytes'

	str := string(reader)

	data:= strings.Split(str,";")

	return data
}

func clock(){
	for{

		time.Sleep(500)
	}
}
