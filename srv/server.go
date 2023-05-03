package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
)

// These constants define the five possible states of the game
const (
	StateWelcomeScreen int = iota // Title screen
	StateChooseRunner             // Player selection screen
	StateLaunchRun                // Countdown before a run
	StateRun                      // Run
	StateResult                   // Results announcement
)

type serverMessage struct {
	message string
	data int
}

func main() {
	gameState := StateWelcomeScreen

	clients := make([]net.Conn, 0)
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

		clients = append(clients, conn)

		if err != nil {
			log.Println("accept error:", err)
			return
		}

		message := "Reponse du serveur"
		_, err = conn.Write([]byte(message))

		if err != nil {
			log.Println("Erreur en envoyant des données au client")
			// return
		}

		log.Println("Message envoyé au client: ", message)

		// Fermer la connexion quand le programme se termine
		defer conn.Close()
	}

	var network bytes.Buffer
	enc := gob.NewEncoder(&network)

	for _, conn := range clients {
		// _, err = conn.Write([]byte("Le jeu va commencer"))

		encodingErr := enc.Encode(serverMessage{"gameState", gameState})

		if encodingErr != nil {
			log.Println("Erreur en encodant les données")
		}

		_, err = conn.Write(network.Bytes())

		if err != nil {
			log.Println("Erreur en envoyant des données au client")
			// return
		}

		log.Println("Message envoyé au client: Le jeu va commencer")
	}
}