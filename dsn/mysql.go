package dsn

import (
	"fmt"
	"net/url"
	"time"
)

type MySQLOption func(*MySQLDSN)

type MySQLDSN struct {
	username string
	password string
	db       string
	protocol string
	host     string
	port     int
	address  string
	params   url.Values
}

func newMySQLDSN() DSNHelper {
	return &MySQLDSN{
		protocol: "tcp",
		host:     "localhost",
		port:     3306,
		params:   url.Values{},
	}
}

func (dsn *MySQLDSN) Apply(opts ...MySQLOption) {
	for _, opt := range opts {
		opt(dsn)
	}
}

// Gen [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...¶mN=valueN]
func (dsn *MySQLDSN) Gen() string {
	var output string
	if dsn.username != "" {
		output += dsn.username
		if dsn.password != "" {
			output += fmt.Sprintf(":%s", dsn.password)
		}
		output += "@"
	}

	if dsn.protocol != "" {
		output += dsn.protocol
	}

	if dsn.address == "" {
		output += fmt.Sprintf("(%s:%d)", dsn.host, dsn.port)
	} else {
		output += fmt.Sprintf("(%s)", dsn.address)
	}

	output += fmt.Sprintf("/%s", dsn.db)
	if len(dsn.params) > 0 {
		output += fmt.Sprintf("?%s", dsn.params.Encode())
	}
	return output
}

func SetMySQLUsername(username, password string) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.username = username
		dsn.password = password
	}
}

func SetMySQLHost(host string, port int) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.address = ""
		dsn.host = host
		dsn.port = port
	}
}

func SetMySQLAddress(addr string) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.address = addr
	}
}

func SetMySQLDatabase(db string) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.db = db
	}
}

func SetMySQLParseTime(b bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		if b {
			dsn.params.Set("parseTime", "true")
		} else {
			dsn.params.Set("parseTime", "false")
		}
	}
}

func SetMySQLCharset(v string) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.params.Set("charset", v)
	}
}

func SetMySQLLoc(v string) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.params.Set("loc", v)
	}
}

func SetMySQLCollation(v string) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.params.Set("collation", v)
	}
}

func SetMySQLAllowCleartextPasswords(ok bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		if ok {
			dsn.params.Set("allowCleartextPasswords", "true")
		} else {
			dsn.params.Set("allowCleartextPasswords", "false")
		}
	}
}

func SetMySQLAutoCommit(ok bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		val := "true"
		if !ok {
			val = "false"
		}
		dsn.params.Set("autocommit", val)
	}
}

func SetMySQLAllowAllFiles(ok bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		val := "true"
		if !ok {
			val = "false"
		}
		dsn.params.Set("allowAllFiles", val)
	}
}

func SetMySQLClientFoundRows(ok bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		val := "true"
		if !ok {
			val = "false"
		}
		dsn.params.Set("clientFoundRows", val)
	}
}

func SetMySQLColumnWithAlias(ok bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		val := "true"
		if !ok {
			val = "false"
		}
		dsn.params.Set("columnsWithAlias", val)
	}
}

func SetMySQLInterpolateParams(ok bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		val := "true"
		if !ok {
			val = "false"
		}
		dsn.params.Set("interpolateParams", val)
	}
}

func SetStrict(ok bool) MySQLOption {
	return func(dsn *MySQLDSN) {
		val := "true"
		if !ok {
			val = "false"
		}
		dsn.params.Set("strict", val)
	}
}

func SetTimeout(timeout time.Duration) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.params.Set("timeout", timeout.String())
	}
}

func SetReadTimeout(timeout time.Duration) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.params.Set("readTimeout", timeout.String())
	}
}

func SetWriteTimeout(timeout time.Duration) MySQLOption {
	return func(dsn *MySQLDSN) {
		dsn.params.Set("writeTimeout", timeout.String())
	}
}

func init() {
	RegisterDSN("mysql", newMySQLDSN)
}
