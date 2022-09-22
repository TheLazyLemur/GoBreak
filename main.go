package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	FPS             int32    = 60
	WindowWidth     int32    = 800
	WindowHeight    int32    = 600
	BackgroundColor rl.Color = rl.Black
	ProjSize        float32  = 25 * 0.8
	ProjSpeed       float32  = 350
	ProjColor       rl.Color = rl.White
	BarLen          float32  = 100
	BarY            float32  = float32(WindowHeight) - ProjSize - 50
	BarSpeed        float32  = ProjSpeed * 1.5
	BarColor        rl.Color = rl.Red
	TargetWidth     float32  = BarLen
	TargetHeight    float32  = ProjSize
	TargetPaddingX  int32    = 20
	TargetPaddingY  int32    = 50
	TargetRows      int32    = 4
	TargetCols      int32    = 5
	TargetGridWidth float32  = (float32(TargetCols)*TargetWidth + (float32(TargetCols)-1)*float32(TargetPaddingX))
	TargetGridX     float32  = float32(WindowWidth)/2 - TargetGridWidth/2
	TargetGridY     int32    = 50
	TargetColor     rl.Color = rl.Green
)

type Target struct {
	x    float32
	y    float32
	dead bool
}

type ProjVel struct {
	x float32
	y float32
}

func initTargets() []Target {
	targets := make([]Target, TargetRows*TargetCols)
	for row := int32(0); row < TargetRows; row++ {
		for col := int32(0); col < TargetCols; col++ {
			targets[row*TargetCols+col].x = (TargetGridX + (float32(TargetWidth)+float32(TargetPaddingX))*float32(col))
			targets[row*TargetCols+col].y = (float32(TargetGridY) + float32(TargetPaddingY)*float32(row))
			targets[row*TargetCols+col].dead = false
		}
	}

	return targets
}

func main() {
	var shouldExit bool

	targets := initTargets()

	var pause bool = false
	var barX float32 = float32(WindowWidth)/2 - BarLen/2
	var barXVel float32 = 0
	var gameStarted bool = false
	var projRec rl.Rectangle
	var projVel ProjVel
	projVel.x = 0
	projVel.y = -1
	projRec.Y = BarY - ProjSize
	projRec.X = float32(WindowWidth)/2 - ProjSize/2

	rl.InitWindow(WindowWidth, WindowHeight, "GoBreaker - Example game")
	rl.SetTargetFPS(FPS)

	for !rl.WindowShouldClose() && !shouldExit {
		rl.BeginDrawing()

		rl.ClearBackground(BackgroundColor)

		for _, target := range targets {
			if !target.dead {
				targetRec := rl.Rectangle{X: target.x, Y: target.y, Width: TargetWidth, Height: TargetHeight}
				rl.DrawRectangleRec(targetRec, TargetColor)
			}
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			pause = !pause
		}

		if rl.IsKeyDown(rl.KeyRight) && !pause {
			barXVel = 1
			if !gameStarted {
				gameStarted = true
				projVel.x = 1
			}
		}

		if rl.IsKeyDown(rl.KeyLeft) && !pause {
			barXVel = -1
			if !gameStarted {
				gameStarted = true
				projVel.x = -1
			}
		}

		if !rl.IsKeyDown(rl.KeyLeft) && !rl.IsKeyDown(rl.KeyRight) {
			barXVel = 0
		}

		barX += barXVel * BarSpeed * rl.GetFrameTime()
		if barX <= 0 {
			barX = 0
		}
		if barX >= float32(WindowWidth)-BarLen {
			barX = float32(WindowWidth) - BarLen
		}

		playerRec := rl.Rectangle{X: barX, Y: BarY, Width: BarLen, Height: 20}
		rl.DrawRectangleRec(playerRec, BarColor)

		if gameStarted && !pause {
			projRec.Y += projVel.y * ProjSpeed * rl.GetFrameTime()
			projRec.X += projVel.x * ProjSpeed * rl.GetFrameTime()
		}

		if projRec.X >= float32(WindowWidth)-ProjSize || projRec.X <= float32(0) {
			projVel.x *= -1
		}

		if projRec.Y < float32(0) {
			projVel.y *= -1
		}

		if projRec.Y > float32(WindowHeight) {
			shouldExit = true
		}

		projRec = rl.Rectangle{X: projRec.X, Y: projRec.Y, Width: ProjSize, Height: ProjSize}

		for i, target := range targets {
			targetRec := rl.Rectangle{X: target.x, Y: target.y, Width: TargetWidth, Height: TargetHeight}

			if rl.CheckCollisionRecs(targetRec, projRec) && !target.dead {
				projVel.y *= -1
				targets[i].dead = true
			}
		}

		if rl.CheckCollisionRecs(playerRec, projRec) {
			projVel.y *= -1
			if barXVel != 0 {
				projVel.x = barXVel
			}
		}

		rl.DrawRectangleRec(projRec, ProjColor)

		rl.EndDrawing()
	}

	rl.CloseWindow()
}
