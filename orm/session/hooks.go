package session

import (
	"goproject/orm/log"
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

//  钩子函数签名 functionName(s *Session) error

func (s *Session) CallMethod(method string, value interface{}) {
	fn := reflect.ValueOf(s.RefTable().Model).MethodByName(method)
	if value != nil {
		fn = reflect.ValueOf(value).MethodByName(method)
	}
	params := []reflect.Value{reflect.ValueOf(s)}
	if fn.IsValid() {
		if r := fn.Call(params); len(r) > 0 {
			if err, ok := r[0].Interface().(error); ok {
				log.Error(err)
			}
		}
	}
}
