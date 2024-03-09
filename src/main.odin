package main

import "core:fmt"
import "core:math"
import rl "vendor:raylib"

SCREEN_WIDTH :: 1280
SCREEN_HEIGHT :: 720

TANK_SIZE :: 64
TERRAIN_HEIGHT :: 200

Bullet :: struct {
	position:  rl.Vector2,
	speed:     rl.Vector2,
	damage:    f32,
	is_active: bool,
}

Tank :: struct {
	position:   rl.Vector2,
	speed:      f32,
	colour:     rl.Color,
	bullets:    [dynamic]Bullet,
	is_dead:    bool,
	health:     f32,
	aim_point:  rl.Vector2,
	aim_angle:  f32,
	aim_power:  f32,
	prev_point: rl.Vector2,
	prev_angle: f32,
	prev_power: f32,
}

player: Tank
enemy: Tank

init :: proc() {
	player.health = 100
	player.position = {340, SCREEN_HEIGHT - TERRAIN_HEIGHT}
	player.speed = 400
	player.colour = rl.DARKBLUE
	player.bullets = make([dynamic]Bullet)

	enemy.health = 100
	enemy.position = {1000, SCREEN_HEIGHT - TERRAIN_HEIGHT}
	enemy.speed = 400
	enemy.colour = rl.MAROON
}

update :: proc() {
	// player movement
	if rl.IsKeyDown(.J) {
		player.position.x -= player.speed * rl.GetFrameTime()
	}
	if rl.IsKeyDown(.L) {
		player.position.x += player.speed * rl.GetFrameTime()
	}

	{ 	// player aiming
		player.aim_power = math.sqrt(
			math.pow(player.position.x - rl.GetMousePosition().x, 2) +
			math.pow(player.position.y - rl.GetMousePosition().y, 2),
		)
		player.aim_angle = math.to_degrees_f32(
			math.asin((player.position.y - rl.GetMousePosition().y) / player.aim_power),
		)
		player.aim_point = rl.GetMousePosition()
	}

	if rl.IsMouseButtonPressed(.LEFT) {
		player.prev_point = player.aim_point
		player.prev_power = player.aim_power
		player.prev_angle = player.aim_angle

		bullet := Bullet{player.position, {1, 1}, 25, false}
		bullet.position = player.position
		append(&player.bullets, bullet)
	}

	// bullet movement
	if len(player.bullets) != 0 {
		for &bullet, index in player.bullets {
			if !bullet.is_active {
				bullet.speed.x =
					math.cos(math.to_radians_f32(player.prev_angle)) * player.prev_power * 3 / 60
				bullet.speed.y =
					-math.sin(math.to_radians_f32(player.prev_angle)) * player.prev_power * 3 / 60
				bullet.is_active = true
			}
			bullet.position.x += bullet.speed.x
			bullet.position.y += bullet.speed.y
			bullet.speed.y += 9.81 / 60
		}
	}
}

render :: proc() {
	rl.BeginDrawing()

	// terrain
	rl.ClearBackground(rl.SKYBLUE)
	rl.DrawRectangleV({0, 520 + TANK_SIZE}, {SCREEN_WIDTH, TERRAIN_HEIGHT}, rl.DARKGREEN)

	for bullet in player.bullets {
		rl.DrawCircleV({bullet.position.x, bullet.position.y}, 10, rl.BLACK)
	}

	// player and enemy rendering
	if !(player.is_dead) {
		rl.DrawRectangleV({player.position.x, player.position.y - 20}, {TANK_SIZE, 10}, rl.RED)
		rl.DrawRectangleV(
			{player.position.x, player.position.y - 20},
			{TANK_SIZE * (player.health / 100), 10},
			rl.GREEN,
		)
		rl.DrawRectangleV(
			{player.position.x, player.position.y},
			{TANK_SIZE, TANK_SIZE},
			player.colour,
		)
	}
	if !(enemy.is_dead) {
		rl.DrawRectangleV({enemy.position.x, enemy.position.y - 20}, {TANK_SIZE, 10}, rl.RED)
		rl.DrawRectangleV(
			{enemy.position.x, enemy.position.y - 20},
			{TANK_SIZE * (enemy.health / 100), 10},
			rl.GREEN,
		)
		rl.DrawRectangleV(
			{enemy.position.x, enemy.position.y},
			{TANK_SIZE, TANK_SIZE},
			enemy.colour,
		)
	}

	// draw aim
	rl.DrawTriangle(
		{player.position.x - TANK_SIZE / 4, player.position.y - TANK_SIZE / 4},
		{player.position.x + TANK_SIZE / 4, player.position.y + TANK_SIZE / 4},
		player.aim_point,
		rl.GRAY,
	)

	rl.EndDrawing()
}

main :: proc() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "tank tank tank")
	rl.SetTargetFPS(60)

	init()

	for !rl.WindowShouldClose() {
		update()
		render()
	}

	rl.CloseWindow()
}
