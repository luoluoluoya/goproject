package dialect

import "reflect"

var dialectMap = make(map[string]Dialect)

type Dialect interface {
	DataTypeOf(typ reflect.Value) string
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) Dialect {
	return dialectMap[name]
}
