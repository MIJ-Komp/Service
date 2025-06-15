package enum

type ComponentType string

const (
	ComponentMotherboard ComponentType = "Motherboard"
	ComponentCPU         ComponentType = "CPU"
	ComponentRAM         ComponentType = "RAM"
	ComponentVGA         ComponentType = "VGA"
)

type ValueType string

const (
	ValueString  ValueType = "string"
	ValueNumber  ValueType = "number"
	ValueEnum    ValueType = "enum"
	ValueBoolean ValueType = "boolean"
)

type ConditionOperator string

const (
	Equal              ConditionOperator = "="
	NotEqual           ConditionOperator = "!="
	GreaterThan        ConditionOperator = ">"
	LessThan           ConditionOperator = "<"
	GreaterThanOrEqual ConditionOperator = ">="
	LessThanOrEqual    ConditionOperator = "<="
	InList             ConditionOperator = "in"
	NotInList          ConditionOperator = "not in"
)
