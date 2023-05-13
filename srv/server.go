package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

// These constants define the five possible states of the game
const (
	StateWelcomeScreen int = iota // Title screen
	StateChooseRunner             // Player selection screen
	StateLaunchRun                // Countdown before a run
	StateRun                      // Run
	StateResult                   // Results announcement
)

type Client struct {
	conn *net.Conn
	id   string
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

func listenClient(conn *net.Conn) (serverGameMessage, error) {
	buffer := make([]byte, 4096)
	n, err := (*conn).Read(buffer)

	if err != nil {
		log.Println("Erreur en lisant les données")
		return serverGameMessage{}, err
	}

	var message serverGameMessage
	err = json.Unmarshal(buffer[:n], &message)

	if err != nil {
		log.Println("Erreur en décodant les données")
		return serverGameMessage{}, err
	}

	log.Println("Message reçu du client: ", message)
	return message, nil
}

func notifyClient(client *Client, gameState *int) {
	data := serverGameMessage{*gameState, client.id, 0, 0, false, 0, 0, false}

	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("Erreur en encodant les données")
	}

	log.Println("Envoi des données au client: ", data)

	_, err = (*client.conn).Write(jsonData)

	if err != nil {
		log.Println("Erreur en envoyant les données")
	}
}

// permet de lire un message envoyé par un client
func readMessage(conn *net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := (*conn).Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

func waitForAllClientsToChooseCharacter(clients []Client) {
	channels := make([]chan bool, len(clients))

	// Créer un canal pour chaque client
	for i := 0; i < len(clients); i++ {
		channels[i] = make(chan bool)
	}

	// Lancer une goroutine pour chaque client qui attend un message booléen
	for i, client := range clients {
		go func(i int, client Client) {
			for {
				message, err := listenClient(client.conn)
				if err != nil {
					log.Println("error reading message from client ", client.id, err)
					continue
				}
				if message.ColorSelected == true {
					channels[i] <- true
					break
				}
			}
		}(i, client)
	}

	// Attendre que tous les canaux aient reçu une valeur true
	for _, ch := range channels {
		<-ch
	}
}

func setState(gameState *int, newState int, clients []Client) {
	*gameState = newState
	for _, client := range clients {
		notifyClient(&client, gameState)
	}
}

func main() {
	gameState := StateWelcomeScreen

	clients := make([]Client, 0)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Println("listen error:", err)
		return
	}

	log.Println("Le serveur est en écoute sur le port 8080")

	// Fermer le listener quand le programme se termine
	defer listener.Close()

	for len(clients) < 2 {
		conn, err := listener.Accept()

		clients = append(clients, Client{&conn, fmt.Sprint("client-", len(clients))})

		if err != nil {
			log.Println("accept error:", err)
			return
		}

		notifyClient(&clients[len(clients)-1], &gameState)
	}

	log.Println("Tous les clients sont connectés")

	setState(&gameState, StateChooseRunner, clients)
	log.Println("Notifier les clients: ", gameState)

	// fonction qui attend de recevoir un message de chaque client pour passer a l'état suivant
	waitForAllClientsToChooseCharacter(clients) // appel synchrone qui bloque le programme

	log.Println("Tous les clients ont choisit leur personnage")
	setState(&gameState, StateLaunchRun, clients)

	for {
		time.Sleep(1 * time.Second)
	}

	for _, client := range clients {
		defer (*client.conn).Close()
	}
}
