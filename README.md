# gorm-shentong
# 使用方式
```go
	db, err := gorm.Open(shentong.New(shentong.Config{
        DSNConfig: &oscar.Config{
            User:   "test",
            Passwd: "testPasswd",
            Host:   "127.0.0.1",
            Port:   2003,
            DBName: "OSRDB",
        },
    }))
```
可以通过DSNConfig进行组装DSN，也可以直接通过DSN进行配置
```go
    db, err := gorm.Open(shentong.New(shentong.Config{
        DSN: "user/password@host:port/dbname",
    }))
```
在某些时候需要适配多种数据库，那么我们会编写统一的model结构体，但是在对应到部分数据库时需要对字段名称进行统一的处理，例如转为大写，那么可以通过 `FieldConvertType` 来进行统一处理
```go
    db, err := gorm.Open(shentong.New(shentong.Config{
        DSN: "user/password@host:port/dbname",
        FieldConvertType: shentong.Upper,
    }))
```
也可以通过 `FieldConvertFunc` 来进行自定义的处理
```go
    db, err := gorm.Open(shentong.New(shentong.Config{
        DSN: "user/password@host:port/dbname",
        FieldConvertFunc: func(field string) string {
            return strings.ToUpper(field)
        },
    }))
```
# 编译说明
TODO