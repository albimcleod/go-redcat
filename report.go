package goredcat

import "encoding/json"

//ReportRequest is the struct to request a report from Redcat
type ReportRequest struct {
	Distinct bool     `json:"Distinct,omitempty"`
	Fields   []string `json:"Fields"`
	//Order       []OrderBy   `json:"Order"`
	Limit       int         `json:"Limit"`
	Start       int         `json:"Start"`
	Constraints Constraints `json:"Constraints,omitempty"`
}

//AddField will add a field to the report request
func (obj *ReportRequest) AddField(field string) {
	obj.Fields = append(obj.Fields, field)
}

//AddField will add a field to the report request
func (obj *ReportRequest) AddConstraint(field string, condition string, value string) {
	cv := ConstraintValue{
		Field:     field,
		Condition: condition,
		Value:     value,
	}
	obj.Constraints.Value = append(obj.Constraints.Value, cv)

}

//Constraints is the struct to add constraints to a report from Redcat
type Constraints struct {
	Operator string            `json:"Operator,omitempty"`
	Value    []ConstraintValue `json:"Value",omitempty`
}

//ConstraintValue is ...
type ConstraintValue struct {
	Field     string `json:"Field"`
	Condition string `json:"Condition"`
	Value     string `json:"Value"`
}

type ReportResult struct {
	Axis    []string        `json:"axis"`
	Count   int             `json:"count"`
	Success bool            `json:"success"`
	Data    [][]json.Number `json:"data"`
}
