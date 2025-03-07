package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 900
	screenHeight = 900
	cellSize     = screenWidth / 10
	emptyCells   = 40
	topMargin    = 50
)

type Game struct {
	grid         [9][9]int
	solution     [9][9]int
	originalGrid [9][9]int
	selectedRow  int
	selectedCol  int
	lastMoveTime time.Time
	gameOver     bool
	win          bool
	startTime    time.Time
	endTime      time.Time
}

var fontFace font.Face
var moveCooldown = 150 * time.Millisecond

// Permet de charger la police
func loadFont() {
	ttfBytes, err := ioutil.ReadFile("details/font/Roboto-Bold.ttf")
	if err != nil {
		log.Fatal(err)
	}
	ttf, err := opentype.Parse(ttfBytes)
	if err != nil {
		log.Fatal(err)
	}
	fontFace, err = opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func loadIcon() ([]image.Image, error) {
	// Liste des fichiers d'icônes
	iconFiles := []string{
		"details/assets/icon64.png",  // 64x64
		"details/assets/icon128.png", // 128x128
		"details/assets/icon256.png", // 256x256
	}

	var icons []image.Image
	// Charger les images à partir des fichiers
	for _, file := range iconFiles {
		img, err := loadImage(file)
		if err != nil {
			return nil, err
		}
		icons = append(icons, img)
	}
	return icons, nil
}

// Fonction pour charger une image à partir d'un fichier PNG
func loadImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", file, err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %v", file, err)
	}
	return img, nil
}

// On génère une grille
func (g *Game) generateSudoku() {
	generateFullSudoku(&g.solution)
	g.grid = g.solution
	removeNumbers(&g.grid, emptyCells)
	g.originalGrid = g.grid
}

// Cette fonction met à jour l'état du jeu
func (g *Game) Update() error {
	if time.Since(g.lastMoveTime) < moveCooldown {
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.selectedRow = (g.selectedRow + 8) % 9
		g.lastMoveTime = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.selectedRow = (g.selectedRow + 1) % 9
		g.lastMoveTime = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.selectedCol = (g.selectedCol + 8) % 9
		g.lastMoveTime = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.selectedCol = (g.selectedCol + 1) % 9
		g.lastMoveTime = time.Now()
	}

	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		g.setCellValue(0)
		g.lastMoveTime = time.Now()
	}

	// Dimensions et position du rectangle
	boxHeight := 100
	boxY := (screenHeight - boxHeight) / 2

	// Dimensions et position du bouton
	buttonWidth := 200
	buttonHeight := 50
	buttonX := (screenWidth - buttonWidth) / 2
	buttonY := boxY + boxHeight + 20

	if g.gameOver && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mouseX, mouseY := ebiten.CursorPosition()

		// Vérifier si le clic est bien dans toute la zone du bouton
		if mouseX >= buttonX && mouseX < buttonX+buttonWidth &&
			mouseY >= buttonY && mouseY < buttonY+buttonHeight {
			g.restartGame()
		}
	}

	for k := ebiten.Key1; k <= ebiten.Key9; k++ {
		if ebiten.IsKeyPressed(k) {
			g.setCellValue(int(k - ebiten.Key0))
			g.lastMoveTime = time.Now()
		}
	}
	for k := ebiten.KeyNumpad1; k <= ebiten.KeyNumpad9; k++ {
		if ebiten.IsKeyPressed(k) {
			g.setCellValue(int(k - ebiten.KeyNumpad0))
			g.lastMoveTime = time.Now()
		}
	}
	return nil
}

// Permet de modifier une cellule
func (g *Game) setCellValue(value int) {
	if g.gameOver || g.originalGrid[g.selectedRow][g.selectedCol] != 0 {
		return
	}

	if value == g.solution[g.selectedRow][g.selectedCol] {
		g.grid[g.selectedRow][g.selectedCol] = value
		if g.isSudokuCompleted() {
			g.win = true
			g.gameOver = true
		}
	} else {
		g.grid[g.selectedRow][g.selectedCol] = -value
	}
}

// Cette fonction vérifie que la grille est complète
func (g *Game) isSudokuCompleted() bool {
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if g.grid[y][x] <= 0 { // Si une case est vide ou fausse
				return false
			}
		}
	}
	return true
}

// Cette fonction génère une grille remplie
func generateFullSudoku(grid *[9][9]int) {
	rand.Seed(time.Now().UnixNano())
	fillGrid(grid)
}

// Cette fonction permet le remplissage de la grille
func fillGrid(grid *[9][9]int) bool {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if grid[row][col] == 0 {
				numbers := rand.Perm(9)
				for _, n := range numbers {
					num := n + 1
					if isValid(grid, row, col, num) {
						grid[row][col] = num
						if fillGrid(grid) {
							return true
						}
						grid[row][col] = 0
					}
				}
				return false
			}
		}
	}
	return true
}

// Cette fonction retire des nombres aléatoires dans la grille
func removeNumbers(grid *[9][9]int, emptyCount int) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < emptyCount; i++ {
		row, col := rand.Intn(9), rand.Intn(9)
		for grid[row][col] == 0 {
			row, col = rand.Intn(9), rand.Intn(9)
		}
		grid[row][col] = 0
	}
}

