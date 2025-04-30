package expr

import (
	"testing"

	"github.com/mailstepcz/go-utils/nocopy"
	"github.com/lib/pq"
	"github.com/mailstepcz/testutils/testcond"
	"github.com/stretchr/testify/require"
)

func TestAppendIdent(t *testing.T) {
	b := appendIdent(nil, "abcd")
	testcond.Equal(t, `"abcd"`, string(b))

	b = appendIdent(nil, "ab.cd")
	testcond.Equal(t, `"ab"."cd"`, string(b))
}

func TestEq(t *testing.T) {
	req := require.New(t)

	e := Eq{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" = $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}

func TestNeq(t *testing.T) {
	req := require.New(t)

	e := Neq{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" <> $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}

func TestEqAny(t *testing.T) {
	req := require.New(t)

	e := EqAny{Ident: "Var", Values: []interface{}{1, 2, 3}}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" = ANY($1)", nocopy.String(b))
	req.Equal([]interface{}{pq.Array([]interface{}{1, 2, 3})}, args)
}

func TestNeqAll(t *testing.T) {
	req := require.New(t)

	e := NeqAll{Ident: "Var", Values: []interface{}{1, 2, 3}}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" <> ALL($1)", nocopy.String(b))
	req.Equal([]interface{}{pq.Array([]interface{}{1, 2, 3})}, args)
}

func TestAnd(t *testing.T) {
	req := require.New(t)

	e := And{[]Expr{
		Eq{Ident: "Var1", Value: 1},
		Eq{Ident: "Var2", Value: 2},
		Eq{Ident: "Var3", Value: 3},
	}}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("(\"Var1\" = $1 AND \"Var2\" = $2 AND \"Var3\" = $3)", nocopy.String(b))
	req.Equal([]interface{}{1, 2, 3}, args)
}

func TestOr(t *testing.T) {
	req := require.New(t)

	e := Or{[]Expr{
		Eq{Ident: "Var1", Value: 1},
		Eq{Ident: "Var2", Value: 2},
		Eq{Ident: "Var3", Value: 3},
	}}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("(\"Var1\" = $1 OR \"Var2\" = $2 OR \"Var3\" = $3)", nocopy.String(b))
	req.Equal([]interface{}{1, 2, 3}, args)
}

func TestContains(t *testing.T) {
	req := require.New(t)

	e := Or{[]Expr{
		Contains{Ident: "Var1", Value: 1},
		Contains{Ident: "Var2", Value: 2},
		Contains{Ident: "Var3", Value: 3},
	}}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("($1 = ANY(\"Var1\") OR $2 = ANY(\"Var2\") OR $3 = ANY(\"Var3\"))", nocopy.String(b))
	req.Equal([]interface{}{1, 2, 3}, args)
}

func TestIsNullAny(t *testing.T) {
	req := require.New(t)

	e := IsNull{Ident: "Var"}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" IS NULL", nocopy.String(b))
	req.Equal([]interface{}(nil), args)
}

func TestIsNotNullAny(t *testing.T) {
	req := require.New(t)

	e := IsNotNull{Ident: "Var"}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" IS NOT NULL", nocopy.String(b))
	req.Equal([]interface{}(nil), args)
}

func TestLt(t *testing.T) {
	req := require.New(t)

	e := Lt{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" < $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}

func TestLte(t *testing.T) {
	req := require.New(t)

	e := Lte{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" <= $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}

func TestGt(t *testing.T) {
	req := require.New(t)

	e := Gt{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" > $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}

func TestGte(t *testing.T) {
	req := require.New(t)

	e := Gte{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("\"Var\" >= $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}
