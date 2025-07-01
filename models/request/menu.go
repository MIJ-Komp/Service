package request

type Menu struct {
	Name     string  `json:"name"`
	ParentId *uint   `json:"parentId"`
	Path     *string `json:"path"`
}

type MenuItem struct {
	ProductCategoryId uint `json:"productCategoryId"`
}
