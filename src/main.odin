package main

import "core:fmt"
import rl "vendor:raylib"

SCREEN_WIDTH :: 1280
SCREEN_HEIGHT :: 720

TANK_SIZE :: 64
TERRAIN_HEIGHT :: 200

Bullet :: struct {
	rect:  rl.Rectangle,
	speed: f32,
}

Tank :: struct {
	rect:    rl.Rectangle,
	speed:   f32,
	colour:  rl.Color,
	bullets: [dynamic]Bullet,
	is_dead: bool,
}

player: Tank
enemy: Tank

init :: proc() {
	player.rect.x = 340
	player.rect.y = SCREEN_HEIGHT - TERRAIN_HEIGHT
	player.speed = 400
	player.colour = rl.DARKBLUE
	player.bullets = make([dynamic]Bullet, 0, 1)

	enemy.rect.x = 1000
	enemy.rect.y = SCREEN_HEIGHT - TERRAIN_HEIGHT
	enemy.speed = 400
	enemy.colour = rl.MAROON
}

update :: proc() {
	// player movement
	if rl.IsKeyDown(.J) {
		player.rect.x -= player.speed * rl.GetFrameTime()
	}
	if rl.IsKeyDown(.L) {
		player.rect.x += player.speed * rl.GetFrameTime()
	}

	// bullet movement
	if len(player.bullets) != 0 {
		player.bullets[0].rect.x += player.bullets[0].speed * rl.GetFrameTime()

		if player.bullets[0].rect.x >= enemy.rect.x {
			fmt.println("hit")
			enemy.is_dead = true
			pop(&player.bullets)
		}
	}

	if rl.IsMouseButtonPressed(.LEFT) {
		bullet_rect := rl.Rectangle {
			player.rect.x + TANK_SIZE,
			SCREEN_HEIGHT - (TERRAIN_HEIGHT),
			10,
			10,
		}
		bullet := Bullet{bullet_rect, 200}
		append(&player.bullets, bullet)
	}
}

render :: proc() {
	rl.BeginDrawing()

	// terrain
	rl.ClearBackground(rl.SKYBLUE)
	rl.DrawRectangleV({0, 520 + TANK_SIZE}, {SCREEN_WIDTH, TERRAIN_HEIGHT}, rl.DARKGREEN)

	if len(player.bullets) != 0 {
		rl.DrawCircleV({player.bullets[0].rect.x, player.bullets[0].rect.y}, 10, rl.BLACK)
	}

	// player and enemy rendering
	if !(player.is_dead) {
		rl.DrawRectangleV({player.rect.x, player.rect.y}, {TANK_SIZE, TANK_SIZE}, player.colour)
	}
	if !(enemy.is_dead) {
		rl.DrawRectangleV({enemy.rect.x, enemy.rect.y}, {TANK_SIZE, TANK_SIZE}, enemy.colour)
	}

	rl.EndDrawing()
}

main :: proc() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "tank tank tank")
	rl.SetTargetFPS(144)

	init()

	for !rl.WindowShouldClose() {
		update()
		render()
	}

	rl.CloseWindow()
}
