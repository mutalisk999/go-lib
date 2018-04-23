package buffer_tcp

import (
	"errors"
	"io"
	"net"
	"strconv"
	"time"
)

type BufferTcpConn struct {
	conn       net.Conn
	sendBuffer string
	readBuffer string
}

type TcpListener struct {
	tcpAddr  *net.TCPAddr
	listener *net.TCPListener
}

func (c *BufferTcpConn) TCPConnect(serverAddr string, serverPort uint16, timeOut float64) error {
	client, err := net.DialTimeout("tcp", serverAddr+":"+strconv.Itoa(int(serverPort)),
		time.Duration(timeOut*1000*1000*1000))
	if err != nil {
		return err
	}

	c.conn = client
	c.sendBuffer = ""
	c.readBuffer = ""
	return nil
}

func (c *BufferTcpConn) TCPRead(nRead uint32) ([]byte, uint32, bool, error) {
	remoteClose := false
	for {
		if uint32(len(c.readBuffer)) >= nRead {
			readBytes := []byte(c.readBuffer[0:nRead])
			c.readBuffer = c.readBuffer[nRead:]
			return readBytes, nRead, remoteClose, nil
		} else if remoteClose == true {
			readBytes := []byte(c.readBuffer)
			c.readBuffer = ""
			return readBytes, uint32(len(readBytes)), remoteClose, nil
		}

		var buf [4096]byte
		n, err := c.conn.Read(buf[0:])
		if err != nil {
			if err == io.EOF {
				remoteClose = true
			} else {
				return nil, 0, remoteClose, err
			}
		}
		c.readBuffer = c.readBuffer + string(buf[0:n])
	}
}

func (c *BufferTcpConn) TCPFlush() error {
	n, err := c.conn.Write([]byte(c.sendBuffer))
	if err != nil {
		return err
	}
	if n != len(c.sendBuffer) {
		return errors.New("can not send completely")
	}
	c.sendBuffer = ""
	return nil
}

func (c *BufferTcpConn) TCPWrite(bytesWrite []byte) error {
	c.sendBuffer = c.sendBuffer + string(bytesWrite)
	if len(c.sendBuffer) > 40960 {
		err := c.TCPFlush()
		if err != nil {
			return err
		}
		c.sendBuffer = ""
	}
	return nil
}

func (c *BufferTcpConn) TCPDisConnect() error {
	if len(c.sendBuffer) > 0 {
		err := c.TCPFlush()
		if err != nil {
			return err
		}
		c.sendBuffer = ""
	}
	c.conn.Close()
	return nil
}

func (c *TcpListener) TCPListen(tcpAddr *net.TCPAddr) error {
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		return err
	}
	c.tcpAddr = tcpAddr
	c.listener = listener
	return nil
}

func (c *TcpListener) TCPAccept() (*BufferTcpConn, error) {
	if c.listener == nil {
		return nil, errors.New("invalid listener")
	}

	conn, err := c.listener.Accept()
	if err != nil {
		return nil, err
	}

	tcpConn := new(BufferTcpConn)
	tcpConn.conn = conn
	tcpConn.sendBuffer = ""
	tcpConn.readBuffer = ""
	return tcpConn, nil
}

func (c *TcpListener) TCPListenClose() {
	if c.listener != nil {
		c.listener.Close()
	}
}
