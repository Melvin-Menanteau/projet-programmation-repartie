package main

import (
	"bufio"
	"log"
	"net"
)

type Client struct {
	conn        net.Conn
	runner      Runner
	globalState int
}

const (
	GlobalWelcomeScreen int = iota
	GlobalChooseRunner
	GlobalLaunchRun
)

func NewClient() *Client {
	return &Client{
		globalState: GlobalWelcomeScreen,
	}
}

func (c *Client) connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) sendMessage(message string) {
	// send message to server using bufio writer
	writer := bufio.NewWriter(c.conn)
	_, err := writer.WriteString(message + "\n")
	if err != nil {
		log.Printf("Error sending message to server: %v", err)
	} else {
		writer.Flush()
	}
}
