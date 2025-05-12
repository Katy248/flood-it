package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Label struct {
	Text            string
	Position        rl.Vector2
	Padding         rl.Vector2
	FontSize        int32
	ForegroundColor rl.Color
}

func (l *Label) Draw() {
	rl.DrawText(l.Text, int32(l.Position.X+l.Padding.X), int32(l.Position.Y+l.Padding.Y), l.FontSize, l.ForegroundColor)
}

func (l *Label) GetHeight() int32 {
	return l.FontSize + int32(l.Padding.Y)*2
}
