package rsql

import (
	"fmt"
	"strconv"
)

/*

or         : and ('OR' and)*
and        : constraint ('AND' constraint)*
constraint : group | comparison
group      : '(' or ')'
comparison : 'IdentifierToken' 'Comparator' arguments
arguments  : '(' listValue ')' | value
value      : 'IntegerToken' | 'DoubleToken' | 'DateToken' | 'DateTimeToken' | 'BooleanToken' | 'StringToken'
listValue  : value (',' value)*
*/

type iterator struct {
	length int
	idx    int
	items  []Token
}

func newIterator(items []Token) *iterator {
	return &iterator{
		length: len(items),
		idx:    0,
		items:  items,
	}
}

func (t *iterator) get(idx int) Token {
	if idx >= t.length {
		return Token{
			Type:  EOFToken,
			Value: "",
			Pos:   t.length,
		}
	}
	return t.items[idx]
}

func (t *iterator) current() Token {
	return t.get(t.idx)
}

func (t *iterator) currentAndMove(potentialCount ...int) Token {
	count := 1
	if len(potentialCount) > 0 {
		count = potentialCount[0]
	}
	result := t.current()
	if result.Type != EOFToken {
		t.idx += count
	}
	return result
}

func Parse(input string) (Expression, error) {
	items, err := NewLexer(input).Parse()
	if err != nil {
		return nil, err
	}
	return or(newIterator(items))
}

func section(tokens *iterator, separator TokenType, apply func(*iterator) (Expression, error)) ([]Expression, error) {
	var result []Expression
	idx := tokens.idx
	cursor := tokens.idx
	opened := 0
	for idx < tokens.length {
		c := tokens.get(idx)
		if c.Type == LeftParenToken {
			opened++
		} else if c.Type == RightParenToken {
			opened--
			if opened < 0 {
				return nil, fmt.Errorf("invalid parentheses")
			}
		} else if c.Type == separator && opened == 0 {
			items := tokens.items[cursor:idx]
			next, err := apply(newIterator(items))
			if err != nil {
				return nil, err
			}
			result = append(result, next)
			cursor = idx + 1
		}
		idx += 1
	}
	if idx > cursor {
		items := tokens.items[cursor:idx]
		next, err := apply(newIterator(items))
		if err != nil {
			return nil, err
		}
		result = append(result, next)
		cursor = idx + 1
	}
	tokens.currentAndMove(cursor - tokens.idx)
	return result, nil
}

func or(tokens *iterator) (Expression, error) {
	items, err := section(tokens, OrToken, and)
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, fmt.Errorf("invalid 'or' section %+v", tokens.current())
	} else if len(items) == 1 {
		return items[0], nil
	}
	return OrExpression{
		Items: items,
	}, nil
}

func and(tokens *iterator) (Expression, error) {
	items, err := section(tokens, AndToken, constraint)
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, fmt.Errorf("invalid 'and' section %+v", tokens.current())
	} else if len(items) == 1 {
		return items[0], nil
	}
	return AndExpression{
		Items: items,
	}, nil
}

func constraint(tokens *iterator) (Expression, error) {
	if tokens.current().Type == LeftParenToken {
		return group(tokens)
	}
	return comparison(tokens)
}

func group(tokens *iterator) (Expression, error) {
	tokens.currentAndMove()
	opened := 1
	idx := tokens.idx
	for idx < tokens.length {
		if tokens.get(idx).Type == LeftParenToken {
			opened++
		} else if tokens.get(idx).Type == RightParenToken {
			opened--
			if opened == 0 {
				break
			}
		}
		idx++
	}
	if opened > 0 {
		return nil, fmt.Errorf("closed parentheses don't match")
	}
	newIterator := newIterator(tokens.items[tokens.idx:idx])
	tokens.currentAndMove(idx - tokens.idx)
	return or(newIterator)
}

func comparison(tokens *iterator) (Expression, error) {
	id, err := identifier(tokens)
	if err != nil {
		return nil, err
	}
	comparator := tokens.currentAndMove()
	args, err := arguments(tokens)
	if err != nil {
		return nil, err
	}
	switch comparator.Type { // TODO Manage that directly to Tokens.
	case EqualsToken:
		return EqualsComparison{Comparison{id, args}}, nil
	case NotEqualsToken:
		return NotEqualsComparison{Comparison{id, args}}, nil
	case LikeToken:
		return LikeComparison{Comparison{id, args}}, nil
	case NotLikeToken:
		return NotLikeComparison{Comparison{id, args}}, nil
	case GreaterToken:
		return GreaterThanComparison{Comparison{id, args}}, nil
	case GreaterOrEqualsToken:
		return GreaterThanOrEqualsComparison{Comparison{id, args}}, nil
	case LessToken:
		return LessThanComparison{Comparison{id, args}}, nil
	case LessOrEqualsToken:
		return LessThanOrEqualsComparison{Comparison{id, args}}, nil
	case InToken:
		tmp, ok := args.(ListValue)
		var lv ListValue
		if !ok {
			lv = ListValue{Value: []Value{tmp}}
		} else {
			lv = tmp
		}
		return InComparison{Comparison{id, lv}}, nil
	case NotInToken:
		tmp, ok := args.(ListValue)
		var lv ListValue
		if !ok {
			lv = ListValue{Value: []Value{tmp}}
		} else {
			lv = tmp
		}
		return NotInComparison{Comparison{id, lv}}, nil
	}
	return nil, fmt.Errorf("'comparator not managed for expression")
}

func identifier(tokens *iterator) (Identifier, error) {
	token := tokens.currentAndMove()
	if token.Type != IdentifierToken {
		return Identifier{}, fmt.Errorf("must be an Identifier")
	}
	return Identifier{token.Value}, nil
}

func arguments(tokens *iterator) (Value, error) {
	if tokens.current().Type == LeftParenToken {
		return valueList(tokens)
	} else {
		return value(tokens)
	}
}

func valueList(tokens *iterator) (Value, error) {
	tokens.currentAndMove() // Remove first (
	var items []Value
	current := tokens.current()
	for current.Type != RightParenToken {
		c, err := value(tokens)
		if err != nil {
			return nil, err
		}
		items = append(items, c)
		current = tokens.currentAndMove()
		if current.Type != RightParenToken && current.Type != CommaToken {
			return nil, fmt.Errorf("invalid list format, next must be comma or Right Parent")
		}
	}
	return ListValue{items}, nil
}

func value(tokens *iterator) (Value, error) {
	v := tokens.currentAndMove()
	switch v.Type {
	case StringToken:
		return StringValue{v.Value}, nil
	case BooleanToken:
		return BooleanValue{v.Value == "true"}, nil
	case DoubleToken:
		c, err := strconv.ParseFloat(v.Value, 64)
		if err != nil {
			return nil, err
		}
		return DoubleValue{c}, nil
	case IntegerToken:
		c, err := strconv.ParseInt(v.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return IntegerValue{c}, nil
	case DateToken:
		return DateValue{v.Value}, nil
	case DateTimeToken:
		return DateTimeValue{v.Value}, nil
	}
	return nil, fmt.Errorf("invalid type")
}
