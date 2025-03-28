package load

type StructInfo struct {
	Name     string
	Fields   []*Field
	Comments []string
}

type Field struct {
	Name            string
	Type            *Type
	DefaultValue    string
	HasDefaultValue bool
	Tags            map[string]string
	IsOverwritten   bool
	IsModifiedCond  string
	MustOverwrite   bool
}

type Type struct {
	Name       string
	IsNillable bool
	IsSlice    bool
}
