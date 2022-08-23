package dto

type OnSuccess[T any] func(data T) error
type OnSuccessList[T any] func(list []T) error
type OnError func(err error) error
type OnIsFond func() error

type FindPagingResult[T any] struct {
	Data        []T    `json:"data"`
	TotalRows   uint64 `json:"totalRows"`
	TotalPages  uint64 `json:"totalPages"`
	PageNum     uint64 `json:"pageNum"`
	PageSize    uint64 `json:"pageSize"`
	Filter      string `json:"filter"`
	Sort        string `json:"sort"`
	Error       error  `json:"-"`
	IsFound     bool   `json:"-"`
	IsTotalRows bool   `json:"isTotalRows"`
}

func NewFindPagingResult[T any](data []T, totalRows uint64, query FindPagingQuery, err error) *FindPagingResult[T] {
	if data != nil && query != nil {
		res := &FindPagingResult[T]{
			Data:        data,
			TotalRows:   0,
			TotalPages:  0,
			PageNum:     query.GetPageNum(),
			PageSize:    query.GetPageSize(),
			Sort:        query.GetSort(),
			Filter:      query.GetFilter(),
			IsFound:     len(data) > 0,
			Error:       err,
			IsTotalRows: query.GetIsTotalRows(),
		}
		if query.GetIsTotalRows() {
			res.TotalRows = totalRows
			res.TotalPages = getTotalPage(totalRows, query.GetPageSize())
		}
	}
	return &FindPagingResult[T]{
		Data:       data,
		TotalRows:  totalRows,
		TotalPages: 0,
		PageNum:    0,
		PageSize:   0,
		Sort:       "",
		Filter:     "",
		IsFound:    false,
		Error:      err,
	}
}

func NewFindPagingResultWithError[T interface{}](err error) *FindPagingResult[T] {
	return &FindPagingResult[T]{
		Data:    nil,
		IsFound: false,
		Error:   err,
	}
}

func getTotalPage(totalRows uint64, pageSize uint64) uint64 {
	if pageSize == 0 {
		return 0
	}
	totalPage := totalRows / pageSize
	if totalRows%pageSize > 1 {
		totalPage++
	}
	return totalPage
}

func (f *FindPagingResult[T]) GetError() error {
	return f.Error
}

func (f *FindPagingResult[T]) GetData() []T {
	return f.Data
}

func (f *FindPagingResult[T]) GetAnyData() any {
	return f.Data
}

func (f *FindPagingResult[T]) GetIsFound() bool {
	return f.IsFound
}

func (f *FindPagingResult[T]) Result() (*FindPagingResult[T], bool, error) {
	return f, f.IsFound, f.Error
}

func (f *FindPagingResult[T]) OnError(onErr OnError) *FindPagingResult[T] {
	if f.Error != nil && onErr != nil {
		f.Error = onErr(f.Error)
	}
	return f
}

func (f *FindPagingResult[T]) OnNotFond(fond OnIsFond) *FindPagingResult[T] {
	if f.Error == nil && !f.IsFound && fond != nil {
		f.Error = fond()
	}
	return f
}

func (f *FindPagingResult[T]) OnSuccess(success OnSuccessList[T]) *FindPagingResult[T] {
	if f.Error == nil && success != nil && f.IsFound {
		f.Error = success(f.Data)
	}
	return f
}

func (f *FindPagingResult[T]) GetTotalRows() uint64 {
	return f.TotalRows
}

func (f *FindPagingResult[T]) GetTotalPages() uint64 {
	return f.TotalPages
}

func (f *FindPagingResult[T]) GetPageNum() uint64 {
	return f.PageNum
}

func (f *FindPagingResult[T]) GetPageSize() uint64 {
	return f.PageSize
}

func (f *FindPagingResult[T]) GetFilter() string {
	return f.Filter
}

func (f *FindPagingResult[T]) GetSort() string {
	return f.Sort
}

type FindPagingResultOptions[T any] struct {
	Data       *[]T
	TotalRows  int64
	TotalPages int64
	PageNum    int64
	PageSize   int64
	Filter     string
	Sort       string
	Error      error
	IsFound    bool
}

func NewFindPagingResultOptions[T any]() *FindPagingResultOptions[T] {
	return &FindPagingResultOptions[T]{}
}

func (f *FindPagingResultOptions[T]) SetData(data *[]T) *FindPagingResultOptions[T] {
	f.Data = data
	return f
}

func (f *FindPagingResultOptions[T]) SetTotalRows(totalRows int64) *FindPagingResultOptions[T] {
	f.TotalRows = totalRows
	return f
}

func (f *FindPagingResultOptions[T]) SetTotalPages(totalPages int64) *FindPagingResultOptions[T] {
	f.TotalPages = totalPages
	return f
}

func (f *FindPagingResultOptions[T]) SetPageNum(pageNum int64) *FindPagingResultOptions[T] {
	f.PageNum = pageNum
	return f
}

func (f *FindPagingResultOptions[T]) SetPageSize(pageSize int64) *FindPagingResultOptions[T] {
	f.PageSize = pageSize
	return f
}

func (f *FindPagingResultOptions[T]) SetFilter(filter string) *FindPagingResultOptions[T] {
	f.Filter = filter
	return f
}

func (f *FindPagingResultOptions[T]) SetSort(sort string) *FindPagingResultOptions[T] {
	f.Sort = sort
	return f
}

func (f *FindPagingResultOptions[T]) SetError(err error) *FindPagingResultOptions[T] {
	f.Error = err
	return f
}

func (f *FindPagingResultOptions[T]) SetIsFound(isFound bool) *FindPagingResultOptions[T] {
	f.IsFound = isFound
	return f
}
