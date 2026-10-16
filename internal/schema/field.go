package schema

// FieldDefinition describes one JSON field to generate.
type FieldDefinition struct {
	Name string
	Type string
	Min  float64
	Max  float64
}
