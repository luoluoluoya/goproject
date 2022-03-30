package schema

import (
	"goproject/orm/dialect"
	"reflect"
	"testing"
)

func TestParseModel(t *testing.T) {
	type person struct {
		Id         int    `orm:"primary key"`
		Name       string `orm:"unique key"`
		Age        uint
		IsDeleted  bool
		unexported bool
	}
	p := &person{}
	m, _ := dialect.GetDialect("mysql")
	s := ParseModel(p, m)
	if s.Name != "person" {
		t.Fatalf("person.TableName excepted person")
	}
	fields := []string{"Id", "Name", "Age", "IsDeleted"}
	if !reflect.DeepEqual(s.FieldNames, fields) {
		t.Fatalf("person.FieldNames excepted %v", fields)
	}
	cases := map[string]Field{
		"Id":        {Name: "Id", Type: "integer", Tag: "primary key"},
		"Name":      {Name: "Name", Type: "text", Tag: "unique key"},
		"Age":       {Name: "Age", Type: "integer", Tag: ""},
		"IsDeleted": {Name: "IsDeleted", Type: "bool", Tag: ""},
	}
	for fieldName, expect := range cases {
		field := s.GetField(fieldName)
		if field.Name != expect.Name || field.Type != expect.Type || field.Tag != expect.Tag {
			//fmt.Println(field.Name, field.Type, field.Tag)
			t.Fatalf("person.%s type:%s tag:%s", expect.Name, expect.Tag, expect.Tag)
		}
	}
}
