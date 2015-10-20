package record

import (
	"fmt"
	
	l "github.com/webdeskltd/log/level"
)

const (
	color_BLACK   ansiColor = iota // 0 Color to black
	color_RED                      // 1 Color to red
	color_GREEN                    // 2 Color to green
	color_YELLOW                   // 3 Color to yellow
	color_BLUE                     // 4 Color to blue
	color_MAGENTA                  // 5 Color to magenta (purple)
	color_CYAN                     // 6 Color to cyan
	color_WHITE                    // 7 Color to white
)

var (
	colorReset string = "\033[0m" // Reset background and foreground colors
	// ANSI colors foreground
	colors map[ansiColor]string = map[ansiColor]string{
		color_BLACK:   fmt.Sprintf("\033[%dm", 30+int(color_BLACK)),
		color_RED:     fmt.Sprintf("\033[%dm", 30+int(color_RED)),
		color_GREEN:   fmt.Sprintf("\033[%dm", 30+int(color_GREEN)),
		color_YELLOW:  fmt.Sprintf("\033[%dm", 30+int(color_YELLOW)),
		color_BLUE:    fmt.Sprintf("\033[%dm", 30+int(color_BLUE)),
		color_MAGENTA: fmt.Sprintf("\033[%dm", 30+int(color_MAGENTA)),
		color_CYAN:    fmt.Sprintf("\033[%dm", 30+int(color_CYAN)),
		color_WHITE:   fmt.Sprintf("\033[%dm", 30+int(color_WHITE)),
	}
	// ANSI colors background
	colorsBackground map[ansiColor]string = map[ansiColor]string{
		color_BLACK:   fmt.Sprintf("\033[%d;1m", 40+int(color_BLACK)),
		color_RED:     fmt.Sprintf("\033[%d;1m", 40+int(color_RED)),
		color_GREEN:   fmt.Sprintf("\033[%d;1m", 40+int(color_GREEN)),
		color_YELLOW:  fmt.Sprintf("\033[%d;1m", 40+int(color_YELLOW)),
		color_BLUE:    fmt.Sprintf("\033[%d;1m", 40+int(color_BLUE)),
		color_MAGENTA: fmt.Sprintf("\033[%d;1m", 40+int(color_MAGENTA)),
		color_CYAN:    fmt.Sprintf("\033[%d;1m", 40+int(color_CYAN)),
		color_WHITE:   fmt.Sprintf("\033[%d;1m", 40+int(color_WHITE)),
	}
	// Colors for error level
	colorLevelMap map[l.Level]ansiStyle = map[l.Level]ansiStyle{
		l.FATAL:    ansiStyle{Background: color_RED, Foreground: color_BLACK},     // Система не стабильна, проолжение работы не возможно
		l.ALERT:    ansiStyle{Background: color_MAGENTA, Foreground: color_WHITE}, // Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
		l.CRITICAL: ansiStyle{Background: color_BLACK, Foreground: color_MAGENTA}, // Критическая ошибка, часть функционала системы работает не корректно
		l.ERROR:    ansiStyle{Background: color_BLACK, Foreground: color_RED},     // Ошибки не прерывающие работу приложения
		l.WARNING:  ansiStyle{Background: color_BLACK, Foreground: color_YELLOW},  // Предупреждения
		l.NOTICE:   ansiStyle{Background: color_BLACK, Foreground: color_GREEN},   // Информационные сообщения
		l.INFO:     ansiStyle{Background: color_BLACK, Foreground: color_WHITE},   // Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
		l.DEBUG:    ansiStyle{Background: color_BLACK, Foreground: color_CYAN},    // Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
	}
)

type (
	ansiColor int16
	ansiStyle struct {
		Background ansiColor // Background color
		Foreground ansiColor // Foreground color
	}
)
