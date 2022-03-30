package session

import (
	"fmt"
	"goproject/orm/log"
	"goproject/orm/schema"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || s.refTable.Name != reflect.TypeOf(value).Name() {
		s.refTable = schema.ParseModel(value, s.dialect)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("session.model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.RefTable()
	sql := strings.Builder{}
	sql.WriteString(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(", table.Name))
	var columns []string
	for _, fieldName := range table.FieldNames {
		field := table.GetField(fieldName)
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	sql.WriteString(strings.Join(columns, ",") + ")")
	fmt.Println(sql)
	_, err := s.Raw(sql.String()).Exec()
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *Session) DropTable() error {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)
	_, err := s.Raw(sql).Exec()
	if err != nil {
		log.Error(err)
	}
	return err
}

func (s *Session) ExistsTable() bool {
	sql, param := s.dialect.TableExistSQL(s.refTable.Name)
	row := s.Raw(sql, param).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == s.refTable.Name
}
