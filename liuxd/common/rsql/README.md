# go-rsql

RSQL like parse (I'm not yet implements all tokens).

## Installation:

- `go get github.com/dohrm/go-rsql`

## Usage:

```go
package main

import (
	"fmt"
	"github.com/dohrm/go-rsql"
)

func main() {
	result, _ := rsql.Parse(`name=='Luke Skywalker' or actor.name=="Mark"`)
	
	
	fmt.Println(fmt.Sprintf("%+v", result))
}

```

## Grammar:

```
or         : and ('OR' | 'or' and)*
and        : constraint ('AND' | 'and' constraint)*
constraint : group | comparison
group      : '(' or ')'
comparison : identifier comparator arguments
identifier : [a-zA-Z0-9]+('.'[a-zA-Z0-9]+)*
comparator : '==' | '!=' | '==~' | '!=~' | '>' | '>=' | '<' | '<=' | '=in=' | '=out='
arguments  : '(' listValue ')' | value
value      : int | double | string | date | datetime | boolean
listValue  : value(','value)*
int        : [0-9]+
double     : [0-9]+'.'[0-9]*
string     : '"'.*'"' | '\''.*'\''
date       : [0-9]{4}'-'[0-9]{2}'-'\[0-9]{2}
datetime   : date'T'[0-9]{2}':'[0-9]{2}':'[0-9]{2}('Z' | (('+'|'-')[0-9]{2}(':')?[0-9]{2}))?
boolean    : 'true' | 'false'

```

Examples :

- `name=='Luke Skywalker' or actor.name=="Mark"`
- `age=in=(1,2,3) or age >= 42`
- `(movie==~'.*H2G2.*' or move=='Seven') and (budget<1500 or rating>=6)`
- `birthDate==1890-08-20`



