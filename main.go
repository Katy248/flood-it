package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var mousePosition rl.Vector2
var hoveredColor rl.Color

const (
	FieldSize1 = 6
	FieldSize2 = 12
	FieldSize3 = 24

	MaxCounter1 = 10
	MaxCounter2 = 22
	MaxCounter3 = 42
)

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(800, 600, "Flood-it")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	field := InitField(FieldSize2)

	for !rl.WindowShouldClose() {
		if rl.IsKeyReleased(rl.KeyQ) {
			break
		}
		if rl.IsKeyReleased(rl.KeyOne) {
			field = InitField(FieldSize1)
		}
		if rl.IsKeyReleased(rl.KeyTwo) {
			field = InitField(FieldSize2)
		}
		if rl.IsKeyReleased(rl.KeyThree) {
			field = InitField(FieldSize3)
		}
		if rl.IsKeyReleased(rl.KeyR) {
			field = InitField(field.Size)
		}

		mousePosition = rl.GetMousePosition()
		field.Update()
		rl.BeginDrawing()
		{
			rl.ClearBackground(ColorWindowBG)
			field.Draw()
			if checkWin(field) {
				field.CanPlay = false
				drawOverlayLabel("You win!", ColorWindowBG, ColorErrorFG)
			} else if checkLoose(field) {
				field.CanPlay = false
				drawOverlayLabel("You lose", ColorErrorBG, ColorErrorFG)
			}
		}
		rl.EndDrawing()

	}
}

const (
	OverlayLabelFontSize = 40
	OverlayLabelPadding  = 10
)

func drawOverlayLabel(text string, bg rl.Color, fg rl.Color) {
	textLength := rl.MeasureText(text, OverlayLabelFontSize)
	rl.DrawRectangleRec(rl.Rectangle{
		X:      float32(int32(rl.GetScreenWidth())/2-textLength/2) - OverlayLabelPadding,
		Y:      float32(rl.GetScreenHeight())/2 - OverlayLabelPadding,
		Width:  float32(textLength) + OverlayLabelPadding*2,
		Height: float32(OverlayLabelFontSize) + OverlayLabelPadding*2,
	}, bg)
	rl.DrawText(
		text,
		int32(rl.GetScreenWidth())/2-textLength/2, int32(rl.GetScreenHeight())/2,
		OverlayLabelFontSize,
		fg,
	)
}

func updateCurrentColor(f *Field, newColor rl.Color) {
	f.CanPlay = false
	updateCellColor(f, 0, 0, newColor, false, false)
	f.currentColor = newColor
	f.CanPlay = true
}

func checkWin(f *Field) bool {
	for i := 0; i < f.Size; i++ {
		for j := 0; j < f.Size; j++ {
			if f.cells[i][j].Color != f.currentColor {
				return false
			}
		}
	}
	return true
}

func checkLoose(f *Field) bool {
	maxClicks := getCurrentMaxClicks(f.Size)
	return f.clicksCount >= maxClicks
}

func getCurrentMaxClicks(fieldSize int) int {
	switch fieldSize {
	case FieldSize1:
		return MaxCounter1
	case FieldSize2:
		return MaxCounter2
	case FieldSize3:
		return MaxCounter3
	default:
		return 0
	}
}

var (
	ColorRed    = rl.Color{224, 27, 36, 255}
	ColorOrange = rl.Color{255, 120, 0, 255}
	ColorYellow = rl.Color{246, 211, 45, 255}
	ColorGreen  = rl.Color{51, 209, 122, 255}
	ColorBlue   = rl.Color{53, 132, 228, 255}
	ColorPurple = rl.Color{145, 65, 172, 255}

	ColorWindowFG = rl.NewColor(255, 255, 255, 255)
	ColorWindowBG = rl.NewColor(34, 34, 38, 255) // rl.Color{36, 31, 49, 255}
	ColorErrorBG  = rl.NewColor(192, 28, 40, 255)
	ColorErrorFG  = rl.NewColor(255, 255, 255, 255)
)
var colors = []rl.Color{
	ColorRed,
	ColorOrange,
	ColorYellow,
	ColorGreen,
	ColorBlue,
	ColorPurple,
}

func getColor() rl.Color {
	index := rl.GetRandomValue(0, int32(len(colors)-1))
	return colors[index]
}
