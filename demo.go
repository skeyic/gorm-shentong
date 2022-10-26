package main

import (
	_ "github.com/Mystery00/go-shentong"
	"github.com/Mystery00/gorm-shentong/shentong"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"os"
	"strings"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
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
		logrus.Error("连接失败：", err)
		os.Exit(0)
	}
	logrus.Infof("连接成功：%s", gormDb.Dialector.Name())
}
