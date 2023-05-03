/*
// Implementation of a main function setting a few characteristics of
// the game window, creating a game, and launching it
*/

package main

import (
	"flag"
	_ "image/png"
	"log"
	"net"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800 // Width of the game window (in pixels)
	screenHeight = 160 // Height of the game window (in pixels)
)

func connectToServer() {
	log.Println("Je me connecte")

	conn, err := net.Dial("tcp", os.Args[1] + ":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return
	}

	defer conn.Close()

	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Erreur en lisant la réponse du serveur:", err)
			return
		}

		// Afficher la réponse du serveur
		log.Println("Réponse du serveur:", string(buffer[:n]))
	}
}

func main() {

	go connectToServer()

	var getTPS bool
	flag.BoolVar(&getTPS, "tps", false, "Afficher le nombre d'appel à Update par seconde")
	flag.Parse()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("LP MiAR -- Programmation répartie (UE03EC2)")

	g := InitGame()
	g.getTPS = getTPS

	err := ebiten.RunGame(&g)
	log.Print(err)

}
