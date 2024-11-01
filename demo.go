package main

import (
	"fmt"
	_ "github.com/Mystery00/go-shentong"
	"github.com/skeyic/gorm-shentong/oscar"
	"github.com/skeyic/gorm-shentong/shentong"
	"gorm.io/gorm"
	"strings"
)

func main() {
	db, err := gorm.Open(shentong.New(shentong.Config{
		DSN: "user/password@host:port/dbname",
		DSNConfig: &oscar.Config{
			User:   "test",
			Passwd: "testPasswd",
			Host:   "127.0.0.1",
			Port:   2003,
			DBName: "OSRDB",
		},
		FieldConvertType: shentong.Custom,
		FieldConvertFunc: func(s string) string {
			return strings.ToUpper(s)
		},
	}))

	if err != nil {
		panic(err)
	}
	fmt.Printf("连接成功：%s", db.Dialector.Name())
}
