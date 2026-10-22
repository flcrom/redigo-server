package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"redigo/cmd"
)

func HandleIO(conn net.Conn, handler *cmd.Handler) {
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		res, err := cmd.HandleRequest(handler, reader)
		if err != nil {
			// system error occurred
			// "stop the connection" kinda errors occured
			if err == io.EOF {
				fmt.Printf("User disconnected EOF error: %s \n", err)
				return
			}

			if err == cmd.NCP {
				return
			}

			fmt.Printf("(server.HandleIO) Error in handling request: %s \n", err)
			return
		}
		_, err = conn.Write([]byte(res))

		if err != nil {
			// system error occurred
			// "stop the connection" kinda errors occured
			fmt.Printf("Err in writing to connection: %s \n", err)
			return
		}
	}
}
