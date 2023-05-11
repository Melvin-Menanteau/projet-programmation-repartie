package main

import (
	"bufio"
	"encoding/json"
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

func (g *Game) listenServer() {
	for {
		buffer := make([]byte, 1024)
		n, err := (*g.serverConnection).Read(buffer)

		if err != nil {
			log.Println("Erreur en lisant les données du server")
		}

		var serverMessage serverMessage
		err = json.Unmarshal(buffer[:n], &serverMessage)

		log.Println("Message reçu du serveur: ", serverMessage)

		if err != nil {
			log.Println("Erreur en décodant les données")
		}

		log.Println("ancien état / nouveau état : ", g.state, "/", serverMessage.State)

		if serverMessage.State == StateChooseRunner {
			log.Println("Serveur prêt a changer d'état, valeur état serveur : ", serverMessage.State)
			g.client.globalState = GlobalChooseRunner
		}

		log.Println("SDflikjhsdafliphsflpsadf")
	}
}
