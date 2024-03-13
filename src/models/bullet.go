package models

import rl "github.com/gen2brain/raylib-go/raylib"

type Bullet struct {
	Position rl.Vector2
	Speed    rl.Vector2
	Damage   float64
	Size     float32
	IsActive bool
}
