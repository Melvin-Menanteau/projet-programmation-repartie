package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	log.Println("Je me connecte")

	conn, err := net.Dial("tcp", "172.18.48.1:8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Erreur en lisant la réponse du serveur:", err)
			return
		}

		// Afficher la réponse du serveur
		fmt.Println("Réponse du serveur:", string(buffer[:n]))
	}
}