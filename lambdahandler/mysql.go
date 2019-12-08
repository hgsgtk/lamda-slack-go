package lambdahandler

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const defaultSQLMode = "TRADITIONAL,NO_AUTO_VALUE_ON_ZERO,ONLY_FULL_GROUP_BY"

const defaultLoc = "Asia%2FTokyo"

const driverName = "mysql"

func NewDBConn(config DBConfig) (*sql.DB, error) {
	ds := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?sql_mode='%s'&parseTime=true&loc=%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
		defaultSQLMode,
		defaultLoc)

	db, err := sql.Open(driverName, ds)
	if err != nil {
		return nil, errors.Wrap(err, "connection database error")
	}

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, errors.Wrap(err, "ping database error")
	}

	return db, nil
}
