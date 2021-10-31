package typedefs

type SourcePosition struct {
	File string
	Line int
	Char int
}

type TaskDefinition struct {
	Pos SourcePosition
	Name string
	NormalizedName string
	Dependencies []string
	Commands []string
}

type MacroDefinition struct {
	Pos SourcePosition
	Name string
	Value string
}
