package shentong

import (
	"database/sql"
	"fmt"
	_ "github.com/Mystery00/go-shentong"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
	"strconv"
	"strings"
)

type Config struct {
	DriverName        string
	DSN               string
	Conn              *sql.DB
	DefaultStringSize uint
	FieldConvertType  FieldConvertType
	FieldConvertFunc  func(string) string
}

type Dialector struct {
	*Config
}

func Open(dsn string) gorm.Dialector {
	return &Dialector{Config: &Config{DSN: dsn}}
}

func New(config Config) gorm.Dialector {
	return &Dialector{Config: &config}
}

func (d Dialector) Name() string {
	return "oracle"
}

func (d Dialector) Initialize(db *gorm.DB) (err error) {
	// register callbacks
	callbacks.RegisterDefaultCallbacks(db, &callbacks.Config{LastInsertIDReversed: true})

	d.DriverName = "aci"
	if d.Conn != nil {
		db.ConnPool = d.Conn
	} else {
		db.ConnPool, err = sql.Open(d.DriverName, d.DSN)
		if err != nil {
			return
		}
	}

	if err = db.Callback().Query().Before("*").Register("shentong_query", queryFix); err != nil {
		return err
	}

	return
}

func (d Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{
		Migrator: migrator.Migrator{
			Config: migrator.Config{
				DB:                          db,
				Dialector:                   d,
				CreateIndexAfterCreateTable: true,
			},
		},
	}
}

func (d Dialector) DataTypeOf(field *schema.Field) string {
	switch field.DataType {
	case schema.Bool:
		return "BOOLEAN"
	case schema.Int, schema.Uint:
		return d.getSchemaIntAndUnitType(field)
	case schema.Float:
		return d.getSchemaFloatType(field)
	case schema.String:
		return d.getSchemaStringType(field)
	case schema.Time:
		return d.getSchemaTimeType(field)
	case schema.Bytes:
		return "BLOB"
	default:
		return d.getSchemaCustomType(field)
	}
}

func (d Dialector) getSchemaFloatType(field *schema.Field) string {
	if field.Precision > 0 {
		return fmt.Sprintf("DECIMAL(%d, %d)", field.Precision, field.Scale)
	}

	if field.Size <= 32 {
		return "FLOAT"
	}

	return "DOUBLE"
}

func (d Dialector) getSchemaStringType(field *schema.Field) string {
	size := field.Size
	if size >= 8000 {
		return "TEXT"
	}

	return fmt.Sprintf("VARCHAR(%d)", size)
}

func (d Dialector) getSchemaTimeType(field *schema.Field) string {
	if field.NotNull || field.PrimaryKey {
		return "TIMESTAMP NOT NULL"
	}
	return "TIMESTAMP NULL"
}

func (d Dialector) getSchemaIntAndUnitType(field *schema.Field) string {
	sqlType := "BIGINT"
	switch {
	case field.Size <= 8:
		sqlType = "TINYINT"
	case field.Size <= 16:
		sqlType = "SMALLINT"
	case field.Size <= 24:
		sqlType = "INT"
	}

	if field.AutoIncrement {
		sqlType += " AUTO_INCREMENT"
	}

	return sqlType
}

func (d Dialector) getSchemaCustomType(field *schema.Field) string {
	sqlType := string(field.DataType)

	if field.AutoIncrement && !strings.Contains(strings.ToLower(sqlType), " auto_increment") {
		sqlType += " AUTO_INCREMENT"
	}

	return sqlType
}

func (d Dialector) DefaultValueOf(*schema.Field) clause.Expression {
	return clause.Expr{SQL: "VALUES (DEFAULT)"}
}

func (d Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteString(":")
	writer.WriteString(strconv.Itoa(len(stmt.Vars)))
}

func (d Dialector) QuoteTo(writer clause.Writer, str string) {
	var (
		underQuoted, selfQuoted bool
		continuousBacktick      int8
		shiftDelimiter          int8
	)

	for _, v := range []byte(str) {
		switch v {
		case '.':
			if continuousBacktick > 0 || !selfQuoted {
				shiftDelimiter = 0
				underQuoted = false
				continuousBacktick = 0
				writer.WriteByte('"')
			}
			writer.WriteByte(v)
			continue
		default:
			if shiftDelimiter-continuousBacktick <= 0 && !underQuoted {
				writer.WriteByte('"')
				underQuoted = true
				if selfQuoted = continuousBacktick > 0; selfQuoted {
					continuousBacktick -= 1
				}
			}

			for ; continuousBacktick > 0; continuousBacktick -= 1 {
				writer.WriteString(`""`)
			}

			writer.WriteByte(v)
		}
		shiftDelimiter++
	}

	if continuousBacktick > 0 && !selfQuoted {
		writer.WriteString(`""`)
	}
	writer.WriteByte('"')
}

func (d Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, nil, `'`, vars...)
}

func (d Dialector) DummyTableName() string {
	return "DUAL"
}
