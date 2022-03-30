package session

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"goproject/orm/dialect"
	"testing"
)

func TestSession_Table(t *testing.T) {
	type person_23u42hkj42 struct {
		Id         int    `orm:"primary key"`
		Name       string `orm:"comment '用户名称'"`
		Age        uint
		IsDeleted  bool
		unexported bool
	}

	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?charset=utf8")
	dailect, _ := dialect.GetDialect("mysql")
	s := New(db, dailect)
	s = s.Model(&person_23u42hkj42{})

	s.DropTable()

	if s.ExistsTable() {
		t.Fatalf("table person_23u42hkj42 not exists")
	}

	err = s.CreateTable()
	if err != nil {
		t.Fatalf("table person_23u42hkj42 create failed")
	}

	//if !s.ExistsTable() {
	//	t.Fatalf("table person_23u42hkj42 exists")
	//}

	err = s.DropTable()
	if err != nil {
		t.Fatalf("table person_23u42hkj42 drop failed")
	}

	if s.ExistsTable() {
		t.Fatalf("table person_23u42hkj42 not exists")
	}
}
