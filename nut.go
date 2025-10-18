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

// ListVar updates NutKeyValMap with the current status of all variables from <upsID>.
func (c *Client) ListVar(upsID string) error {
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

outer:
	for {
		l, err := c.read()
		if err != nil {
			return fmt.Errorf("in-loop error: %w", err)
		}
		lSplit := strings.Split(l, " ")
		switch lSplit[0] {
		case "VAR":
			NutKeyValMap[strings.Trim(lSplit[2], "\"")] = strings.Trim(strings.Join(lSplit[3:], " "), "\"")
		default:
			break outer
		}
	}
	return nil
}

// GetVar returns the value of the specified variable for <upsID>.
func (c *Client) GetVar(upsID, varName string) (string, error) {
	cmd := "GET VAR \"" + upsID + "\" \"" + varName + "\""
	if err := c.write(cmd); err != nil {
		return "", err
	}
	l, err := c.read()
	if err != nil {
		return "", err
	}
	value := strings.Trim(strings.Join(strings.Split(l, " ")[3:], " "), "\"")
	return value, nil
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
