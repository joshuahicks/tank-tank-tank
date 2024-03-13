package models

import rl "github.com/gen2brain/raylib-go/raylib"

type Tank struct {
	Position rl.Vector2
	Speed    float32
	Colour   rl.Color
	Bullets  []*Bullet
	IsDead   bool
	Health   float32

	AimPoint rl.Vector2
	AimPower float64
	AimAngle float64

	PrevPoint rl.Vector2
	PrevPower float64
	PrevAngle float64
}
