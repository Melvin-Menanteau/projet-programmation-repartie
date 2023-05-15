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
	conn           *net.Conn
	id             string
	gameState      int
	xpos           float64
	ypos           float64
	arrived        bool
	runTime        time.Duration
	colorScheme    int
	colorSelected  bool
	speed 		   float64
	animationFrame int
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

func listenClient(conn *net.Conn) (serverGameMessage, error) {
	buffer := make([]byte, 4096)
	n, err := (*conn).Read(buffer)

	if err != nil {
		log.Println("[ListenClient] Erreur en lisant les données")
		return serverGameMessage{}, err
	}

	var message serverGameMessage
	err = json.Unmarshal(buffer[:n], &message)

	if err != nil {
		log.Println("[ListenClient] Erreur en décodant les données")
		return serverGameMessage{}, err
	}

	log.Println("[ListenClient] Message reçu du client: ", message)

	return message, nil
}

func notifyClient(client *Client, message serverGameMessage) {
	jsonData, err := json.Marshal(message)

	if err != nil {
		log.Println("[NotifyClient] Erreur en encodant les données")
	}

	log.Println(fmt.Sprintf("[NotifyClient] Envoi des données au client (%s): ", client.id), message)

	_, err = (*client.conn).Write(jsonData)

	if err != nil {
		log.Println("[NotifyClient] Erreur en envoyant les données")
	}
}

func notifyAllClients(clients []Client, sourceClient Client) {
	for _, client := range clients {
		notifyClient(&client, buildServerGameMessage(&sourceClient, client.id == sourceClient.id))
	}
}

func buildServerGameMessage(client *Client, isSelf bool) serverGameMessage {
	log.Println("[BuildServerGameMessage] Construction du message pour le client ", client)
	
	return serverGameMessage{
		client.gameState,
		client.id,
		client.xpos,
		client.ypos,
		client.arrived,
		client.runTime,
		client.colorScheme,
		client.colorSelected,
		client.speed,
		client.animationFrame,
		isSelf,
		client.nbPlayersReady}
}

// Attendre que tous les clients aient choisi leur personnage
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

				clients[i].colorScheme = message.ColorScheme
				clients[i].colorSelected = true

				notifyAllClients(clients, clients[i])

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

// Attends que tous les clients aient finis la course
func waitForAllClientsToFinishRun(clients []Client) {
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

				clients[i].xpos = message.Xpos
				clients[i].animationFrame = message.AnimationFrame
				clients[i].runTime = message.RunTime
				clients[i].arrived = message.Arrived

				notifyAllClients(clients, clients[i])
				
				if message.Arrived == true {
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
		client.gameState = newState
		notifyClient(&client, buildServerGameMessage(&client, true))
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

	for len(clients) < 4 {
		conn, err := listener.Accept()

		clients = append(clients, Client{&conn, fmt.Sprintf("player%d", len(clients)), StateWelcomeScreen, 0, 0, false, 0, 0, false, 0, 0, len(clients)})

		if err != nil {
			log.Println("accept error:", err)
			return
		}

		notifyClient(&clients[len(clients)-1], buildServerGameMessage(&clients[len(clients)-1], true))

		for i, client := range clients {
			clients[i].nbPlayersReady = len(clients)
			notifyAllClients(clients, client)
		}
	}

	log.Println("Tous les clients sont connectés")

	setState(&gameState, StateChooseRunner, clients)
	log.Println("Notifier les clients: ", gameState)

	// fonction qui attend de recevoir un message de chaque client pour passer a l'état suivant
	waitForAllClientsToChooseCharacter(clients) // appel synchrone qui bloque le programme

	log.Println("Tous les clients ont choisit leur personnage")

	setState(&gameState, StateLaunchRun, clients)

	// Remettre le nombre de joueurs prêt à 0 pour ne pas pouvoir relancer la course directement après les résultats
	// for i, client := range clients {
	// 	clients[i].nbPlayersReady = 0
	// 	notifyAllClients(clients, client)
	// }

	// Attends que tous le clients aient finis la course
	waitForAllClientsToFinishRun(clients)

	log.Println("Tous les clients ont finis la course")
	setState(&gameState, StateResult, clients)

	for {
		time.Sleep(1 * time.Second)
	}

	for _, client := range clients {
		defer (*client.conn).Close()
	}
}
