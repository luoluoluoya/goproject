package orm

import (
	"database/sql"
	"fmt"
	"goproject/orm/dialect"
	"goproject/orm/log"
	"goproject/orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver string, source string) (e *Engine, err error) {
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()
	db, err := sql.Open(driver, source)
	if err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}
	dialect, ok := dialect.GetDialect(driver)
	if !ok {
		err = fmt.Errorf("driver %s not registered", driver)
		return
	}
	log.Infof("Connect [%s] %s success\n", driver, source)
	return &Engine{db: db, dialect: dialect}, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Errorf("Database failed to close: %v", err)
		return
	}
	log.Info("Database closed")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
