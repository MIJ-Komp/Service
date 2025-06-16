package request

type Menu struct {
	Name     string `json:"name"`
	ParentId *uint  `json:"parentId"`
}

type MenuItem struct {
	ProductCategoryId uint `json:"productCategoryId"`
}
