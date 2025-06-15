package request

type CompatibilityRule struct {
	SourceComponentTypeCode string `json:"sourceComponentTypeCode"`
	TargetComponentTypeCode string `json:"targetComponentTypeCode"`
	SourceKey               string `json:"sourceKey"`
	TargetKey               string `json:"targetKey"`
	Condition               string `json:"condition"`
	ValueType               string `json:"valueType"`
	ErrorMessage            string `json:"errorMessage"`
	IsActive                bool   `json:"isActive"`
}
