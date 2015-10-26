package level

func init() {
	Map2Level = make(map[LevelName]Level)
	for n := range Map {
		Map2Level[Map[n]] = n
	}
}

func New(l Level) *LevelObject {
	var self *LevelObject = new(LevelObject)
	self.Level = l
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
