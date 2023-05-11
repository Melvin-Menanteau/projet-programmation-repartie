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
	state         int
	idPlayer      string
	xpos          float64
	ypos          float64
	arrived       bool
	runTime       time.Duration
	colorScheme   int
	colorSelected bool
}

func listenClient(conn *net.Conn) {
	for {
		buffer := make([]byte, 1024)
		n, err := (*conn).Read(buffer)

		if n == 0 {
			continue
		}

		if err != nil {
			log.Println("Erreur en lisant les données")
			continue
		}

		var message serverGameMessage
		err = json.Unmarshal(buffer[:n], &message)

		if err != nil {
			continue
		}

		log.Println("Message reçu du client: ", message)
	}
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

		go listenClient(&conn)
	}

	log.Println("Tous les clients sont connectés")

	gameState++

	log.Println("Notifier les clients: ", gameState)

	for _, client := range clients {
		notifyClient(&client, &gameState)
	}

	for {
		time.Sleep(1 * time.Second)
	}

	for _, client := range clients {
		defer (*client.conn).Close()
	}
}
