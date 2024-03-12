package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	Title         = "Tank Tank Tank"
	ScreenWidth   = 1280
	ScreenHeight  = 720
	TerrainHeight = 200
	TankSize      = 64
)

type Bullet struct {
	position rl.Vector2
	speed    rl.Vector2
	damage   float64
	size     float32
	isActive bool
}

type Tank struct {
	position rl.Vector2
	speed    float32
	colour   rl.Color
	bullets  []*Bullet
	isDead   bool
	health   float32

	aimPoint rl.Vector2
	aimPower float64
	aimAngle float64

	prevPoint rl.Vector2
	prevPower float64
	prevAngle float64
}

var (
	player Tank
	enemy  Tank

	mousePosition    rl.Vector2
	spawnEnemyButton rl.Rectangle
	spawnEnemy       bool
)

func initGame() {
	player.health = 100
	player.position = rl.Vector2{X: 340, Y: ScreenHeight - TerrainHeight}
	player.speed = 400
	player.colour = rl.DarkBlue
	player.bullets = make([]*Bullet, 0)

	enemy.health = 100
	enemy.position = rl.Vector2{X: 1000, Y: ScreenHeight - TerrainHeight}
	enemy.speed = 400
	enemy.colour = rl.Maroon

	spawnEnemyButton = rl.NewRectangle(1160, 20, 100, 30)
	spawnEnemy = false
}

func update() {
	mousePosition = rl.GetMousePosition()

	{ // spawn enemy button
		if rl.CheckCollisionPointRec(mousePosition, spawnEnemyButton) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				spawnEnemy = true
				enemy.health = 100
				enemy.isDead = false
			}
		}
	}

	{ // player movement
		if rl.IsKeyDown(rl.KeyJ) {
			player.position.X -= player.speed * rl.GetFrameTime()
		}

		if rl.IsKeyDown(rl.KeyL) {
			player.position.X += player.speed * rl.GetFrameTime()
		}
	}

	{ // player aiming
		player.aimPower = math.Sqrt(
			math.Pow(
				float64(player.position.X-mousePosition.X),
				2,
			) + math.Pow(
				float64(player.position.Y-mousePosition.Y),
				2,
			),
		)
		player.aimAngle = rl.Rad2deg * math.Asin(
			float64(player.position.Y-mousePosition.Y)/player.aimPower,
		)
		player.aimPoint = mousePosition
	}

	{ // shooting
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			player.prevPoint = player.aimPoint
			player.prevAngle = player.aimAngle
			player.prevPower = player.aimPower

			bullet := Bullet{position: player.position, damage: 25, size: 10, isActive: false}
			player.bullets = append(player.bullets, &bullet)
		}
	}

	{ // bullet movement
		if len(player.bullets) != 0 {
			for index, bullet := range player.bullets {
				if !bullet.isActive {
					bullet.speed.X = float32(
						math.Cos(rl.Deg2rad*player.prevAngle) * player.prevPower * 3 / 60,
					)
					bullet.speed.Y = float32(
						-math.Sin(rl.Deg2rad*player.prevAngle) * player.prevPower * 3 / 60,
					)
					bullet.isActive = true
				}
				bullet.position.X += bullet.speed.X
				bullet.position.Y += bullet.speed.Y
				bullet.speed.Y += 9.81 / 60

				// hit detection
				if rl.CheckCollisionCircleRec(
					bullet.position,
					bullet.size,
					rl.NewRectangle(enemy.position.X, enemy.position.Y, TankSize, TankSize),
				) && !enemy.isDead {
					enemy.health -= float32(bullet.damage)

					player.bullets = append(player.bullets[:index], player.bullets[index+1:]...)

					if enemy.health <= 0 {
						enemy.isDead = true
					}
				}
			}
		}
	}
}

func render() {
	rl.BeginDrawing()

	{ // terrain
		rl.ClearBackground(rl.SkyBlue)
		rl.DrawRectangleV(
			rl.Vector2{X: 0, Y: 520 + TankSize},
			rl.Vector2{X: ScreenWidth, Y: TerrainHeight},
			rl.DarkGreen,
		)
	}

	{ // spawn enemy button
		rl.DrawRectangleRec(spawnEnemyButton, rl.White)
	}

	{ // bullets
		for _, bullet := range player.bullets {
			rl.DrawCircleV(
				rl.Vector2{X: bullet.position.X, Y: bullet.position.Y},
				bullet.size,
				rl.Black,
			)
		}
	}

	{ // player and enemy rendering
		if !(player.isDead) {
			rl.DrawRectangleV(
				rl.Vector2{X: player.position.X, Y: player.position.Y - 20},
				rl.Vector2{X: TankSize, Y: 10},
				rl.Red,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: player.position.X, Y: player.position.Y - 20},
				rl.Vector2{X: TankSize * (player.health / 100), Y: 10},
				rl.Green,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: player.position.X, Y: player.position.Y},
				rl.Vector2{X: TankSize, Y: TankSize},
				player.colour,
			)
		}
		if !(enemy.isDead) {
			rl.DrawRectangleV(
				rl.Vector2{X: enemy.position.X, Y: enemy.position.Y - 20},
				rl.Vector2{X: TankSize, Y: 10},
				rl.Red,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: enemy.position.X, Y: enemy.position.Y - 20},
				rl.Vector2{X: TankSize * (enemy.health / 100), Y: 10},
				rl.Green,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: enemy.position.X, Y: enemy.position.Y},
				rl.Vector2{X: TankSize, Y: TankSize},
				enemy.colour,
			)
		}
	}

	{ // player aiming
		rl.DrawTriangle(
			rl.Vector2{X: player.position.X - TankSize/4, Y: player.position.Y - TankSize/4},
			rl.Vector2{X: player.position.X + TankSize/4, Y: player.position.Y + TankSize/4},
			player.aimPoint,
			rl.Gray,
		)
	}

	rl.EndDrawing()
}

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, Title)
	rl.SetTargetFPS(60)

	initGame()

	for !rl.WindowShouldClose() {
		update()
		render()
	}

	rl.CloseWindow()
}
