package shentong

import (
	"strings"
)

type FieldConvertType int8

const (
	None        FieldConvertType = iota //不转换
	ToUpperCase                         //转换为大写
	ToLowerCase                         //转换为小写
	Custom                              //自定义转换
)

func (t FieldConvertType) convert(columnName string) string {
	switch t {
	case None:
		return columnName
	case ToUpperCase:
		return strings.ToUpper(columnName)
	case ToLowerCase:
		return strings.ToLower(columnName)
	default:
		return columnName
	}
}
