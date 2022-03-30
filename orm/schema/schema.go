package schema

import (
	"goproject/orm/dialect"
	"reflect"
)

// 字段
type Field struct {
	Name string
	Type string
	Tag  string
}

// 表
type Schema struct {
	Model        interface{}
	Name         string
	Fields       []*Field
	FieldNames   []string
	fieldMapping map[string]*Field
}

func (s *Schema) GetField(fieldName string) *Field {
	return s.fieldMapping[fieldName]
}

func ParseModel(model interface{}, dialect dialect.Dialect) *Schema {
	s := &Schema{fieldMapping: make(map[string]*Field)}
	typ := reflect.Indirect(reflect.ValueOf(model)).Type()
	s.Model = model
	s.Name = typ.Name()
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Anonymous || !f.IsExported() {
			continue
		}
		field := &Field{
			Name: f.Name,
			Type: dialect.DataTypeOf(reflect.Indirect(reflect.New(f.Type))),
			Tag:  f.Tag.Get("orm"),
		}
		s.Fields = append(s.Fields, field)
		s.FieldNames = append(s.FieldNames, field.Name)
		s.fieldMapping[field.Name] = field
	}
	return s
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
