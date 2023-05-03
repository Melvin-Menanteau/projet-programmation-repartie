package main

import (
	"log"
	"net"
)

func main() {
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
			return
		}

		log.Println("Message envoyé au client: ", message)

		// Fermer la connexion quand le programme se termine
		defer conn.Close()
	}

	for _, conn := range clients {
		_, err = conn.Write([]byte("Le jeu va commencer"))

		if err != nil {
			log.Println("Erreur en envoyant des données au client")
			return
		}

		log.Println("Message envoyé au client: Le jeu va commencer")
	}
}