package types

import "api.mijkomp.com/models/enum"

type ProductWithSpecs struct {
	ID            string
	ComponentType enum.ComponentType
	Specs         map[string]string
}
