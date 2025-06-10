package response

type PageResult struct {
	Items      interface{} `json:"items"`
	TotalCount int64       `json:"totalCount"`
	PageSize   int64       `json:"pageSize"`
}
