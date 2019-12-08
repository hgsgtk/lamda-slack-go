package lambdahandler

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type dbConfig struct {
	host     string
	user     string
	password string
	name     string
	sqlMode  string
	location string
}

func (c *dbConfig) GetDS() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:3306)/%s?sql_mode='%s'&parseTime=true&loc=%s",
		c.user,
		c.password,
		c.host,
		c.name,
		c.sqlMode,
		c.location)
}

const defaultSQLMode = "TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY"

const defaultLoc = "Asia%2FTokyo"

const driverName = "mysql"

func NewDBConn(host, user, password, name string) (*sql.DB, error) {
	dc := &dbConfig{
		host:     host,
		user:     user,
		password: password,
		name:     name,
		sqlMode:  defaultSQLMode,
		location: defaultLoc,
	}

	db, err := sql.Open(driverName, dc.GetDS())
	if err != nil {
		return nil, errors.Wrap(err, "connection database error")
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, errors.Wrap(err, "ping database error")
	}

	return db, nil
}
