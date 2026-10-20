package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"
)

const (
	Width  = 40
	Height = 15
)

type Paddle struct {
	Y, Length int
}

type Ball struct {
	X, Y, DirX, DirY int
}

// Очистка экрана стандартными средствами ОС
func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	padLength := 3
	p1 := Paddle{Y: Height/2 - padLength/2, Length: padLength}
	p2 := Paddle{Y: Height/2 - padLength/2, Length: padLength}
	
	ball := Ball{
		X:    Width / 2,
		Y:    Height / 2,
		DirX: 1,
		DirY: 1,
	}
	
	score1, score2 := 0, 0

	// Упрощенный неблокирующий ввод для демонстрации в консоли.
	// Так как стандартный os.Stdin блокирует поток, сделаем авто-движение ракетки Игрока 2 (ИИ),
	// а Игрок 1 будет двигаться автоматически для тестов сборки, либо игра будет просто идти сама.
	
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		<-ticker.C
		
		// ЛОГИКА ДВИЖЕНИЯ МЯЧА
		ball.X += ball.DirX
		ball.Y += ball.DirY

		// Простейший ИИ для ракеток, чтобы они следили за мячом (играет сама в себя)
		if ball.Y > p1.Y+p1.Length/2 && p1.Y < Height-p1.Length-1 {
			p1.Y++
		} else if ball.Y < p1.Y+p1.Length/2 && p1.Y > 1 {
			p1.Y--
		}
		
		if ball.Y > p2.Y+p2.Length/2 && p2.Y < Height-p2.Length-1 {
			p2.Y++
		} else if ball.Y < p2.Y+p2.Length/2 && p2.Y > 1 {
			p2.Y--
		}

		// Отскок от верхней и нижней границ
		if ball.Y <= 1 || ball.Y >= Height-2 {
			ball.DirY = -ball.DirY
		}

		// Отскок от левой ракетки
		if ball.X == 2 && ball.Y >= p1.Y && ball.Y < p1.Y+p1.Length {
			ball.DirX = -ball.DirX
		}

		// Отскок от правой ракетки
		if ball.X == Width-3 && ball.Y >= p2.Y && ball.Y < p2.Y+p2.Length {
			ball.DirX = -ball.DirX
		}

		// Гол
		if ball.X <= 1 {
			score2++
			ball.X, ball.Y = Width/2, Height/2
			ball.DirX = 1
		} else if ball.X >= Width-2 {
			score1++
			ball.X, ball.Y = Width/2, Height/2
			ball.DirX = -1
		}

		// ОТРИСОВКА КАДРА В БУФЕР СТРОКИ (чтобы экран не мерцал)
		clearScreen()
		var frame string

		for y := 0; y < Height; y++ {
			for x := 0; x < Width; x++ {
				if y == 0 || y == Height-1 {
					frame += "#"
				} else if x == 0 || x == Width-1 {
					frame += "#"
				} else if x == 1 && y >= p1.Y && y < p1.Y+p1.Length {
					frame += "|"
				} else if x == Width-2 && y >= p2.Y && y < p2.Y+p2.Length {
					frame += "|"
				} else if x == ball.X && y == ball.Y {
					frame += "O"
				} else {
					frame += " "
				}
			}
			frame += "\n"
		}

		fmt.Print(frame)
		fmt.Printf(" Игрок 1: %d  |  Игрок 2: %d\n", score1, score2)
		
		if score1 >= 5 || score2 >= 5 {
			fmt.Println(" Игра окончена! ")
			return
		}
	}
}
