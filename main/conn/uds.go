package conn

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const sockAddr = "/tmp/netbreaker-uds.sock"

func echoServer(c net.Conn) {
	io.Copy(c, c)
	c.Close()
}

func ListenUDS() {
	fmt.Printf("hi\n")
	if err := os.RemoveAll(sockAddr); err != nil {
		fmt.Print(err)
	}
	l, err := net.Listen("unix", sockAddr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()
	for {
		// Accept new connections, dispatching them to echoServer
		// in a goroutine.
		conn, err := l.Accept()
		if err != nil {
			log.Fatal("accept error:", err)
		}
		go echoServer(conn)
	}
}
