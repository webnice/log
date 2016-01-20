package level // import "github.com/webdeskltd/log/level"

func init() {
	Map2Level = make(map[LevelName]Level)
	for n := range Map {
		Map2Level[Map[n]] = n
	}
}

// New - Creating an object, level logging
func New(l Level) *LevelObject {
	var self *LevelObject = new(LevelObject)
	self.Level = l
	return self
}

// NewFromMesssage - Determines the level of logging on the first word of the message
func NewFromMesssage(str string, deflt Level) *LevelObject {
	var self *LevelObject = new(LevelObject)
	self.Level = deflt

	// Опредяеляем уровень

	return self
}

func (self *LevelObject) String() (ret string) {
	resp, _ := Map[self.Level]
	ret = string(resp)
	return
}

func (self *LevelObject) Int8() int8 {
	return int8(self.Level)
}
