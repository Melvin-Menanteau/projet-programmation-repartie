package main

import (
	// "bytes"
	// "encoding/gob"
	"encoding/binary"
	"encoding/json"
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

type serverMessage struct {
	State     int
	Time      int
	Position  float64
	Character int
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

		var message serverMessage
		err = json.Unmarshal(buffer[:n], &message)

		if err != nil {
			continue
		}

		log.Println("Message reçu du client: ", message)

		log.Println("Notifier le client: ", message.State)

		// notifyClient(conn, &message.State)
	}
}

func notifyClient(conn *net.Conn, gameState *int) {
	jsonData, err := json.Marshal(serverMessage{*gameState, 0, 0, 0})

	if err != nil {
		log.Println("Erreur en encodant les données")
	}

	_, err = (*conn).Write(jsonData)

	if err != nil {
		log.Println("Erreur en envoyant les données")
	}
}

// Envoie l'état du jeu au client
func notifyClientGameState(conn *net.Conn, gameState *int) {
	// Conversion de l'entier représentant l'état du jeu en un slice de bytes
	stateBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(stateBytes, uint32(*gameState))

	// Envoi du slice de bytes sur la connexion
	_, err := (*conn).Write(stateBytes)
	if err != nil {
		log.Println("Erreur en envoyant l'état du jeu au client:", err)
	}
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

		go listenClient(&conn)
		// go notifyClient(&conn, &gameState)
	}

	log.Println("Tous les clients sont connectés")

	gameState++

	log.Println("Notifier les clients: ", gameState)

	for _, conn := range clients {
		notifyClient(&conn, &gameState)
	}

	for {
		time.Sleep(1 * time.Second)
	}

	// for gameState == StateWelcomeScreen {

	// 	// recoit les messages des clients si ils sont prêts

	// 	for _, conn := range clients {
	// 		notifyClientGameState(&conn, &gameState)
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }

	// var network bytes.Buffer
	// enc := gob.NewEncoder(&network)

	// for _, conn := range clients {
	// 	jsonData, err := json.Marshal(serverMessage{gameState, 0, 0, 0})

	// 	if err != nil {
	// 		log.Println("Erreur en encodant les données")
	// 	}

	// 	_, err = conn.Write(jsonData)

	// 	buffer := make([]byte, 1024)
	// 	n, err := conn.Read(buffer)

	// 	if err != nil {
	// 		log.Println("Erreur en lisant les données du client")
	// 		return
	// 	}

	// 	var message serverMessage
	// 	err = json.Unmarshal(buffer[:n], &message)

	// 	if err != nil {
	// 		log.Println("Erreur en décodant les données")
	// 	}

	// 	log.Println("Message reçu du client: ", message)
	// 	log.Println("Message reçu du client: ", message.State)

	// _, err = conn.Write([]byte("Le jeu va commencer"))

	// encodingErr := enc.Encode(serverMessage{"gameState", gameState})

	// if encodingErr != nil {
	// 	log.Println("Erreur en encodant les données")
	// }

	// _, err = conn.Write(network.Bytes())

	// 	if err != nil {
	// 		log.Println("Erreur en envoyant des données au client")
	// 		// return
	// 	}

	// 	log.Println("Message envoyé au client: Le jeu va commencer")
	// }

	for client := range clients {
		defer clients[client].Close()
	}
}
