package shentong

import (
	"gorm.io/gorm"
)

type fieldConverter func(string) string

func queryFix(db *gorm.DB) {
	config := db.Dialector.(*Dialector).Config
	if config.FieldConvertType == None {
		//不转换，直接返回
		return
	}
	var converter fieldConverter
	if config.FieldConvertType == Custom {
		//自定义转换
		converter = config.FieldConvertFunc
	} else {
		converter = config.FieldConvertType.convert
	}
	for _, field := range db.Statement.Schema.Fields {
		//将新的映射关系存入db.Statement.Schema.FieldsByDBName
		db.Statement.Schema.FieldsByDBName[converter(field.DBName)] = field
	}
}
