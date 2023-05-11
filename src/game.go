/*
//  Data structure for representing a game. Implements the ebiten.Game
//  interface (Update in game-update.go, Draw in game-draw.go, Layout
//  in game-layout.go). Provided with a few utilitary functions:
//    - initGame
*/

package main

import (
	"bytes"
	"course/assets"
	"image"
	"log"
	"net"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	state            int           // Current state of the game
	stateLock        sync.Mutex    // Lock for the state
	runnerImage      *ebiten.Image // Image with all the sprites of the runners
	runners          [4]Runner     // The four runners used in the game
	f                Field         // The running field
	launchStep       int           // Current step in StateLaunchRun state
	resultStep       int           // Current step in StateResult state
	getTPS           bool          // Help for debug
	serverConnection *net.Conn
	client           *Client // Client associated with the runner
}

// These constants define the five possible states of the game
const (
	StateWelcomeScreen int = iota // Title screen
	StateChooseRunner             // Player selection screen
	StateLaunchRun                // Countdown before a run
	StateRun                      // Run
	StateResult                   // Results announcement
)

// setter synchronisé pour l'état du jeu
func (g *Game) SetState(state int) {
	g.stateLock.Lock()         // acquisition du verrou
	defer g.stateLock.Unlock() // libération du verrou après la fin de la méthode

	g.state = state
}

// InitGame builds a new game ready for being run by ebiten
func InitGame(serverConnection *net.Conn) *Game {

	if serverConnection == nil {
		log.Fatal("No server connection")
	}

	g := &Game{}
	g.serverConnection = serverConnection

	g.client = NewClient()

	// go g.notifyServer()
	go g.listenServer()

	// Open the png image for the runners sprites
	img, _, err := image.Decode(bytes.NewReader(assets.RunnerImage))
	if err != nil {
		log.Fatal(err)
	}
	g.runnerImage = ebiten.NewImageFromImage(img)

	// Define game parameters
	start := 50.0
	finish := float64(screenWidth - 50)
	frameInterval := 20

	// Create the runners
	for i := range g.runners {
		interval := 0
		if i == 0 {
			interval = frameInterval
		}
		g.runners[i] = Runner{
			xpos: start, ypos: 50 + float64(i*20),
			maxFrameInterval: interval,
			colorScheme:      0,
		}
	}

	// Create the field
	g.f = Field{
		xstart:   start,
		xarrival: finish,
		chrono:   time.Now(),
	}

	return g
}
