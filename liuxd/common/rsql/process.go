package rsql

import (
	"errors"
	"fmt"
)

type Process interface {
	OnAndItem()
	OnAndStart()
	OnAndEnd()
	OnOrItem()
	OnOrStart()
	OnOrEnd()
	OnEquals(name string, value interface{}, rValue Value)
	OnNotEquals(name string, value interface{}, rValue Value)
	OnLike(name string, value interface{}, rValue Value)
	OnNotLike(name string, value interface{}, rValue Value)
	OnGreaterThan(name string, value interface{}, rValue Value)
	OnGreaterThanOrEquals(name string, value interface{}, rValue Value)
	OnLessThan(name string, value interface{}, rValue Value)
	OnLessThanOrEquals(name string, value interface{}, rValue Value)
	OnIn(name string, value interface{}, rValue Value)
	OnNotIn(name string, value interface{}, rValue Value)
	GetFilter(tenantId string) interface{}
}

type process struct {
	str string
}

func NewSqlProcess() Process {
	return &process{}
}

func (p *process) OnNotEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s != (%v)", p.str, name, value)
}

func (p *process) OnLike(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s like (%v)", p.str, name, value)
}

func (p *process) OnNotLike(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s not like %v", p.str, name, value)
}

func (p *process) OnGreaterThan(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s>%v", p.str, name, value)
}

func (p *process) OnGreaterThanOrEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s>=%v", p.str, name, value)
}

func (p *process) OnLessThan(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s<%v", p.str, name, value)
}

func (p *process) OnLessThanOrEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s <= %v", p.str, name, value)
}

func (p *process) OnIn(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s in %v", p.str, name, value)
}

func (p *process) OnNotIn(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s not in %v", p.str, name, value)
}

func (p *process) OnEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s=%v", p.str, name, value)
}

func (p *process) NotEquals(name string, value interface{}, rValue Value) {
	p.str = fmt.Sprintf("%s %s=%v", p.str, name, value)
}

func (p *process) OnAndItem() {
	p.str = fmt.Sprintf("%s and ", p.str)
}
func (p *process) OnAndStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}
func (p *process) OnAndEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}
func (p *process) OnOrItem() {
	p.str = fmt.Sprintf("%s or ", p.str)
}
func (p *process) OnOrStart() {
	p.str = fmt.Sprintf("%s(", p.str)
}
func (p *process) OnOrEnd() {
	p.str = fmt.Sprintf("%s)", p.str)
}

func (p *process) GetFilter(tenantId string) interface{} {
	return p.str
}
func (p *process) Print() {
	fmt.Print(p.str)
}

func ParseProcess(input string, process Process) error {
	if len(input) == 0 {
		return nil
	}
	expr, err := Parse(input)
	if err != nil {
		return errors.New(fmt.Sprintf("rsql %s expression error, %s", input, err.Error()))
	}
	err = parseProcess(expr, process)
	if err != nil {
		return errors.New(fmt.Sprintf("rsql %s parseProcess error, %s", input, err.Error()))
	}
	return nil
}

func parseProcess(expr Expression, process Process) error {
	switch expr.(type) {
	case AndExpression:
		ex, _ := expr.(AndExpression)
		process.OnAndStart()
		for i, e := range ex.Items {
			_ = parseProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnAndItem()
			}
		}
		process.OnAndEnd()
		break
	case OrExpression:
		ex, _ := expr.(OrExpression)
		process.OnOrStart()
		for i, e := range ex.Items {
			_ = parseProcess(e, process)
			if i < len(ex.Items)-1 {
				process.OnOrItem()
			}
		}
		process.OnOrEnd()
		break
	case NotEqualsComparison:
		ex, _ := expr.(NotEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotEquals(name, value, ex.Comparison.Val)
		break
	case EqualsComparison:
		ex, _ := expr.(EqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnEquals(name, value, ex.Comparison.Val)
		break
	case LikeComparison:
		ex, _ := expr.(LikeComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLike(name, value, ex.Comparison.Val)
		break
	case NotLikeComparison:
		ex, _ := expr.(NotLikeComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotLike(name, value, ex.Comparison.Val)
		break
	case GreaterThanComparison:
		ex, _ := expr.(GreaterThanComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThan(name, value, ex.Comparison.Val)
		break
	case GreaterThanOrEqualsComparison:
		ex, _ := expr.(GreaterThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnGreaterThanOrEquals(name, value, ex.Comparison.Val)
		break
	case LessThanComparison:
		ex, _ := expr.(LessThanComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLessThan(name, value, ex.Comparison.Val)
		break
	case LessThanOrEqualsComparison:
		ex, _ := expr.(LessThanOrEqualsComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnLessThanOrEquals(name, value, ex.Comparison.Val)
		break
	case InComparison:
		ex, _ := expr.(InComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnIn(name, value, ex.Comparison.Val)
		break
	case NotInComparison:
		ex, _ := expr.(NotInComparison)
		name := ex.Comparison.Identifier.Val
		value := getValue(ex.Comparison.Val)
		process.OnNotIn(name, value, ex.Comparison.Val)
		break
	}
	return nil
}

func getValue(val Value) interface{} {
	var value interface{}
	switch val.(type) {
	case IntegerValue:
		value = val.(IntegerValue).Value
		break
	case BooleanValue:
		value = val.(BooleanValue).Value
		break
	case StringValue:
		value = fmt.Sprintf(`"%v"`, val.(StringValue).Value)
		break
	case DateTimeValue:
		value = fmt.Sprintf(`"%v"`, val.(DateTimeValue).Value)
		break
	case DoubleValue:
		value = val.(DoubleValue).Value
		break
	}
	return value
}
