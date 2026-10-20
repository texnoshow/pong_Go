package main

import (
	"math/rand"
	"os"
	"time"

	"://github.com"
)

const (
	Width  = 60
	Height = 20
)

type Paddle struct {
	Y, Length int
}

type Ball struct {
	X, Y, DirX, DirY int
}

func main() {
	// Инициализация экрана tcell
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	// Начальное состояние игры
	rand.Seed(time.Now().UnixNano())
	padLength := 4
	p1 := Paddle{Y: Height/2 - padLength/2, Length: padLength}
	p2 := Paddle{Y: Height/2 - padLength/2, Length: padLength}
	
	ball := Ball{
		X:    Width / 2,
		Y:    Height / 2,
		DirX: 1,
		DirY: 1,
	}
	if rand.Intn(2) == 0 {
		ball.DirX = -1
	}

	score1, score2 := 0, 0
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorWhite)

	// Канал для отслеживания нажатий клавиш
	inputChan := make(chan tcell.Event)
	go func() {
		for {
			inputChan <- screen.PollEvent()
		}
	}()

	// Игровой цикл (примерно 30 FPS)
	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case ev := <-inputChan:
			switch e := ev.(type) {
			case *tcell.EventResize:
				screen.Sync()
			case *tcell.EventKey:
				// Выход из игры
				if e.Key() == tcell.KeyEscape || e.Rune() == 'q' || e.Rune() == 'Q' {
					return
				}
				// Управление Игроком 1 (W/S)
				if e.Rune() == 'w' || e.Rune() == 'W' {
					if p1.Y > 1 {
						p1.Y--
					}
				}
				if e.Rune() == 's' || e.Rune() == 'S' {
					if p1.Y < Height-p1.Length-1 {
						p1.Y++
					}
				}
				// Управление Игроком 2 (I/K)
				if e.Rune() == 'i' || e.Rune() == 'I' {
					if p2.Y > 1 {
						p2.Y--
					}
				}
				if e.Rune() == 'k' || e.Rune() == 'K' {
					if p2.Y < Height-p2.Length-1 {
						p2.Y++
					}
				}
			}

		case <-ticker.C:
			// ЛОГИКА ДВИЖЕНИЯ
			ball.X += ball.DirX
			ball.Y += ball.DirY

			// Отскок от верхней и нижней границ
			if ball.Y <= 1 || ball.Y >= Height-2 {
				ball.DirY = -ball.DirY
			}

			// Отскок от левой ракетки (Игрок 1)
			if ball.X == 2 && ball.Y >= p1.Y && ball.Y < p1.Y+p1.Length {
				ball.DirX = -ball.DirX
			}

			// Отскок от правой ракетки (Игрок 2)
			if ball.X == Width-3 && ball.Y >= p2.Y && ball.Y < p2.Y+p2.Length {
				ball.DirX = -ball.DirX
			}

			// Пропуск мяча (Гол)
			if ball.X <= 1 {
				score2++
				ball.X, ball.Y = Width/2, Height/2
				ball.DirX = 1
			} else if ball.X >= Width-2 {
				score1++
				ball.X, ball.Y = Width/2, Height/2
				ball.DirX = -1
			}

			// ОТРИСОВКА ЭКРАНА
			screen.Clear()

			// Верхняя и нижняя стены
			for x := 0; x < Width; x++ {
				screen.SetContent(x, 0, '#', nil, style)
				screen.SetContent(x, Height-1, '#', nil, style)
			}
			// Левая и правая стены
			for y := 0; y < Height; y++ {
				screen.SetContent(0, y, '#', nil, style)
				screen.SetContent(Width-1, y, '#', nil, style)
			}

			// Ракетка 1
			for y := p1.Y; y < p1.Y+p1.Length; y++ {
				screen.SetContent(1, y, '|', nil, style)
			}
			// Ракетка 2
			for y := p2.Y; y < p2.Y+p2.Length; y++ {
				screen.SetContent(Width-2, y, '|', nil, style)
			}

			// Мяч
			screen.SetContent(ball.X, ball.Y, 'O', nil, style)

			// Вывод счета под полем
			scoreStr := " Игрок 1 (W/S): %d   |   Игрок 2 (I/K): %d   [Q - Выход]"
			runes := []rune(string([]byte(scoreStr))) // простая конвертация для вывода текста
			// Для вывода форматированной строки нарисуем её посимвольно
			drawText(screen, 2, Height+1, style, "Игрок 1 (W/S): "+string(rune(score1+48))+"   |   Игрок 2 (I/K): "+string(rune(score2+48))+"   [Q - Выход]")

			screen.Show()
		}
	}
}

func drawText(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for i, r := range str {
		s.SetContent(x+i, y, r, nil, style)
	}
}
