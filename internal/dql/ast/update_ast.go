package ast

type UpdateQueryAST struct {
	Table      string
	SetFields  map[string]string
	Conditions map[string]string
	Where      map[string]string // ✅ Add WHERE conditions here
}
