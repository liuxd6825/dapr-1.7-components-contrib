package rsql

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ParseProcess(t *testing.T) {
	input := "((toto==32 and userId=='001' ) or (user=='admin' and sex==1)) and user==~'000'"
	p := &process{}
	err := ParseProcess(input, p)
	p.Print()
	assert.Error(t, err)
}
