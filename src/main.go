package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	models "github.com/joshuahicks/tank-tank-tank/src/models"
)

const (
	Title         = "Tank Tank Tank"
	ScreenWidth   = 1280
	ScreenHeight  = 720
	TerrainHeight = 200
	TankSize      = 64
)

var (
	player models.Tank
	enemy  models.Tank

	mousePosition    rl.Vector2
	spawnEnemyButton rl.Rectangle
	spawnEnemy       bool
)

func initGame() {
	player.Health = 100
	player.Position = rl.Vector2{X: 340, Y: ScreenHeight - TerrainHeight}
	player.Speed = 400
	player.Colour = rl.DarkBlue
	player.Bullets = make([]*models.Bullet, 0)

	enemy.Health = 100
	enemy.Position = rl.Vector2{X: 1000, Y: ScreenHeight - TerrainHeight}
	enemy.Speed = 400
	enemy.Colour = rl.Maroon

	spawnEnemyButton = rl.NewRectangle(1160, 20, 100, 30)
}

func update() {
	mousePosition = rl.GetMousePosition()

	{ // spawn enemy button
		if rl.CheckCollisionPointRec(mousePosition, spawnEnemyButton) {
			if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
				enemy.Health = 100
				enemy.IsDead = false
			}
		}
	}

	{ // player movement
		if rl.IsKeyDown(rl.KeyJ) {
			player.Position.X -= player.Speed * rl.GetFrameTime()
		}

		if rl.IsKeyDown(rl.KeyL) {
			player.Position.X += player.Speed * rl.GetFrameTime()
		}
	}

	{ // player aiming
		player.AimPower = math.Sqrt(
			math.Pow(
				float64(player.Position.X-mousePosition.X),
				2,
			) + math.Pow(
				float64(player.Position.Y-mousePosition.Y),
				2,
			),
		)
		player.AimAngle = rl.Rad2deg * math.Asin(
			float64(player.Position.Y-mousePosition.Y)/player.AimPower,
		)
		player.AimPoint = mousePosition
	}

	{ // shooting
		if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
			player.PrevPoint = player.AimPoint
			player.PrevAngle = player.AimAngle
			player.PrevPower = player.AimPower

			bullet := models.Bullet{
				Position: player.Position,
				Damage:   25,
				Size:     10,
				IsActive: false,
			}
			player.Bullets = append(player.Bullets, &bullet)
		}
	}

	{ // bullet movement
		if len(player.Bullets) != 0 {
			for index, bullet := range player.Bullets {
				if !bullet.IsActive {
					bullet.Speed.X = float32(
						math.Cos(rl.Deg2rad*player.PrevAngle) * player.PrevPower * 3 / 60,
					)
					bullet.Speed.Y = float32(
						-math.Sin(rl.Deg2rad*player.PrevAngle) * player.PrevPower * 3 / 60,
					)
					bullet.IsActive = true
				}
				bullet.Position.X += bullet.Speed.X
				bullet.Position.Y += bullet.Speed.Y
				bullet.Speed.Y += 9.81 / 60

				// hit detection
				if rl.CheckCollisionCircleRec(
					bullet.Position,
					bullet.Size,
					rl.NewRectangle(enemy.Position.X, enemy.Position.Y, TankSize, TankSize),
				) && !enemy.IsDead {
					enemy.Health -= float32(bullet.Damage)

					player.Bullets = append(player.Bullets[:index], player.Bullets[index+1:]...)

					if enemy.Health <= 0 {
						enemy.IsDead = true
					}
				}
			}
		}
	}
}

func render() {
	rl.BeginDrawing()
	defer rl.EndDrawing()

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
		for _, bullet := range player.Bullets {
			rl.DrawCircleV(
				rl.Vector2{X: bullet.Position.X, Y: bullet.Position.Y},
				bullet.Size,
				rl.Black,
			)
		}
	}

	{ // player and enemy rendering
		if !(player.IsDead) {
			rl.DrawRectangleV(
				rl.Vector2{X: player.Position.X, Y: player.Position.Y - 20},
				rl.Vector2{X: TankSize, Y: 10},
				rl.Red,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: player.Position.X, Y: player.Position.Y - 20},
				rl.Vector2{X: TankSize * (player.Health / 100), Y: 10},
				rl.Green,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: player.Position.X, Y: player.Position.Y},
				rl.Vector2{X: TankSize, Y: TankSize},
				player.Colour,
			)
		}
		if !(enemy.IsDead) {
			rl.DrawRectangleV(
				rl.Vector2{X: enemy.Position.X, Y: enemy.Position.Y - 20},
				rl.Vector2{X: TankSize, Y: 10},
				rl.Red,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: enemy.Position.X, Y: enemy.Position.Y - 20},
				rl.Vector2{X: TankSize * (enemy.Health / 100), Y: 10},
				rl.Green,
			)
			rl.DrawRectangleV(
				rl.Vector2{X: enemy.Position.X, Y: enemy.Position.Y},
				rl.Vector2{X: TankSize, Y: TankSize},
				enemy.Colour,
			)
		}
	}

	{ // player aiming
		rl.DrawTriangle(
			rl.Vector2{X: player.Position.X - TankSize/4, Y: player.Position.Y - TankSize/4},
			rl.Vector2{X: player.Position.X + TankSize/4, Y: player.Position.Y + TankSize/4},
			player.AimPoint,
			rl.Gray,
		)
	}
}

func main() {
	rl.InitWindow(ScreenWidth, ScreenHeight, Title)
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	initGame()

	for !rl.WindowShouldClose() {
		update()
		render()
	}
}
