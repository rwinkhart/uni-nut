package nut

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// A Client wraps a connection to a NUT server.
type Client struct {
	conn        net.Conn
	br          *bufio.Reader
	upsID       string
	upsIDLength int
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

// Authenticate logs into credential-protected sessions.
// If authentication is enabled on your NUT server, run
// this immediately after dialing.
func (c *Client) Authenticate(username, password string) error {
	var err error
	if err = c.write("USERNAME " + username + "\nPASSWORD " + password); err != nil {
		return err
	}
	_, err = c.clearBuffer("OK L", "")
	if err != nil {
		return err
	}
	return nil
}

// AutomaticallySetID detects the ID of the connected UPS.
// This only works when there is only one UPS per
// NUT server, which is the case with UniFi UPS units.
func (c *Client) AutomaticallySetID() error {
	if err := c.write("LIST UPS"); err != nil {
		return err
	}
	l, err := c.clearBuffer("END", "UPS")
	if err != nil {
		return err
	}
	lSplit := strings.Split(l, " ")
	c.upsID = strings.Join(lSplit[1:len(lSplit)-2], " ")
	c.upsIDLength = len(strings.Split(c.upsID, " "))
	return nil
}

// ManuallySetID allows the user to specify the UPS ID
// manually if auto-detection is not desired.
func (c *Client) ManuallySetID(upsID string) {
	c.upsID = upsID
	c.upsIDLength = len(strings.Split(upsID, " "))
}

// ListVar updates NutKeyValMap with the current status of all variables from the target UPS.
func (c *Client) ListVar() error {
	cmd := "LIST VAR \"" + c.upsID + "\""
	if err := c.write(cmd); err != nil {
		return err
	}
	l, err := c.read()
	if err != nil {
		return err
	}
	expectedPrefix := "BEGIN LIST VAR"
	if !strings.HasPrefix(l, expectedPrefix) {
		return fmt.Errorf("pre-loop error: expected prefix %q, got line %q", expectedPrefix, l)
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
			NutKeyValMap[lSplit[1+c.upsIDLength]] = strings.Trim(strings.Join(lSplit[2+c.upsIDLength:], " "), "\"")
		default:
			break outer
		}
	}
	return nil
}

// GetVar returns the value of the specified variable for the target UPS.
func (c *Client) GetVar(varName string) (string, error) {
	if err := c.write("GET VAR \"" + c.upsID + "\" \"" + varName + "\""); err != nil {
		return "", err
	}
	l, err := c.read()
	if err != nil {
		return "", err
	}
	lSplit := strings.Split(l, " ")
	if len(lSplit) < 4 {
		return "", fmt.Errorf("invalid response to GET VAR; check your UPS ID")
	}
	value := strings.Trim(strings.Join(lSplit[2+c.upsIDLength:], " "), "\"")
	return value, nil
}

func newClient(conn net.Conn) *Client {
	return &Client{conn, bufio.NewReader(conn), "", 0}
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

func (c *Client) clearBuffer(tillPrefix, storePrefix string) (string, error) {
	var stored string
	for {
		l, err := c.read()
		if err != nil {
			return "", err
		}
		if storePrefix != "" && strings.HasPrefix(l, storePrefix) {
			stored = l
		} else if strings.HasPrefix(l, tillPrefix) {
			break
		}
	}
	return stored, nil
}
