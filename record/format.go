package record

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/webdeskltd/debug"
)

const (
	tagName string = `fmt`
)

const test string = ` 1: %{id}
 2: %{pid:8d}                  - (int      ) Process id
 3: %{application}             - (string   ) Application name basename of os.Args[0]
 4: %{hostname}                - (string   ) Server host name
 5: %{time}                    - (time.Time) Time when log occurred
 6: %{level:-8d}               - (int8     ) Log level
 7: %{message}                 - (string   ) Message
 8: %{color}                   - %{begcolor}(bool     ) ANSI color based on log level%{endcolor}
 9: %{longfile}                - (string   ) Full file name and line number: /a/b/c/d.go
10: %{shortfile}               - (string   ) Final file name element and line number: d.go
11: %{line}                    - (int      ) Line number in file
12: %{package}                 - (string   ) Full package path, eg. github.com/webdeskltd/log
13: %{module} or %{shortpkg}   - (string   ) Module name base package path, eg. log
14: %{function} or %{facility} - (string   ) Full function name, eg. PutUint32
15: %{callstack}               - (string   ) Full call stack

"%{color}[%{module:-10s}] %{time:2006-01-02T15:04:05.000Z07:00} (%{level:7s}): %{message} (%{package}) (%{function}:%{line}) (%{shortfile}:%{line}) (%{longfile})"
%{level:.1s}
`
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
const (
	level_FATAL    logLevel = iota // 0 Fatal: system is unusable
	level_ALERT                    // 1 Alert: action must be taken immediately
	level_CRITICAL                 // 2 Critical: critical conditions
	level_ERROR                    // 3 Error: error conditions
	level_WARNING                  // 4 Warning: warning conditions
	level_NOTICE                   // 5 Notice: normal but significant condition
	level_INFO                     // 6 Informational: informational messages
	level_DEBUG                    // 7 Debug: debug-level messages
)

var (
	rexFormat     *regexp.Regexp    = regexp.MustCompile(`%{([a-z]+)(?::(.*?[^\\]))?}`) // Регулярное выражение поиска констант шаблона
	templateNames map[string]recDic                                                     // Справочник доступных констант шаблона
	colorReset    string            = "\033[0m"                                         // Сброс цветов
)
var (
	// ANSI colors
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
	colorLevelMap map[logLevel]ansiStyle = map[logLevel]ansiStyle{
		level_FATAL:    ansiStyle{Background: color_RED, Foreground: color_BLACK},     // Система не стабильна, проолжение работы не возможно
		level_ALERT:    ansiStyle{Background: color_MAGENTA, Foreground: color_WHITE}, // Система не стабильна но может частично продолжить работу (например запусился один из двух серверов - что-то работает а что-то нет)
		level_CRITICAL: ansiStyle{Background: color_BLACK, Foreground: color_MAGENTA}, // Критическая ошибка, часть функционала системы работает не корректно
		level_ERROR:    ansiStyle{Background: color_BLACK, Foreground: color_RED},     // Ошибки не прерывающие работу приложения
		level_WARNING:  ansiStyle{Background: color_BLACK, Foreground: color_YELLOW},  // Предупреждения
		level_NOTICE:   ansiStyle{Background: color_BLACK, Foreground: color_GREEN},   // Информационные сообщения
		level_INFO:     ansiStyle{Background: color_BLACK, Foreground: color_WHITE},   // Сообщения информационного характера описывающие шаги выполнения алгоритмов приложения
		level_DEBUG:    ansiStyle{Background: color_BLACK, Foreground: color_CYAN},    // Режим отладки, аналогичен INFO но с подробными данными и дампом переменных
	}
)

type (
	logLevel  int8
	ansiColor int16
	ansiStyle struct {
		Background ansiColor // Цвет фона
		Foreground ansiColor // Цвет текста
	}
	recDic struct {
		Index  int
		Format string
		Type   string
		Name   string
	}
)

func init() {
	makeDictionary()
	debug.Nop()
}

// Создание на основе структуры констант используемых при работе
func makeDictionary() {
	var v *Record = new(Record)
	var rv reflect.Value
	var rt reflect.Type
	var rs reflect.StructField
	var i, n int
	var s string
	var names, attr []string
	templateNames = make(map[string]recDic)
	rv = reflect.Indirect(reflect.ValueOf(v))
	rt = rv.Type()

	for i = 0; i < rt.NumField(); i++ {
		rs = rt.Field(i)
		s = rs.Tag.Get(tagName)
		if s == `-` || s == `` {
			continue
		}
		names = strings.Split(s, `,`)
		for n = range names {
			attr = strings.Split(names[n], `:`)
			if len(attr) == 1 {
				attr = append(attr, `v`)
			}
			if len(attr) == 2 {
				templateNames[attr[0]] = recDic{
					Index:  i,
					Format: `%` + attr[1],
					Type:   rt.Field(i).Type.String(),
					Name:   rt.Field(i).Name,
				}
			}
		}
	}
	v = nil
	debug.Dumper(templateNames)
}

