package main

import (
	"fmt"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var mousePosition rl.Vector2
var currentColor rl.Color
var clicksCount = 0
var hoveredColor rl.Color

var canPlay = false

func main() {
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(800, 600, "Flood-it")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	initField()

	fontSize := int32(40)
	for !rl.WindowShouldClose() {
		if rl.IsKeyReleased(rl.KeyQ) {
			break
		}
		mousePosition = rl.GetMousePosition()
		updateField()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		drawField()
		if checkWin() {
			canPlay = false
			rl.DrawRectangleRec(rl.Rectangle{
				X:      float32(int32(rl.GetScreenWidth())/2 - rl.MeasureText("You win!", 40)/2),
				Y:      float32(rl.GetScreenHeight()) / 2,
				Width:  float32(rl.MeasureText("You win!", fontSize)),
				Height: float32(fontSize),
			}, rl.RayWhite)
			rl.DrawText(
				"You win!", int32(rl.GetScreenWidth())/2-rl.MeasureText("You win!", 40)/2, int32(rl.GetScreenHeight())/2, fontSize, rl.Red,
			)
		}
		rl.EndDrawing()

	}
}

func initField() {
	field = make([][]*Cell, FieldSize)
	for i := 0; i < FieldSize; i++ {
		field[i] = make([]*Cell, FieldSize)
		for j := 0; j < FieldSize; j++ {
			field[i][j] = &Cell{
				Color:     getColor(),
				Rectangle: &rl.Rectangle{Width: CellWidth, Height: CellHeight},
			}
			fmt.Println(field[i][j])
		}
	}
	currentColor = field[0][0].Color
	clicksCount = 0
	canPlay = true
}

type Cell struct {
	*rl.Rectangle
	Color    rl.Color
	Hover    bool
	Captured bool
	Clicked  bool
}

var field [][]*Cell

func (c *Cell) Update(pos rl.Vector2) {
	c.Hover = rl.CheckCollisionPointRec(mousePosition, *c.Rectangle)
	if c.Hover {
		hoveredColor = c.Color
	}
	c.Clicked = c.Hover && rl.IsMouseButtonReleased(rl.MouseButtonLeft)
	c.Captured = c.Hover && rl.IsMouseButtonDown(rl.MouseButtonLeft)
	c.Rectangle.X = pos.X
	c.Rectangle.Y = pos.Y
}

func (c *Cell) Draw() {
	visibleColor := c.Color

	// rect := rl.Rectangle{X: float32(x), Y: float32(y), Width: CellWidth, Height: CellHeight}
	if c.Captured && canPlay && c.Color != currentColor {
		visibleColor.A -= 50
	}

	rl.DrawRectangleRec(*c.Rectangle, visibleColor)
	if c.Hover && canPlay && c.Color != currentColor {
		rl.DrawRectangleLinesEx(*c.Rectangle, 1, rl.Black)
	}
	// if c.Hover && canPlay && c.Color != currentColor || hoveredColor == c.Color {
	// 	rl.DrawCircle(int32(c.Rectangle.X+c.Rectangle.Width/2), int32(c.Rectangle.Y+c.Rectangle.Height/2), 5, rl.Black)
	// }
}

const (
	FieldSize1 = 6
	FieldSize2 = 12
	FieldSize3 = 24
)

var FieldSize = FieldSize2

const CellHeight = 40
const CellWidth = 40

var FieldWidth int
var FieldHeight int

func updateField() {
	var FieldWidth = FieldSize * CellWidth
	var FieldHeight = FieldSize * CellHeight
	if rl.IsKeyReleased(rl.KeyR) {
		initField()
		return
	}

	if rl.IsKeyReleased(rl.KeyOne) {
		FieldSize = FieldSize1
		initField()
	}
	if rl.IsKeyReleased(rl.KeyTwo) {
		FieldSize = FieldSize2
		initField()
	}
	if rl.IsKeyReleased(rl.KeyThree) {
		FieldSize = FieldSize3
		initField()
	}

	hoveredColor = rl.Blank
	startPos := rl.Vector2{
		X: float32(rl.GetScreenWidth())/2 - float32(FieldWidth)/2,
		Y: float32(rl.GetScreenHeight())/2 - float32(FieldHeight)/2,
	}
	for i := 0; i < FieldSize; i++ {
		for j := 0; j < FieldSize; j++ {
			field[i][j].Update(rl.Vector2{
				X: float32(int(startPos.X) + CellWidth*i),
				Y: float32(startPos.Y) + float32(CellHeight*j),
			})
		}
	}
}

func drawField() {
	for i := 0; i < FieldSize; i++ {
		for j := 0; j < FieldSize; j++ {
			field[i][j].Draw()
			clicked := field[i][j].Clicked
			if clicked && canPlay && field[i][j].Color != currentColor {
				go updateCurrentColor(field[i][j].Color)
				clicksCount++
			}
		}
	}
	rl.DrawText(fmt.Sprintf("Clicks %d", clicksCount), 0, 0, 40, currentColor)
}
func updateCurrentColor(newColor rl.Color) {
	canPlay = false
	updateCellColor(0, 0, newColor, false, false)
	currentColor = newColor
	canPlay = true
}

func checkWin() bool {
	for i := 0; i < FieldSize; i++ {
		for j := 0; j < FieldSize; j++ {
			if field[i][j].Color != currentColor {
				return false
			}
		}
	}

	return true

}

func updateCellColor(i, j int, newColor rl.Color, goLeft bool, goUp bool) {
	cell := field[i][j]
	if cell.Color != currentColor {
		return
	}
	if cell.Color == newColor {
		return
	}

	field[i][j].Color = newColor
	time.Sleep(20 * time.Millisecond)

	var wg sync.WaitGroup
	if i < FieldSize-1 {
		wg.Add(1)
		go func() {
			updateCellColor(i+1, j, newColor, false, true)
			wg.Done()
		}()
	}
	if j < FieldSize-1 {
		wg.Add(1)
		go func() {
			updateCellColor(i, j+1, newColor, true, false)
			wg.Done()
		}()
	}
	if i > 0 && goLeft {
		wg.Add(1)
		go func() {
			updateCellColor(i-1, j, newColor, true, true)
			wg.Done()
		}()
	}
	if j > 0 && goUp {
		wg.Add(1)
		go func() {
			updateCellColor(i, j-1, newColor, true, true)
			wg.Done()
		}()
	}
	wg.Wait()

}

var colors = []rl.Color{
	rl.Red,
	rl.Yellow,
	rl.Orange,
	rl.Green,
	rl.Blue,
	rl.Violet,
}

func getColor() rl.Color {
	index := rl.GetRandomValue(0, int32(len(colors)-1))
	return colors[index]
}
