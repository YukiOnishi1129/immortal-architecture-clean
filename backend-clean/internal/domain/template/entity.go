package template

import "time"

type Template struct {
	ID        string
	Name      string
	OwnerID   string
	Fields    []Field
	UpdatedAt time.Time
}

type Field struct {
	ID         string
	Label      string
	Order      int
	IsRequired bool
}
