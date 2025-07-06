package enum

type EProductType string

const (
	ProductTypeSingle EProductType = "single"
	ProductTypeGroup  EProductType = "group"
)

func (e EProductType) DisplayString() string {
	switch e {
	case ProductTypeSingle:
		return "Satuan"
	default:
		return "Bundle"
	}
}
