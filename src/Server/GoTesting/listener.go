package main

//todo
//Listener -> Server communication. Skapa en loop som lyssnar efter input och output från klienten till listener. Send/Receive, Accept connection,
//Listener ID, Port
//Tog emot en port och spawnade en connection

import(
	"net"
	"fmt"
	"os"
	"bufio"
)



func createListener(port string) net.Listener{
	//Bör inte ignorera err m.h.a _

	fmt.Println("Creating listener...")
	connection, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return connection
}

func main() {
	connection := createListener(":9001")
	//reader := bufio.NewReader(os.Stdin)
	var conn net.Conn
	for {
		conn, _ := connection.Accept()
		
		fmt.Println("Listening...")

		message, err := bufio.NewReader(conn).ReadString('\n')

		if err != nil {
			conn.Close()
		}

		fmt.Print("Received: ", string(message))

		bytes := make([]byte, 1024)
		bytes = []byte("Succesfully sent")
		conn.Write(bytes)
		
		fmt.Println()
	}
	
	defer conn.Close()
}