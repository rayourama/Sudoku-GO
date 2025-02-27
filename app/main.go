package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
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
	cellSize     = screenWidth / 9
	emptyCells   = 40
)

type Game struct {
	grid         [9][9]int
	solution     [9][9]int
	originalGrid [9][9]int
	selectedRow  int
	selectedCol  int
	lastMoveTime time.Time
}

var fontFace font.Face
var moveCooldown = 150 * time.Millisecond

func loadFont() {
	ttfBytes, err := ioutil.ReadFile("font/Roboto-Bold.ttf")
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

func (g *Game) generateSudoku() {
	generateFullSudoku(&g.solution)
	g.grid = g.solution
	removeNumbers(&g.grid, emptyCells)
	g.originalGrid = g.grid
}

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

func (g *Game) setCellValue(value int) {
	if g.originalGrid[g.selectedRow][g.selectedCol] == 0 { // Vérifie que la cellule est modifiable
		currentValue := g.grid[g.selectedRow][g.selectedCol]

		// Effacer si Backspace est pressé ou si le même numéro est entré
		if value == 0 || abs(currentValue) == value {
			g.grid[g.selectedRow][g.selectedCol] = 0
		} else if value == g.solution[g.selectedRow][g.selectedCol] {
			g.grid[g.selectedRow][g.selectedCol] = value // Correct -> affiché en vert
		} else {
			g.grid[g.selectedRow][g.selectedCol] = -value // Faux -> affiché en rouge
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{240, 240, 240, 255})
	lineColor := color.RGBA{0, 0, 0, 255}

	for i := 0; i <= 9; i++ {
		thickness := 1.0
		if i%3 == 0 {
			thickness = 3.0
		}
		vector.StrokeLine(screen, float32(i*cellSize), 0, float32(i*cellSize), float32(screenHeight), float32(thickness), lineColor, true)
		vector.StrokeLine(screen, 0, float32(i*cellSize), float32(screenWidth), float32(i*cellSize), float32(thickness), lineColor, true)
	}

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
				}
				text.Draw(screen, num, fontFace, x*cellSize+30, y*cellSize+50, colorToUse)
			}
		}
	}
	selectedColor := color.RGBA{255, 0, 0, 255} // Rouge pour la sélection
	vector.StrokeRect(screen, float32(g.selectedCol*cellSize), float32(g.selectedRow*cellSize), float32(cellSize), float32(cellSize), 2, selectedColor, true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	loadFont()
	game := &Game{}
	game.generateSudoku()
	game.selectedRow, game.selectedCol = 0, 0

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Sudoku")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func generateFullSudoku(grid *[9][9]int) {
	rand.Seed(time.Now().UnixNano())
	fillGrid(grid)
}

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

func isValid(grid *[9][9]int, row, col, num int) bool {
	for i := 0; i < 9; i++ {
		if grid[row][i] == num || grid[i][col] == num || grid[row/3*3+i/3][col/3*3+i%3] == num {
			return false
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
