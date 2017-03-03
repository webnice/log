package formater

//import "gopkg.in/webnice/debug.v1"
import (
	"bytes"
	"fmt"
	"strconv"
	"unicode/utf8"

	s "gopkg.in/webnice/log.v2/sender"
	t "gopkg.in/webnice/log.v2/trace"
)

// New Create new object
func New() Interface {
	var ftr = new(impl)
	return ftr
}

// CheckFormat Проверка шаблона на корректность
func (ftr *impl) CheckFormat(tpl string) (matches [][]int, err error) {
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
				err = fmt.Errorf("%s:%s", errUnknownVariable.Error(), name)
				return
			}
			pre = end
		}
	}
	return
}

// TruncateString Обрезание строки в соответствии с форматом
func (ftr *impl) TruncateString(src, layout string) (ret string) {
	var chanks []string
	var err error
	var l int64

	ret = src
	chanks = rexTruncate.FindStringSubmatch(layout)
	if len(chanks) == 3 {
		l, err = strconv.ParseInt(chanks[2], 0, 64)
		if err == nil {
			for i, w, c := 0, 0, int64(0); i < len(src); i += w {
				_, w = utf8.DecodeRuneInString(src[i:])
				c++
				if c >= l {
					ret = src[:i+w]
					break
				}
			}
		}
	}
	return
}

// FormatedElement Получение форматированного элемента
func (ftr *impl) FormatedElement(rcd *t.Info, elm recDic, layout string) (ret string) {
	var frm, timeFormat string
	var parts []string

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
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.Id.String(), layout))
		default:
			ret = fmt.Sprintf(layout, rcd.Id)
		}
	case "Pid":
		if elm.Format == frm {
			ret = fmt.Sprintf(layout, rcd.Pid)
		} else {
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "AppName":
		if elm.Format == frm {
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.AppName, layout))
		} else {
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "HostName":
		if elm.Format == frm {
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.HostName, layout))
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
			ret = fmt.Sprintf("%s", rcd.TodayAndNow.Format(timeFormat))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Level":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, rcd.Level.String())
			ret = fmt.Sprintf(layout, ftr.TruncateString(ret, layout))
		case `d`:
			ret = fmt.Sprintf(layout, rcd.Level)
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Message":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.Message, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Color":
		rcd.Color = true
	case "ColorBeg":
		ret += fmt.Sprint(colorsBackground[colorLevelMap[rcd.Level].Background])
		ret += fmt.Sprint(colors[colorLevelMap[rcd.Level].Foreground])
	case "ColorEnd":
		if rcd.Color == false {
			ret += fmt.Sprint(colorReset)
		}
	case "FileNameLong":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.FileNameLong, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "FileNameShort":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.FileNameShort, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "FileLine":
		switch frm {
		case `d`:
			ret = fmt.Sprintf(layout, rcd.FileLine)
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Package":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.Package, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Module":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.Module, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "Function":
		switch frm {
		case `s`:
			ret = fmt.Sprintf(layout, ftr.TruncateString(rcd.Function, layout))
		default:
			ret = fmt.Sprintf("BAD_FORMAT_'%s',_USE_'%%%s'", layout, elm.Format)
		}
	case "CallStack":
		ret = fmt.Sprintf(layout, rcd.CallStack)
	}
	return
}

// Text Формарирует сообщение в текстовую строку согласно шаблону
func (ftr *impl) Text(msg s.Message, tpl string) (ret *bytes.Buffer, err error) {
	var matches [][]int
	var r []int
	var pre, start, end int
	var name, layout string

	matches, err = ftr.CheckFormat(tpl)
	if err != nil {
		return
	}

	ret = bytes.NewBufferString(``)
	for _, r = range matches {
		start, end = r[0], r[1]
		if start > pre {
			ret.WriteString(tpl[pre:start])
		}
		name = tpl[r[2]:r[3]]
		layout = ""
		if r[4] != -1 {
			layout = `%` + tpl[r[4]:r[5]]
		}
		ret.WriteString(ftr.FormatedElement(msg.Trace, templateNames[name], layout))
		pre = end
	}
	if tpl[pre:] != "" {
		ret.WriteString(tpl[pre:])
	}

	if msg.Trace.Color {
		var tmp = fmt.Sprintf("%s%s%s",
			colorReset+
				colorsBackground[colorLevelMap[msg.Level].Background]+
				colors[colorLevelMap[msg.Level].Foreground],
			ret.String(),
			colorReset,
		)
		ret = bytes.NewBufferString(tmp)
	}
	return
}
