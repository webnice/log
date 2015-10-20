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

var (
	rexFormat          *regexp.Regexp    = regexp.MustCompile(`%{([a-z]+)(?::(.*?[^\\]))?}`) // Регулярное выражение поиска констант шаблона
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
					Format: `%` + attr[1],
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
	case "color":
		this.color = true
	case "colorBeg":
		ret += fmt.Sprint(colorsBackground[colorLevelMap[this.Level].Background])
		ret += fmt.Sprint(colors[colorLevelMap[this.Level].Foreground])
	case "colorEnd":
		if this.color == false {
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

func (this *Record) Format(tpl string) (ret string, err error) {
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
		result.WriteString(this.getFormatedElement(templateNames[name], layout))
		pre = end
	}
	if tpl[pre:] != "" {
		result.WriteString(tpl[pre:])
	}

	ret = result.String()

	if this.color {
		fmt.Print(colorReset)
		fmt.Print(colorsBackground[colorLevelMap[this.Level].Background])
		fmt.Print(colors[colorLevelMap[this.Level].Foreground])
	}
	fmt.Print(ret)
	if this.color {
		fmt.Print(colorReset)
	}
	fmt.Println()
	return
}
