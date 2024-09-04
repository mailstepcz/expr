package expr

import (
	"testing"

	"github.com/fealsamh/go-utils/nocopy"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func TestEq(t *testing.T) {
	req := require.New(t)

	e := Eq{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("Var = $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}

func TestNeq(t *testing.T) {
	req := require.New(t)

	e := Neq{Ident: "Var", Value: 1234}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("Var <> $1", nocopy.String(b))
	req.Equal([]interface{}{1234}, args)
}

func TestEqAny(t *testing.T) {
	req := require.New(t)

	e := EqAny{Ident: "Var", Values: []interface{}{1, 2, 3}}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("Var = ANY($1)", nocopy.String(b))
	req.Equal([]interface{}{pq.Array([]interface{}{1, 2, 3})}, args)
}

func TestNeqAll(t *testing.T) {
	req := require.New(t)

	e := NeqAll{Ident: "Var", Values: []interface{}{1, 2, 3}}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("Var <> ALL($1)", nocopy.String(b))
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
	req.Equal("(Var1 = $1 AND Var2 = $2 AND Var3 = $3)", nocopy.String(b))
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
	req.Equal("(Var1 = $1 OR Var2 = $2 OR Var3 = $3)", nocopy.String(b))
	req.Equal([]interface{}{1, 2, 3}, args)
}

func TestIsNullAny(t *testing.T) {
	req := require.New(t)

	e := IsNull{Ident: "Var"}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("Var IS NULL", nocopy.String(b))
	req.Equal([]interface{}(nil), args)
}

func TestIsNotNullAny(t *testing.T) {
	req := require.New(t)

	e := IsNotNull{Ident: "Var"}
	b, args := new(PostgresHandler).ToSQL(e, nil, nil)
	req.Equal("Var IS NOT NULL", nocopy.String(b))
	req.Equal([]interface{}(nil), args)
}
