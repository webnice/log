package record

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	l "github.com/webdeskltd/log/level"

	"github.com/webdeskltd/debug"
)

const (
	tagName           string = `fmt`
	defaultTimeFormat string = `2006-01-02T15:04:05.999Z07:00`
)

var (
	rexFormat          *regexp.Regexp    = regexp.MustCompile(`%{([a-z]+)(?::(.*?[^\\]))?}`) // Регулярное выражение поиска констант шаблона
	rexTruncate        *regexp.Regexp    = regexp.MustCompile(`(.*?)(\d+?)s`)                // Регулярное выражение разбора формата строки
	rexTime            *regexp.Regexp    = regexp.MustCompile(`^%(.*)t$`)                    // Регулярное выражение разбора формата времени
	templateNames      map[string]recDic                                                     // Справочник доступных констант шаблона
	errWrongTag        error             = errors.New(`Wrong tag`)                           // return if tag is incorrect
	errUnknownVariable error             = errors.New(`Unknown variable`)                    // return if found unknown variable as prefix
	errInvalidFormat   error             = errors.New(`Invalid log format`)                  // return if log format is empty or not one variable found
)

type (
	recDic struct {
		Index  int
		Format string
		Type   string
		Name   string
	}
)

func init() {
	makeDictionary(new(Record))
	debug.Nop()
}

// Создание на основе структуры констант используемых при работе
func makeDictionary(v interface{}) (err error) {
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
					Format: attr[1],
					Type:   rt.Field(i).Type.String(),
					Name:   rt.Field(i).Name,
				}
			}
			if len(attr) > 2 {
				err = errors.New(errWrongTag.Error() + `:` + s)
				return
			}
		}
	}
	v = nil
	return
}

// Проверка шаблона на корректность
func CheckFormat(tpl string) (matches [][]int, err error) {
	var r []int
	var pre, start, end int
	var name string
	matches = rexFormat.FindAllStringSubmatchIndex(tpl, -1)
	if len(matches) == 0 {
		err = errInvalidFormat
		return
	}
	for _, r = range matches {
		start, end = r[0], r[1]
		if start > pre {
			name = tpl[r[2]:r[3]]
			if _, ok := templateNames[name]; ok == false {
				err = errors.New(errUnknownVariable.Error() + ":" + name)
				return
			}
			pre = end
		}
	}
	return
}

// Обрезание строки в соответствии с форматом
func TruncateString(src, layout string) (ret string) {
	var chanks []string
	var err error
	var l int64

	ret = src
	chanks = rexTruncate.FindStringSubmatch(layout)
	if len(chanks) == 3 {
		l, err = strconv.ParseInt(chanks[2], 0, 64)
		if err == nil && int(l) <= len(src) {
			ret = src[:int(l)]
		}
	}
	return
}

func (self *Record) getFormatedElement(elm recDic, layout string) (ret string) {
	var frm, timeFormat string
	var parts []string
	var ok bool

	if layout == "" {
		layout = `%` + elm.Format
	}
	if len(layout) > 0 {
		frm = layout[len(layout)-1:]
	}

	// Вариант1 - Через reflect
	// ... удалён

	// Вариант2 - без использования reflect
	// ... быстрее, но подразумевает что структура меняться не будет
	switch elm.Name {
	case "Id":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, TruncateString(self.Id.String(), layout))
		default:
			ret = fmt.Sprintf(layout, self.Id)
		}
	case "Pid":
		if elm.Format == frm {
			ret = fmt.Sprintf(layout, self.Pid)
		} else {
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "AppName":
		if elm.Format == frm {
			ret = fmt.Sprintf(layout, TruncateString(self.AppName, layout))
		} else {
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "HostName":
		if elm.Format == frm {
			ret = fmt.Sprintf(layout, TruncateString(self.HostName, layout))
		} else {
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "TodayAndNow":
		switch frm {
		case `t`:
			parts = rexTime.FindStringSubmatch(layout)
			if len(parts) == 2 {
				timeFormat = parts[1]
			}
			if timeFormat == "" {
				timeFormat = defaultTimeFormat
			}
			ret = fmt.Sprintf("%s", self.TodayAndNow.Format(timeFormat))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Level":
		switch frm {
		case `s`:
			if _, ok = l.Map[self.Level]; ok {
				ret = fmt.Sprintf(layout, l.Map[self.Level])
			} else {
				ret = `-`
			}
			ret = fmt.Sprintf(layout, TruncateString(ret, layout))
		case `d`:
			ret = fmt.Sprintf(layout, self.Level)
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Message":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, TruncateString(self.Message, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "color":
		self.color = true
	case "colorBeg":
		ret += fmt.Sprint(colorsBackground[colorLevelMap[self.Level].Background])
		ret += fmt.Sprint(colors[colorLevelMap[self.Level].Foreground])
	case "colorEnd":
		if self.color == false {
			ret += fmt.Sprint(colorReset)
		}
	case "FileNameLong":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, TruncateString(self.FileNameLong, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "FileNameShort":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, TruncateString(self.FileNameShort, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "FileLine":
		switch frm {
		case `d`:
			ret = fmt.Sprintf(layout, self.FileLine)
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Package":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, TruncateString(self.Package, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Module":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, TruncateString(self.Module, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Function":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, TruncateString(self.Function, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "CallStack":
		ret = fmt.Sprintf(layout, self.CallStack)
	}
	return
}

// Forming a formatted log message based on Record
func (self *Record) Format(tpl string) (ret string, err error) {
	var matches [][]int
	var r []int
	var pre, start, end int
	var name, layout string
	var resultTmp []byte
	var result *bytes.Buffer
	matches, err = CheckFormat(tpl)
	if err != nil {
		return
	}
	if self.resolver != nil {
		self.resolver(self)
	}

	result = bytes.NewBuffer(resultTmp)
	for _, r = range matches {
		start, end = r[0], r[1]
		if start > pre {
			result.WriteString(tpl[pre:start])
		}
		name = tpl[r[2]:r[3]]
		layout = ""
		if r[4] != -1 {
			layout = `%` + tpl[r[4]:r[5]]
		}
		result.WriteString(self.getFormatedElement(templateNames[name], layout))
		pre = end
	}
	if tpl[pre:] != "" {
		result.WriteString(tpl[pre:])
	}
	if self.color {
		ret = fmt.Sprintf("%s%s%s",
			colorReset+
				colorsBackground[colorLevelMap[self.Level].Background]+
				colors[colorLevelMap[self.Level].Foreground],
			result.String(),
			colorReset,
		)
	} else {
		ret = result.String()
	}
	return
}
