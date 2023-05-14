package main

import (
	"encoding/json"
	"fmt"
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
	IsSelf        bool
}

const (
	GlobalWelcomeScreen int = iota
	GlobalChooseRunner
	GlobalLaunchRun
	GlobalStateRun
	GlobalResult
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

		if serverMessage.IsSelf {
			if g.state != serverMessage.State {
				log.Println("Changement d'état de ", g.state, " => ", serverMessage.State)

				switch serverMessage.State {
				case StateWelcomeScreen:
					g.client.globalState = GlobalWelcomeScreen
					break
				case StateChooseRunner:
					g.client.globalState = GlobalChooseRunner
					break
				case StateLaunchRun:
					g.client.globalState = GlobalLaunchRun
					break
				case StateRun:
					g.client.globalState = GlobalStateRun
					break
				case StateResult:
					g.client.globalState = GlobalResult
					break
				}
			}

			if g.client.idPlayer != serverMessage.IdPlayer {
				log.Println("Changement du nom à", serverMessage.IdPlayer)
				g.client.idPlayer = serverMessage.IdPlayer
				g.runners[0].playerName = serverMessage.IdPlayer
			}
		} else {
			for i := 0; i < len(g.runners); i++ {
				if !g.runners[i].hasBeenAttributed {
					log.Println(fmt.Sprintf("Le runner %d a été attribué à %s", i, serverMessage.IdPlayer))
					g.runners[i].playerName = serverMessage.IdPlayer
					g.runners[i].hasBeenAttributed = true
				}

				if g.runners[i].playerName == serverMessage.IdPlayer {
					g.runners[i].playerName = serverMessage.IdPlayer
					g.runners[i].xpos = serverMessage.Xpos
					// g.runners[i].ypos = serverMessage.Ypos
					g.runners[i].arrived = serverMessage.Arrived
					g.runners[i].runTime = serverMessage.RunTime
					g.runners[i].colorScheme = serverMessage.ColorScheme
					g.runners[i].colorSelected = serverMessage.ColorSelected

					break
				}
			}
		}
	}
}

func (g *Game) notifyServer() {
	jsonData, err := json.Marshal(serverGameMessage{
		g.state,
		g.client.idPlayer,
		g.client.runner.xpos,
		g.client.runner.ypos,
		g.client.runner.arrived,
		g.client.runner.runTime,
		g.client.runner.colorScheme,
		g.client.runner.colorSelected,
		true})

	if g.state == StateChooseRunner || g.state == StateRun {
		log.Println("Envoi des données au serveur: ", string(jsonData))
	}

	if err != nil {
		log.Println("Erreur en encodant les données")
	}

	_, err = (*g.serverConnection).Write(jsonData)

	if err != nil {
		log.Println("Erreur en envoyant les données")
	}
}
