package oscar

import (
	"bytes"
	"errors"
	"net/url"
	"sort"
	"strconv"
)

type Config struct {
	User   string            // Username
	Passwd string            // Password (requires User)
	Host   string            // Host
	Port   int               // Port
	DBName string            // Database name
	Params map[string]string // Connection parameters
}

func NewConfig() *Config {
	return &Config{
		Params: make(map[string]string),
	}
}

func (cfg *Config) normalize() error {
	if cfg.User == "" {
		return errors.New("user is empty")
	}
	if cfg.Passwd == "" {
		return errors.New("passwd is empty")
	}
	if cfg.Host == "" {
		cfg.Host = "127.0.0.1"
	}
	if cfg.Port == 0 {
		cfg.Port = 2003
	}
	if cfg.DBName == "" {
		return errors.New("dbName is empty")
	}
	return nil
}

func writeDSNParam(buf *bytes.Buffer, hasParam *bool, name, value string) {
	buf.Grow(1 + len(name) + 1 + len(value))
	if !*hasParam {
		*hasParam = true
		buf.WriteByte('?')
	} else {
		buf.WriteByte('&')
	}
	buf.WriteString(name)
	buf.WriteByte('=')
	buf.WriteString(value)
}

// FormatDSN formats the given Config into a DSN string which can be passed to the driver.
// USER/PASSWD@HOST:PORT/DB_NAME
func (cfg *Config) FormatDSN() string {
	err := cfg.normalize()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer

	// username/password
	buf.WriteString(cfg.User)
	buf.WriteByte('/')
	buf.WriteString(cfg.Passwd)

	buf.WriteByte('@')

	// host:port
	buf.WriteString(cfg.Host)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(cfg.Port))

	buf.WriteByte('/')

	// dbname
	buf.WriteString(cfg.DBName)

	// ?param1=value1&...&paramN=valueN
	hasParam := false
	if cfg.Params != nil {
		var params []string
		for param := range cfg.Params {
			params = append(params, param)
		}
		sort.Strings(params)
		for _, param := range params {
			writeDSNParam(&buf, &hasParam, param, url.QueryEscape(cfg.Params[param]))
		}
	}

	return buf.String()
}
