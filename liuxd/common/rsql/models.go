package rsql

type Identifier struct {
	Val string
}

type Value interface{ ValueName() string }

type StringValue struct{ Value string }

func (StringValue) ValueName() string { return "string" }

type BooleanValue struct{ Value bool }

func (BooleanValue) ValueName() string { return "bool" }

type DateValue struct{ Value string }

func (DateValue) ValueName() string { return "date" }

type DateTimeValue struct{ Value string }

func (DateTimeValue) ValueName() string { return "datetime" }

type IntegerValue struct{ Value int64 }

func (IntegerValue) ValueName() string { return "int" }

type DoubleValue struct{ Value float64 }

func (DoubleValue) ValueName() string { return "double" }

type ListValue struct{ Value []Value }

func (ListValue) ValueName() string { return "list" }

type Expression interface{ ExpressionName() string }

type OrExpression struct{ Items []Expression }

func (OrExpression) ExpressionName() string { return "Or" }

type AndExpression struct{ Items []Expression }

func (AndExpression) ExpressionName() string { return "And" }

type Comparison struct {
	Identifier Identifier
	Val        Value
}

func (Comparison) ExpressionName() string { return "Comparison" }

type EqualsComparison struct{ Comparison }

func (EqualsComparison) ExpressionName() string { return "==" }

type NotEqualsComparison struct{ Comparison }

func (NotEqualsComparison) ExpressionName() string { return "!=" }

type LikeComparison struct{ Comparison }

func (LikeComparison) ExpressionName() string { return "~=" }

type NotLikeComparison struct{ Comparison }

func (NotLikeComparison) ExpressionName() string { return "!~=" }

type GreaterThanComparison struct{ Comparison }

func (GreaterThanComparison) ExpressionName() string { return ">" }

type GreaterThanOrEqualsComparison struct{ Comparison }

func (GreaterThanOrEqualsComparison) ExpressionName() string { return ">=" }

type LessThanComparison struct{ Comparison }

func (LessThanComparison) ExpressionName() string { return "<" }

type LessThanOrEqualsComparison struct{ Comparison }

func (LessThanOrEqualsComparison) ExpressionName() string { return "<=" }

type InComparison struct{ Comparison }

func (InComparison) ExpressionName() string { return "=in=" }

type NotInComparison struct{ Comparison }

func (NotInComparison) ExpressionName() string { return "=out=" }
