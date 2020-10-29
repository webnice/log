package formater // import "github.com/webnice/log/v2/formater"

import (
	"errors"
	"reflect"
	"strings"

	t "github.com/webnice/log/v2/trace"
)

func init() {
	_ = makeDictionary(new(t.Info))
}

// Создание на основе структуры констант используемых в качестве формата
func makeDictionary(v interface{}) (err error) {
	var (
		rv          reflect.Value
		rt          reflect.Type
		rs          reflect.StructField
		i, n        int
		s           string
		names, attr []string
	)

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