// Cette fonction vérifie l'entrée utilisateur
func isValid(grid *[9][9]int, row, col, num int) bool {
	for i := 0; i < 9; i++ {
		if grid[row][i] == num || grid[i][col] == num || grid[row/3*3+i/3][col/3*3+i%3] == num {
			return false
		}
	}
	return true
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{240, 240, 240, 255})

	var timerText string
	var chrono string

	if !g.gameOver && !g.win {
		// Calcul du temps écoulé avant la victoire
		écoulé := time.Since(g.startTime)
		minutes := int(écoulé.Seconds()) / 60
		seconds := int(écoulé.Seconds()) % 60
		chrono = fmt.Sprintf("Temps: %02d:%02d", minutes, seconds)
		text.Draw(screen, chrono, fontFace, screenWidth/2-110, 40, color.Black)

	} else if g.gameOver && g.win {
		if g.endTime.IsZero() {
			g.endTime = time.Now()
		}

		// Calcul du temps entre le début et la fin du jeu
		écoulé := g.endTime.Sub(g.startTime)
		minutes := int(écoulé.Seconds()) / 60
		seconds := int(écoulé.Seconds()) % 60
		timerText = fmt.Sprintf("Votre temps est de : %02d:%02d", minutes, seconds)
	}

	// Affichage du timer en haut
	text.Draw(screen, timerText, fontFace, screenWidth/2-200, 40, color.Black)

	lineColor := color.RGBA{0, 0, 0, 255}

	gridWidth := cellSize * 9
	gridLeft := (screenWidth - gridWidth) / 2

	for i := 0; i <= 9; i++ {
		thickness := 1.0
		if i%3 == 0 {
			thickness = 3.0
		}
		// Lignes verticales
		vector.StrokeLine(screen, float32(gridLeft+i*cellSize), float32(topMargin), float32(gridLeft+i*cellSize), float32(screenHeight-40), float32(thickness), lineColor, true)
		// Lignes horizontales
		vector.StrokeLine(screen, float32(gridLeft), float32(i*cellSize)+float32(topMargin), float32(gridLeft+gridWidth), float32(i*cellSize)+float32(topMargin), float32(thickness), lineColor, true)
	}

	// Affichage des chiffres
	for y := 0; y < 9; y++ {
		for x := 0; x < 9; x++ {
			if g.grid[y][x] != 0 {
				num := fmt.Sprintf("%d", abs(g.grid[y][x]))
				var colorToUse color.Color

				if g.originalGrid[y][x] != 0 {
					colorToUse = color.Black // Chiffre initial
				} else if g.grid[y][x] > 0 {
					colorToUse = color.RGBA{0, 200, 0, 255} // Vert pour correct
				} else {
					colorToUse = color.RGBA{255, 0, 0, 255} // Rouge pour erreur
					num = fmt.Sprintf("%d", -g.grid[y][x])
				}
				text.Draw(screen, num, fontFace, int(float32(gridLeft)+float32(x*cellSize)+30), int(float32(y*cellSize)+float32(topMargin)+50), colorToUse)
			}
		}
	}

	if g.gameOver {
		var msg string
		var bgColor color.RGBA
		if g.win {
			msg = "Vous avez gagné !"
			bgColor = color.RGBA{50, 205, 50, 200}
		}

		// Dimensions et position du rectangle
		boxWidth := 500
		boxHeight := 100
		boxX := (screenWidth - boxWidth) / 2
		boxY := (screenHeight - boxHeight) / 2

		// Dimensions et position du bouton
		buttonWidth := 200
		buttonHeight := 50
		buttonX := (screenWidth - buttonWidth) / 2
		buttonY := boxY + boxHeight + 20

		// Dessiner le bouton
		buttonColor := color.RGBA{30, 144, 255, 255} // Bleu
		vector.DrawFilledRect(screen, float32(buttonX), float32(buttonY), float32(buttonWidth), float32(buttonHeight), buttonColor, true)
		vector.StrokeRect(screen, float32(buttonX), float32(buttonY), float32(buttonWidth), float32(buttonHeight), 2, color.White, true)

		// Afficher le texte du bouton
		btnTextX := buttonX + 40
		btnTextY := buttonY + 40
		text.Draw(screen, "Rejouer", fontFace, btnTextX, btnTextY, color.White)

		// Dessiner un rectangle arrondi
		vector.DrawFilledRect(screen, float32(boxX), float32(boxY), float32(boxWidth), float32(boxHeight), bgColor, true)

		// Afficher le texte du résultat
		textX := boxX + 110
		textY := boxY + 40
		text.Draw(screen, msg, fontFace, textX, textY, color.Black)

		// Afficher le message avec le temps
		textX = boxX + 40
		textY = boxY + 80
		text.Draw(screen, timerText, fontFace, textX, textY, color.Black)
	}

	// Dessiner le rectangle de sélection
	if !g.gameOver { // Afficher la sélection uniquement si le jeu n'est pas terminé
		selectedColor := color.RGBA{255, 0, 0, 255} // Rouge pour la sélection
		vector.StrokeRect(screen, float32(g.selectedCol*cellSize)+float32(gridLeft), float32(g.selectedRow*cellSize)+float32(topMargin), float32(cellSize), float32(cellSize), 2, selectedColor, true)
	}

}

// Cette fonction permet de relancer une nouvelle partie
func (g *Game) restartGame() {
	g.generateSudoku() // Générer une nouvelle grille
	g.gameOver = false
	g.win = false
	g.selectedRow, g.selectedCol = 0, 0
	g.startTime = time.Now() // Réinitialiser le chrono à zéro
}

// Cette fonction définit la taille de la fenêtre
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Fonction principale
func main() {
	icons, err := loadIcon()
	if err != nil {
		log.Fatalf("Erreur lors du chargement de l'icône: %v", err)
	}
	loadFont()
	game := &Game{}
	game.generateSudoku()
	game.startTime = time.Now()
	game.selectedRow, game.selectedCol = 0, 0

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sudoku")
	ebiten.SetWindowIcon(icons)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
