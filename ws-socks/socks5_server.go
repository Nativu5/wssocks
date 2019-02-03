package ws_socks

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/url"
	"strconv"
)

type ClientConfig struct {
	LocalAddr  string // "host:port"
	ServerAddr *url.URL
}

// client part of socks5 server
type Client struct {
	Config ClientConfig
}

// response to socks5 client and start to exchange data between socks5 client and
func (client *Client) Reply(conn net.Conn, onDial func(conn *net.TCPConn, addr string) error) error {
	defer conn.Close()
	var buffer [1024]byte

	n, err := conn.Read(buffer[:])
	if err != nil {
		return err
	}
	var addr string
	//sock5 proxy
	if buffer[0] != 0x05 {
		return errors.New("only socks5 supported")
	}

	// response to socks5 client
	// see rfc 1982 for more details (https://tools.ietf.org/html/rfc1928)
	n, err = conn.Write([]byte{0x05, 0x00})
	if err != nil {
		return err
	}

	n, err = conn.Read(buffer[:])
	if err != nil {
		return err
	}

	var host string
	switch buffer[3] {
	case 0x01:
		// ipv4 address
		ipv4 := make([]byte, 4)
		if _, err := io.ReadAtLeast(bytes.NewReader(buffer[4:]), ipv4, len(ipv4)); err != nil {
			return err
		}
		host = net.IP(ipv4).String()
	case 0x04:
		// ipv6
		ipv6 := make([]byte, 16)
		if _, err := io.ReadAtLeast(bytes.NewReader(buffer[4:]), ipv6, len(ipv6)); err != nil {
			return err
		}
		host = net.IP(ipv6).String()
	case 0x03:
		// domain
		addrLen := int(buffer[4])
		domain := make([]byte, addrLen)
		if _, err := io.ReadAtLeast(bytes.NewReader(buffer[5:]), domain, addrLen); err != nil {
			return err
		}
		host = string(domain)
	}

	port := make([]byte, 2)
	err = binary.Read(bytes.NewReader(buffer[n-2:n]), binary.BigEndian, &port)
	if err != nil {
		return err
	}

	addr = net.JoinHostPort(host, strconv.Itoa((int(port[0])<<8)|int(port[1])))
	if err := onDial(conn.(*net.TCPConn), addr); err != nil {
		return err
	}

	return nil
}
