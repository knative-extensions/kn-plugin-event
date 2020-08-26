package event

// Spec holds specification of event to be created
type Spec struct {
	Type   string
	ID     string
	Source string
	Fields []FieldSpec
}

// FieldSpec holds a specification of a event's data field
type FieldSpec struct {
	Path  string
	Value interface{}
}
