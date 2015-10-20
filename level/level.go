package level

import ()

func New(l Level) *LevelObject {
	var self *LevelObject = new(LevelObject)
	self.Level = l
	return self
}

func (self *LevelObject) String() (ret string) {
	ret, _ = Map[self.Level]
	return
}

func (self *LevelObject) Int8() int8 {
	return int8(self.Level)
}
