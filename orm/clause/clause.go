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
	sql, vars := generators[typ](vars...)
	c.sql[typ] = sql
	c.sqlVars[typ] = vars
}

func (c *Clause) Build(orders ...Type) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
