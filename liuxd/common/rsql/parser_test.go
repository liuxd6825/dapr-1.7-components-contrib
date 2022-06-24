package rsql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParse_Invalid(t *testing.T) {
	_, err := Parse("")
	assert.Error(t, err)
	_, err = Parse("toto==toto")
	assert.Error(t, err)
	_, err = Parse(`(toto=="titi"`)
	assert.Error(t, err)
	_, err = Parse(`toto=="titi")`)
	assert.Error(t, err)
	_, err = Parse(`toto=="titi" and and`)
	assert.Error(t, err)

	// TODO More test
}

func TestParse_Comparison(t *testing.T) {
	v, _ := Parse("toto==42")
	assert.IsType(t, EqualsComparison{}, v)
	v, _ = Parse("toto!=42")
	assert.IsType(t, NotEqualsComparison{}, v)
	v, _ = Parse("toto!=~42")
	assert.IsType(t, NotLikeComparison{}, v)
	v, _ = Parse("toto>=42")
	assert.IsType(t, GreaterThanOrEqualsComparison{}, v)
	v, _ = Parse("toto<=42")
	assert.IsType(t, LessThanOrEqualsComparison{}, v)
	v, _ = Parse("toto=in=42")
	assert.IsType(t, InComparison{}, v)
	v, _ = Parse("toto=out=(42,43)")
	assert.IsType(t, NotInComparison{}, v)
}

func TestParse_And(t *testing.T) {
	v, _ := Parse("toto==42 and titi==42")
	assert.IsType(t, AndExpression{}, v)
	assert.Len(t, v.(AndExpression).Items, 2)
}

func TestParse_Or(t *testing.T) {
	v, _ := Parse("toto==42 or titi==42")
	assert.IsType(t, OrExpression{}, v)
	assert.Len(t, v.(OrExpression).Items, 2)

}

func TestParse_AndOrPriority(t *testing.T) {
	v, _ := Parse(`toto==42 and titi==42 or tutu==666 and titi==666 and tutu==42 or bob=="bibi"`)
	assert.IsType(t, OrExpression{}, v)
	orExp := v.(OrExpression)
	assert.Len(t, orExp.Items, 3)
	assert.IsType(t, AndExpression{}, orExp.Items[0])
	assert.Len(t, orExp.Items[0].(AndExpression).Items, 2)
	assert.IsType(t, AndExpression{}, orExp.Items[1])
	assert.Len(t, orExp.Items[1].(AndExpression).Items, 3)
	assert.IsType(t, EqualsComparison{}, orExp.Items[2])
}

func TestParse_WithGroups(t *testing.T) {
	v, _ := Parse(`((toto==42 and titi!=42) and (tutu=in=(666,777) or titi<666)) or toto==1`)
	assert.IsType(t, OrExpression{}, v)
	assert.Len(t, v.(OrExpression).Items, 2)
	v1 := v.(OrExpression).Items[0]
	v2 := v.(OrExpression).Items[1]
	assert.IsType(t, AndExpression{}, v1)
	assert.IsType(t, EqualsComparison{}, v2)
	v11 := v1.(AndExpression).Items[0]
	v12 := v1.(AndExpression).Items[1]
	assert.IsType(t, AndExpression{}, v11)
	assert.IsType(t, OrExpression{}, v12)

	assert.IsType(t, EqualsComparison{}, v11.(AndExpression).Items[0])
	assert.IsType(t, NotEqualsComparison{}, v11.(AndExpression).Items[1])
	assert.IsType(t, InComparison{}, v12.(OrExpression).Items[0])
	assert.IsType(t, LessThanComparison{}, v12.(OrExpression).Items[1])

	assert.Len(t, v12.(OrExpression).Items[0].(InComparison).Val.(ListValue).Value, 2)

}

func BenchmarkParse(b *testing.B) {
	query := `(((toto==42 and titi!=42) and (tutu=in=(666,777) or titi<666)) or toto==1) and (tata<=12 or teet==1890-08-20)`

	for n := 0; n < b.N; n++ {
		_, e := Parse(query)
		if e != nil {
			panic(e)
		}
	}
}
