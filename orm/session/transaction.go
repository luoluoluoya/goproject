package session

import "goproject/orm/log"

func (s *Session) Begin() (err error) {
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (s *Session) Commit() (err error) {
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	return
}
