package formater

//import "gopkg.in/webnice/debug.v1"
import (
	"errors"
	"reflect"
	"strings"

	t "gopkg.in/webnice/log.v2/trace"
)

func init() {
	makeDictionary(new(t.Info))
}

// Создание на основе структуры констант используемых в качестве формата
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
