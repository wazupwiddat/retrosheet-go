package db

import (
	"fmt"
	"time"

	"github.com/gocraft/dbr"
)

func Open(driver, dsn string) (*dbr.Connection, error) {
	if dsn == "" {
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?parseTime=true",
			"root",
			"",
			"localhost",
			"3306",
			"baseball",
		)
	}

	conn, err := dbr.Open(driver, dsn, nil)
	if err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(20)
	conn.SetMaxIdleConns(20)
	conn.SetConnMaxLifetime(time.Second * 60)

	return conn, nil
}
