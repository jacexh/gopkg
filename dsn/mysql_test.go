package dsn

import (
	"fmt"
	"testing"
	"time"
)

func TestDSN(t *testing.T) {
	dsn, err := GetGenerator("mysql")
	if err != nil {
		t.FailNow()
	}
	mysqlDSN := dsn.(*MySQLDSN)
	mysqlDSN.Apply(
		SetMySQLUsername("admin", "admin1234"),
		SetMySQLDatabase("example"),
		SetMySQLCharset("utf8mb4"),
		SetMySQLParseTime(true),
	)
	if dsn.Gen() != "admin:admin1234@tcp(localhost:3306)/example?charset=utf8mb4&parseTime=true" {
		t.FailNow()
	}
}

func TestTimeout(t *testing.T) {
	cDSNFormat := "%s&"
	fmt.Printf(cDSNFormat, time.Second*20)
}
