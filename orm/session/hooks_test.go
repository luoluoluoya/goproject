package session

import (
	"database/sql"
	"fmt"
	"goproject/orm/dialect"
	"goproject/orm/log"
	"testing"
)

type account struct {
	ID       int `geeorm:"PRIMARY KEY"`
	Password string
}

func (account *account) BeforeInsert(s *Session) error {
	log.Info("before insert", account)
	account.ID += 1000
	return nil
}

func (account *account) AfterQuery(s *Session) error {
	log.Info("after query", account)
	account.Password = "******"
	return nil
}

func TestSession_CallMethod(t *testing.T) {
	db, _ := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	dailect, _ := dialect.GetDialect("mysql")
	s := New(db, dailect)
	s = s.Model(&account{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_, _ = s.Insert(&account{1, "123456"}, &account{2, "qwerty"})

	u := &account{}

	err := s.First(u)
	fmt.Println(u)
	if err != nil || u.ID != 1001 || u.Password != "******" {
		t.Fatal("Failed to call hooks after query, got", u)
	}
	_ = s.DropTable()
}
