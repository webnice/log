package formater

//import "gopkg.in/webnice/debug.v1"
import (
	"bytes"
	"errors"
	"regexp"

	s "gopkg.in/webnice/log.v2/sender"
)

const tagName = `fmt`
const defaultTimeFormat = `2006-01-02T15:04:05.999Z07:00`

var (
	rexFormat          *regexp.Regexp    = regexp.MustCompile(`%{([a-z]+)(?::(.*?[^\\]))?}`) // Регулярное выражение поиска констант шаблона
	rexTruncate        *regexp.Regexp    = regexp.MustCompile(`(.*?)(\d+?)s`)                // Регулярное выражение разбора формата строки
	rexTime            *regexp.Regexp    = regexp.MustCompile(`^%(.*)t$`)                    // Регулярное выражение разбора формата времени
	templateNames      map[string]recDic                                                     // Справочник доступных констант шаблона
	errWrongTag        error             = errors.New(`Wrong tag`)                           // return if tag is incorrect
	errUnknownVariable error             = errors.New(`Unknown variable`)                    // return if found unknown variable as prefix
	errInvalidFormat   error             = errors.New(`Invalid log format`)                  // return if log format is empty or not one variable found
)

type recDic struct {
	Index  int
	Format string
	Type   string
	Name   string
}

// Interface is an interface of package
type Interface interface {
	// Text Формарирует сообщение в текстовую строку согласно шаблону
	Text(s.Message, string) (*bytes.Buffer, error)
}

// impl is an implementation of package
type impl struct {
}
