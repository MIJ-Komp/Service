package enum

type EProductType string

const (
	ProductTypeSimple  EProductType = "simple"
	ProductTypeVariant EProductType = "variant"
	ProductTypeGroup   EProductType = "bundle"
)

func (e EProductType) DisplayString() string {
	switch e {
	case ProductTypeSimple:
		return "Sederhana"
	case ProductTypeVariant:
		return "Bervariasi"
	default:
		return "Bundle"
	}
}
