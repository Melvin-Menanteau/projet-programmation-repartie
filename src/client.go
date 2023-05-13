package main

import (
	"encoding/json"
	"log"
	"net"
	"time"
)

type Client struct {
	idPlayer    string
	conn        net.Conn
	runner      *Runner
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

		if g.client.idPlayer != serverMessage.IdPlayer {
			if g.client.idPlayer == "" && serverMessage.IdPlayer != "" {
				log.Println("Changement du nom à", serverMessage.IdPlayer)
				g.client.idPlayer = serverMessage.IdPlayer
				g.runners[0].playerName = serverMessage.IdPlayer
				log.Println("Nom du runner à", g.runners[0].playerName)
			} else {
				for i := 0; i < len(g.runners); i++ {
					if g.runners[i].playerName == "" {
						g.runners[i].playerName = serverMessage.IdPlayer
					}

					if g.runners[i].playerName == serverMessage.IdPlayer {
						g.runners[i].playerName = serverMessage.IdPlayer
						g.runners[i].xpos = serverMessage.Xpos
						g.runners[i].ypos = serverMessage.Ypos
						g.runners[i].arrived = serverMessage.Arrived
						g.runners[i].runTime = serverMessage.RunTime
						g.runners[i].colorScheme = serverMessage.ColorScheme
						g.runners[i].colorSelected = serverMessage.ColorSelected

						break
					}
				}
			}
		}

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
	jsonData, err := json.Marshal(serverGameMessage{g.state, g.client.idPlayer, g.client.runner.xpos, g.client.runner.ypos, g.client.runner.arrived, g.client.runner.runTime, g.client.runner.colorScheme, g.client.runner.colorSelected})

	if err != nil {
		log.Println("Erreur en encodant les données")
	}

	_, err = (*g.serverConnection).Write(jsonData)

	if err != nil {
		log.Println("Erreur en envoyant les données")
	}
}
