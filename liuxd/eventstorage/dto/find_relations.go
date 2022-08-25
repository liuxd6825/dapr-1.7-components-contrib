package dto

type FindRelationsRequest struct {
	TenantId      string `json:"tenantId"`
	AggregateType string `json:"aggregateType"`
	Filter        string `json:"filter"`
	Sort          string `json:"sort"`
	PageNum       uint64 `json:"pageNum"`
	PageSize      uint64 `json:"pageSize"`
	IsTotalRows   bool   `json:"isTotalRows"`
}

func (g *FindRelationsRequest) GetTenantId() string {
	return g.TenantId
}

func (g *FindRelationsRequest) GetFilter() string {
	return g.Filter
}

func (g *FindRelationsRequest) GetSort() string {
	return g.Sort
}

func (g *FindRelationsRequest) GetPageNum() uint64 {
	return g.PageNum
}

func (g *FindRelationsRequest) GetPageSize() uint64 {
	return g.PageSize
}

func (g *FindRelationsRequest) GetIsTotalRows() bool {
	return g.IsTotalRows
}

type FindRelationsResponse struct {
	Data       []*Relation      `json:"data"`
	Headers    *ResponseHeaders `json:"headers"`
	TotalRows  uint64           `json:"totalRows"`
	TotalPages uint64           `json:"totalPages"`
	PageNum    uint64           `json:"pageNum"`
	PageSize   uint64           `json:"pageSize"`
	Filter     string           `json:"filter"`
	Sort       string           `json:"sort"`
	Error      string           `json:"error"`
	IsFound    bool             `json:"isFound"`
}

type Relation struct {
	Id            string `json:"id"`
	TenantId      string `json:"tenantId"`
	TableName     string `json:"tableName"`
	AggregateId   string `json:"aggregateId"`
	AggregateType string `json:"aggregateType"`
	IsDeleted     bool   `json:"isDeleted"`
	RelName       string `json:"relName"`
	RelValue      string `json:"relValue"`
}

/*func (r *Relation) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})
	data["id"] = r.Id
	data["tenantId"] = r.TenantId
	data["tableName"] = r.TableName
	data["aggregateId"] = r.AggregateId
	data["isDeleted"] = r.IsDeleted
	data["aggregateType"] = r.AggregateType
	data["re"]
	for k, v := range r.Items {
		name := utils.AsJsonName(k)
		data[name] = v
	}
	return json.Marshal(data)
}
*/
