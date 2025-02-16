package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// Dimensions de la grille
const (
	screenWidth  = 900
	screenHeight = 900
	cellSize     = screenWidth / 9
)

// Game structure
type Game struct{}

// Cette fonction met à jour l'état du jeu
func (g *Game) Update() error {
	return nil
}

// Interface graphique
func (g *Game) Draw(screen *ebiten.Image) {
	// Couleurs
	bgColor := color.RGBA{240, 240, 240, 255} // Gris clair
	screen.Fill(bgColor)
	lineColor := color.RGBA{0, 0, 0, 255}

	// Dessiner la grille 9x9
	for i := 0; i <= 9; i++ {
		thickness := 1.0
		if i%3 == 0 {
			thickness = 3.0 // Traits plus épais pour séparer les blocs 3x3
		}
		// Lignes verticales
		vector.StrokeLine(screen, float32(i*cellSize), 0, float32(i*cellSize), screenHeight, float32(thickness), lineColor, true)
		// Lignes horizontales
		vector.StrokeLine(screen, 0, float32(i*cellSize), screenWidth, float32(i*cellSize), float32(thickness), lineColor, true)
	}
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
	ebiten.SetWindowTitle("Sudoku")

	// Lancement du jeu
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
