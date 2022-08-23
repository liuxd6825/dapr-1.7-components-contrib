/*  Filter
// - name=="Kill Bill";year=gt=2003
// - name=="Kill Bill" and year>2003
// - genres=in=(sci-fi,action);(director=='Christopher Nolan',actor==*Bale);year=ge=2000
// - genres=in=(sci-fi,action) and (director=='Christopher Nolan' or actor==*Bale) and year>=2000
// - director.lastName==Nolan;year=ge=2000;year=lt=2010
// - director.lastName==Nolan and year>=2000 and year<2010
// - genres=in=(sci-fi,action);genres=out=(romance,animated,horror),director==Que*Tarantino
// - genres=in=(sci-fi,action) and genres=out=(romance,animated,horror) or director==Que*Tarantino
// or         : and ('OR' | 'or' and) *
// and        : constraint ('AND' | 'and' constraint)*
// constraint : group | comparison
// group      : '(' or ')'
// comparison : identifier comparator arguments
// identifier : [a-zA-Z0-9]+('.'[a-zA-Z0-9]+)*
// comparator : '==' | '!=' | '==~' | '!=~' | '>' | '>=' | '<' | '<=' | '=in=' | '=out='
// arguments  : '(' listValue ')' | value
// value      : int | double | string | date | datetime | boolean
// listValue  : value(','value)*
// int        : [0-9]+
// double     : [0-9]+'.'[0-9]*
// string     : '"'.*'"' | '\''.*'\''
// date       : [0-9]{4}'-'[0-9]{2}'-'\[0-9]{2}
// datetime   : date'T'[0-9]{2}':'[0-9]{2}':'[0-9]{2}('Z' | (('+'|'-')[0-9]{2}(':')?[0-9]{2}))?
// boolean    : 'true' | 'false'
*/

package dto

type FindPagingQuery interface {
	GetTenantId() string
	GetFilter() string
	GetSort() string
	GetPageNum() uint64
	GetPageSize() uint64
	GetIsTotalRows() bool
}

func NewFindPagingQuery() FindPagingQuery {
	query := &findPagingQuery{PageSize: 20}
	return query
}

type findPagingQuery struct {
	TenantId       string
	Filter         string
	Sort           string
	PageNum        uint64
	PageSize       uint64
	IsTotalRecords bool
}

func (q *findPagingQuery) SetTenantId(value string) {
	q.TenantId = value
}

func (q *findPagingQuery) SetFilter(value string) {
	q.Filter = value
}

func (q *findPagingQuery) SetSort(value string) {
	q.Sort = value
}

func (q *findPagingQuery) SetPageNum(value uint64) {
	q.PageNum = value
}

func (q *findPagingQuery) SetPageSize(value uint64) {
	q.PageSize = value
}

func (q *findPagingQuery) SetIsTotalRows(value bool) {
	q.IsTotalRecords = value
}

func (q *findPagingQuery) GetTenantId() string {
	return q.TenantId
}

func (q *findPagingQuery) GetFilter() string {
	return q.Filter
}

func (q *findPagingQuery) GetSort() string {
	return q.Sort
}

func (q *findPagingQuery) GetPageNum() uint64 {
	return q.PageNum
}

func (q *findPagingQuery) GetPageSize() uint64 {
	return q.PageSize
}

func (q *findPagingQuery) GetIsTotalRows() bool {
	return q.IsTotalRecords
}
