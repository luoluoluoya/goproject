package dialect

import (
	"reflect"
	"testing"
)

func TestMysql_DataTypeOf(t *testing.T) {
	mysql := &mysql{}
	cases := map[interface{}]string{
		interface{}(32):        "integer",
		interface{}(uint8(19)): "integer",
		interface{}(int64(32)): "bigint",
		interface{}(true):      "bool",
		interface{}(12.1):      "real",
		interface{}("lisi"):    "text",
		//interface{}([]int{1, 2, 3}): "blob",
	}
	for k, v := range cases {
		if r := mysql.DataTypeOf(reflect.ValueOf(k)); r != v {
			t.Fatalf("%v except type: %v. returend: %v\n", k, v, r)
		}
	}
}
