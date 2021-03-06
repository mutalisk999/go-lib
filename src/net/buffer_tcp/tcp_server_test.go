package buffer_tcp

import (
	"fmt"
	. "github.com/mutalisk999/go-lib/src/sched/goroutine_mgr"
	"net"
	"testing"
	"time"
)

func TestTcpServer(t *testing.T) {
	manager := new(GoroutineManager)
	manager.Initialise("mgr1")

	listener := new(TcpListener)
	err := listener.TCPListen(&net.TCPAddr{IP: net.ParseIP("0.0.0.0"), Port: 8888})
	if err != nil {
		fmt.Println(err)
		return
	}

	manager.GoroutineCreatePn("connector", connectorCallback, nil)

	for {
		conn, err := listener.TCPAccept()
		if err != nil {
			fmt.Println(err)
			return
		}
		manager.GoroutineCreateP1("receiver", receiverCallback, conn)
	}
}

func connectorCallback(g Goroutine, args ...interface{}) {
	defer g.OnQuit()

	fmt.Println("connectorCallback")
	conn := new(BufferTcpConn)
	err := conn.TCPConnect("127.0.0.1", 8888, 1)

	if err != nil {
		fmt.Println("connect error")
		return
	}

	for i := 0; i < 100; i++ {
		_ = conn.TCPWrite([]byte("1234567890abcdefg"))
		_ = conn.TCPFlush()
		time.Sleep(1000)
	}
	_ = conn.TCPDisConnect()
}

func receiverCallback(g Goroutine, conn interface{}) {
	defer g.OnQuit()

	fmt.Println("receiverCallback")
	for {
		c, _ := conn.(*BufferTcpConn)
		buffer, _, flag, err := c.TCPRead(1000)
		fmt.Println("receiver", buffer, flag, err)
		if err == nil && flag != true {
		} else {
			break
		}
	}
}
