package request

type ProductCategory struct {
	Name     string `json:"name"`
	ParentId *uint  `json:"parentId"`
}