// Проверка шаблона на корректность
func CheckFormat(tpl string) (err error) {
	var matches [][]int
	var r []int
	var pre, start, end int
	var name string
	matches = rexFormat.FindAllStringSubmatchIndex(tpl, -1)
	for _, r = range matches {
		start, end = r[0], r[1]
		if start > pre {
			name = tpl[r[2]:r[3]]
			if _, ok := templateNames[name]; ok == false {
				err = errors.New("Unknown variable: " + name)
				return
			}
			pre = end
		}
	}
	return
}

func (this *Record) getFormatedElement(elm recDic, layout string) (ret string) {
	if layout == "" {
		layout = elm.Format
	}
	// Вариант1 - Через reflect
	// ... Удалён

	// Вариант2 - без использования reflect
	switch elm.Name {
	case "Id":
		ret = fmt.Sprintf(layout, this.Id)
	case "Pid":
		ret = fmt.Sprintf(layout, this.Pid)
	case "AppName":
		ret = fmt.Sprintf(layout, this.AppName)
	case "HostName":
		ret = fmt.Sprintf(layout, this.HostName)
	case "TodayAndNow":
		//ret = fmt.Sprintf(layout, this.TodayAndNow)
	case "Level":
		ret = fmt.Sprintf(layout, this.Level)
	case "Message":
		ret = fmt.Sprintf(layout, this.Message)
	case "Color":
		this.Color = true
	case "BegColor":
		ret += fmt.Sprint(colorsBackground[colorLevelMap[logLevel(this.Level)].Background])
		ret += fmt.Sprint(colors[colorLevelMap[logLevel(this.Level)].Foreground])
	case "EndColor":
		if this.Color == false {
			ret += fmt.Sprint(colorReset)
		}
	case "FileNameLong":
		ret = fmt.Sprintf(layout, this.FileNameLong)
	case "FileNameShort":
		ret = fmt.Sprintf(layout, this.FileNameShort)
	case "FileLine":
		ret = fmt.Sprintf(layout, this.FileLine)
	case "Package":
		ret = fmt.Sprintf(layout, this.Package)
	case "Module":
		ret = fmt.Sprintf(layout, this.Module)
	case "Function":
		ret = fmt.Sprintf(layout, this.Function)
	case "CallStack":
		ret = fmt.Sprintf(layout, this.CallStack)
	}
	return
}

func (this *Record) Format() (ret string, err error) {
	var matches [][]int
	var r []int
	var pre, start, end int
	var name, layout string
	var resultTmp []byte
	var result *bytes.Buffer

	err = CheckFormat(test)
	if err != nil {
		return
	}
	result = bytes.NewBuffer(resultTmp)

	matches = rexFormat.FindAllStringSubmatchIndex(test, -1)
	fmt.Println(test)

	//	this.Level = int8(level_FATAL)
	//	this.Level = int8(level_ALERT)
	//	this.Level = int8(level_CRITICAL)
	//	this.Level = int8(level_ERROR)
	//	this.Level = int8(level_WARNING)
	//	this.Level = int8(level_NOTICE)
	//	this.Level = int8(level_INFO)
	//	this.Level = int8(level_DEBUG)

	pre = 0
	for _, r = range matches {
		start, end = r[0], r[1]
		if start > pre {
			result.WriteString(test[pre:start])
		}
		name = test[r[2]:r[3]]
		layout = ""
		if r[4] != -1 {
			layout = `%` + test[r[4]:r[5]]
		}
		result.WriteString(this.getFormatedElement(templateNames[name], layout))
		pre = end
	}
	if test[pre:] != "" {
		result.WriteString(test[pre:])
	}

	ret = result.String()

	if this.Color {
		fmt.Print(colorReset)
		fmt.Print(colorsBackground[colorLevelMap[logLevel(this.Level)].Background])
		fmt.Print(colors[colorLevelMap[logLevel(this.Level)].Foreground])
	}
	fmt.Print(ret)
	if this.Color {
		fmt.Print(colorReset)
	}
	fmt.Println()
	return
}
