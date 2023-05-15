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

func connectToServer() *net.Conn {
	log.Println("Je me connecte")

	conn, err := net.Dial("tcp", os.Args[1]+":8080")
	if err != nil {
		log.Println("Dial error:", err)
		return nil
	}

	return &conn

}

func main() {
	conn := connectToServer()

	defer (*conn).Close()

	var getTPS bool
	flag.BoolVar(&getTPS, "tps", false, "Afficher le nombre d'appel à Update par seconde")
	flag.Parse()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Rapide et furieux")

	g := InitGame(conn)
	g.getTPS = getTPS

	err := ebiten.RunGame(g)
	log.Print(err)

}
