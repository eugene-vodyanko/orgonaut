package oracle

import (
	"database/sql"
	_ "github.com/sijms/go-ora/v2"
	"time"
)

type Oracle struct {
	Db *sql.DB
}

func New(username, password, url, schema string, maxOpenConns, maxIdleConns, maxLifetime, maxIdleTime int) (*Oracle, error) {

	db, err := sql.Open("oracle", makeDatabaseURL(username, password, url))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(maxIdleTime) * time.Second)
	db.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	// Option not used: negatively affects performance
	// db.SetMaxIdleConns(maxIdleConns)

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Oracle{db}, nil
}

func (o *Oracle) Ping() error {
	return o.Db.Ping()
}

func (o *Oracle) Close() error {
	return o.Db.Close()
}

func makeDatabaseURL(username, password, url string) string {
	return "oracle://" + username + ":" + password + "@" + url
}
