package main

import "core:fmt"
import rl "vendor:raylib"

SCREEN_WIDTH :: 1280
SCREEN_HEIGHT :: 720

TANK_SIZE :: 64
TERRAIN_HEIGHT :: 200

Bullet :: struct {
	rect:   rl.Rectangle,
	damage: f32,
	speed:  f32,
}

Tank :: struct {
	rect:    rl.Rectangle,
	speed:   f32,
	colour:  rl.Color,
	bullets: [dynamic]Bullet,
	is_dead: bool,
	health:  f32,
}

player: Tank
enemy: Tank

init :: proc() {
	player.health = 100
	player.rect.x = 340
	player.rect.y = SCREEN_HEIGHT - TERRAIN_HEIGHT
	player.speed = 400
	player.colour = rl.DARKBLUE
	player.bullets = make([dynamic]Bullet)

	enemy.health = 100
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
		for &bullet, index in player.bullets {
			bullet.rect.x += bullet.speed * rl.GetFrameTime()

			if bullet.rect.x >= enemy.rect.x && !enemy.is_dead {
				enemy.health -= bullet.damage

				if enemy.health <= 0 {
					enemy.is_dead = true
				}
				ordered_remove(&player.bullets, index)
			}

			if bullet.rect.x > SCREEN_WIDTH {
				ordered_remove(&player.bullets, index)
			}
		}
	}

	if rl.IsMouseButtonPressed(.LEFT) {
		bullet_rect := rl.Rectangle {
			player.rect.x + TANK_SIZE,
			SCREEN_HEIGHT - (TERRAIN_HEIGHT),
			10,
			10,
		}
		bullet := Bullet{bullet_rect, 25, 300}
		append(&player.bullets, bullet)
	}
}

render :: proc() {
	rl.BeginDrawing()

	// terrain
	rl.ClearBackground(rl.SKYBLUE)
	rl.DrawRectangleV({0, 520 + TANK_SIZE}, {SCREEN_WIDTH, TERRAIN_HEIGHT}, rl.DARKGREEN)

	for bullet in player.bullets {
		rl.DrawCircleV({bullet.rect.x, bullet.rect.y}, 10, rl.BLACK)
	}

	// player and enemy rendering
	if !(player.is_dead) {
		rl.DrawRectangleV({player.rect.x, player.rect.y - 20}, {TANK_SIZE, 10}, rl.RED)
		rl.DrawRectangleV(
			{player.rect.x, player.rect.y - 20},
			{TANK_SIZE * (player.health / 100), 10},
			rl.GREEN,
		)
		rl.DrawRectangleV({player.rect.x, player.rect.y}, {TANK_SIZE, TANK_SIZE}, player.colour)
	}
	if !(enemy.is_dead) {
		rl.DrawRectangleV({enemy.rect.x, enemy.rect.y - 20}, {TANK_SIZE, 10}, rl.RED)
		rl.DrawRectangleV(
			{enemy.rect.x, enemy.rect.y - 20},
			{TANK_SIZE * (enemy.health / 100), 10},
			rl.GREEN,
		)
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
