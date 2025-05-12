package main

import (
	"fmt"
	"sync"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Field struct {
	cells        [][]*Cell
	currentColor rl.Color
	clicksCount  int
	CanPlay      bool
	hoveredColor rl.Color

	// not updated fields
	FieldWidth  int32
	FieldHeight int32
	Size        int

	clicksLabel *Label
}

func (f *Field) Update() {

	f.hoveredColor = rl.Blank
	startPos := rl.Vector2{
		X: float32(rl.GetScreenWidth())/2 - float32(f.FieldWidth)/2,
		Y: float32(rl.GetScreenHeight())/2 - float32(f.FieldHeight)/2,
	}
	if startPos.Y < float32(f.clicksLabel.GetHeight()) {
		startPos.Y = float32(f.clicksLabel.GetHeight())
	}
	for i := int32(0); int(i) < f.Size; i++ {
		for j := int32(0); int(j) < f.Size; j++ {
			f.cells[i][j].Update(rl.Vector2{
				X: startPos.X + float32(CellWidth*i),
				Y: startPos.Y + float32(CellHeight*j),
			})
		}
	}
	f.clicksLabel.Text = fmt.Sprintf("Clicks %d/%d", f.clicksCount, getCurrentMaxClicks(f.Size))
	f.clicksLabel.ForegroundColor = f.currentColor

}
func (f *Field) Draw() {
	for i := 0; i < f.Size; i++ {
		for j := 0; j < f.Size; j++ {
			f.cells[i][j].Draw()
			clicked := f.cells[i][j].Clicked
			if clicked && f.CanPlay && f.cells[i][j].Color != f.currentColor {
				go updateCurrentColor(f, f.cells[i][j].Color)
				f.clicksCount++
			}
		}
	}
	f.clicksLabel.Draw()

}

func InitField(fieldSize int) *Field {
	field := &Field{
		clicksCount: 0,
		CanPlay:     true,
		cells:       make([][]*Cell, fieldSize),
		FieldWidth:  int32(fieldSize) * CellWidth,
		FieldHeight: int32(fieldSize) * CellHeight,
		clicksLabel: &Label{
			Position: rl.Vector2{0, 0},
			Padding:  rl.Vector2{X: 4, Y: 2},
			FontSize: 40,
		},
		Size: fieldSize,
	}
	rl.SetWindowMinSize(
		fieldSize*int(CellWidth),
		fieldSize*int(CellHeight)+int(field.clicksLabel.GetHeight()),
	)
	for i := 0; i < fieldSize; i++ {
		field.cells[i] = make([]*Cell, fieldSize)
		for j := 0; j < fieldSize; j++ {
			field.cells[i][j] = &Cell{
				Field: field,
				Color: getColor(),
				Rectangle: &rl.Rectangle{
					Width:  float32(CellWidth),
					Height: float32(CellHeight),
				},
			}
		}
	}
	field.currentColor = field.cells[0][0].Color
	return field
}

type Cell struct {
	*rl.Rectangle
	Color    rl.Color
	Hover    bool
	Captured bool
	Clicked  bool

	Field *Field
}

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
	if c.Captured && c.Field.CanPlay && c.Color != c.Field.currentColor {
		visibleColor.A -= 50
	}

	rl.DrawRectangleRec(*c.Rectangle, visibleColor)
	if c.Hover && c.Field.CanPlay && c.Color != c.Field.currentColor {
		rl.DrawRectangleLinesEx(*c.Rectangle, 1, rl.Black)
	}
}

const CellHeight int32 = 40
const CellWidth int32 = 40

func updateCellColor(f *Field, i, j int, newColor rl.Color, goLeft bool, goUp bool) {
	cell := f.cells[i][j]
	if cell.Color != f.currentColor {
		return
	}
	if cell.Color == newColor {
		return
	}

	f.cells[i][j].Color = newColor
	time.Sleep(20 * time.Millisecond)

	var wg sync.WaitGroup
	if i < f.Size-1 {
		wg.Add(1)
		go func() {
			updateCellColor(f, i+1, j, newColor, false, true)
			wg.Done()
		}()
	}
	if j < f.Size-1 {
		wg.Add(1)
		go func() {
			updateCellColor(f, i, j+1, newColor, true, false)
			wg.Done()
		}()
	}
	if i > 0 && goLeft {
		wg.Add(1)
		go func() {
			updateCellColor(f, i-1, j, newColor, true, true)
			wg.Done()
		}()
	}
	if j > 0 && goUp {
		wg.Add(1)
		go func() {
			updateCellColor(f, i, j-1, newColor, true, true)
			wg.Done()
		}()
	}
	wg.Wait()
}
