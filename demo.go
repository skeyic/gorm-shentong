package main

import (
	"fmt"
	_ "github.com/Mystery00/go-shentong"
	"github.com/Mystery00/gorm-shentong/shentong"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
)

func main() {
	gormDb, err := gorm.Open(shentong.New(shentong.Config{
		DSN:              "test/testPasw@127.0.0.1:2003/OSRDB",
		FieldConvertType: shentong.Custom,
		FieldConvertFunc: func(s string) string {
			return strings.ToUpper(s)
		},
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, //表名后面不加s
		},
	})

	if err != nil {
		panic(err)
	}
	fmt.Printf("连接成功：%s", gormDb.Dialector.Name())
}
