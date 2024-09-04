package expr

import (
	"strconv"

	"github.com/lib/pq"
)

// Handler is a handler for generating expressions from an AST.
type Handler interface {
	NextPlaceholder([]byte) []byte
	WrapSlice([]interface{}) interface{}
}

// SQLHandler is a handler for generating SQL expressions.
type SQLHandler interface {
	Handler
	ToSQL([]byte, []interface{}) ([]byte, []interface{})
}

// PostgresHandler is a handler for PostgreSQL.
type PostgresHandler struct {
	idx int64
}

// NextPlaceholder returns the next placeholer for an expression's argument.
func (h *PostgresHandler) NextPlaceholder(b []byte) []byte {
	h.idx++
	b = append(b, '$')
	return strconv.AppendInt(b, h.idx, 10)
}

// WrapSlice wraps an expression's slice argument.
func (h *PostgresHandler) WrapSlice(sl []interface{}) interface{} {
	return pq.Array(sl)
}

// ToSQL generates the SQL representation of the AST.
func (h *PostgresHandler) ToSQL(e Expr, b []byte, args []interface{}) ([]byte, []interface{}) {
	return e.Linearise(h, b, args)
}

// Expr is an AST node representing an expression.
type Expr interface {
	Linearise(Handler, []byte, []interface{}) ([]byte, []interface{})
}

// Eq is an AST node for equality.
type Eq struct {
	Ident string
	Value interface{}
}

// Linearise linearises the AST.
func (e Eq) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, e.Ident...)
	b = append(b, " = "...)
	return h.NextPlaceholder(b), append(args, e.Value)
}

// Neq is an AST node for inequality.
type Neq struct {
	Ident string
	Value interface{}
}

// Linearise linearises the AST.
func (e Neq) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, e.Ident...)
	b = append(b, " <> "...)
	return h.NextPlaceholder(b), append(args, e.Value)
}

// EqAny is an AST node for equality with ANY.
type EqAny struct {
	Ident  string
	Values []interface{}
}

// Linearise linearises the AST.
func (e EqAny) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, e.Ident...)
	b = append(b, " = ANY("...)
	b = h.NextPlaceholder(b)
	return append(b, ')'), append(args, h.WrapSlice(e.Values))
}

// NeqAll is an AST node for inequality with ALL.
type NeqAll struct {
	Ident  string
	Values []interface{}
}

// Linearise linearises the AST.
func (e NeqAll) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, e.Ident...)
	b = append(b, " <> ALL("...)
	b = h.NextPlaceholder(b)
	return append(b, ')'), append(args, h.WrapSlice(e.Values))
}

// And is an AST node for conjunction.
type And struct {
	Exprs []Expr
}

// Linearise linearises the AST.
func (e And) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, '(')
	for i, expr := range e.Exprs {
		if i > 0 {
			b = append(b, " AND "...)
		}
		b, args = expr.Linearise(h, b, args)
	}
	return append(b, ')'), args
}

// Or is an AST node for disjunction.
type Or struct {
	Exprs []Expr
}

// Linearise linearises the AST.
func (e Or) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, '(')
	for i, expr := range e.Exprs {
		if i > 0 {
			b = append(b, " OR "...)
		}
		b, args = expr.Linearise(h, b, args)
	}
	return append(b, ')'), args
}

// IsNull is an AST node for nullity testing.
type IsNull struct {
	Ident string
}

// Linearise linearises the AST.
func (e IsNull) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, e.Ident...)
	b = append(b, " IS NULL"...)
	return b, args
}

// IsNotNull is an AST node for non-nullity testing.
type IsNotNull struct {
	Ident string
}

// Linearise linearises the AST.
func (e IsNotNull) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, e.Ident...)
	b = append(b, " IS NOT NULL"...)
	return b, args
}

// Lt is an AST node for representing the "less than" comparison operation.
type Lt struct {
	Ident string
	Value interface{}
}

// Linearise linearises the AST.
func (e Lt) Linearise(h Handler, b []byte, args []interface{}) ([]byte, []interface{}) {
	b = append(b, e.Ident...)
	b = append(b, " < "...)
	return h.NextPlaceholder(b), append(args, e.Value)
}

var (
	_ Expr = Eq{}
	_ Expr = Neq{}
	_ Expr = EqAny{}
	_ Expr = NeqAll{}
	_ Expr = And{}
	_ Expr = IsNull{}
	_ Expr = IsNotNull{}
	_ Expr = Lt{}

	_ Handler = (*PostgresHandler)(nil)
)
