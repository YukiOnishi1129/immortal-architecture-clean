package note

import "time"

type NoteStatus string

const (
	StatusDraft   NoteStatus = "Draft"
	StatusPublish NoteStatus = "Publish"
)

type Note struct {
	ID         string
	Title      string
	TemplateID string
	OwnerID    string
	Status     NoteStatus
	Sections   []Section
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type Section struct {
	ID      string
	NoteID  string
	FieldID string
	Content string
}
