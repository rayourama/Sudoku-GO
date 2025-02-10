package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Dimensions de la grille
const (
	screenWidth  = 900
	screenHeight = 900
)

// Game structure
type Game struct{}

// Cette fonction met à jour l'état du jeu
func (g *Game) Update() error {
	return nil
}

// Interface graphique
func (g *Game) Draw(screen *ebiten.Image) {
	texte := "Bienvenue dans Sudoku en Go!"
	ebitenutil.DebugPrint(screen, texte)
}

// Layout gère la disposition de l'écran
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// Création du jeu
	game := &Game{}

	// Configuration de la fenêtre
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sudoku en Go")

	// Lancement du jeu
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
