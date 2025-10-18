package nut

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// A Client wraps a connection to a NUT server.
type Client struct {
	conn net.Conn
	br   *bufio.Reader
}

// Global map to be updated by GetListVar and read by the importing program.
var NutKeyValMap = make(map[string]string)

// Dial dials a NUT server using TCP. If the address does not contain
// a port number, it will default to 3493.
func Dial(addr string) (*Client, error) {
	_, _, err := net.SplitHostPort(addr)
	if err != nil {
		addr = net.JoinHostPort(addr, "3493")
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return newClient(conn), nil
}

// Close closes the connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetListVar(upsID string) error {
	cmd := "LIST VAR \"" + upsID + "\""
	if err := c.write(cmd); err != nil {
		return err
	}
	l, err := c.read()
	if err != nil {
		return err
	}
	expected := "BEGIN " + cmd
	if l != expected {
		return fmt.Errorf("pre-loop error: expected %q, got %q", expected, l)
	}

	for {
		l, _ := c.read()
		if !strings.HasPrefix(l, "VAR \""+upsID+"\" ") {
			break
		}
		lSplit := strings.Split(l, " ")
		NutKeyValMap[strings.Trim(lSplit[2], "\"")] = strings.Trim(strings.Join(lSplit[3:], " "), "\"")
	}
	return nil
}

// newClient wraps an existing net.Conn.
func newClient(conn net.Conn) *Client {
	return &Client{conn, bufio.NewReader(conn)}
}

func (c *Client) write(s string) error {
	_, err := c.conn.Write([]byte(s + "\n"))
	return err
}

func (c *Client) read() (string, error) {
	l, err := c.br.ReadString('\n')
	if err != nil {
		return "", err
	}
	if len(l) > 0 {
		l = l[:len(l)-1]
	}
	return l, nil
}
