package models

type TodoExtend struct {
	Todo
	Tags []*Tag
}

type Tag struct {
	Name string
}
