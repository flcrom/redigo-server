package main

import (
	"fmt"
	"net"
	"os"
	"redigo/cmd"
	"redigo/internal/server"
	"redigo/store"
)

func main() {

	fmt.Println("Listening on port 6379!")

	//listener listens for clients
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	consumeListener(listener)

}

func consumeListener(listener net.Listener) {
	//initialize store and the handler
	store := store.NewStore()
	handler := cmd.NewHandler(store)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		fmt.Println("User connected")
		go server.HandleIO(conn, handler)
	}
}

// there are 2 kinds of errors that are propogating back to HandleIO
// one of them is coming from the envelope, the other one comes from the system
// when there is an error of the system, i need to recognise it other wise
