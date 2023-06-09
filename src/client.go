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
	nbPlayersReady int
}

type serverGameMessage struct {
	State          int
	IdPlayer       string
	Xpos           float64
	Ypos           float64
	Arrived        bool
	RunTime        time.Duration
	ColorScheme    int
	ColorSelected  bool
	Speed		   float64
	AnimationFrame int
	IsSelf         bool
	NbPlayersReady int
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
			// On met à jour les données du client
			g.client.nbPlayersReady = serverMessage.NbPlayersReady

			// Si l'état du jeu a changé, on change l'état du client
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

			// Si le nom du client a changé, on change le nom du client
			if g.client.idPlayer != serverMessage.IdPlayer {
				log.Println("Changement du nom à", serverMessage.IdPlayer)
				g.client.idPlayer = serverMessage.IdPlayer
				g.runners[0].playerName = serverMessage.IdPlayer
				g.runners[0].hasBeenAttributed = true
			}
		} else {
			// Si aucun client n'a été attribué, on attribue le premier client qui se présente
			if (func() bool {
				for i := 0; i < len(g.runners); i++ {
					if g.runners[i].playerName == serverMessage.IdPlayer {
						return false
					}
				}
				return true
			}()) {
				for i := 0; i < len(g.runners); i++ {
					if !g.runners[i].hasBeenAttributed {
						log.Println(fmt.Sprintf("Le runner %d a été attribué à %s", i, serverMessage.IdPlayer))
						g.runners[i].playerName = serverMessage.IdPlayer
						g.runners[i].hasBeenAttributed = true

						break
					}
				}
			}

			for i := 0; i < len(g.runners); i++ {
				if g.runners[i].playerName == serverMessage.IdPlayer {
					g.runners[i].playerName = serverMessage.IdPlayer
					g.runners[i].xpos = serverMessage.Xpos
					g.runners[i].arrived = serverMessage.Arrived
					g.runners[i].runTime = serverMessage.RunTime
					g.runners[i].colorScheme = serverMessage.ColorScheme
					g.runners[i].colorSelected = serverMessage.ColorSelected
					g.runners[i].speed = serverMessage.Speed
					g.runners[i].animationFrame = serverMessage.AnimationFrame
					g.runners[i].hasBeenAttributed = true

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
		g.client.runner.speed,
		g.client.runner.animationFrame,
		true,
		g.client.nbPlayersReady})

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
