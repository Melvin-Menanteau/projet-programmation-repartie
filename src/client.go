package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type Client struct {
	conn        net.Conn
	runner      Runner
	globalState int
}

type serverGameMessage struct {
	State         int
	IdPlayer      string
	Xpos          float64
	Ypos          float64
	Arrived       bool
	RunTime       time.Duration
	ColorScheme   int
	ColorSelected bool
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

func (g *Game) listenServer() {
	for {
		buffer := make([]byte, 4096)
		n, err := (*g.serverConnection).Read(buffer)

		if err != nil {
			log.Println("Erreur en lisant les données du server")
		}

		var serverMessage serverGameMessage
		err = json.Unmarshal(buffer[:n], &serverMessage)

		log.Println("Message reçu du serveur: ", serverMessage)

		if err != nil {
			log.Println("Erreur en décodant les données")
		}

		log.Println("ancien état / nouveau état : ", g.state, "/", serverMessage.State)

		switch serverMessage.State {
		case StateChooseRunner:
			g.client.globalState = GlobalChooseRunner
			break
		case StateLaunchRun:
			g.client.globalState = GlobalLaunchRun
			break
		}
	}
}

func (g *Game) notifyServer() {
	jsonData, err := json.Marshal(serverGameMessage{g.state, "aaaa", 0, 0, false, 0, 0, false})

	if err != nil {
		log.Println("Erreur en encodant les données")
	}

	_, err = (*g.serverConnection).Write(jsonData)

	if err != nil {
		log.Println("Erreur en envoyant les données")
	}
}
