package ft

import "strings"

type Tag struct {
	ID     uint `gorm:"primary_key;autoIncrement"`
	FileID uint
	Name   string
	Value  string
}

func ParseTagKeyValue(s string) (string, string) {
	var name, value string
	p := strings.IndexRune(s, '=')
	if p == -1 {
		name = s
	} else {
		name = s[:p]
		value = s[p+1:]
	}
	return name, value
}
