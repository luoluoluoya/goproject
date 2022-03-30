package clause

import "strings"

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
)

type Clause struct {
	sql     map[Type]string
	sqlVars map[Type][]interface{}
}

func (c *Clause) Set(typ Type, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[Type]string)
	}
	if c.sqlVars == nil {
		c.sqlVars = make(map[Type][]interface{})
	}
	c.sql[typ], c.sqlVars[typ] = generators[typ](vars...)
}

func (c *Clause) Build(types ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, typ := range types {
		if sql, ok := c.sql[typ]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[typ]...)

		}
	}
	return strings.Join(sqls, " "), vars
}
